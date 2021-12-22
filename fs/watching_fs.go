package fs

import (
	"io/fs"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

type Event struct {
	File      string
	Timestamp time.Time
}

func NewWatching(backing FS) (*Watching, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watching{
		watcher: watcher,
		backing: backing,
	}, nil
}

type Watching struct {
	events  chan Event
	watcher *fsnotify.Watcher
	backing FS
}

func (w *Watching) Open(name string) (file fs.File, err error) {
	file, err = w.backing.Open(name)
	if err == nil {
		dir := filepath.Dir(name)
		if watchErr := w.watcher.Add(dir); watchErr != nil {
			log.Errorf("Failed to watch %s: %v", dir, watchErr)
		}
		return file, nil
	}
	return nil, err
}

func (w *Watching) Events() chan Event {
	if w.events == nil {
		w.events = make(chan Event)
		go transportEvents(w.watcher.Events, w.events)
	}
	return w.events
}

func (w *Watching) Close() error {
	return w.watcher.Close()
}

func transportEvents(in <-chan fsnotify.Event, out chan<- Event) {
	for ev := range in {
		ev.Name = filepath.Base(ev.Name)
		out <- Event{
			File:      ev.Name,
			Timestamp: time.Now(),
		}
	}
}
