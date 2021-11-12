package main

import "fmt"
import "github.com/C-Qingpi/GO-BlockChain-Test-2021/blockchain"
import "strconv"

func PrintChain(chain *blockchain.BlockChain) {
	for _, block := range chain.Blocks {

		fmt.Printf("-----------\n")
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Printf("nonce: %v\n", block.Nonce)
		fmt.Printf("hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
	fmt.Printf("==========================\n")
	fmt.Printf("difficulty: %d\n", blockchain.Difficulty)
	fmt.Printf("==========================\n")
}
func main() {
	chain := blockchain.InitBlockChain()
	chain.AddBlock("first block after genesis")
	chain.AddBlock("second block after genesis")
	chain.AddBlock("third block after genesis")

	PrintChain(chain)
}
