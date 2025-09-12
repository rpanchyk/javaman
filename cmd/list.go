package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/javaman/internal/globals"
	"github.com/rpanchyk/javaman/internal/models"
	"github.com/rpanchyk/javaman/internal/services/lister"
	"github.com/spf13/cobra"
)

var (
	number  int
	version string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows list of available Java versions",
	Run: func(cmd *cobra.Command, args []string) {
		listFetcher := lister.NewFilteredListFetcher(
			&globals.Config,
			&models.ListFilter{
				Number:  number,
				Version: version,
			},
			&globals.DefaultListFetcher,
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
			fmt.Printf("%s %-30s %-10s %-6s %-15s %s\n",
				defaultMarker, sdk.Vendor+"-"+sdk.Version, sdk.Os, sdk.Arch, downloadedMarker, installedMarker)
		}

		if len(sdks) == 0 {
			fmt.Println("No available SDK versions found")
		}
	},
}

func init() {
	listCmd.Flags().IntVarP(&number, "number", "n", 0, "Number of the last available SDK versions to show")
	listCmd.Flags().StringVarP(&version, "version", "v", "", "SDK version to show")
	RootCmd.AddCommand(listCmd)
}
