package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


func isNodeOnline(node string) bool {
    client := http.Client{
        Timeout: 2 * time.Second, // Tempo limite para a conexão
    }

    response, err := client.Get(fmt.Sprintf("%s/blockchain", node))
    if err != nil {
        return false // O nó está offline
    }
    defer response.Body.Close()

    return true // O nó está online
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockchain)
}

func startServer(port string) {
    http.HandleFunc("/blockchain", handleGetBlockchain)
	http.HandleFunc("/addblock", handleAddBlock)
	http.HandleFunc("/registerNode", handleRegisterNode)
	http.HandleFunc("/removeNode", handleRemoveNode)
    http.ListenAndServe(":"+port, nil)
}

func handleAddBlock(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
	}
	var requestData struct {
		Data string `json:"data"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Erro ao processar os dados", http.StatusBadRequest)
		return
	}

	newBlock := addBlock(requestData.Data, 2)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBlock)
}

func handleRegisterNode(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

	var requestData struct {
        Node string `json:"node"`
    }

	err := json.NewDecoder(r.Body).Decode(&requestData)

	if err != nil || requestData.Node == "" {
        http.Error(w, "Dados inválidos", http.StatusBadRequest)
        return
    }

	if requestData.Node == selfNode {
        http.Error(w, "Não é permitido registrar o próprio nó", http.StatusBadRequest)
        return
    }

	for _, node := range nodes {
        if node == requestData.Node {
            w.WriteHeader(http.StatusConflict) // 409 - Nó já registrado
            fmt.Fprintf(w, "Nó já registrado: %s", requestData.Node)
            return
        }
    }

	nodes = append(nodes, requestData.Node)
	saveNodesToFile()
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Nó registrado com sucesso: %s", requestData.Node)
	fmt.Println("Nós registrados atualmente:", nodes)
}

func handleRemoveNode(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    var requestData struct {
        Node string `json:"node"`
    }
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil || requestData.Node == "" {
        http.Error(w, "Dados inválidos", http.StatusBadRequest)
        return
    }

    // Remover o nó da lista
    for i, node := range nodes {
        if node == requestData.Node {
            nodes = append(nodes[:i], nodes[i+1:]...)
            saveNodesToFile() // Salva a lista atualizada
            fmt.Fprintf(w, "Nó removido com sucesso: %s", requestData.Node)
            fmt.Println("Nós registrados atualmente:", nodes)
            return
        }
    }

    http.Error(w, "Nó não encontrado", http.StatusNotFound)
}

func syncWithNetwork() {
    for _, node := range nodes {
		
		if node == selfNode {
            //fmt.Printf("Ignorando sincronização com o próprio nó: %s\n", selfNode)
			continue
        }

		if !isNodeOnline(node) {
            //fmt.Printf("Nó offline, ignorando: %s\n", node)
            continue
        }

        // Faz uma requisição GET ao endpoint /blockchain do nó
        response, err := http.Get(fmt.Sprintf("%s/blockchain", node))
        if err != nil {
            fmt.Printf("Erro ao conectar ao nó %s: %v\n", node, err)
            continue
        }
        defer response.Body.Close()

        var receivedBlockchain []Block
        err = json.NewDecoder(response.Body).Decode(&receivedBlockchain)
        if err != nil {
            fmt.Printf("Erro ao decodificar a blockchain do nó %s: %v\n", node, err)
            continue
        }

        // Verificar se a blockchain recebida é válida e mais longa
        if len(receivedBlockchain) > len(blockchain) && isBlockchainValid(receivedBlockchain) {
            blockchain = receivedBlockchain
			saveBlockchainToFile()
            fmt.Printf("Blockchain atualizada a partir do nó: %s\n", node)
        } 
    }
}

func startAutoSync(interval time.Duration) {
    ticker := time.NewTicker(interval) // Cria um ticker que dispara a cada intervalo definido
    go func() {
        for range ticker.C {
            //fmt.Println("Iniciando sincronização automática...")
            syncWithNetwork() // Chama a função de sincronização
        }
    }()
}
