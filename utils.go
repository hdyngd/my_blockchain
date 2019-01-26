package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

// IntToHex int64型をビッグエンディアンの順でバイト配列に変換
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer) // 空のバッファーの生成
	// int64型の値をバイト配列化し、バッファーに追記
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err) // エラーが発生した場合強制終了
	}
	return buff.Bytes() // バッファーをバイト配列に変換し返す
}
