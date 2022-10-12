package internal

import (
	"bytes"
	"cobra_bitcoin/configs"
	"cobra_bitcoin/utils"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

//Block 定义区块结构
type Block struct {
	Version  uint64 // 区块版本号
	PrevHash []byte // 前区块哈希

	MerkleRoot []byte // 梅克尔根 先填写为空

	TimeStamp  uint64 // 从1970.1.1至今的秒数
	Difficulty uint64 // 挖矿难度值
	Nonce      uint64 // 随机数

	Hash []byte // 哈希, 为了方便，我们将当前区块的哈希放入Block中

	//Data []byte // 数据
	Transactions []*Transaction // v4改成交易数据
}

func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{
		Version:    00,
		PrevHash:   prevHash,
		MerkleRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: configs.Bits, //
		//Nonce:      10, // 后期调整
		Hash:       []byte{},
		//Data:       []byte(data),
		Transactions: txs,
	}

	//block.SetHash()
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce

	block.HashTransaction() // 赋值 merkleroot

	return block
}

// SetHash 为了生成区块哈希，实现一个简单的函数，来计算哈希值，没有随机数，没有难度值
func (b *Block) SetHash() {
	//var data []byte
	//data = append(data, utils.Uint64ToByte(b.Version)...)
	//data = append(data, b.PrevHash...)
	//data = append(data, b.MerkleRoot...)
	//data = append(data, utils.Uint64ToByte(b.TimeStamp)...)
	//data = append(data, utils.Uint64ToByte(b.Difficulty)...)
	//data = append(data, utils.Uint64ToByte(b.Nonce)...)
	//data = append(data, b.Data...)

	tmp := [][]byte{
		utils.Uint64ToByte(b.Version),
		b.PrevHash,
		b.MerkleRoot,
		utils.Uint64ToByte(b.TimeStamp),
		utils.Uint64ToByte(b.Difficulty),
		utils.Uint64ToByte(b.Nonce),
		//b.Data, // todo
	}
	data := bytes.Join(tmp, []byte{})

	hash := sha256.Sum256(data)
	b.Hash = hash[:]
}


// Serialize 序列化，将区块转换成字节流
func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if err != nil {
		log.Panicln(err)
	}

	return buffer.Bytes()
}

// DeSerialize 反序列化
func DeSerialize(data []byte) *Block {
	//fmt.Printf("解码传入的数据： %x \n", data)
	b := &Block{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(b)
	if err != nil {
		log.Panicf("DeSerialize: %s ", err)
	}

	return b
}

// HashTransaction 模拟梅克尔根，做一个简单的处理
func (b *Block) HashTransaction() {
	// 我们的交易的ID就是交易的哈希值，所以我们可以将交易id拼接起来，整体做一次哈希运算
	// 作为 merkleroot

	var hashes []byte
	for _, v := range b.Transactions {
		hashes = append(hashes, v.TxID...)
	}
	hash := sha256.Sum256(hashes)
	b.MerkleRoot = hash[:]
}