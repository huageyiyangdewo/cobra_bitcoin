package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer

	// buffer 这里为什么需要加 &，因为 bytes.Buffer 实现Write方法时，
	// 使用的是指针接收者，所以再给 io.Writer 的变量赋值时，需要用 变量的地址，而不能用 变量
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		panic(fmt.Sprintf("Uint64ToByte error: %s", err))
	}

	return buffer.Bytes()

}
