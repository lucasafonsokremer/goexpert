package loadtest

import (
	"net/http"
	"sync"
	"time"
)

// RequestResult representa o resultado de uma requisição HTTP
type RequestResult struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

// Report contém as métricas do teste de carga
type Report struct {
	TotalTime      time.Duration
	TotalRequests  int
	Status200      int
	StatusCodes    map[int]int
	SuccessRate    float64
	FailedRequests int
}

// LoadTester é responsável por executar os testes de carga
type LoadTester struct {
	url         string
	requests    int
	concurrency int
	client      *http.Client
}

// New cria uma nova instância de LoadTester
func New(url string, requests int, concurrency int) *LoadTester {
	return &LoadTester{
		url:         url,
		requests:    requests,
		concurrency: concurrency,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Run executa o teste de carga e retorna o relatório
func (lt *LoadTester) Run() Report {
	startTime := time.Now()

	results := make(chan RequestResult, lt.requests)
	var wg sync.WaitGroup

	// Controlar o número de goroutines simultâneas
	semaphore := make(chan struct{}, lt.concurrency)

	// Criar worker pool
	for i := 0; i < lt.requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Adquirir slot no semáforo
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Fazer request
			result := lt.makeRequest()
			results <- result
		}()
	}

	// Aguardar todas as goroutines terminarem
	go func() {
		wg.Wait()
		close(results)
	}()

	// Coletar resultados
	report := lt.collectResults(results)
	report.TotalTime = time.Since(startTime)

	return report
}

// makeRequest realiza uma requisição HTTP
func (lt *LoadTester) makeRequest() RequestResult {
	start := time.Now()
	resp, err := lt.client.Get(lt.url)
	duration := time.Since(start)

	if err != nil {
		return RequestResult{
			Duration: duration,
			Error:    err,
		}
	}
	defer resp.Body.Close()

	return RequestResult{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
	}
}

// collectResults coleta e processa os resultados das requisições
func (lt *LoadTester) collectResults(results chan RequestResult) Report {
	report := Report{
		StatusCodes: make(map[int]int),
	}

	for result := range results {
		report.TotalRequests++

		if result.Error != nil {
			report.FailedRequests++
			report.StatusCodes[0]++ // Código 0 para erros
		} else {
			report.StatusCodes[result.StatusCode]++
			if result.StatusCode == 200 {
				report.Status200++
			}
		}
	}

	if report.TotalRequests > 0 {
		report.SuccessRate = float64(report.Status200) / float64(report.TotalRequests) * 100
	}

	return report
}
