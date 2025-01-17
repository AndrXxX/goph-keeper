package itemsloader

import (
	"bytes"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums/contenttypes"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

const fetchRoute = "/api/updates/{type}"

type ItemsLoader[T any] struct {
	Sender     requestSender
	URLBuilder urlBuilder
	Fetcher    sliceFetcher[T]
}

func (c *ItemsLoader[T]) Download(itemType string) (statusCode int, l []T) {
	url := c.URLBuilder.Build(fetchRoute, map[string]string{"type": itemType})

	resp, sErr := c.Sender.Get(url, contenttypes.ApplicationJSON)
	if sErr != nil {
		logger.Log.Error("failed to send request", zap.Error(sErr), zap.String("itemType", itemType))
		return resp.StatusCode, l
	}

	l, fErr := c.Fetcher.FetchSlice(resp.Body)
	if fErr != nil {
		logger.Log.Error("failed to fetch slice", zap.Error(fErr), zap.String("itemType", itemType))
		return resp.StatusCode, l
	}
	return resp.StatusCode, l
}

func (c *ItemsLoader[T]) Upload(itemType string, list []T) (statusCode int, err error) {
	url := c.URLBuilder.Build(fetchRoute, map[string]string{"type": itemType})

	data, mErr := json.Marshal(list)
	if mErr != nil {
		return 0, fmt.Errorf("marshal data with itemType (%s) %w", itemType, mErr)
	}

	resp, sErr := c.Sender.Post(url, contenttypes.ApplicationJSON, bytes.NewBuffer(data))
	if sErr != nil {
		return 0, fmt.Errorf("post request with itemType (%s) %w", itemType, sErr)
	}
	if resp != nil {
		if resp.Body != nil {
			return resp.StatusCode, resp.Body.Close()
		}
		return resp.StatusCode, nil
	}
	return 0, nil
}
