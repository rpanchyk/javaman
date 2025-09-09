package remover

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/rpanchyk/javaman/internal/models"
	"github.com/rpanchyk/javaman/internal/services/lister"
	"github.com/rpanchyk/javaman/internal/utils"
)

type DefaultRemover struct {
	config      *models.Config
	listFetcher lister.ListFetcher
}

func NewDefaultRemover(
	config *models.Config,
	listFetcher lister.ListFetcher) *DefaultRemover {

	return &DefaultRemover{
		config:      config,
		listFetcher: listFetcher,
	}
}

func (r DefaultRemover) Remove(version string, removeDownloaded, removeInstalled bool) error {
	sdks, err := r.listFetcher.Fetch()
	if err != nil {
		return fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	sdk, err := utils.FindByVersion(version, sdks, r.config)
	if err != nil {
		return fmt.Errorf("cannot find specified SDK: %w", err)
	}
	if sdk.IsDefault {
		return fmt.Errorf("cannot remove SDK version %s since it is used as default", version)
	}
	fmt.Printf("Found SDK: %+v\n", *sdk)

	if removeDownloaded {
		if err := r.removeDownloaded(sdk); err != nil {
			return err
		}
	}

	if removeInstalled {
		if err := r.removeInstalled(sdk); err != nil {
			return err
		}
	}

	return nil
}

func (r DefaultRemover) findSdk(version string, sdks []models.Sdk) (*models.Sdk, error) {
	for _, sdk := range sdks {
		if sdk.Version == version {
			return &sdk, nil
		}
	}
	return nil, fmt.Errorf("version %s not found", version)
}

func (r DefaultRemover) removeDownloaded(sdk *models.Sdk) error {
	version := sdk.Vendor + "-" + sdk.Version

	filePath := filepath.Join(r.config.DownloadDir, path.Base(sdk.URL))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("SDK %s version is not downloaded\n", version)
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("cannot remove downloaded archive of SDK %s version: %w", version, err)
	}
	fmt.Printf("File %s has been removed\n", filePath)

	fmt.Printf("Downloaded archive of SDK %s version has been removed\n", version)
	return nil
}

func (r DefaultRemover) removeInstalled(sdk *models.Sdk) error {
	version := sdk.Vendor + "-" + sdk.Version
	installDir := filepath.Join(r.config.InstallDir, version)
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		fmt.Printf("SDK %s version is not installed\n", version)
		return nil
	}

	if err := os.RemoveAll(installDir); err != nil {
		return fmt.Errorf("cannot remove %s: %w", installDir, err)
	}
	fmt.Printf("Directory %s has been removed\n", installDir)

	fmt.Printf("Installation directories of SDK %s version has been removed\n", version)
	return nil
}
