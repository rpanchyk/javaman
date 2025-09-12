package lister

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/rpanchyk/javaman/internal/models"
)

type FilteredListFetcher struct {
	config      *models.Config
	listFilter  *models.ListFilter
	listFetcher ListFetcher
}

func NewFilteredListFetcher(
	config *models.Config,
	listFilter *models.ListFilter,
	listFetcher ListFetcher) *FilteredListFetcher {

	return &FilteredListFetcher{
		config:      config,
		listFilter:  listFilter,
		listFetcher: listFetcher,
	}
}

func (f FilteredListFetcher) Fetch() ([]models.Sdk, error) {
	sdks, err := f.listFetcher.Fetch()
	if err != nil {
		return nil, fmt.Errorf("error fetching SDKs: %w", err)
	}

	filtered, err := f.filterSdks(sdks)
	if err != nil {
		return nil, fmt.Errorf("error filtering SDKs: %w", err)
	}

	return f.enrichSdks(filtered)
}

func (f FilteredListFetcher) filterSdks(sdks []models.Sdk) ([]models.Sdk, error) {
	r := regexp.MustCompile("[^0-9]+")
	sort.Slice(sdks, func(i, j int) bool {
		firstVersion := r.ReplaceAllString(sdks[i].Version, ".")
		secondVersion := r.ReplaceAllString(sdks[j].Version, ".")

		first := strings.Split(firstVersion, ".")
		second := strings.Split(secondVersion, ".")

		// 1.9.1 vs 1.8
		// 1.8 vs 1.9.1
		diff := len(first) - len(second)
		if diff != 0 {
			if diff > 0 {
				for i := 0; i < diff; i++ {
					second = append(second, "0")
				}
			} else {
				for i := 0; i < -diff; i++ {
					first = append(first, "0")
				}
			}
		}

		for k := 0; k < len(first); k++ {
			if first[k] != second[k] { // 1.9.1 vs 1.9.2
				f, err := strconv.Atoi(first[k])
				if err != nil {
					panic(err)
				}
				s, err := strconv.Atoi(second[k])
				if err != nil {
					panic(err)
				}
				return f > s
			}
		}

		return false
	})
	//fmt.Printf("Sorted sdks: %v\n", sdks)

	maxCount := f.config.ListLimit
	if f.listFilter.Number > 0 {
		maxCount = f.listFilter.Number
	}
	filterVersion := ""
	if f.listFilter.Version != "" {
		filterVersion = f.listFilter.Version
	}

	res := make([]models.Sdk, 0)
	count := 0
	ver := ""
	for _, sdk := range sdks {
		if filterVersion != "" && strings.Index(sdk.Version, filterVersion) == -1 {
			continue
		}
		if ver != sdk.Version {
			count++
		}
		if count > maxCount {
			break
		}

		res = append(res, sdk)
		ver = sdk.Version
	}

	clear(sdks)
	return res, nil
}

func (f FilteredListFetcher) enrichSdks(sdks []models.Sdk) ([]models.Sdk, error) {
	for i := 0; i < len(sdks); i++ {
		filePath := filepath.Join(f.config.DownloadDir, path.Base(sdks[i].URL))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			sdks[i].IsDownloaded = true
			sdks[i].FilePath = filePath
		}

		installDir := filepath.Join(f.config.InstallDir, sdks[i].Vendor+"-"+sdks[i].Version)
		if info, err := os.Stat(installDir); err == nil && info.IsDir() {
			sdks[i].IsInstalled = true
		}

		if envVar, ok := os.LookupEnv("JAVA_HOME"); ok && envVar == installDir {
			sdks[i].IsDefault = true
		}
	}
	return sdks, nil
}
