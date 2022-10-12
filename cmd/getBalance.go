package cmd

import (
	"github.com/spf13/cobra"
)

var address string

// getBalanceCmd represents the getBalance command
var getBalanceCmd = &cobra.Command{
	Use:   "getBalance",
	Short: "get balance",
	Long: `getBalance 地址  获取余额`,
	Run: func(cmd *cobra.Command, args []string) {
		blockChain.GetBalance(address)
	},
}

func init() {
	rootCmd.AddCommand(getBalanceCmd)

	getBalanceCmd.Flags().StringVarP(&address, "address", "a", "", "btc address")
	getBalanceCmd.MarkFlagRequired(address)
}
