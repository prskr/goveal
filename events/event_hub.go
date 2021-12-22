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

	"github.com/baez90/goveal/fs"
)

const (
	baseDecimal = 10
)

type (
	ReloadTrigger interface {
		Triggers(ev fs.Event) bool
	}
	ContentEvent struct {
		File         string `json:"file"`
		FileNameHash string `json:"fileNameHash"`
		Timestamp    string `json:"ts"`
		ForceReload  bool   `json:"forceReload"`
	}
	EventSource interface {
		io.Closer
		fs.FS
		Events() chan fs.Event
	}
	EventHandler interface {
		OnEvent(ev ContentEvent) error
	}
	EventHandlerFunc func(ev ContentEvent) error
	subscription     struct {
		EventHandler
		OnError chan error
	}
	FileNameTrigger   string
	FileSuffixTrigger string
)

func (t FileNameTrigger) Triggers(ev fs.Event) bool {
	fileBase := filepath.Base(ev.File)
	return strings.EqualFold(fileBase, string(t))
}

func (t FileSuffixTrigger) Triggers(ev fs.Event) bool {
	return strings.HasSuffix(strings.ToLower(filepath.Base(ev.File)), strings.ToLower(string(t)))
}

func (f EventHandlerFunc) OnEvent(ev ContentEvent) error {
	return f(ev)
}

func NewEventHub(eventSource EventSource, fileNameHash hash.Hash, triggers ...ReloadTrigger) *EventHub {
	hub := &EventHub{
		FileNameHash:   fileNameHash,
		reloadTriggers: triggers,
		source:         eventSource,
		subscriptions:  make(map[uuid.UUID]*subscription),
		done:           make(chan struct{}),
	}

	go hub.processEvents()

	return hub
}

type EventHub struct {
	FileNameHash   hash.Hash
	reloadTriggers []ReloadTrigger
	lock           sync.RWMutex
	done           chan struct{}
	source         EventSource
	subscriptions  map[uuid.UUID]*subscription
}

func (h *EventHub) Subscribe(handler EventHandler) (id uuid.UUID, onError <-chan error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	s := &subscription{
		EventHandler: handler,
		OnError:      make(chan error),
	}
	clientID := uuid.New()
	h.subscriptions[clientID] = s

	return clientID, s.OnError
}

func (h *EventHub) Unsubscribe(id uuid.UUID) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.subscriptions, id)
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

	var triggerReload bool
	for idx := range h.reloadTriggers {
		if triggerReload = h.reloadTriggers[idx].Triggers(ev); triggerReload {
			break
		}
	}

	ce := ContentEvent{
		File:         fmt.Sprintf("/%s", ev.File),
		Timestamp:    strconv.FormatInt(ev.Timestamp.Unix(), baseDecimal),
		ForceReload:  triggerReload,
		FileNameHash: hex.EncodeToString(h.FileNameHash.Sum([]byte(path.Base(ev.File)))),
	}

	for _, handler := range h.subscriptions {
		if err := handler.OnEvent(ce); err != nil {
			handler.OnError <- err
		}
	}
}
