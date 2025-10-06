# Weather API Service with CEP Integration

Esta API fornece informações de temperatura baseadas em CEPs brasileiros. O serviço converte o CEP em nome da cidade e retorna a temperatura atual em Celsius, Fahrenheit e Kelvin.

## Pré-requisitos

- Docker e Docker Compose
- Go 1.21+ (apenas para desenvolvimento)
- Chave de API do WeatherAPI (cadastre-se em https://www.weatherapi.com/)

## Configuração

1. Crie o arquivo `.env` na raiz do projeto com sua chave da WeatherAPI:
```bash
echo "WEATHER_API_KEY=sua_chave_aqui" > .env
```

## Executando com Docker Compose

Para iniciar a aplicação:
```bash
docker-compose up --build
```

A API estará disponível em `http://localhost:8080`

## Exemplos de Uso

1. **CEP Válido** (retorna 200):
```bash
curl -i http://localhost:8080/weather/89221370  # Joinville/SC
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
curl -i http://localhost:8080/weather/1234567  # CEP com formato inválido
```
Resposta:
```json
{
    "message": "invalid zipcode"
}
```

3. **CEP Inexistente** (retorna 404):
```bash
curl -i http://localhost:8080/weather/99999999  # CEP com formato válido mas inexistente
```
Resposta:
```json
{
    "message": "can not find zipcode"
}
```

## Executando Testes

### Usando Go diretamente:

```bash
# Executar todos os testes
go test -v ./...
```

### Usando Dockerfile.test:

```bash
# Construir e executar testes em um container
docker build -f Dockerfile.test -t weather-api-tests .
docker run --env-file .env weather-api-tests
```

## Estrutura do Projeto

```
.
├── cmd
│   └── api
│       └── main.go          # Ponto de entrada da aplicação
├── internal
│   ├── handlers
│   │   └── weather_handler.go
│   └── services
│       ├── cep_service.go
│       ├── cep_service_test.go
│       └── weather_service.go
├── Dockerfile
├── Dockerfile.test
├── docker-compose.yml
└── README.md
```

## Tecnologias Utilizadas

- Go 1.21
- Gin Web Framework
- Docker & Docker Compose
- ViaCEP API
- WeatherAPI