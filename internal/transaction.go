package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type TxInput struct {
	TxID    []byte // 交易id
	Index   int64  // output 的索引
	Address string // 地址
}

type TxOutput struct {
	Value   float64 // 转账金额
	Address string  // 地址
}

type Transaction struct {
	TxID      []byte     // 交易id
	TxInputs  []TxInput  // 所有的inputs
	TxOutputs []TxOutput // 所有的outputs
}

// SetTxID 设置交易ID
func (tx *Transaction) SetTxID() {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panicf("Transaction SetTxID err:%s", err)
	}

	hash := sha256.Sum256(buffer.Bytes())
	tx.TxID = hash[:]

}


// 实现挖矿交易
// 特点：只有输出，没有有效的输入（不需要交易ID，不需要索引，不需要签名）
// 把挖矿的人传递进来，因为有奖励

func NewCoinbaseTx(address string, data string) *Transaction {

	// 我们在后面的交易中，需要识别一个交易是否为coinbase,所以需要设置一些特殊的值，用于判断
	inputs := []TxInput{{
		TxID: nil,
		Index: -1,
		Address: data,
	}}

	outputs := []TxOutput{{
		Value: 12.5,
		Address: address,
	}}

	tx := &Transaction{
		TxInputs: inputs,
		TxOutputs: outputs,
	}
	tx.SetTxID()

	return tx
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	utxos := make(map[string][]int64) // 标识能用的 utxo
	var resValue float64 // 这些 utxo 存储的金额

	// 1、遍历账本，找到属于付款人的合适的金额，把这个 outputs 找到
	utxos, resValue = bc.FindNeedUtxos(from, amount)

	//2、如果找到的钱不足，则创建交易失败
	if resValue < amount {
		fmt.Printf("余额不足，创建交易失败！\n")
		return nil
	}

	var inputs []TxInput
	var outputs []TxOutput

	//3、将 outputs 转成 inputs
	for txid, indexes := range utxos {
		for _, i := range indexes {
			input := TxInput{
				TxID: []byte(txid),
				Index: i,
				Address: from,
			}

			inputs = append(inputs, input)
		}
	}

	// 4、创建输出，创建一个属于收款人的 output
	output := TxOutput{
		Value: amount,
		Address: to,
	}
	outputs  = append(outputs, output)


	//5、如果有找零，创建输入收款人的 output
	if resValue > amount {
		output1 := TxOutput{resValue-amount, from}
		outputs  = append(outputs, output1)
	}

	tx := Transaction{
		TxID: nil,
		TxInputs: inputs,
		TxOutputs: outputs,
	}
	// 设置交易ID
	tx.SetTxID()
	return &tx
}