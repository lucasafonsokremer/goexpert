package reporter

import (
	"fmt"

	"github.com/lucasafonsokremer/goexpert/desafio-stress-test/internal/loadtest"
)

// Reporter é responsável por exibir os relatórios
type Reporter struct{}

// New cria uma nova instância de Reporter
func New() *Reporter {
	return &Reporter{}
}

// Print exibe o relatório formatado no console
func (r *Reporter) Print(report loadtest.Report) {
	fmt.Println("==========================================")
	fmt.Println("          RESULTADOS DO TESTE DE CARGA")
	fmt.Println("==========================================")
	fmt.Printf("\nTempo total gasto: %v\n", report.TotalTime)
	fmt.Printf("Quantidade total de requests realizados: %d\n", report.TotalRequests)
	fmt.Printf("Requests com status HTTP 200: %d\n", report.Status200)
	fmt.Printf("Taxa de sucesso: %.2f%%\n", report.SuccessRate)

	if report.FailedRequests > 0 {
		fmt.Printf("Requests com erro: %d\n", report.FailedRequests)
	}

	fmt.Println("\nDistribuição de códigos de status HTTP:")
	for statusCode, count := range report.StatusCodes {
		if statusCode == 0 {
			fmt.Printf("  Erros de conexão: %d\n", count)
		} else {
			fmt.Printf("  HTTP %d: %d\n", statusCode, count)
		}
	}

	fmt.Println("\n==========================================")
}
