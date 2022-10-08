package internal

import (
	"bytes"
	"cobra_bitcoin/configs"
	"cobra_bitcoin/utils"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	block *Block

	// 存储哈希值，它内置了一些方法
	// Cmp: 比较方法
	// SetBytes: 把 bytes 转成 big.Int 类型
	// SetString: 把 string 转成 big.Int 类型
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	pow := &ProofOfWork{
		block: block,
	}

	// 写难度值，难度值应该是推导出来的，但是为了简化，先把难度值固定，后面再改
	// 16进制格式的字符串
	// 0000100000000000000000000000000000000000000000000000000000000000

	//targetString := "0000100000000000000000000000000000000000000000000000000000000000"
	//var bigIntTmp big.Int
	////bigIntTmp.SetBytes([]byte(targetString))
	//bigIntTmp.SetString(targetString, 16)
	//
	//pow.target = &bigIntTmp

	bigIntTmp := big.NewInt(1)
	//bigIntTmp.Lsh(bigIntTmp, 256)
	//bigIntTmp.Rsh(bigIntTmp, 16)
	bigIntTmp.Lsh(bigIntTmp, 256 - configs.Bits) // 简写方式
	pow.target = bigIntTmp

	return pow
}

// Run 这是 pow 的运算函数，为了获取挖矿的随机数，同时返回区块的哈希值
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	// 1、获取block数据
	// 2、拼接 nonce
	// 3、sha256
	// 4、与难度值比较
	// a. 哈希值大于难度值，nonce++
	// b.哈希值小于难度值，挖矿成功，退出

	var nonce uint64
	var hash [32]byte

	for  {

		hash = sha256.Sum256(pow.prepareData(nonce))

		// 将hash(数组类型) 转换为big.Int，然后与pow.target比较
		var bigIntTmp big.Int
		bigIntTmp.SetBytes(hash[:])
		if  bigIntTmp.Cmp(pow.target) == -1{
			fmt.Printf("挖矿成功！nonce:%d, hash: %x\n", nonce, hash)
			break
		} else {
			nonce++
		}
	}

	return hash[:], nonce
}

// prepareData 拼接数据
func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	b := pow.block

	tmp := [][]byte{
		utils.Uint64ToByte(b.Version),
		b.PrevHash,
		b.MerkleRoot,
		utils.Uint64ToByte(b.TimeStamp),
		utils.Uint64ToByte(b.Difficulty),
		utils.Uint64ToByte(nonce),
		b.Data,
	}
	data := bytes.Join(tmp, []byte{})
	return data
}


// IsValid 校验是否合法
func (pow *ProofOfWork) IsValid() bool {
	// 在校验的时候，block的数据是完整的
	hash := sha256.Sum256(pow.prepareData(pow.block.Nonce))

	var tmp big.Int
	tmp.SetBytes(hash[:])

	//if tmp.Cmp(pow.target) == -1 {
	//	return true
	//}
	//return false

	return tmp.Cmp(pow.target) == -1
}