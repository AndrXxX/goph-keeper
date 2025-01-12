package helpers

import (
	"sync/atomic"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"golang.org/x/tools/container/intsets"
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
	l.List.SetItems([]list.Item{})
	for i := range items {
		l.List.InsertItem(intsets.MaxInt, items[i])
	}
	l.refreshing.Store(false)
}

func (l *ListRefresher[T]) RefreshIn(t time.Duration) {
	go func() {
		time.Sleep(t)
		l.Refresh()
	}()
}
