package main

import (
	"fmt"
)

func main() {
	fmt.Println("Welcome to My Blockchain Program v0.04! ")
	// // bc := NewBlockchain()
	// bc := CreateBlockchain("")
	// // 終了時データベースをクローズ
	// defer bc.db.Close()
	// cli := CLI{bc}
	cli := CLI{}
	cli.Run()
}
