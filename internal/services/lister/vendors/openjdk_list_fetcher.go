package vendors

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/rpanchyk/javaman/internal/clients"
	"github.com/rpanchyk/javaman/internal/models"
)

type OpenJdkListFetcher struct {
	config     *models.Config
	httpClient clients.HttpClient
}

func NewOpenJdkListFetcher(
	config *models.Config,
	httpClient clients.HttpClient) *OpenJdkListFetcher {

	return &OpenJdkListFetcher{
		config:     config,
		httpClient: httpClient,
	}
}

func (f OpenJdkListFetcher) Fetch() ([]models.Sdk, error) {
	fmt.Printf("Fetching openjdk SDKs ... ")

	status, response, err := f.httpClient.Get("https://jdk.java.net/archive/")
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}
	if status < 200 || status > 299 {
		return nil, fmt.Errorf("error making http request, status: %d", status)
	}

	// https://download.java.net/java/GA/jdk23.0.2/6da2a6609d6e406f85c491fcb119101b/7/GPL/openjdk-23.0.2_linux-x64_bin.tar.gz
	// https://download.java.net/java/GA/jdk23.0.2/6da2a6609d6e406f85c491fcb119101b/7/GPL/openjdk-23.0.2_macos-aarch64_bin.tar.gz
	// https://download.java.net/java/GA/jdk23.0.2/6da2a6609d6e406f85c491fcb119101b/7/GPL/openjdk-23.0.2_windows-x64_bin.zip

	r, err := regexp.Compile(`href=['"](.*/openjdk-([0-9._\-]+)_(\w+)-(\w+)_bin\.(?:tar\.gz|zip)+)['"]`)
	if err != nil {
		return nil, fmt.Errorf("error compile regexp: %w", err)
	}

	sdks := make([]models.Sdk, 0)
	for _, parts := range r.FindAllStringSubmatch(response, -1) {
		//fmt.Printf("parts: %+v\n", parts)

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
			Vendor:  "openjdk",
			URL:     strings.TrimSpace(parts[1]),
			Version: strings.TrimSpace(parts[2]),
			Os:      os,
			Arch:    arch,
		}

		if !slices.Contains(sdks, sdk) {
			sdks = append(sdks, sdk)
		}
	}

	fmt.Println("OK")
	return sdks, nil
}
