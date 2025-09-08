package lister

import (
	"fmt"
	"regexp"
	"runtime"
	"slices"
	"strings"

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
	sdks = []models.Sdk{} // TODO: remove this line

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
	sdks := []models.Sdk{}
	for _, vendor := range f.config.Vendors {
		if strings.ToLower(strings.TrimSpace(vendor)) == "microsoft" {
			microsoftSdks, err := f.getMicrosoftSdks()
			if err != nil {
				return nil, fmt.Errorf("error getting microsoft sdks: %w", err)
			}
			sdks = append(sdks, microsoftSdks...)
		}
	}

	//headers := map[string]string{
	//	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:142.0) Gecko/20100101 Firefox/142.0",
	//	"Accept":     "application/json",
	//}
	//status, response, err := f.httpClient.GetWithHeaders(f.config.ReleaseURL, headers)
	//if err != nil {
	//	return nil, fmt.Errorf("error making http request: %w", err)
	//}
	//if status < 200 || status > 299 {
	//	return nil, fmt.Errorf("error making http request, status: %d", status)
	//}
	////fmt.Printf("response: %s\n", response)
	//
	//sdks, err := f.parsePage(response)
	//if err != nil {
	//	return nil, fmt.Errorf("error parsing response: %w", err)
	//}
	//fmt.Printf("SDKs: %+v\n", sdks)

	return sdks, nil
}

//func (f DefaultListFetcher) parsePage(body string) ([]models.Sdk, error) {
//	sdks := make([]models.Sdk, 0)
//
//	var releases models.Releases
//	err := json.Unmarshal([]byte(body), &releases)
//	if err != nil {
//		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
//	}
//	fmt.Printf("releases: %+v\n", releases)
//
//	for _, release := range releases.Data.Releases {
//		if release.Family > 8 && release.Status == "DELIVERED" {
//			version := release.Version
//			if release.Type == "MAJOR" || release.Type == "FEATURE" {
//				version = version + ".0.0"
//			}
//
//			sdks = append(sdks, models.Sdk{
//				Version: version,
//			})
//		}
//	}
//
//	return sdks, nil
//}

func (f DefaultListFetcher) getMicrosoftSdks() ([]models.Sdk, error) {
	urls := []string{
		"https://learn.microsoft.com/en-us/java/openjdk/download",
		"https://learn.microsoft.com/en-us/java/openjdk/older-releases",
	}

	sdks := make([]models.Sdk, 0)
	for _, url := range urls {
		status, response, err := f.httpClient.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error making http request: %w", err)
		}
		if status < 200 || status > 299 {
			return nil, fmt.Errorf("error making http request, status: %d", status)
		}

		// <a href="https://aka.ms/download-jdk/microsoft-jdk-11.0.14.1_1-31205-linux-x64.tar.gz" data-linktype="external">microsoft-jdk-11.0.14.1_1-31205-linux-x64.tar.gz</a>
		// Group 1: https://aka.ms/download-jdk/microsoft-jdk-11.0.14.1_1-31205-linux-x64.tar.gz
		// Group 2: 11.0.14.1_1-31205
		// Group 3: linux
		// Group 4: x64

		// <a href="https://aka.ms/download-jdk/microsoft-jdk-21.0.8-macos-x64.tar.gz" data-linktype="external">microsoft-jdk-21.0.8-macos-x64.tar.gz</a>
		// Group 1: https://aka.ms/download-jdk/microsoft-jdk-21.0.8-macos-x64.tar.gz
		// Group 2: 21.0.8
		// Group 3: macos
		// Group 4: x64

		// <a href="https://aka.ms/download-jdk/microsoft-jdk-21.0.8-windows-x64.zip" data-linktype="external">microsoft-jdk-21.0.8-windows-x64.zip</a>
		// Group 1: https://aka.ms/download-jdk/microsoft-jdk-21.0.8-windows-x64.zip
		// Group 2: 21.0.8
		// Group 3: windows
		// Group 4: x64

		r, err := regexp.Compile(`href=['"](.*/microsoft-jdk-([0-9._\-]+)-(\w+)-(\w+)\.(?:tar\.gz|zip)+)['"]`)
		if err != nil {
			return nil, fmt.Errorf("error compile regexp: %w", err)
		}

		for _, parts := range r.FindAllStringSubmatch(response, -1) {

			var os models.Os
			switch strings.ToLower(parts[3]) {
			case "linux":
				os = models.Linux
			case "macos":
				os = models.Macos
			case "windows":
				os = models.Windows
			default:
				continue
			}

			var arch models.Arch
			switch strings.ToLower(parts[4]) {
			case "x64":
				arch = models.X64
			case "aarch64":
				arch = models.ARM
			default:
				continue
			}

			sdk := models.Sdk{
				Vendor:  "microsoft",
				URL:     parts[1],
				Version: parts[2],
				Os:      os,
				Arch:    arch,
			}

			if f.config.ListFilterOs && runtime.GOOS != sdk.Os.GoOs() {
				continue
			}
			if f.config.ListFilterArch && runtime.GOARCH != sdk.Arch.GoArch() {
				continue
			}

			if !slices.Contains(sdks, sdk) {
				sdks = append(sdks, sdk)
			}
		}
	}

	return sdks, nil
}
