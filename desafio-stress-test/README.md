# Stress Test CLI

Objetivo: Criar um sistema CLI em Go para realizar testes de carga em um serviço web. O usuário deverá fornecer a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas.

## Funcionalidades

- Execução de testes de carga configuráveis via CLI
- Controle de concorrência para simular múltiplos usuários
- Relatório detalhado com métricas de performance
- Distribuição de códigos de status HTTP
- Taxa de sucesso e tempo total de execução

## Parâmetros

- `--url`: URL do serviço a ser testado (obrigatório)
- `--requests`: Número total de requisições (obrigatório)
- `--concurrency`: Número de chamadas simultâneas (obrigatório)

## Como Usar

### Build da Imagem Docker

```bash
docker build -t stress-test .
```

### Executando o Teste

**Exemplo básico:**
```bash
docker run stress-test --url=http://google.com --requests=1000 --concurrency=10
```

## Relatório

O sistema gera um relatório com as seguintes informações:

- Tempo total gasto na execução
- Quantidade total de requisições realizadas
- Número de requisições com status HTTP 200
- Distribuição de outros códigos de status HTTP (como 404, 500, etc.)

## Exemplo de Saída

```
Iniciando teste de carga...
URL: http://google.com
Total de Requests: 1000
Concorrência: 10

==========================================
          RESULTADOS DO TESTE DE CARGA
==========================================

Tempo total gasto: 46.981953342s
Quantidade total de requests realizados: 1000
Requests com status HTTP 200: 1000
Taxa de sucesso: 100.00%

Distribuição de códigos de status HTTP:
  HTTP 200: 1000

==========================================
```