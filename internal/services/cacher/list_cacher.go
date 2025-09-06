package cacher

import (
	"github.com/rpanchyk/javaman/internal/models"
)

type ListCacher interface {
	Get() ([]models.Sdk, error)
	Save(sdks []models.Sdk) error
}
