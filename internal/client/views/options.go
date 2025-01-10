package views

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/messages"
)

type Option func(c *container)

func WithShowMessage() Option {
	return func(c *container) {
		key := fmt.Sprintf("%T", messages.ShowMessage{})
		c.uo[key] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.ShowMessage)
			c.messages.Store(msg.Message, msg.Message)
			go func() {
				time.Sleep(errorsTimeout)
				c.messages.Delete(msg.Message)
			}()
			return c, nil
		}
	}
}
