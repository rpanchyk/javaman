package lister

import (
	"fmt"

	"github.com/rpanchyk/javaman/internal/models"
	"github.com/rpanchyk/javaman/internal/services/cacher"
)

type DefaultListFetcher struct {
	fetchers   []ListFetcher
	listCacher cacher.ListCacher
}

func NewDefaultListFetcher(
	fetchers []ListFetcher,
	listCacher cacher.ListCacher) *DefaultListFetcher {

	return &DefaultListFetcher{
		fetchers:   fetchers,
		listCacher: listCacher,
	}
}

func (f DefaultListFetcher) Fetch() ([]models.Sdk, error) {
	sdks, err := f.listCacher.Get()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of SDKs from cache: %w", err)
	}

	if len(sdks) == 0 {
		sdks, err = f.fetchSdks()
		if err != nil {
			return nil, fmt.Errorf("cannot fetch list of SDKs: %w", err)
		}
		if err := f.listCacher.Save(sdks); err != nil {
			return nil, fmt.Errorf("cannot save list of SDKs to cache: %w", err)
		}
	}

	return sdks, nil
}

func (f DefaultListFetcher) fetchSdks() ([]models.Sdk, error) {
	sdks := make([]models.Sdk, 0)

	for _, fetcher := range f.fetchers {
		vendorSdks, err := fetcher.Fetch()
		if err != nil {
			return nil, fmt.Errorf("error getting sdks: %w", err)
		}
		sdks = append(sdks, vendorSdks...)
	}

	//fmt.Printf("SDKs: %+v\n", sdks)
	return sdks, nil
}
