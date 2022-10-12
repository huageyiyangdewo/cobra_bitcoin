package cmd

import (
	"github.com/spf13/cobra"
)

var data string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "addBlock",
	Short: "add block",
	Long: `addBlock xxxx 添加数据到区块链`,
	Run: func(cmd *cobra.Command, args []string) {
		//blockChain.AddBlock(data) todo
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&data, "data", "d", "", "block data")

	createCmd.MarkFlagRequired("data")
}
