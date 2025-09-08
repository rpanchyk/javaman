package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rpanchyk/javaman/internal/clients"
	"github.com/rpanchyk/javaman/internal/services/cacher"
	"github.com/rpanchyk/javaman/internal/services/lister"
	"github.com/rpanchyk/javaman/internal/services/lister/vendors"
	"github.com/rpanchyk/javaman/internal/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows list of available Java versions",
	Run: func(cmd *cobra.Command, args []string) {
		fetchers := make([]lister.ListFetcher, 0)
		for _, vendor := range utils.Config.Vendors {
			vendorName := strings.ToLower(strings.TrimSpace(vendor))
			if vendorName == "microsoft" {
				fetchers = append(fetchers, vendors.NewMicrosoftListFetcher(
					&utils.Config,
					&clients.SimpleHttpClient{},
				))
			}
		}

		listFetcher := lister.NewFilteredListFetcher(
			&utils.Config,
			lister.NewDefaultListFetcher(
				fetchers,
				cacher.NewDefaultListCacher(&utils.Config)),
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
