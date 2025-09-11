package vendors

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/rpanchyk/javaman/internal/clients"
	"github.com/rpanchyk/javaman/internal/models"
)

type AmazonListFetcher struct {
	config     *models.Config
	httpClient clients.HttpClient
}

func NewAmazonListFetcher(
	config *models.Config,
	httpClient clients.HttpClient) *AmazonListFetcher {

	return &AmazonListFetcher{
		config:     config,
		httpClient: httpClient,
	}
}

func (f AmazonListFetcher) Fetch() ([]models.Sdk, error) {
	fmt.Printf("Fetching corretto SDKs ... ")

	status, response, err := f.httpClient.Get("https://raw.githubusercontent.com/corretto/corretto-downloads/refs/heads/main/latest_links/version-info.json")
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}
	if status < 200 || status > 299 {
		return nil, fmt.Errorf("error making http request, status: %d", status)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}
	//fmt.Printf("Data: %+v\n", data)

	majors := make([]string, 0)
	for key := range data {
		for _, v := range data[key].([]interface{}) {
			majors = append(majors, fmt.Sprintf("%v", v))
		}
	}
	slices.Sort(majors)
	//fmt.Printf("Majors: %+v\n", majors)

	// https://corretto.aws/downloads/resources/11.0.27.6.1/amazon-corretto-11.0.27.6.1-linux-aarch64.tar.gz
	// https://corretto.aws/downloads/resources/11.0.27.6.1/amazon-corretto-11.0.27.6.1-macosx-aarch64.tar.gz
	// https://corretto.aws/downloads/resources/11.0.28.6.1/amazon-corretto-11.0.28.6.1-windows-x64-jdk.zip

	r, err := regexp.Compile(`href=['"](.*/amazon-corretto-([0-9._\-]+)-(\w+)-(\w+)(?:-jdk)*\.(?:tar\.gz|zip)+)['"]`)
	if err != nil {
		return nil, fmt.Errorf("error compile regexp: %w", err)
	}

	sdks := make([]models.Sdk, 0)
	for _, major := range majors {
		page := 0
		for {
			page++
			url := fmt.Sprintf("https://github.com/corretto/corretto-%s/releases?page=%d", major, page)
			//fmt.Printf("Fetching %s\n", url)

			status, response, err := f.httpClient.Get(url)
			if err != nil {
				return nil, fmt.Errorf("error making http request: %w", err)
			}
			if status < 200 || status > 299 {
				break
			}

			matches := r.FindAllStringSubmatch(response, -1)
			if len(matches) == 0 {
				break
			}
			for _, parts := range matches {
				//fmt.Printf("parts: %+v\n", parts)

				var os models.Os
				switch strings.ToLower(parts[3]) {
				case "linux":
					os = models.Linux
				case "macosx":
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
					Vendor:  "corretto",
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
	}

	fmt.Println("OK")
	return sdks, nil
}
