package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
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

func (h ContentEventHandler) Close() (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("failed to close event handler: %v", rec)
		}
	}()

	close(h)
	return
}

type Events struct {
	logger *log.Logger
	hub    *events.EventHub
}

func RegisterEventsAPI(router *httprouter.Router, hub *events.EventHub, logger *log.Logger) {
	ev := &Events{hub: hub, logger: logger}
	router.GET("/api/v1/events", ev.EventHandler)

}

func (e *Events) EventHandler(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Transfer-Encoding", "chunked")
	writer.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
	writer.Header().Set("Access-Control-Allow-Credentials", "true")

	var (
		handler  = make(ContentEventHandler)
		clientID = e.hub.Subscribe(handler)
		buf      = new(bytes.Buffer)
		enc      = json.NewEncoder(buf)
	)

	defer func() {
		if err := e.hub.Unsubscribe(clientID); err != nil {
			e.logger.Warnf("Error occurred while unsubscribing: %v", err)
		}
	}()

	for {
		select {
		case ev := <-handler:
			if err := enc.Encode(ev); err != nil {
				e.logger.Errorf("Failed to marshal to JSON: %v", err)
				continue
			} else if _, err = fmt.Fprintf(writer, "data: %s\n\n", buf.String()); err != nil {
				e.logger.Errorf("Failed to write to client: %v", err)
				return
			} else if f, ok := writer.(http.Flusher); !ok {
				e.logger.Errorf("Cannot flush data")
				writer.WriteHeader(http.StatusBadRequest)
				return
			} else {
				f.Flush()
			}
		case <-req.Context().Done():
			return
		}
	}
}
