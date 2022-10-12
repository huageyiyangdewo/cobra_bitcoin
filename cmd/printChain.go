package cmd

import (
	"bytes"
	"cobra_bitcoin/internal"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

// printChainCmd represents the printChain command
var printChainCmd = &cobra.Command{
	Use:   "printChain",
	Short: "print block chain info",
	Long: `printChain 打印区块链`,
	Run: func(cmd *cobra.Command, args []string) {
		it := blockChain.NewBlockChainIterator()
		for {
			v := it.Next()
			fmt.Printf("Version: %d \n", v.Version)
			fmt.Printf("PrevHash: %x \n", v.PrevHash)
			fmt.Printf("MerkleRoot: %x \n", v.MerkleRoot)
			t := time.Unix(int64(v.TimeStamp), 0).Format("2006-01-02 15:04:05")
			fmt.Printf("TimeStamp: %v \n", t)
			fmt.Printf("Difficulty: %d \n", v.Difficulty)
			fmt.Printf("Nonce: %d \n", v.Nonce)
			fmt.Printf("Hash: %x \n", v.Hash)
			fmt.Printf("Data: %s \n", v.Transactions[0].TxInputs[0].Address)

			pow := internal.NewProofOfWork(v)
			fmt.Printf("IsValid: %v\n", pow.IsValid())
			fmt.Println("---")

			if bytes.Equal(v.PrevHash, []byte{}) {
				fmt.Println("区块遍历结束！")
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(printChainCmd)

}
