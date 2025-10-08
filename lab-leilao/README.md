# Auction System - Full Cycle Go Expert

Sistema de leilão com fechamento automático implementado em Go.

## Funcionalidades

- Criação de leilões
- Criação de lances (bids)
- **Fechamento automático de leilões após tempo configurável**
- Consulta de leilões e lances

## Nova Funcionalidade: Fechamento Automático de Leilões

O sistema agora possui uma rotina automática que monitora e fecha leilões que ultrapassaram o tempo definido.

### Implementação

A solução foi implementada no arquivo `internal/infra/database/auction/create_auction.go` e inclui:

1. **Função `getAuctionInterval()`**: Calcula o tempo de duração do leilão baseado na variável de ambiente `AUCTION_INTERVAL`
2. **Goroutine `monitorExpiredAuctions()`**: Monitora continuamente os leilões ativos e verifica se expiraram
3. **Função `checkAndCloseExpiredAuctions()`**: Busca leilões ativos e fecha aqueles que ultrapassaram o tempo definido
4. **Método `UpdateAuctionStatus()`**: Atualiza o status de um leilão no banco de dados
5. **Controle de concorrência**: Utiliza `sync.Mutex` para evitar race conditions ao acessar mapas compartilhados

### Como Funciona

1. Quando o `AuctionRepository` é criado, uma goroutine é iniciada automaticamente
2. A goroutine verifica a cada 10 segundos se existem leilões vencidos
3. Quando um leilão ultrapassa o tempo definido em `AUCTION_INTERVAL`, seu status é automaticamente alterado para `Completed`
4. Um mapa de controle evita múltiplas atualizações do mesmo leilão

## Requisitos

- Docker e Docker Compose
- (Opcional) Go 1.20+ para desenvolvimento local

## Variáveis de Ambiente

Configure as seguintes variáveis no arquivo `cmd/auction/.env`:

```env
BATCH_INSERT_INTERVAL=10s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=5m

MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

### Variável Principal

- **`AUCTION_INTERVAL`**: Define o tempo de duração dos leilões (ex: `20s`, `5m`, `1h`)
  - Padrão: `5m` (5 minutos) caso não seja definida ou seja inválida

## Como Executar

### 1. Subir o ambiente com Docker Compose

```bash
docker-compose up --build
```

Isso irá:
- Construir a imagem da aplicação
- Iniciar o MongoDB
- Iniciar a aplicação na porta 8080

### 2. A API estará disponível em

```
http://localhost:8080
```

## Endpoints da API

### Criar Leilão

Utilize o arquivo `api/create_auction.http` para testar as chamadas HTTP (pode usar o VS Code com extensão REST Client)

Após executar o POST, você deveria receber a seguinte resposta:

```
HTTP/1.1 201 Created
Date: Wed, 08 Oct 2025 11:33:58 GMT
Content-Length: 0
Connection: close
```

**Condições (condition):**
- `1`: New (Novo)
- `2`: Used (Usado)
- `3`: Refurbished (Recondicionado)

### Listar Leilões

**Importante**: Por padrão, o sistema está configurado para aguardar 5 minutos até colocar como inativo o leilão. Você confere se o leilão está ativo, através do Status 0.

Utilize o arquivo `api/get_auction.http` para testar as chamadas HTTP (pode usar o VS Code com extensão REST Client)

Após executar o GET, você deveria receber a seguinte resposta:

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Oct 2025 16:23:31 GMT
Content-Length: 205
Connection: close

[
  {
    "id": "26e4b886-4f86-484c-b9de-7b776ac51e73",
    "product_name": "Notebook Dell",
    "category": "Electronics",
    "description": "Notebook Dell Inspiron 15",
    "condition": 1,
    "status": 0,
    "timestamp": "2025-10-08T16:17:14Z"
  }
]
```

**Status:**
- `0`: Active (Ativo)
- `1`: Completed (Finalizado)

### Buscar Leilão por ID

```bash
curl localhost:8080/auction/26e4b886-4f86-484c-b9de-7b776ac51e73
```

Após executar o GET, você deveria receber a seguinte resposta:

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Oct 2025 16:23:36 GMT
Content-Length: 205
Connection: close

[
  {
    "id": "26e4b886-4f86-484c-b9de-7b776ac51e73",
    "product_name": "Notebook Dell",
    "category": "Electronics",
    "description": "Notebook Dell Inspiron 15",
    "condition": 1,
    "status": 0,
    "timestamp": "2025-10-08T16:17:14Z"
  }
]
```


### Criar Lance

Utilize o arquivo `api/create_bid.http` para testar as chamadas HTTP (pode usar o VS Code com extensão REST Client)

**Você precisa informar o auctionId no arquivo**

Após executar o POST, você deveria receber a seguinte resposta:

```bash
HTTP/1.1 201 Created
Date: Wed, 08 Oct 2025 16:24:02 GMT
Content-Length: 0
Connection: close
```

**Nota**: O sistema valida automaticamente se o leilão está aberto antes de aceitar o lance.

### Listar Lances de um Leilão

Após criar o o lance, aguarde 20s até o processo de atualização em lotes ocorrer. Utilize o arquivo `api/get_bid.http` para testar as chamadas HTTP (pode usar o VS Code com extensão REST Client)

**Você precisa informar o auctionId no arquivo**

Após executar o GET, você deveria receber a seguinte resposta:

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 08 Oct 2025 16:26:44 GMT
Content-Length: 197
Connection: close

[
  {
    "id": "65ac1a11-4b87-4461-a2dc-f51996399377",
    "user_id": "9d90c058-294c-47a8-8f81-229202359424",
    "auction_id": "26e4b886-4f86-484c-b9de-7b776ac51e73",
    "amount": 1500,
    "timestamp": "2025-10-08T16:20:37Z"
  }
]
```

## Executar Testes

Para executar os testes do fechamento automático:

```bash
# Com o MongoDB rodando
docker-compose up -d mongodb

# Executar os testes
go test ./internal/infra/database/auction/... -v

# Ou executar dentro do container
docker-compose exec app go test ./internal/infra/database/auction/... -v
```

### Testes Implementados

1. **TestAuctionAutoClose**: Verifica se o leilão é fechado automaticamente após o tempo expirar
2. **TestAuctionNotClosedBeforeExpiry**: Valida que o leilão permanece ativo antes do tempo expirar
3. **TestUpdateAuctionStatus**: Testa a função de atualização de status
4. **TestGetAuctionInterval**: Testa a leitura e parsing da variável de ambiente

```bash
=== RUN   TestAuctionAutoClose
--- PASS: TestAuctionAutoClose (14.04s)
=== RUN   TestAuctionNotClosedBeforeExpiry
--- PASS: TestAuctionNotClosedBeforeExpiry (5.05s)
=== RUN   TestUpdateAuctionStatus
--- PASS: TestUpdateAuctionStatus (0.04s)
=== RUN   TestGetAuctionInterval
--- PASS: TestGetAuctionInterval (0.00s)
PASS
ok      fullcycle-auction_go/internal/infra/database/auction    19.133s
```

## Estrutura do Projeto

```
.
├── cmd/
│   └── auction/
│       ├── main.go
│       └── .env
├── configuration/
│   ├── database/
│   ├── logger/
│   └── rest_err/
├── internal/
│   ├── entity/
│   │   ├── auction_entity/
│   │   └── bid_entity/
│   ├── infra/
│   │   ├── api/
│   │   └── database/
│   │       ├── auction/
│   │       │   ├── create_auction.go       # ⭐ IMPLEMENTAÇÃO PRINCIPAL
│   │       │   ├── create_auction_test.go  # ⭐ TESTES
│   │       │   └── find_auction.go
│   │       └── bid/
│   ├── usecase/
│   └── internal_error/
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Detalhes Técnicos

### Gerenciamento de Concorrência

O sistema utiliza as seguintes técnicas para garantir segurança em ambientes concorrentes:

1. **Mutex para mapas compartilhados**: `closedAuctionsMapMutex` protege o acesso ao mapa de leilões fechados
2. **Ticker com intervalo fixo**: Verifica leilões a cada 10 segundos para evitar overhead
3. **Mapa de cache**: Evita múltiplas atualizações do mesmo leilão no banco de dados

### Fluxo de Fechamento

```
[Criação do Repository]
        ↓
[Inicia goroutine de monitoramento]
        ↓
[A cada 10s: busca leilões ativos]
        ↓
[Calcula tempo de expiração]
        ↓
[Leilão expirado?] → Sim → [Atualiza status para Completed]
        ↓                           ↓
       Não                   [Adiciona ao mapa de cache]
        ↓                           ↓
   [Aguarda próximo ciclo]    [Loga fechamento]
```

## Parar o Ambiente

```bash
docker-compose down
```

Para remover os volumes (dados do MongoDB):

```bash
docker-compose down -v
```