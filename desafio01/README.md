# Challenge 01

## 1. Server

- First, go to the server directory:

```bash
cd server
```

- Install the necessary dependencies:

```bash
go mod tidy
```

- Start the server:

```bash
go run main.go
```

The server will run on port 8080 and automatically generate a cotacoes.db SQLite database file in the same folder. This file will hold the exchange rate data pulled from the external API.

## 2. Client

- Now, open a new session on your VScode terminal and go to the client directory:

```bash
cd client
```

- Install the necessary dependencies:

```bash
go mod tidy
```

- Start the client:

```bash
go run main.go
```

## 3. Check results

- The client will send a request to the server to fetch the dollar exchange rate and save it to a file named cotacao.txt in the same directory. The file will contain something like:

```bash
Dólar: 5.5725
```

- The client will response:

```bash
2025/06/10 21:26:19 Cotação salva com sucesso: Dólar: 5.5725
```