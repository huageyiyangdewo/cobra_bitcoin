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
func NewBlockChain() *BlockChain {
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

			genesisBlock := NewBlock(configs.GenesisInfo, []byte{})
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
func (bc *BlockChain) AddBlock(data string) {
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

		block := NewBlock(data, bc.Tail)
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