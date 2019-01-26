package main

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

// Blockchainに追加本体の定義
type Blockchain struct {
	//blocks []*Block // Blockのポインタ配列
	// ブロックチェーンをblocksの代わりに
	// 最後のブロックのハッシュと
	// データベースのポインタを保存
	tip []byte   // 最後のブロックのハッシュ値
	db  *bolt.DB // データベース
}

// Blockchanの列挙構造体
type BlockchainIterator struct {
	currentHash []byte   // 現在のハッシュ
	db          *bolt.DB // データベース
}

// 新しいBlockを生成し、Blockchainに追加
// func (bc *Blockchain) AddBlock(data string) {
// 	// prevBlock := bc.blocks[len(bc.blocks)-1]
// 	// newBlock := NewBlock(data, prevBlock.Hash)
// 	// bc.blocks = append(bc.blocks, newBlock)

// 	var lastHash []byte // 最後のハッシュ値
// 	// キー"l"から最後のハッシュ値を取得。View=読み取り
// 	err := bc.db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte(blocksBucket))
// 		lastHash = b.Get([]byte("l"))
// 		return nil
// 	})
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	// 新規ブロックの生成
// 	newBlock := NewBlock(data, lastHash)
// 	err = bc.db.Update(func(tx *bolt.Tx) error {
// 		//blocksBucketを取得し、newBlockを追加
// 		b := tx.Bucket([]byte(blocksBucket))
// 		err := b.Put(newBlock.Hash, newBlock.Serialize())
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		// DB上の最終ブロックのハッシュを更新
// 		err = b.Put([]byte("l"), newBlock.Hash)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		// このプログラム上の最終ハッシュ値を更新
// 		bc.tip = newBlock.Hash
// 		return nil
// 	})
// }

// 新しいBlockchainの生成
// func NewBlockchain() *Blockchain {
// 	// // 初期ブロックを格納した新しいBlockchainを生成
// 	// blockchain := &Blockchain{[]*Block{NewGenesisBlock()}}
// 	// return blockchain

// 	var tip []byte // 最後のブロックのハッシュ
// 	// BoltDBのオープン
// 	db, err := bolt.Open(dbFile, 0600, nil)
// 	if err != nil {
// 		log.Panic(err) // エラー時は強制終了
// 	}

// 	// データベースを更新用に開く
// 	err = db.Update(func(tx *bolt.Tx) error {
// 		// 開いた後に、キーblocksBucketのバケットを取得
// 		b := tx.Bucket([]byte(blocksBucket))
// 		// 存在しない場合＝初めてのアクセス
// 		if b == nil {
// 			fmt.Println("ブロックチェーンが存在しません。新しいものを生成します...")
// 			genesis := NewGenesisBlock()
// 			// キーblocksBucketのバケットを生成
// 			b, err := tx.CreateBucket([]byte(blocksBucket))
// 			if err != nil {
// 				log.Panic(err)
// 			}
// 			// 生成した初期ブロックを、ハッシュ値をキーにしてDBに保存
// 			err = b.Put(genesis.Hash, genesis.Serialize())
// 			if err != nil {
// 				log.Panic(err)
// 			}
// 			// 生成した最後のハッシュ値を、"l"をキーにしてDBに保存
// 			err = b.Put([]byte("l"), genesis.Hash)
// 			if err != nil {
// 				log.Panic(err)
// 			}
// 			// tip = 最後のハッシュ値
// 			tip = genesis.Hash
// 		} else {
// 			// すでに存在している場合は最後のハッシュ値を読み込む
// 			tip = b.Get([]byte("l"))
// 		}
// 		return nil // エラーがない場合
// 	}) // Updateメソッドの終了
// 	// エラー処理
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	// Blockchain構造体を生成し返す
// 	bc := Blockchain{tip, db}
// 	return &bc
// }

// イテレータ
func (bc *Blockchain) Iterator() *BlockchainIterator {
	// 最終ブロックのハッシュ値とデータベースを格納し
	// 返す
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

// 指定されたブロックを取得し、ポインタを返す
// また、イテレータのcurrentHashについて、
// 一つ前のブロックのハッシュ値となるように更新
func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock) // ブロックの復元
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	// イテレータのcurrentHashを一つ前のブロックのハッシュ値に
	i.currentHash = block.PrevBlockHash
	return block
}

// データベースの存在を確認
// importに"os"を追加
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// 初期ブロックに書き込まれるデータ
const genesisCoinbaseData = "初期ブロックのデータ"

// ブロックチェーンの生成
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("ブロックチェーンはすでに存在します。")
		os.Exit(1)
	}

	var tip []byte // 最後のブロックのハッシュ
	// BoltDBのオープン
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err) // エラー時は強制終了
	}
	// データベースを更新用に開く
	err = db.Update(func(tx *bolt.Tx) error {
		// コインベーストランザクションを生成
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		// 初期ブロックの生成
		genesis := NewGenesisBlock(cbtx)
		// ブロック格納用バケットの生成
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}
		// 初期ブロックのシリアライズ化
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		// 初期ブロックのハッシュの記録
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		// グローバル変数tipに初期ブロックのハッシュを代入
		tip = genesis.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	// ブロックチェーン構造体を生成
	bc := Blockchain{tip, db}
	return &bc
}

// ブロックチェーンの初期化（復元）
func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("ブロックチェーンがありません。最初に作成してください。")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := Blockchain{tip, db}
	return &bc
}

// 未使用のトランザクションの探索
//	"encoding/hex"をimportに含める
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	// 未使用のトランザクション
	var unspentTXs []Transaction
	// 使用済みのトランザクション出力
	spentTXOs := make(map[string][]int)
	// 最後のブロックから逆順に探索するためのイテレータ
	bci := bc.Iterator()

	for {
		block := bci.Next() // 一つ前のブロック
		for _, tx := range block.Transactions {
			// それぞれのトランザクションのIDを取得
			txID := hex.EncodeToString(tx.ID)

		Outputs: // ループのラベル
			for outIdx, out := range tx.Vout {
				// 出力アウトプットは使用済みか？
				if spentTXOs[txID] != nil {
					// 既に対象トランザクションで使用済みの
					// 出力がある場合
					for _, spentOut := range spentTXOs[txID] {
						// 対象トランザクションの出力を確認し、
						// 一致するものがあるならば、次の出力に
						// 処理を移す
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// 未使用の場合かつ指定アドレスでアンロック
				// でいる場合
				if out.CanBeUnlockedWith(address) {
					// 未使用のトランザクションとして登録
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// トランザクションがコインベースでない時
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					// それぞれの入力を確認
					if in.CanUnlockOutputWith(address) {
						// 指定したアドレスで出力をアンロック可能な場合
						// 入力のトランザクションidを16進数で文字列化
						inTxID := hex.EncodeToString(in.Txid)
						// 入力にリンクした出力を使用済み出力として登録
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		// 初期ブロックならば終了
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

// 未使用のトランザクション出力を探す
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		// 未使用のトランザクションから
		// アドレスで使用可能なトランザクション出力を
		// UTXOs配列に追記
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// 支払可能なトランザクション出力の探索
// 支払い対象となるトランザクション出力を返す
func (bc *Blockchain) FindSpendableOutputs(
	address string, amount int) (int, map[string][]int) {
	// 未使用の出力
	unspentOutputs := make(map[string][]int)
	// アドレスから未使用のトランザクションを取得
	unspentTXs := bc.FindUnspentTransactions(address)
	// 累積値
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		// 未使用トランザクションごとの処理
		// IDを16進数文字列化
		txID := hex.EncodeToString(tx.ID)
		// それぞれの出力についてくり返す
		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				// 出力が指定アドレスでアンロック可能
				// かつ、累積値が指定数量未満の時だけ追加
				accumulated += out.Value
				// 未使用トランザクション出力をunspentOutputsに追加
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				// 指定数量よりも累積値が多くなったら終了
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs
}

// ブロックのマイニング
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte
	// データベースからトランザクションデータと
	// 最終ブロックのハッシュを取得
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	// 新規ブロックを作成
	newBlock := NewBlock(transactions, lastHash)
	// シリアライズ化を行い、データベースに保存
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		// 最終ブロックのハッシュを記録
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		// 最終ハッシュをbc上で参照できるように代入
		bc.tip = newBlock.Hash
		return nil
	})
}
