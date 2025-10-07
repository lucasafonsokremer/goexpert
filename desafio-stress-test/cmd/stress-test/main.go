package main

import (
	"flag"
	"fmt"

	"github.com/lucasafonsokremer/goexpert/desafio-stress-test/internal/loadtest"
	"github.com/lucasafonsokremer/goexpert/desafio-stress-test/internal/reporter"
)

func main() {
	// Definir flags CLI
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")

	flag.Parse()

	// Validar parâmetros
	if err := validateParams(*url, *requests, *concurrency); err != nil {
		fmt.Println(err)
		flag.Usage()
		return
	}

	// Exibir informações do teste
	printTestInfo(*url, *requests, *concurrency)

	// Executar teste de carga
	tester := loadtest.New(*url, *requests, *concurrency)
	report := tester.Run()

	// Exibir relatório
	rep := reporter.New()
	rep.Print(report)
}

func validateParams(url string, requests int, concurrency int) error {
	if url == "" {
		return fmt.Errorf("erro: --url é obrigatório")
	}

	if requests <= 0 {
		return fmt.Errorf("erro: --requests deve ser maior que 0")
	}

	if concurrency <= 0 {
		return fmt.Errorf("erro: --concurrency deve ser maior que 0")
	}

	return nil
}

func printTestInfo(url string, requests int, concurrency int) {
	fmt.Printf("Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Total de Requests: %d\n", requests)
	fmt.Printf("Concorrência: %d\n\n", concurrency)
}
