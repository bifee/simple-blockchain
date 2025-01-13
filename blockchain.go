package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)
var blockchain []Block

type Block struct {
	Index int
	Timestamp string
	Data string 
	PrevHash string
	Hash string
}

func calculateHash(block Block) string{
	record := fmt.Sprintf("%d%s%s%s", block.Index,  block.Timestamp, block.Data, block.PrevHash)
	hash := sha256.Sum256([]byte(record))

	return hex.EncodeToString(hash[:])
}

func createGenesisBlock() Block{
	genesisBlock := Block{
		Index : 0,
		Timestamp: time.Now().String(),
		Data: "Bloco Gênesis",
		PrevHash: "0",
	}

	genesisBlock.Hash = calculateHash(genesisBlock)
	blockchain = append(blockchain, genesisBlock)

	return genesisBlock
}

func addBlock(data string) Block{
	newBlock := Block{
		Index : (blockchain[len(blockchain)-1].Index + 1),
		Timestamp: time.Now().String(),
		Data: data,
		PrevHash: blockchain[len(blockchain)-1].Hash,
	}
	newBlock.Hash = calculateHash(newBlock)
	blockchain = append(blockchain, newBlock)

	return newBlock
}

func isBlockchainValid() bool {
	for i := 1; i < len(blockchain); i++{
		currentBlock := blockchain[i]
		prevBlock := blockchain[i - 1]
		if currentBlock.PrevHash != prevBlock.Hash {return false}
		if calculateHash(currentBlock) != currentBlock.Hash {return false}
	}
	return true
}
func main(){
	createGenesisBlock()

    // Adicionando novos blocos
    addBlock("Bloco 1: Transação A -> B")
    addBlock("Bloco 2: Transação C -> D")
    addBlock("Bloco 3: Transação E -> F")

    // Imprimindo todos os blocos na blockchain
    for i, block := range blockchain {
        fmt.Printf("Bloco %d: %+v\n", i, block)
    }
	fmt.Println("Blockchain é válida?", isBlockchainValid())

}