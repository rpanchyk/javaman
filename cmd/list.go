package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/javaman/internal/services/lister"
	"github.com/rpanchyk/javaman/internal/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows list of available Java versions",
	Run: func(cmd *cobra.Command, args []string) {
		listFetcher := lister.NewFilteredListFetcher(
			&utils.Config,
			&utils.DefaultListFetcher,
		)
		sdks, err := listFetcher.Fetch()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, sdk := range sdks {
			defaultMarker := " "
			if sdk.IsDefault {
				defaultMarker = "*"
			}
			downloadedMarker := "            "
			if sdk.IsDownloaded {
				downloadedMarker = "[downloaded]"
			}
			installedMarker := ""
			if sdk.IsInstalled {
				installedMarker = "[installed]"
			}
			fmt.Printf("%s %s-%-20s %-10s %-10s %-15s %s\n",
				defaultMarker, sdk.Vendor, sdk.Version, sdk.Os, sdk.Arch, downloadedMarker, installedMarker)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
