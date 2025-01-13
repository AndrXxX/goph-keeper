package itemsloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/enums/contenttypes"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

const (
	sendfileUpdateUrl = "/api/files/update"        // POST
	uploadFileUrl     = "/api/files/upload/{id}"   // POST
	downloadFileUrl   = "/api/files/download/{id}" // GET
)

type FilesLoader struct {
	Sender     requestSender
	URLBuilder urlBuilder
	Fetcher    sliceFetcher[entities.FileItem]
	FS         fileStorage
}

func (c *FilesLoader) Download() (statusCode int, l []entities.FileItem) {
	url := c.URLBuilder.Build(fetchRoute, map[string]string{"type": datatypes.Files})

	resp, sErr := c.Sender.Get(url, contenttypes.ApplicationJSON)
	if sErr != nil {
		logger.Log.Error("failed to send request", zap.Error(sErr))
		return resp.StatusCode, l
	}
	if resp.StatusCode != http.StatusOK {
		logger.Log.Error("failed to send request", zap.Int("status_code", resp.StatusCode))
	}
	l, fErr := c.Fetcher.FetchSlice(resp.Body)
	if fErr != nil {
		logger.Log.Error("failed to fetch slice", zap.Error(fErr))
		return resp.StatusCode, nil
	}
	var res []entities.FileItem
	for i := range l {
		if c.FS.IsExist(l[i].ID) {
			res = append(res, l[i])
			continue
		}
		if err := c.DownloadFile(l[i]); err != nil {
			logger.Log.Error("failed to download file", zap.Error(err))
			continue
		}
	}
	return resp.StatusCode, res
}

func (c *FilesLoader) DownloadFile(item entities.FileItem) error {
	url := c.URLBuilder.Build(downloadFileUrl, map[string]string{"id": item.ID.String()})
	resp, sErr := c.Sender.Get(url, contenttypes.OctetStream)
	if sErr != nil {
		return fmt.Errorf("DownloadFile request %w", sErr)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DownloadFile request status %s", resp.Status)
	}
	err := c.FS.Store(resp.Body, item.ID)
	if err != nil {
		return fmt.Errorf("store file after download: %w", err)
	}
	return nil
}

func (c *FilesLoader) Upload(list []entities.FileItem) (statusCode int, err error) {
	var res []entities.FileItem
	for i := range list {
		if list[i].TempPath == "" {
			return 0, fmt.Errorf("file not selected")
		}
		f, err := os.OpenFile(list[i].TempPath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return 0, fmt.Errorf("failed to open file: %w", err)
		}
		data, mErr := json.Marshal(list[i])
		if mErr != nil {
			return 0, fmt.Errorf("marshal data %w", mErr)
		}

		url := c.URLBuilder.Build(sendfileUpdateUrl, map[string]string{})
		resp, sErr := c.Sender.Post(url, contenttypes.ApplicationJSON, bytes.NewReader(data))
		if sErr != nil {
			return 0, fmt.Errorf("post request %w", sErr)
		}
		if resp.StatusCode != http.StatusOK {
			return resp.StatusCode, fmt.Errorf("upload file request status %s", resp.Status)
		}
		rawId, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return resp.StatusCode, fmt.Errorf("read response %w", sErr)
		}
		id, err := uuid.ParseBytes(rawId)
		if err != nil {
			return resp.StatusCode, fmt.Errorf("uuid parse %w", sErr)
		}
		list[i].ID = id
		err = c.FS.Store(f, id)
		if err != nil {
			return 0, fmt.Errorf("store file after upload: %w", err)
		}
		res = append(res, list[i])
	}
	return http.StatusOK, nil
}

func (c *FilesLoader) UploadFile(item entities.FileItem, data io.Reader) error {
	url := c.URLBuilder.Build(uploadFileUrl, map[string]string{"id": item.ID.String()})
	resp, sErr := c.Sender.Post(url, contenttypes.OctetStream, data)
	if sErr != nil {
		return fmt.Errorf("UploadFile request %w", sErr)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("UploadFile request status %s", resp.Status)
	}
	return nil
}
