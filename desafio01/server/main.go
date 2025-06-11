package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

type APIResponse struct {
	USDBRL Cotacao `json:"USDBRL"`
}

func main() {
	db, err := sql.Open("sqlite", "./cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY, bid TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctxAPI, cancelAPI := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancelAPI()

		req, _ := http.NewRequestWithContext(ctxAPI, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			// Verifica se o erro foi causado por timeout do contexto
			if ctxAPI.Err() == context.DeadlineExceeded {
				log.Println("Erro: contexto encerrado por timeout ao chamar a API")
				http.Error(w, "Erro ao obter cotação: timeout", http.StatusRequestTimeout)
				return
			}
			log.Println("Erro na chamada da API:", err)
			http.Error(w, "Erro ao obter cotação", http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		var apiResp APIResponse
		if err := json.NewDecoder(res.Body).Decode(&apiResp); err != nil {
			log.Println("Erro ao decodificar JSON:", err)
			http.Error(w, "Erro ao processar resposta", http.StatusInternalServerError)
			return
		}

		// Salvar no banco com timeout de 10ms
		ctxDB, cancelDB := context.WithTimeout(r.Context(), 50*time.Millisecond)
		defer cancelDB()

		query := "INSERT INTO cotacoes (bid) VALUES (?)"
		_, err = db.ExecContext(ctxDB, query, apiResp.USDBRL.Bid)
		if err != nil {
			// Verifica se o erro foi causado por timeout do contexto
			if ctxDB.Err() == context.DeadlineExceeded {
				log.Println("Erro: contexto encerrado por timeout ao salvar no banco")
				http.Error(w, "Erro ao obter cotação: timeout", http.StatusRequestTimeout)
				return
			}
			log.Println("Erro ao salvar no banco:", err)
			http.Error(w, "Erro ao salvar no banco", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"bid": apiResp.USDBRL.Bid})
	})

	log.Println("Servidor iniciado em :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
