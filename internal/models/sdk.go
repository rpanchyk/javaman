package models

type Sdk struct {
	URL          string `json:"url"`
	Vendor       string `json:"vendor"`
	Version      string `json:"version"`
	Os           Os     `json:"os"`
	Arch         Arch   `json:"arch"`
	FilePath     string `json:"-"`
	IsDownloaded bool   `json:"-"`
	IsInstalled  bool   `json:"-"`
	IsDefault    bool   `json:"-"`
}
