package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"strconv"
	"testing"
)

func TestHello(t *testing.T) {
	fmt.Println("Hello")
}

func TestErrorFatal(t *testing.T) {
	fmt.Println("TestErrorFatal")
	count := 0
	t.Error("Error, but will continue...") // 継続
	count++
	t.Fatalf("Fatal Successed? count = %d", count) // 処理の停止
	// これ以降は実行されない
	count++
	fmt.Printf("This is never executed. count = %d", count)
}

func TestBlock(t *testing.T) {
	fmt.Println("TestBlock")
	block := Block{} // 構造体を生成しblockに代入
	fmt.Println(block)
	pBlock := &Block{} // 構造体を生成しそのポインタを返す
	fmt.Println(pBlock)
}

// strconv.FormatIntのテスト
// 要import strconv
func TestStrconv(t *testing.T) {
	dataBinary := strconv.FormatInt(12345, 2) // 2進数
	fmt.Println(dataBinary)
	dataDecimal := strconv.FormatInt(12345, 10) // 10進数
	fmt.Println(dataDecimal)
	dataHexadecimal := strconv.FormatInt(12345, 16) // 16進数
	fmt.Println(dataHexadecimal)
}

// bytes.Joinのテスト
// 要import bytes
func TestBytesJoin(t *testing.T) {
	data1 := "hogehoge"
	byteArrayData1 := []byte(data1)
	fmt.Println(byteArrayData1)
	byteArrayData2 := []byte{1, 2, 3, 4, 5}
	fmt.Println(byteArrayData2)
	result := bytes.Join(
		[][]byte{byteArrayData1, byteArrayData2},
		[]byte{0x00, 0x00, 0x00, 0x00})
	fmt.Println(result)
}

// crypto.ハッシュ関数のテスト
// 要import crypto/md5, crypto/sha1, crypto/sha256, crypto/sha512
func TestHash(t *testing.T) {
	data := []byte("Hash me!")
	fmt.Printf("MD5 = % X\n", md5.Sum(data))
	fmt.Printf("SHA1 = % X\n", sha1.Sum(data))
	fmt.Printf("SHA256 = % X\n", sha256.Sum256(data))
	fmt.Printf("SHA512 = % X\n", sha512.Sum512(data))
}

// SetHashメソッドの動作検証
func TestSetHash(t *testing.T) {
	block := Block{} // 構造体を生成しblockに代入
	block.SetHash()
	fmt.Println(block.Hash)
}

func TestNewBlock(t *testing.T) {
	newBlock := NewBlock("hogehoge", []byte{})
	fmt.Println(newBlock.Hash)
}
