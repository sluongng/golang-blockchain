package main

import "fmt"
import "strconv"

func main() {
	bc := NewBlockChain()

	bc.AddBlock("Send 1 BTC to Son Luong")
	bc.AddBlock("Send 2 more to Son Luong")

	for _, block := range bc.blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

		fmt.Println()
	}
}
