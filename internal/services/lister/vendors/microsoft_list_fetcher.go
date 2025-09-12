package vendors

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/rpanchyk/javaman/internal/clients"
	"github.com/rpanchyk/javaman/internal/models"
	"github.com/rpanchyk/javaman/internal/utils"
)

type MicrosoftListFetcher struct {
	config     *models.Config
	httpClient clients.HttpClient
}

func NewMicrosoftListFetcher(
	config *models.Config,
	httpClient clients.HttpClient) *MicrosoftListFetcher {

	return &MicrosoftListFetcher{
		config:     config,
		httpClient: httpClient,
	}
}

func (f MicrosoftListFetcher) Fetch() ([]models.Sdk, error) {
	fmt.Printf("Fetching microsoft SDKs ...")
	pb := utils.NewDotProgressBar()
	pb.Start()
	defer pb.Stop()

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
				URL:     strings.TrimSpace(parts[1]),
				Version: strings.TrimSpace(parts[2]),
				Os:      os,
				Arch:    arch,
			}

			if !slices.Contains(sdks, sdk) {
				sdks = append(sdks, sdk)
			}
		}
	}

	fmt.Println(" OK")
	return sdks, nil
}
