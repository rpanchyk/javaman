package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/javaman/internal/services/defaulter"
	"github.com/rpanchyk/javaman/internal/utils"
	"github.com/spf13/cobra"
)

var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Set specified Java version as default",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defaulter := defaulter.NewDefaultDefaulter(
			&utils.Config,
			&utils.DefaultListFetcher,
		)
		if err := defaulter.Default(args[0]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(defaultCmd)
}
