package helpers

import (
	"sync/atomic"
	"time"

	"github.com/charmbracelet/bubbles/list"
)

type ListRefresher[T list.Item] struct {
	S          Storage[T]
	List       *list.Model
	refreshing atomic.Bool
}

func (l *ListRefresher[T]) Refresh() {
	if l.refreshing.Load() || l.S == nil {
		return
	}
	l.refreshing.Store(true)
	items := l.S.FindAll(nil)
	conv := make([]list.Item, len(items))
	for i := range items {
		conv[i] = items[i]
	}
	l.List.SetItems(conv)
	l.refreshing.Store(false)
}

func (l *ListRefresher[T]) RefreshIn(t time.Duration) {
	go func() {
		time.Sleep(t)
		l.Refresh()
	}()
}
