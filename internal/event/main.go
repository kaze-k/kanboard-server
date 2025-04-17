package event

import (
	"server/internal/constant"
	"sync"
)

type Event struct {
	EventType *constant.EventType
	ProjectID *uint
	TaskID    *uint
	UserID    *uint
	Content   *string
}

type (
	Handler func(Event)
)

var (
	AdminHandlers    []Handler
	KanboardHandlers []Handler
	mu               sync.Mutex
)

func AdminSubscribe(handler Handler) {
	mu.Lock()
	defer mu.Unlock()
	AdminHandlers = append(AdminHandlers, handler)
}

func AdminPublish(event Event) {
	mu.Lock()
	defer mu.Unlock()
	for _, handler := range AdminHandlers {
		go handler(event)
	}
}

func KanboardSubscribe(handler Handler) {
	mu.Lock()
	defer mu.Unlock()
	KanboardHandlers = append(KanboardHandlers, handler)
}

func KanboardPublish(event Event) {
	mu.Lock()
	defer mu.Unlock()
	for _, handler := range KanboardHandlers {
		go handler(event)
	}
}
