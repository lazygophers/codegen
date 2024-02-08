package cli

import "github.com/spf13/cobra"

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
