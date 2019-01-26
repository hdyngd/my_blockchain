package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CLI responsible for processing command line arguments
type CLI struct {
	bc *Blockchain
}

// 使用方法についての出力
func (cli *CLI) printUsage() {
	fmt.Println("使用方法:")
	fmt.Println("  getbalance -address ADDRESS" +
		"- アドレスの残高を表示する")
	fmt.Println("  createblockchain -address ADDRESS " +
		"- ブロックチェーンを生成し初期ブロック報酬をアドレスに送信する")
	fmt.Println("  printchain - ブロックチェーンの全てのブロックを出力する")
	fmt.Println("  send -from 送信元アドレス " +
		"-to 送信先アドレス -amount 送金額 - fromからtoへコインを送金する")
}

// パラメータの検証
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// ブロックの追加コマンドの処理
// func (cli *CLI) addBlock(data string) {
// 	cli.bc.AddBlock(data)
// 	fmt.Println("Success!")
// }

// チェーンの出力

func (cli *CLI) printChain() {
	bc := NewBlockchain("")
	defer bc.db.Close()
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

// ブロックチェーンを生成する
func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("Done!")
}

// 残高を取得し出力する
func (cli *CLI) getBalance(address string) {
	// ブロックチェーンを生成
	bc := NewBlockchain(address)
	defer bc.db.Close()
	// 残高を初期化し、対象アドレスについての
	// 未使用トランザクション出力全てを取得
	balance := 0
	UTXOs := bc.FindUTXO(address)
	// 全てのValueの値の合計を求める
	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

// 送金処理
func (cli *CLI) send(from, to string, amount int) {
	// ブロックチェーンを取得
	bc := NewBlockchain(from)
	defer bc.db.Close()
	// 未使用トランザクション出力を用いて送金する
	tx := NewUTXOTransaction(from, to, amount, bc)
	// 新しいトランザクションを用いてブロックをマイニング
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("成功しました!")
}

// CLIの実行
func (cli *CLI) Run() {
	// 引数の検証
	cli.validateArgs()
	// addBlockコマンドの解析
	// addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	// printChainコマンドの解析
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 追加するブロックデータ
	// 第１引数: オプション名 第２引数: デフォルト値
	// 第３引数: 説明
	// addBlockData := addBlockCmd.String("data", "", "Block data")

	createBlockchainCmd :=
		flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createBlockchainAddress :=
		createBlockchainCmd.String(
			"address", "", "初期ブロックの報酬を送信するアドレス")

	// sendコマンドの対応
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "送信元ウォレットアドレス")
	sendTo := sendCmd.String("to", "", "送信先ウォレットアドレス")
	sendAmount := sendCmd.Int("amount", 0, "送金額")

	// getbalanceコマンドの対応
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "指定アドレスの残高表示")

	// コマンドラインの第１引数からコマンド名を判別
	switch os.Args[1] {
	// case "addblock": // ブロックの追加
	// 	err := addBlockCmd.Parse(os.Args[2:]) // 第２引数以降を取得
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	case "printchain": // チェーンの表示
		err := printChainCmd.Parse(os.Args[2:]) // 第２引数以降を取得
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain": // ブロックチェーンの作成
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send": // 送金処理
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default: // それ以外
		cli.printUsage() // 使い方の表示
		os.Exit(1)       // 終了
	}

	// if addBlockCmd.Parsed() { // addblockコマンドか？
	// 	if *addBlockData == "" { // パラメータがない場合
	// 		addBlockCmd.Usage() // 使い方を表示
	// 		os.Exit(1)
	// 	}
	// 	cli.addBlock(*addBlockData) // addblockコマンドを実行
	// }

	if printChainCmd.Parsed() { // printchainコマンドか？
		cli.printChain() // チェーンの表示
	}

	if createBlockchainCmd.Parsed() { // createBlockchainコマンドか？
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() { // sendコマンドの場合
		// 検証
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		// 送金処理
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
	if getBalanceCmd.Parsed() { // getbalanceコマンドの場合
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

}
