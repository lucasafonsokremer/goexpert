package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal("Erro ao criar request:", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Erro: contexto encerrado por timeout")
			return
		}
		log.Println("Erro ao fazer requisição:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Erro ao ler resposta:", err)
		return
	}

	var cot CotacaoResponse
	if err := json.Unmarshal(body, &cot); err != nil {
		log.Println("Erro ao decodificar JSON:", err)
		log.Println("Resposta:", string(body))
		return
	}

	content := "Dólar: " + cot.Bid
	err = os.WriteFile("cotacao.txt", []byte(content), 0644)
	if err != nil {
		log.Println("Erro ao salvar arquivo:", err)
	}

	log.Println("Cotação salva com sucesso:", content)
}
