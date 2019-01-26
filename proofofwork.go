package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

// Int64の最大値
var maxNonce = math.MaxInt64

// マイニング難易度
// 先頭何ビットが0となるようなnonceを採掘
const targetBits = 8

// 要import math/big
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// PoWの生成
func NewProofOfWork(b *Block) *ProofOfWork {
	// 1をbig.Int型に変換
	target := big.NewInt(1)
	// 256からtargetBits引いた分だけ左算術シフト
	target.Lsh(target, uint(256-targetBits))

	// ProofOfWorkの構造体を生成し、
	// そのポインタをpowに代入
	pow := &ProofOfWork{b, target}
	return pow
}

// nonceを代入してPoW比較対象の元データを作成
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	// 前ブロックのハッシュ、タイムスタンプ、データ、
	// マイニング難易度とnonceを連結したバイト配列の生成
	data := bytes.Join( //2次元バイト配列を連結し一つのバイト配列に
		[][]byte{
			pow.block.PrevBlockHash,
			//pow.block.Data,
			pow.block.HashTransactions(), // 追加
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{}, // 区切りデータ（空データ）
	)
	return data
}

// PoW作業
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int // big.Int型のハッシュ
	var hash [32]byte   // バイト配列のハッシュ
	nonce := 0

	fmt.Printf("ブロックの採掘 データ＝ \"%s\"\n", pow.block.Transactions)

	for nonce < maxNonce { // ０から最大値まで繰り返す
		data := pow.prepareData(nonce)                // nonceを含めたデータの結合
		hash = sha256.Sum256(data)                    // 結合データのハッシュ値算出(32バイト)
		fmt.Printf("\rnonce=%d:hash=%x", nonce, hash) // ハッシュ値出力（同一行に上書き出力）
		hashInt.SetBytes(hash[:])                     // バイト配列で値を代入
		if hashInt.Cmp(pow.target) == -1 {            // ハッシュ値とtargetを比較
			break // ハッシュ値の方が小さい場合終了
		} else {
			nonce++ // ハッシュ値の方が大きい場合次の数にトライ
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

// PoWのブロックの検証
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int // big.Int型の変数
	//
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data) // ハッシュ値を求める
	hashInt.SetBytes(hash[:])   // big.Intにハッシュ値を代入
	// 基準値よりも低ければtrue, そうだなければfalse
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
