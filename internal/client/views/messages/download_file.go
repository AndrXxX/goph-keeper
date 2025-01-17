package messages

import "github.com/AndrXxX/goph-keeper/internal/client/entities"

type DownloadFile struct {
	Item *entities.FileItem
	Path string
}
