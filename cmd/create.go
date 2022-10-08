package cmd

import (
	"github.com/spf13/cobra"
)

var data string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "addBlock",
	Short: "create block",
	Long: `create new block`,
	Run: func(cmd *cobra.Command, args []string) {
		blockChain.AddBlock(data)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&data, "data", "d", "", "block data")

	createCmd.MarkFlagRequired("data")
}
