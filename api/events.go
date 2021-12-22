package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/baez90/goveal/events"
)

type ContentEventHandler chan events.ContentEvent

func (h ContentEventHandler) OnEvent(ce events.ContentEvent) error {
	const enqueueTimeout = 50 * time.Millisecond
	select {
	case h <- ce:
		return nil
	case <-time.After(enqueueTimeout):
		return errors.New("failed to enqueue due to timeout")
	}
}

type Events struct {
	logger *log.Logger
	hub    *events.EventHub
}

func RegisterEventsAPI(app *fiber.App, hub *events.EventHub, logger *log.Logger) {
	ev := &Events{hub: hub, logger: logger}
	app.Get("/api/v1/events", ev.EventHandler)
}

func (e *Events) EventHandler(fc *fiber.Ctx) error {
	var (
		ctx               = fc.Context()
		handler           = make(ContentEventHandler)
		clientID, onError = e.hub.Subscribe(handler)
	)

	ctx.SetContentType("text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		for ev := range handler {
			if msg, err := json.Marshal(ev); err != nil {
				e.logger.Errorf("Failed to marshal to JSON: %v", err)
				continue
			} else if _, err = fmt.Fprintf(w, "data: %s\n\n", string(msg)); err != nil {
				e.logger.Errorf("Failed to write to client: %v", err)
				continue
			} else if err = w.Flush(); err != nil {
				e.hub.Unsubscribe(clientID)
			}
		}
	})

	go func() {
		for err := range onError {
			e.logger.Errorf("Error while sending events to client: %v", err)
		}
	}()

	return nil
}
