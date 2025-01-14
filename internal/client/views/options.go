package views

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/config"
	"github.com/AndrXxX/goph-keeper/internal/client/jobs"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/messages"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/queue"
)

type Option func(c *container)

func WithMap(m Map) Option {
	return func(c *container) {
		c.views = m
	}
}

func WithRepeatableJob(qr contract.QueueRunner, ri time.Duration, job queue.Job) Option {
	return func(c *container) {
		go func() {
			for {
				if c.quitting.Load() {
					return
				}
				time.Sleep(ri)
				if err := qr.AddJob(job); err != nil {
					logger.Log.Error(err.Error())
					return
				}
			}
		}()
	}
}

func WithBuildInfo(cfg *config.Config) Option {
	return func(c *container) {
		c.bi = &contract.BuildInfo{
			Version: cfg.BuildVersion,
			Date:    cfg.BuildDate,
		}
	}
}

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
			return c, helpers.GenCmd(messages.Quit{})
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
			c.quitting.Store(true)
			handler()
			return c, tea.Quit
		}
	}
}

func WithDownloadFile(fs contract.FileStorage) Option {
	return func(c *container) {
		c.uo[getKeyType(messages.DownloadFile{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.DownloadFile)
			src, err := fs.Get(msg.Item.ID)
			if err != nil {
				return c, helpers.GenCmd(messages.ShowError{Err: "Ошибка при скачивании файла: " + err.Error()})
			}
			dst, err := os.OpenFile(path.Join(msg.Path, msg.Item.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return c, helpers.GenCmd(messages.ShowError{Err: "Ошибка при скачивании файла: " + err.Error()})
			}
			_, err = io.Copy(dst, src)
			if err != nil {
				return c, helpers.GenCmd(messages.ShowError{Err: "Ошибка при скачивании файла: " + err.Error()})
			}
			return c, helpers.GenCmd(messages.ShowMessage{Message: "Файл " + msg.Item.Name + " сохранен"})
		}
	}
}

func getKeyType(v tea.Msg) string {
	return fmt.Sprintf("%T-Handler", v)
}
