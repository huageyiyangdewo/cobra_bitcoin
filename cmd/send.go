package cmd

import (
	"cobra_bitcoin/internal"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var from, to, miner, value, sendData string

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send btc",
	Long: `send FROM TO AMOUNT MINER 转账命令`,
	Run: func(cmd *cobra.Command, args []string) {
		amount, _ := strconv.ParseFloat(value, 64)


		// 创建挖矿交易
		coinbase := internal.NewCoinbaseTx(miner, sendData)
		txs := []*internal.Transaction{coinbase}

		// 创建普通交易
		tx := internal.NewTransaction(from, to, amount, blockChain)
		if tx != nil {
			txs = append(txs, tx)
		} else {
			fmt.Println("发现无效交易!过滤")
		}
		// 添加区块
		blockChain.AddBlock(txs)

		fmt.Println("挖矿成功！")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVarP(&from, "from", "f", "", "from address")
	sendCmd.Flags().StringVarP(&to, "to", "t", "", "to address")
	sendCmd.Flags().StringVarP(&miner, "miner", "m", "", "miner address")
	sendCmd.Flags().StringVarP(&value, "amount", "a", "0.0", "send amount")
	sendCmd.Flags().StringVarP(&sendData, "sendData", "d", "", "sendData")

	sendCmd.MarkFlagRequired(from)
	sendCmd.MarkFlagRequired(to)
	sendCmd.MarkFlagRequired(miner)
	sendCmd.MarkFlagRequired(value)
}
