package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/javaman/internal/clients"
	"github.com/rpanchyk/javaman/internal/globals"
	"github.com/rpanchyk/javaman/internal/services/downloader"
	"github.com/rpanchyk/javaman/internal/services/installer"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install specified Java version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		installer := installer.NewDefaultInstaller(
			&globals.Config,
			&globals.DefaultListFetcher,
			downloader.NewDefaultDownloader(
				&globals.Config,
				&globals.DefaultListFetcher,
				&clients.SimpleHttpSaver{}),
		)
		if err := installer.Install(args[0]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
