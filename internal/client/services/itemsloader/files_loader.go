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
	sendfileUpdateUrl = "/api/files/update"         // POST
	uploadFileUrl     = "/api/files/upload/{id}/"   // POST
	downloadFileUrl   = "/api/files/download/{id}/" // GET
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
		if err := c.downloadFile(l[i]); err != nil {
			logger.Log.Error("failed to download file", zap.Error(err))
			continue
		}
	}
	return resp.StatusCode, res
}

func (c *FilesLoader) downloadFile(item entities.FileItem) error {
	url := c.URLBuilder.Build(downloadFileUrl, map[string]string{"id": item.ID.String()})
	resp, sErr := c.Sender.Get(url, contenttypes.OctetStream)
	if sErr != nil {
		return fmt.Errorf("download file request %w", sErr)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download file request status %s", resp.Status)
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
			return 0, fmt.Errorf("open file: %w", err)
		}
		resp, uErr := c.uploadFileUpdate(list[i])
		if uErr != nil {
			return 0, fmt.Errorf("upload file update %w", uErr)
		}
		id, err := c.fetchId(resp.Body)
		if err != nil {
			return resp.StatusCode, fmt.Errorf("fetch id %w", err)
		}
		list[i].ID = id
		err = c.uploadFile(list[i], f)
		if err != nil {
			return 0, fmt.Errorf("upload file after upload data: %w", err)
		}
		res = append(res, list[i])
	}
	return http.StatusOK, nil
}

func (c *FilesLoader) uploadFileUpdate(item entities.FileItem) (*http.Response, error) {
	data, mErr := json.Marshal(item)
	if mErr != nil {
		return nil, fmt.Errorf("marshal data %w", mErr)
	}
	url := c.URLBuilder.Build(sendfileUpdateUrl, map[string]string{})
	resp, err := c.Sender.Post(url, contenttypes.ApplicationJSON, bytes.NewReader(data))
	if resp != nil && resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("upload file request status %s", resp.Status)
	}
	return resp, err
}

func (c *FilesLoader) fetchId(data io.ReadCloser) (uuid.UUID, error) {
	rawId, rErr := io.ReadAll(data)
	defer func(data io.ReadCloser) {
		_ = data.Close()
	}(data)
	if rErr != nil {
		return uuid.UUID{}, fmt.Errorf("read response %w", rErr)
	}
	id, pErr := uuid.ParseBytes(rawId)
	if pErr != nil {
		return uuid.UUID{}, fmt.Errorf("uuid parse %w", pErr)
	}
	return id, nil
}

func (c *FilesLoader) uploadFile(item entities.FileItem, data io.Reader) error {
	url := c.URLBuilder.Build(uploadFileUrl, map[string]string{"id": item.ID.String()})
	resp, sErr := c.Sender.Post(url, contenttypes.OctetStream, data)
	if sErr != nil {
		return fmt.Errorf("upload file request %w", sErr)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload file request status %s", resp.Status)
	}
	return nil
}
