package defaulter

import (
	"fmt"

	"github.com/rpanchyk/javaman/internal/models"
	"github.com/rpanchyk/javaman/internal/services/lister"
	"github.com/rpanchyk/javaman/internal/utils"
)

type DefaultDefaulter struct {
	config      *models.Config
	listFetcher lister.ListFetcher
}

func NewDefaultDefaulter(
	config *models.Config,
	listFetcher lister.ListFetcher) *DefaultDefaulter {

	return &DefaultDefaulter{
		config:      config,
		listFetcher: listFetcher,
	}
}

func (d DefaultDefaulter) Default(version string) error {
	sdks, err := d.listFetcher.Fetch()
	if err != nil {
		return fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	sdk, err := utils.FindByVersion(version, sdks)
	if err != nil {
		return fmt.Errorf("cannot find specified SDK: %w", err)
	}

	if !sdk.IsInstalled {
		return fmt.Errorf("SDK version %s is not installed", version)
	}

	if sdk.IsDefault {
		fmt.Printf("SDK version %s is already used as default\n", version)
		return nil
	}

	platformDefaulter := &PlatformDefaulter{Config: d.config}
	return platformDefaulter.Default(version)
}
