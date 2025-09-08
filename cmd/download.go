package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/javaman/internal/clients"
	"github.com/rpanchyk/javaman/internal/services/downloader"
	"github.com/rpanchyk/javaman/internal/utils"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download specified Java version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		downloader := downloader.NewDefaultDownloader(
			&utils.Config,
			&utils.DefaultListFetcher,
			&clients.SimpleHttpSaver{},
		)
		if _, err := downloader.Download(args[0]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)
}
