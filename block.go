package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

// データ本体
type Block struct {
	Timestamp int64 // 作成日時
	// Data          []byte // 保存データ
	Transactions  []*Transaction // 保存トランザクション
	PrevBlockHash []byte         // 一つ前のブロックのハッシュ値
	Hash          []byte         // 上記を結合した結果のハッシュ値
	Nonce         int            // 採掘用のデータ
}

// ポインタレシーバを用いたハッシュ値の代入メソッド
// func (b *Block) SetHash() {
// 	// Int型は簡単にバイト配列に変換できないので、strconv.FormatIntメソッドを用いて
// 	// 10進数の文字列化を行い、その上でバイト配列化
// 	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
// 	// ハッシュ算出の対象を連結
// 	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
// 	// ハッシュアルゴリズムSHA256によるハッシュ値の算出
// 	hash := sha256.Sum256(headers)
// 	// メンバーHashに結果を代入
// 	// ハッシュ関数の結果は固定長のバイト配列[32]byteのため
// 	// 全要素スライス[:]を使わないとデータ型が不一致となる
// 	b.Hash = hash[:]
// }

// 新規ブロックを生成し、構造体Blockのポインタを返す
// タイムスタンプを割り当て、ハッシュを算出
//func NewBlock(data string, prevBlockHash []byte) *Block {
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	// Block構造体を初期化し、そのポインタを代入
	// 第1メンバー: 現在の時刻をTime型で取得し、Unixタイム型に変換
	//// 第2メンバー: 文字列をバイト配列に変換
	// 第2メンバー: transactionsをそのまま代入
	// 第3メンバー: 引数によって渡された一つ前のブロックのハッシュを代入
	// 第4メンバー: 最初は空の状態で生成
	// 第5メンバー: 最初は0を代入
	// block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	// ハッシュ値の代入処理
	//block.SetHash()
	pow := NewProofOfWork(block) // PoWを用いて採掘した上でハッシュとnonceを格納
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

// 前のハッシュを持たない初期ブロックの生成
// func NewGenesisBlock(coinbase *Transaction) *Block {
func NewGenesisBlock(coinbase *Transaction) *Block {
	// return NewBlock("初期ブロック", []byte{})
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// ブロックのシリアライゼーション
func (b *Block) Serialize() []byte {
	// バッファーとしてresultを宣言
	var result bytes.Buffer
	// バッファーのエンコーダーを生成
	encoder := gob.NewEncoder(&result)
	// ブロックのエンコード＝シリアライズ化
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// ブロックのデシアライゼーション
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

// ブロックハッシュを算出するために
// トランザクションデータ全てをハッシュ化
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte //  トランザクションごとのハッシュを格納
	var txHash [32]byte   //  結果
	// Transactionsの要素全てについて繰り返す
	for _, tx := range b.Transactions {
		// ID＝ハッシュ値を追記
		txHashes = append(txHashes, tx.ID)
	}
	// 全結合したデータをハッシュ化
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:] // 固定長のバイト配列から可変バイト配列に変換
}
