package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/rpanchyk/javaman/internal/clients"
	"github.com/rpanchyk/javaman/internal/globals"
	"github.com/rpanchyk/javaman/internal/models"
	"github.com/rpanchyk/javaman/internal/services/cacher"
	"github.com/rpanchyk/javaman/internal/services/lister"
	"github.com/rpanchyk/javaman/internal/services/lister/vendors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "javaman",
	Short: "Java version manager",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	go catchSignal()

	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func catchSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV)
	sig := <-sigs
	fmt.Println("Signal obtained:", sig)
}

func init() {
	cobra.OnInitialize(initConfig, initListFetcher)
}

func initConfig() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Cannot get user home directory, error:", err)
		os.Exit(1)
	}

	viper.AddConfigPath(filepath.Join(userHomeDir, ".javaman"))
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Cannot read config, error:", err)
		os.Exit(1)
	}

	globals.Config = getConfig()
	fmt.Printf("Config: %+v\n", globals.Config)
}

func getConfig() models.Config {
	var config models.Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Cannot unmarshal config, error:", err)
		os.Exit(1)
	}

	config.DownloadDir = toAbsPath(config.DownloadDir)
	config.InstallDir = toAbsPath(config.InstallDir)
	config.ListCacheFile = toAbsPath(config.ListCacheFile)
	return config
}

func toAbsPath(path string) string {
	if !filepath.IsAbs(path) {
		configFile := viper.GetViper().ConfigFileUsed()
		return filepath.Join(filepath.Dir(configFile), path)
	}
	return path
}

func initListFetcher() {
	httpClient := &clients.SimpleHttpClient{}

	fetchers := make([]lister.ListFetcher, 0)
	for _, vendor := range globals.Config.Vendors {
		switch strings.ToLower(strings.TrimSpace(vendor)) {
		case "corretto":
			fetchers = append(fetchers, vendors.NewAmazonListFetcher(
				&globals.Config,
				httpClient,
			))
		case "microsoft":
			fetchers = append(fetchers, vendors.NewMicrosoftListFetcher(
				&globals.Config,
				httpClient,
			))
		case "openjdk":
			fetchers = append(fetchers, vendors.NewOpenJdkListFetcher(
				&globals.Config,
				httpClient,
			))
		}
	}

	globals.DefaultListFetcher = *lister.NewDefaultListFetcher(
		fetchers,
		cacher.NewDefaultListCacher(&globals.Config),
	)
}
