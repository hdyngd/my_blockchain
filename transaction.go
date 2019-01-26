package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID   []byte     // トランザクションのID
	Vin  []TXInput  // トランザクション入力
	Vout []TXOutput // トランザクション出力
}

// トランザクション出力
type TXOutput struct {
	// 値（通貨量、ビットコインの場合satoshi数）
	Value int
	//公開鍵Script
	ScriptPubKey string
}

// トランザクション入力
type TXInput struct {
	// トランザクションID
	Txid []byte
	// トランザクション出力のインデックス
	Vout int
	// 署名Script
	ScriptSig string
}

// トランザクションのIDとしてハッシュ値を代入
// 	"bytes" "crypto/sha256" "encoding/gob" "log" をimportに追加
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx) // Transactionをバイト配列化
	if err != nil {
		log.Panic(err)
	}
	// SHA256でハッシュ値を求める
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// 最初のトランザクション
// importに"fmt"を追加
const subsidy = 10 // 報酬額
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" { // data = ""の時、最初のトランザクションとして扱う
		data = fmt.Sprintf("'%s'に対する報酬", to)
	}
	txin := TXInput{[]byte{}, -1, data} // トランザクション入力の生成
	txout := TXOutput{subsidy, to}      // トランザクション出力の生成
	// トランザクションの生成
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID() // IDの割り当て
	return &tx
}

// そのアドレスがトランザクションを作成したか否かをチェック
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// 出力が提供データによってアンロックすることができたかをチェック
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// Coinbaseを判定
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && // 入力が1つしかない
		len(tx.Vin[0].Txid) == 0 && // 入力のTxidが0である
		tx.Vin[0].Vout == -1 // 入力のリンク元出力が-1である
}

// 送金処理のトランザクションの生成
// encoding/hexをimportに追加
func NewUTXOTransaction(
	from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	// 送金可能な金額を算出
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	// 送金可能額accが送金しようとしている
	// 金額amountよりも小さい場合、例外発生
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// トランザクション入力リストの生成
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		// 未使用出力それぞれにについて
		for _, out := range outs {
			// トランザクション入力構造体を生成
			// 出力と入力の間のリンクをはるリンク
			input := TXInput{txID, out, from}
			//
			inputs = append(inputs, input)
		}
	}
	//出力リストの作成
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		// ぴったりの金額出ない場合、最後の出力は差分値を代入
		outputs = append(outputs, TXOutput{acc - amount, from}) // 変更
	}
	// トランザクションの生成
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}
