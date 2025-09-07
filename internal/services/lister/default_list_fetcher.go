package lister

import (
	"encoding/json"
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
	sdks = []models.Sdk{}

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
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:142.0) Gecko/20100101 Firefox/142.0",
		"Accept":     "application/json",
	}
	status, response, err := f.httpClient.GetWithHeaders(f.config.ReleaseURL, headers)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}
	if status < 200 || status > 299 {
		return nil, fmt.Errorf("error making http request, status: %d", status)
	}
	//fmt.Printf("response: %s\n", response)

	sdks, err := f.parsePage(response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	fmt.Printf("SDKs: %+v\n", sdks)

	return sdks, nil
}

func (f DefaultListFetcher) parsePage(body string) ([]models.Sdk, error) {
	sdks := make([]models.Sdk, 0)

	var releases models.Releases
	err := json.Unmarshal([]byte(body), &releases)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	fmt.Printf("releases: %+v\n", releases)

	for _, release := range releases.Data.Releases {
		if release.Family > 8 && release.Status == "DELIVERED" {
			version := release.Version
			if release.Type == "MAJOR" || release.Type == "FEATURE" {
				version = version + ".0.0"
			}

			sdks = append(sdks, models.Sdk{
				Version: version,
			})
		}
	}

	return sdks, nil
}
