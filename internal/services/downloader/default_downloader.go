package downloader

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/rpanchyk/javaman/internal/clients"
	"github.com/rpanchyk/javaman/internal/models"
	"github.com/rpanchyk/javaman/internal/services/lister"
)

type DefaultDownloader struct {
	config      *models.Config
	listFetcher lister.ListFetcher
	httpSaver   clients.HttpSaver
}

func NewDefaultDownloader(
	config *models.Config,
	listFetcher lister.ListFetcher,
	httpSaver clients.HttpSaver) *DefaultDownloader {

	return &DefaultDownloader{
		config:      config,
		listFetcher: listFetcher,
		httpSaver:   httpSaver,
	}
}

func (d DefaultDownloader) Download(version string) (*models.Sdk, error) {
	sdks, err := d.listFetcher.Fetch()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	sdk, err := d.findSdk(version, sdks)
	if err != nil {
		return nil, fmt.Errorf("cannot find specified SDK: %w", err)
	}
	fmt.Printf("Found SDK: %+v\n", *sdk)

	filePath, err := d.downloadSdk(sdk.URL, d.config.DownloadDir)
	if err != nil {
		return nil, fmt.Errorf("cannot download specified SDK: %w", err)
	}

	sdk.FilePath = filePath
	sdk.IsDownloaded = true

	fmt.Printf("Downloaded SDK: %+v\n", *sdk)
	return sdk, nil
}

func (d DefaultDownloader) findSdk(version string, sdks []models.Sdk) (*models.Sdk, error) {
	r, err := regexp.Compile(`(\w+)-([0-9._\-]+)`)
	if err != nil {
		return nil, fmt.Errorf("error compile regexp: %w", err)
	}

	var sdkVendor string
	var sdkVersion string
	for _, parts := range r.FindAllStringSubmatch(version, -1) {
		sdkVendor = parts[1]
		sdkVersion = parts[2]
	}

	sdkOs := models.Linux
	err = sdkOs.UnmarshalJSON([]byte(runtime.GOOS))
	if err != nil {
		return nil, fmt.Errorf("error resolve os: %w", err)
	}

	sdkArch := models.X64
	err = sdkArch.UnmarshalJSON([]byte(runtime.GOARCH))
	if err != nil {
		return nil, fmt.Errorf("error resolve arch: %w", err)
	}

	for _, sdk := range sdks {
		if sdk.Vendor == sdkVendor && sdk.Version == sdkVersion && sdk.Os == sdkOs && sdk.Arch == sdkArch {
			return &sdk, nil
		}
	}
	return nil, fmt.Errorf("version %s not found", version)
}

func (d DefaultDownloader) downloadSdk(url, dir string) (string, error) {
	filePath := filepath.Join(dir, path.Base(url))
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("SDK %s has been already downloaded\n", filePath)
		return filePath, nil
	}

	if err := d.httpSaver.Save(url, filePath); err != nil {
		return "", fmt.Errorf("cannot save file to: %s error: %w", filePath, err)
	}

	fmt.Printf("SDK %s has been downloaded\n", filePath)
	return filePath, nil
}
