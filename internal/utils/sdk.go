package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/rpanchyk/javaman/internal/models"
)

func FindByVersion(version string, sdks []models.Sdk, config *models.Config) (*models.Sdk, error) {
	r, err := regexp.Compile(`([\w\-]+)-([0-9._\-]+)`)
	if err != nil {
		return nil, fmt.Errorf("error compile regexp: %w", err)
	}

	var sdkVendor string
	var sdkVersion string
	for _, parts := range r.FindAllStringSubmatch(version, -1) {
		sdkVendor = parts[1]
		sdkVersion = parts[2]
	}

	currentOs := CurrentOs()
	currentArch := CurrentArch()

	for _, sdk := range sdks {
		if sdk.Vendor == sdkVendor && sdk.Version == sdkVersion && sdk.Os == currentOs && sdk.Arch == currentArch {
			return enrichSdk(&sdk, config)
		}
	}
	return nil, fmt.Errorf("version %s not found", version)
}

func enrichSdk(sdk *models.Sdk, config *models.Config) (*models.Sdk, error) {
	filePath := filepath.Join(config.DownloadDir, path.Base(sdk.URL))
	if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
		sdk.IsDownloaded = true
		sdk.FilePath = filePath
	}

	installDir := filepath.Join(config.InstallDir, sdk.Vendor+"-"+sdk.Version)
	if info, err := os.Stat(installDir); err == nil && info.IsDir() {
		sdk.IsInstalled = true
	}

	if envVar, ok := os.LookupEnv("JAVA_HOME"); ok && envVar == installDir {
		sdk.IsDefault = true
	}

	return sdk, nil
}
