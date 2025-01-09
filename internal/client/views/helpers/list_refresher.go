package helpers

import (
	"time"

	"github.com/charmbracelet/bubbles/list"

	"github.com/AndrXxX/goph-keeper/internal/client/interfaces"
)

type ListRefresher[T list.Item] struct {
	S          interfaces.Storage[T]
	List       *list.Model
	refreshing bool
}

func (l *ListRefresher[T]) Refresh() {
	if l.refreshing || l.S == nil {
		return
	}
	l.refreshing = true
	items := l.S.FindAll(nil)
	l.List.SetItems([]list.Item{})
	for i := range items {
		l.List.InsertItem(-1, items[i])
	}
	l.refreshing = false
}

func (l *ListRefresher[T]) RefreshIn(t time.Duration) {
	go func() {
		time.Sleep(t)
		l.Refresh()
	}()
}
