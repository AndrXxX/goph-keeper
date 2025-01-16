package helpers

import (
	"slices"
	"strings"
	"sync"
	"time"
)

type MsgList struct {
	vals []string
	m    sync.Mutex
}

func (l *MsgList) DeleteAfter(val string, timeout time.Duration) {
	go func() {
		time.Sleep(timeout)
		l.m.Lock()
		defer l.m.Unlock()
		pos := slices.Index(l.vals, val)
		if pos != -1 {
			l.vals = append(l.vals[:pos], l.vals[pos+1:]...)
		}
	}()
}

func (l *MsgList) Join(separator string) string {
	l.m.Lock()
	defer l.m.Unlock()
	return strings.Join(l.vals, separator)
}

func (l *MsgList) Store(v string) {
	l.m.Lock()
	l.vals = append(l.vals, v)
	l.m.Unlock()
}
