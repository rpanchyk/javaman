package lister

import (
	"github.com/rpanchyk/javaman/internal/models"
)

type ListFetcher interface {
	Fetch() ([]models.Sdk, error)
}
