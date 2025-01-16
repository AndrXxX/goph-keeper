package helpers

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type MsgList struct {
	sync.Map
}

func (l *MsgList) DeleteAfter(msg string, timeout time.Duration) {
	go func() {
		time.Sleep(timeout)
		l.Delete(msg)
	}()
}

func (l *MsgList) Join(separator string) string {
	b := strings.Builder{}
	l.Range(func(_, v any) bool {
		b.WriteString(fmt.Sprintf("%s%s", v.(string), separator))
		return true
	})
	return b.String()
}
