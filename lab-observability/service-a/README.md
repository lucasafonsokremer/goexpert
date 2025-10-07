# Service A - Serviço de Entrada de CEP

Este serviço é responsável por receber e validar CEPs (Código de Endereçamento Postal) brasileiros antes de encaminhá-los para o Serviço B para obter informações meteorológicas.

## API Endpoints

### POST /cep

Recebe um CEP como entrada, valida o formato e encaminha para o Serviço B.

## Pré-requisitos

- Docker e Docker Compose
- Go 1.21+ (apenas para desenvolvimento)

## Configuração

As variáveis de ambiente necessárias já estão configuradas no docker-compose.yml.

## Executando com Docker Compose

Para iniciar a aplicação:
```bash
docker-compose up --build
```

A API estará disponível em `http://localhost:8080`

## Exemplos de Uso

1. **CEP Válido** (retorna 200):
```bash
curl -i -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'
```
Resposta:
```json
{
    "temp_C": 25.0,
    "temp_F": 77.0,
    "temp_K": 298.15
}
```

2. **CEP Inválido** (retorna 422):
```bash
curl -i -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "1234567"}'  # CEP com formato inválido
```
Resposta:
```json
{
    "message": "invalid zipcode"
}
```

3. **CEP Não Encontrado** (retorna 404):
```bash
curl -i -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "99999999"}'  # CEP inexistente
```
Resposta:
```json
{
    "message": "zipcode not found"
}
```

## Variáveis de Ambiente

- `WEATHER_SERVICE_URL`: URL do Serviço B (padrão: http://service-b:8080)
- `PORT`: Porta do servidor (padrão: 8080)
