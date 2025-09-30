package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func buscaAPI(cep string, apiFmtAddr string, apiCh chan<- string) {
	url := fmt.Sprintf(apiFmtAddr, cep)
	resp, err := http.Get(url)
	if err != nil {
		apiCh <- fmt.Sprintf("Erro ao acessar %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		apiCh <- fmt.Sprintf("Erro ao ler resposta de %s: %v", url, err)
		return
	}

	apiCh <- fmt.Sprintf("Resposta de %s: %s", url, string(body))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <CEP>")
		fmt.Println("Exemplo de formato: 01001-000 ou 01001000")
		return
	}

	cep := os.Args[1]
	fmt.Println("CEP informado:", cep)

	brasilApiCh := make(chan string)
	viaCepCh := make(chan string)

	go buscaAPI(cep, "https://brasilapi.com.br/api/cep/v1/%s", brasilApiCh)
	go buscaAPI(cep, "http://viacep.com.br/ws/%s/json/", viaCepCh)

	select {
	case res := <-brasilApiCh:
		fmt.Println("Primeira resposta:", res)
	case res := <-viaCepCh:
		fmt.Println("Primeira resposta:", res)
	case <-time.After(time.Second):
		fmt.Println("Timeout: Nenhuma API respondeu em 1 segundo.")
	}
}
