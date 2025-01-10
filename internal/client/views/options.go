package views

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/jobs"
	"github.com/AndrXxX/goph-keeper/internal/client/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

type Option func(c *container)

func WithStartView(view names.ViewName) Option {
	return func(c *container) {
		c.current = view
	}
}

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

func WithUpdateUser(a contract.UserAccessor) Option {
	return func(c *container) {
		c.uo[getKeyType(messages.UpdateUser{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.UpdateUser)
			a.SetUser(msg.User)
			return c, nil
		}
	}
}

func WithAuth(a contract.UserAccessor) Option {
	return func(c *container) {
		c.uo[getKeyType(messages.Auth{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.Auth)
			a.SetMasterPass(msg.MasterPass)
			err := a.Auth()
			if err != nil {
				return c, helpers.GenCmd(messages.ShowError{Err: fmt.Sprintf(err.Error())})
			}
			return c, tea.Batch(helpers.GenCmd(messages.ChangeView{Name: names.MainMenu}))
		}
	}
}

func WithUploadItemUpdates(sm contract.SyncManager, qr contract.QueueRunner) Option {
	return func(c *container) {
		c.uo[getKeyType(messages.UploadItemUpdates{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.UploadItemUpdates)
			err := qr.AddJob(&jobs.UploadItemsUpdatesJob{
				Type:        msg.Type,
				Items:       msg.Items,
				SyncManager: sm,
			})
			if err != nil {
				return c, helpers.GenCmd(messages.ShowError{Err: fmt.Sprintf("Ошибка при синхронизации: %s", err)})
			}
			return c, nil
		}
	}
}

func WithQuit(handler func()) Option {
	return func(c *container) {
		c.uo[getKeyType(messages.Quit{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			c.quitting = true
			handler()
			return c, tea.Quit
		}
	}
}

func getKeyType(v tea.Msg) string {
	return fmt.Sprintf("%T", v)
}
