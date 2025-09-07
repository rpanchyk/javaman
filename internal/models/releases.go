package models

type Releases struct {
	ApiVersion string `mapstructure:"apiVersion"`
	Generated  string `mapstructure:"generated"`
	Data       Data   `mapstructure:"data"`
}

type Data struct {
	Count    int       `mapstructure:"count"`
	Releases []Release `mapstructure:"releases"`
}

type Release struct {
	Version string `mapstructure:"version"`
	Family  int    `mapstructure:"family"`
	Status  string `mapstructure:"status"`
	Type    string `mapstructure:"type"`
}
