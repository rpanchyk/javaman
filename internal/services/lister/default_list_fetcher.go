package lister

import (
	"fmt"

	"github.com/rpanchyk/javaman/internal/clients"
	"github.com/rpanchyk/javaman/internal/models"
	"github.com/rpanchyk/javaman/internal/services/cacher"
)

type DefaultListFetcher struct {
	config     *models.Config
	httpClient clients.HttpClient
	listCacher cacher.ListCacher
}

func NewDefaultListFetcher(
	config *models.Config,
	httpClient clients.HttpClient,
	listCacher cacher.ListCacher) *DefaultListFetcher {

	return &DefaultListFetcher{
		config:     config,
		httpClient: httpClient,
		listCacher: listCacher,
	}
}

func (f DefaultListFetcher) Fetch() ([]models.Sdk, error) {
	sdks, err := f.listCacher.Get()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of SDKs from cache: %w", err)
	}

	if len(sdks) == 0 {
		sdks, err = f.downloadSdks()
		if err != nil {
			return nil, fmt.Errorf("cannot download list of SDKs: %w", err)
		}
		if err := f.listCacher.Save(sdks); err != nil {
			return nil, fmt.Errorf("cannot save list of SDKs to cache: %w", err)
		}
	}

	return sdks, nil
}

func (f DefaultListFetcher) downloadSdks() ([]models.Sdk, error) {
	response, err := f.httpClient.Get(f.config.ReleaseURL)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}

	sdks, err := f.parsePage(response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return sdks, nil
}

func (f DefaultListFetcher) parsePage(body string) ([]models.Sdk, error) {
	sdks := make([]models.Sdk, 0)

	// TODO: parse

	return sdks, nil
}
