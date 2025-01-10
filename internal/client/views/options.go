package views

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

type Option func(c *container)

func WithShowMessage(timeout time.Duration) Option {
	return func(c *container) {
		c.uo[getKeyType(messages.ShowMessage{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.ShowMessage)
			c.messages.Store(msg.Message, msg.Message)
			go func() {
				time.Sleep(timeout)
				c.messages.Delete(msg.Message)
			}()
			return c, nil
		}
	}
}

func WithShowError(timeout time.Duration) Option {
	return func(c *container) {
		c.uo[getKeyType(messages.ShowError{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.ShowError)
			c.errors.Store(msg.Err, msg.Err)
			go func() {
				time.Sleep(timeout)
				c.errors.Delete(msg.Err)
			}()
			return c, nil
		}
	}
}

func WithUpdateUser(as *state.AppState) Option {
	return func(c *container) {
		c.uo[getKeyType(messages.UpdateUser{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.UpdateUser)
			as.User = msg.User
			return c, nil
		}
	}
}

func WithAuth(as *state.AppState) Option {
	return func(c *container) {
		c.uo[getKeyType(messages.Auth{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.Auth)
			as.User.MasterPassword = msg.MasterPass
			err := as.Auth()
			if err != nil {
				return c, helpers.GenCmd(messages.ShowError{Err: fmt.Sprintf(err.Error())})
			}
			return c, tea.Batch(helpers.GenCmd(messages.ChangeView{Name: names.MainMenu}))
		}
	}
}

func getKeyType(v tea.Msg) string {
	return fmt.Sprintf("%T", v)
}
