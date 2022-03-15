package events

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/baez90/goveal/fs"
)

const (
	baseDecimal = 10
)

type (
	EventMutation interface {
		OnEvent(in ContentEvent, ev fs.Event) (out ContentEvent)
	}
	ContentEvent struct {
		File         string `json:"file"`
		FileNameHash string `json:"fileNameHash"`
		Timestamp    string `json:"ts"`
		ForceReload  bool   `json:"forceReload"`
		ReloadConfig bool   `json:"reloadConfig"`
	}
	EventSource interface {
		io.Closer
		fs.FS
		Events() chan fs.Event
	}
	EventHandler interface {
		OnEvent(ev ContentEvent) error
		Close() error
	}
	EventHandlerFunc            func(ev ContentEvent) error
	MutationReloadForFile       string
	MutationConfigReloadForFile string
)

func (t MutationReloadForFile) OnEvent(in ContentEvent, ev fs.Event) (out ContentEvent) {
	if strings.EqualFold(filepath.Base(ev.File), filepath.Base(string(t))) {
		in.ForceReload = true
	}
	return in
}

func (t MutationConfigReloadForFile) OnEvent(in ContentEvent, ev fs.Event) (out ContentEvent) {
	if strings.EqualFold(filepath.Base(ev.File), string(t)) {
		in.ReloadConfig = true
	}
	return in
}

func (f EventHandlerFunc) OnEvent(ev ContentEvent) error {
	return f(ev)
}

func NewEventHub(logger *log.Logger, eventSource EventSource, fileNameHash hash.Hash, mutations ...EventMutation) *EventHub {
	hub := &EventHub{
		Logger:        logger,
		FileNameHash:  fileNameHash,
		mutations:     mutations,
		source:        eventSource,
		subscriptions: make(map[uuid.UUID]EventHandler),
		done:          make(chan struct{}),
	}

	go hub.processEvents()

	return hub
}

type EventHub struct {
	FileNameHash  hash.Hash
	Logger        *log.Logger
	mutations     []EventMutation
	lock          sync.RWMutex
	done          chan struct{}
	source        EventSource
	subscriptions map[uuid.UUID]EventHandler
}

func (h *EventHub) Subscribe(handler EventHandler) (id uuid.UUID) {
	h.lock.Lock()
	defer h.lock.Unlock()

	clientID := uuid.New()
	h.subscriptions[clientID] = handler

	return clientID
}

func (h *EventHub) Unsubscribe(id uuid.UUID) (err error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if handler, ok := h.subscriptions[id]; ok {
		err = handler.Close()
		delete(h.subscriptions, id)
		return err
	}

	return nil
}

func (h *EventHub) Close() error {
	close(h.done)
	return h.source.Close()
}

func (h *EventHub) processEvents() {
	events := h.source.Events()
	for {
		select {
		case ev, more := <-events:
			if !more {
				return
			}
			h.notifySubscribers(ev)
		case _, more := <-h.done:
			if !more {
				return
			}
		}
	}
}

func (h *EventHub) notifySubscribers(ev fs.Event) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	ce := ContentEvent{
		File:         fmt.Sprintf("/%s", ev.File),
		Timestamp:    strconv.FormatInt(ev.Timestamp.Unix(), baseDecimal),
		FileNameHash: hex.EncodeToString(h.FileNameHash.Sum([]byte(path.Base(ev.File)))),
	}

	for idx := range h.mutations {
		ce = h.mutations[idx].OnEvent(ce, ev)
	}

	for _, handler := range h.subscriptions {
		if err := handler.OnEvent(ce); err != nil {
			h.Logger.Errorf("Failed to propagate event: %v", err)
		}
	}
}
