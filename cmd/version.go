package cmd

import (
	"fmt"

	"github.com/rpanchyk/javaman/internal/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows version of javaman",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("javaman version %s %s/%s\n", "1.0.0", utils.CurrentOs(), utils.CurrentArch())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
