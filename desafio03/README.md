# Order System

Este projeto implementa um sistema de pedidos com três interfaces: HTTP, GraphQL e gRPC.

## Importante

A migration está configurada para executar 30s após a inicialização dos containers

## Como executar tudo com Docker

1. **Build e subida dos containers (incluindo aplicação Go):**
   ```sh
   docker-compose up --build -d
   ```
   Isso irá subir MySQL, RabbitMQ, rodar as migrações e iniciar a aplicação Go (`ordersystem`).

2. **A aplicação Go estará disponível na porta 8080.**

## Endpoints disponíveis

- **HTTP:** Porta `8000`
- **GraphQL:** Porta `8080`
- **gRPC:** Porta `50051`

## Como testar

### HTTP

Utilize o arquivo `api/create_order.http` para testar as chamadas HTTP (pode usar o VS Code com extensão REST Client).

Após executar o POST, você deveria receber a seguinte resposta:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 30 Sep 2025 22:39:58 GMT
Content-Length: 66
Connection: close

{
  "orders": [
    {
      "id": "a",
      "price": 100.5,
      "tax": 0.5,
      "final_price": 101
    }
  ]
}
```

### gRPC

Use um cliente como [Evans](https://github.com/ktr0731/evans):

```sh
evans -r repl
```
Conecte em `localhost:50051` e utilize os métodos disponíveis.

### GraphQL

Acesse no navegador: [http://localhost:8080](http://localhost:8080)

#### Exemplo de mutation para criar um pedido:

```graphql
mutation {
  createOrder(input: {
    id:"1",
    Price: 100.0,
    Tax: 10.0
  }) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

#### Exemplo de query para listar pedidos da página 1:

```graphql
query {
  listOrders(page: 1) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

---

Pronto! Com esses passos você consegue executar e testar todas as interfaces do projeto.
