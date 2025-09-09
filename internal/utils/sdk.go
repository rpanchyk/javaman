package utils

import (
	"fmt"
	"regexp"
	"runtime"

	"github.com/rpanchyk/javaman/internal/models"
)

func FindByVersion(version string, sdks []models.Sdk) (*models.Sdk, error) {
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
