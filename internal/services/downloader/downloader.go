package downloader

import (
	"github.com/rpanchyk/javaman/internal/models"
)

type Downloader interface {
	Download(version string) (*models.Sdk, error)
}
