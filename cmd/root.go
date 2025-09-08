package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/rpanchyk/javaman/internal/clients"
	"github.com/rpanchyk/javaman/internal/models"
	"github.com/rpanchyk/javaman/internal/services/cacher"
	"github.com/rpanchyk/javaman/internal/services/lister"
	"github.com/rpanchyk/javaman/internal/services/lister/vendors"
	"github.com/rpanchyk/javaman/internal/utils"
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
	cobra.OnInitialize(initConfig, initFetcher)
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

	utils.Config = getConfig()
	fmt.Printf("Config: %+v\n", utils.Config)
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

func initFetcher() {
	fetchers := make([]lister.ListFetcher, 0)
	for _, vendor := range utils.Config.Vendors {
		switch strings.ToLower(strings.TrimSpace(vendor)) {
		case "microsoft":
			fetchers = append(fetchers, vendors.NewMicrosoftListFetcher(
				&utils.Config,
				&clients.SimpleHttpClient{},
			))
		}
	}

	utils.DefaultListFetcher = *lister.NewDefaultListFetcher(
		fetchers,
		cacher.NewDefaultListCacher(&utils.Config),
	)
}
