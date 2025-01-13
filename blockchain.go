package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)
var blockchain []Block

type Block struct {
	Index int
	Timestamp string
	Data string 
	PrevHash string
	Hash string
	Nonce int
	Difficulty int
}

func calculateHash(block Block) string{
	record := fmt.Sprintf("%d%s%s%s%d", block.Index,  block.Timestamp, block.Data, block.PrevHash, block.Nonce)
	hash := sha256.Sum256([]byte(record))

	return hex.EncodeToString(hash[:])
}

func createGenesisBlock() Block{
	genesisBlock := Block{
		Index : 0,
		Timestamp: time.Now().String(),
		Data: "Bloco Gênesis",
		PrevHash: "0",
		Nonce: 0,
		Difficulty: 0,
	}

	genesisBlock.Hash = calculateHash(genesisBlock)
	blockchain = append(blockchain, genesisBlock)

	return genesisBlock
}

func addBlock(data string, difficulty int) Block{
	
	if difficulty > 5 {
		fmt.Println("Dificuldade muito alta, ajustando para 5.")
		difficulty = 5
	}
	
	newBlock := Block{
		Index : (blockchain[len(blockchain)-1].Index + 1),
		Timestamp: time.Now().String(),
		Data: data,
		PrevHash: blockchain[len(blockchain)-1].Hash,
		Nonce: 0,
		Difficulty: difficulty,
	}
	runProofOfWork(&newBlock)
	blockchain = append(blockchain, newBlock)

	return newBlock
}

func isBlockchainValid() bool {
	for i := 1; i < len(blockchain); i++{
		currentBlock := blockchain[i]
		prevBlock := blockchain[i - 1]
		if currentBlock.PrevHash != prevBlock.Hash || calculateHash(currentBlock) != currentBlock.Hash {return false}
	}
	return true
}

func saveBlockchainToFile(){
	file, err := os.Create("data.json") // Tenta criar o arquivo
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return // Encerra a função em caso de erro
	}
	defer file.Close()
	jsonData, err := json.MarshalIndent(blockchain, "", "  ")
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo:", err)
		return
	}
	fmt.Println("Blockchain salva com sucesso em data.json")
}

func repairBlockchain(){
	originalLength := len(blockchain)
	for i := 1; i < len(blockchain); i++{
		currentBlock := blockchain[i]
		prevBlock := blockchain[i - 1]
		if currentBlock.PrevHash != prevBlock.Hash || calculateHash(currentBlock) != currentBlock.Hash {
			fmt.Printf("Bloco %d é inválido. Removendo blocos subsequentes...\n", currentBlock.Index)
			blockchain = blockchain[:i] // Remove todos os blocos após o inválido
			fmt.Printf("Blockchain reparada. %d blocos removidos.\n", originalLength-i)
			saveBlockchainToFile()     // Salva a blockchain corrigida
			return
		}
	}
	fmt.Println("Blockchain está válida.")
}

func loadBlockchainFromFile(){
	file, err := os.Open("data.json")
	if err != nil {
		fmt.Println("Arquivo não encontrado. Criando nova blockchain com bloco gênesis...")
		createGenesisBlock()
		return
	}
	defer file.Close()
	blockchainData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}
	err = json.Unmarshal(blockchainData, &blockchain)
	if err != nil {
		fmt.Println("Erro ao desserializar JSON:", err)
		return
	}
	fmt.Println("Blockchain carregada com sucesso!")
	if isBlockchainValid() {
        fmt.Println("Blockchain é válida.")
    } else {
        fmt.Println("Erro: Blockchain carregada é inválida.")
		repairBlockchain()
    }
}

func runProofOfWork(block *Block){
	prefix := ""
    for i := 0; i < block.Difficulty; i++ {
        prefix += "0"
    }
	start := time.Now()
	for{
		block.Hash = calculateHash(*block)
		if block.Hash[:block.Difficulty] == prefix{
			fmt.Printf("Bloco minerado em %s! Nonce: %d, Hash: %s\n", time.Since(start), block.Nonce, block.Hash)
            break	
		}
		block.Nonce++
		if block.Nonce%10000 == 0 { // Log a cada 10.000 tentativas
            fmt.Printf("Tentativa atual: Nonce=%d, Hash=%s\n", block.Nonce, block.Hash)
        }
	}
}


func main(){
	loadBlockchainFromFile()

	addBlock("aaa", 2)
	addBlock("bbb", 3)
 
    saveBlockchainToFile()
}