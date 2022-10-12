package internal

import (
	"cobra_bitcoin/configs"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

// BlockChain 创建区块链
// 1、bolt 数据库的句柄
// 2、最后一个区块的哈希值
type BlockChain struct {
	//Blocks []*Block
	Db *bolt.DB

	Tail []byte  // 最后一个区块的哈希值
}

// NewBlockChain 实现创建区块链的方法
func NewBlockChain(miner string) *BlockChain {
	//genesisBlock := NewBlock(configs.GenesisInfo, []byte{0x00000000000000})
	//
	//bc := &BlockChain{
	//	Blocks: []*Block{genesisBlock},
	//}
	//
	//return bc
	db, err := bolt.Open(configs.BlockChainName, 0600, nil)
	if err != nil {
		log.Panicln(err)
	}

	var tail []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(configs.BlockBucketName))
		if b == nil {
			fmt.Printf("bucket不存在，准备创建！\n")

			b, err = tx.CreateBucket([]byte(configs.BlockBucketName))
			if err != nil {
				log.Panicln(err)
			}

			// 创世块中只有一个挖矿交易，只有 coinbase
			coinbase := NewCoinbaseTx(miner, configs.GenesisInfo)
			genesisBlock := NewBlock([]*Transaction{coinbase}, []byte{})
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化，转成字节流*/ )
			if err != nil {
				fmt.Printf("db put err: %s \n", err)
			}
			err = b.Put([]byte(configs.LastHashKey), genesisBlock.Hash)
			if err != nil {
				fmt.Printf("db put err: %s \n", err)
			}

			tail = genesisBlock.Hash

		} else {
			tail = b.Get([]byte(configs.LastHashKey))
		}
		return nil

	})

	if err != nil {
		fmt.Printf("db update err: %s \n", err)
	}

	return &BlockChain{db, tail}
}

// AddBlock 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//// 1、创建一个区块
	//lastBlock := bc.Blocks[len(bc.Blocks)-1] // 最后一个区块就是新区块的 prevHash
	//hash := lastBlock.Hash
	//// 2、添加到 bc.Blocks 数组中
	//block := NewBlock(data, hash)
	//bc.Blocks = append(bc.Blocks, block)

	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(configs.BlockBucketName))
		if b == nil {
			fmt.Printf("bucket不存在，请检查！\n")
			os.Exit(1)

		}

		block := NewBlock(txs, bc.Tail)
		err := b.Put(block.Hash, block.Serialize() /*将区块序列化，转成字节流*/ )
		if err != nil {
			fmt.Printf("db put err: %s \n", err)
		}
		err = b.Put([]byte(configs.LastHashKey), block.Hash)
		if err != nil {
			fmt.Printf("db put err: %s \n", err)
		}

		bc.Tail = block.Hash
		return nil

	})

	if err != nil {
		fmt.Printf("db update err: %s \n", err)
	}
}

// BlockChainIterator 定义一个区块链的迭代器，包括 db, current
type BlockChainIterator struct {
	Db *bolt.DB
	Current []byte // 储存当前区块的哈希值
}

// NewBlockChainIterator 创建迭代器，
func (bc *BlockChain) NewBlockChainIterator() *BlockChainIterator {
	return &BlockChainIterator{
		Db: bc.Db,
		Current: bc.Tail,
	}
}

func (it *BlockChainIterator) Next() *Block {
	var block *Block

	err := it.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(configs.BlockBucketName))
		if b == nil {
			fmt.Printf("bucket不存在，请检查！\n")
			os.Exit(1)

		}

		blockInfo := b.Get(it.Current)
		block = DeSerialize(blockInfo)

		it.Current = block.PrevHash

		return nil
	})

	if err != nil {
		fmt.Printf("db BlockChainIterator next err: %s \n", err)
	}

	return block
}

// FindMyUtxos 找到所有的 utxo
func (bc *BlockChain) FindMyUtxos(address string) []TxOutput {

	var utxos []TxOutput // 未消耗的

	// 这里是标识已经消耗过的 utxo 的结构，key是交易id，value是这个id在input中的索引
	spentUTXOs := make(map[string][]int64)

	it := bc.NewBlockChainIterator()
	// 遍历账本
	for {
		block := it.Next()
		// 遍历交易
		for _, tx := range block.Transactions {
			// 遍历 inputs
			for _, input := range tx.TxInputs {
				if input.Address == address {
					key := string(input.TxID)
					spentUTXOs[key] = append(spentUTXOs[key], input.Index)
				}
			}

			// 遍历 output
			OUTPUT:
				for i, output := range tx.TxOutputs {
					key := string(tx.TxID)
					indexes := spentUTXOs[key] // 当前交易有被消耗过的 output
					if len(indexes) != 0 {
						for _, j  := range indexes {
							if int64(i) == j { // 这笔交易已经被消耗过了
								continue OUTPUT
							}
						}
					}

					// 找到属于我的所有的 output
					if output.Address == address {
						fmt.Printf("找到了属于 %s 的output，i: %d\n", address, i)
						utxos = append(utxos, output)
					}
				}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}


	return utxos
}

func (bc *BlockChain) GetBalance(address string) {
	utxos := bc.FindMyUtxos(address)

	var total = 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("%s 的余额为 %f\n", address, total)
}

// FindNeedUtxos 遍历账本，找到属于付款人的合适的金额，把这个 outputs 找到
func (bc *BlockChain) FindNeedUtxos(from string, value float64) (map[string][]int64, float64) {
	utxos := make(map[string][]int64) // 标识能用的 utxo
	var resValue float64 // 这些 utxo 存储的金额


	// ++++++++++++++++++++
	// 这里是标识已经消耗过的 utxo 的结构，key是交易id，value是这个id在input中的索引
	spentUTXOs := make(map[string][]int64)

	it := bc.NewBlockChainIterator()
	// 遍历账本
	for {
		block := it.Next()
		// 遍历交易
		for _, tx := range block.Transactions {
			// 遍历 inputs
			for _, input := range tx.TxInputs {
				if input.Address == from {
					key := string(input.TxID)
					spentUTXOs[key] = append(spentUTXOs[key], input.Index)
				}
			}

			// 遍历 output
		OUTPUT:
			for i, output := range tx.TxOutputs {
				key := string(tx.TxID)
				indexes := spentUTXOs[key] // 当前交易有被消耗过的 output
				if len(indexes) != 0 {
					for _, j  := range indexes {
						if int64(i) == j { // 这笔交易已经被消耗过了
							continue OUTPUT
						}
					}
				}

				// 找到属于我的所有的 output
				if output.Address == from {
					fmt.Printf("找到了属于 %s 的output，i: %d\n", from, i)

					// 找到符合条件的 output,添加到返回逻辑中
					utxos[key] = append(utxos[key], int64(i))
					resValue += output.Value
					// 判断一下金额是否满足
					if resValue >= value {
						// 满足，直接返回，
						return utxos, resValue
					}
					// 不满足，继续遍历直到遍历完整个 链
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}


	// ++++++++++++++++++++
	return utxos, 0.0
}