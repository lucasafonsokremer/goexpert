# Desafio 2 - Desafio Multi Threading

Este projeto é um programa Go que demonstra o uso de concorrência para otimizar a busca por informações de CEP. Ele envia requisições simultâneas para a BrasilAPI e a ViaCEP, utilizando goroutines para as consultas e canais para gerenciar as respostas. A lógica é simples: a primeira API que responder tem seu resultado retornado. Se, após 1 segundo, nenhuma das APIs responder, o processo é cancelado com um timeout.

## Como executar

```sh
go run main.go <CEP>
```

Exemplo:

```sh
go run main.go 01001-000
```

ou

```sh
go run main.go 01001000
```

## Requisitos
- Go 1.18 ou superior
- Acesso à internet