package views

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/asaskevich/govalidator"
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

func WithViews(m Map) Option {
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
				c.syncCnt.Add(1)
				go func() {
					time.Sleep(time.Second)
					c.syncCnt.Add(-1)
				}()
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
		c.uo[helpers.GenMsgKey(messages.ShowMessage{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.ShowMessage)
			c.messages.Store(msg.Message)
			c.messages.DeleteAfter(msg.Message, timeout)
			return c, nil
		}
	}
}

func WithShowError(timeout time.Duration) Option {
	return func(c *container) {
		c.uo[helpers.GenMsgKey(messages.ShowError{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.ShowError)
			c.errors.Store(msg.Err)
			c.errors.DeleteAfter(msg.Err, timeout)
			return c, nil
		}
	}
}

func WithValidityError(timeout time.Duration) Option {
	return func(c *container) {
		c.uo[helpers.GenMsgKey(messages.ValidityError{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.ValidityError)
			var errs govalidator.Errors
			if errors.As(msg.Error, &errs) {
				err := "Ошибка валидации"
				c.errors.Store(err)
				c.errors.DeleteAfter(err, timeout)
				for _, e := range errs {
					c.errors.Store(e.Error())
					c.errors.DeleteAfter(e.Error(), timeout)
				}
				return c, nil
			}
			c.errors.Store(msg.Error.Error())
			c.errors.DeleteAfter(msg.Error.Error(), timeout)
			return c, nil
		}
	}
}

func WithUpdateUser(a contract.UserAccessor) Option {
	return func(c *container) {
		c.uo[helpers.GenMsgKey(messages.UpdateUser{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.UpdateUser)
			a.SetUser(msg.User)
			return c, nil
		}
	}
}

func WithAuth(a contract.UserAccessor) Option {
	return func(c *container) {
		c.uo[helpers.GenMsgKey(messages.Auth{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
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
		c.uo[helpers.GenMsgKey(messages.UploadItemUpdates{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			msg := v.(messages.UploadItemUpdates)
			c.syncCnt.Add(1)
			go func() {
				time.Sleep(time.Second)
				c.syncCnt.Add(-1)
			}()
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
		c.uo[helpers.GenMsgKey(messages.Quit{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
			c.quitting.Store(true)
			handler()
			return c, tea.Quit
		}
	}
}

func WithDownloadFile(fs contract.FileStorage) Option {
	return func(c *container) {
		c.uo[helpers.GenMsgKey(messages.DownloadFile{})] = func(v tea.Msg) (tea.Model, tea.Cmd) {
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

func WithUpdateInterval(i time.Duration) Option {
	return func(c *container) {
		c.updateInterval = i
	}
}
