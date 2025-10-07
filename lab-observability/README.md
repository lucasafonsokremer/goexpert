
# Demo de Serviços de Clima com Observabilidade Distribuída

Este projeto contém dois microsserviços que trabalham juntos para fornecer informações de clima baseadas em CEPs brasileiros. O sistema implementa rastreamento distribuído usando OpenTelemetry e Zipkin para observabilidade.

## Arquitetura

- **Service A**: Serviço de validação de entrada, que recebe e valida CEPs.
- **Service B**: Serviço de clima, que busca dados de localização no ViaCEP e dados de clima na WeatherAPI.

## Pré-requisitos

- Docker e Docker Compose
- Go 1.21+ (apenas para desenvolvimento)
- Chave de API do WeatherAPI (cadastre-se em https://www.weatherapi.com/)

## Configuração do Ambiente

1. Crie um arquivo `.env` na pasta service-b com sua chave da WeatherAPI:
   ```bash
   echo "WEATHER_API_KEY=sua_chave_aqui" > .env
   ```
   > **Atenção:** O serviço B depende dessa variável para funcionar corretamente.

## Como Executar

1. Clone o repositório
2. Inicie ambos os serviços com Docker Compose:

   ```bash
   docker-compose up --build
   ```

Isso irá subir:
- Service A na porta 8080
- Service B na porta 8081

## Testando os Serviços

### Exemplo de CEP Válido

```bash
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'
```
Resposta esperada:
```json
{
    "city":"Sao Caetano do Sul",
    "temp_C": 30.3,
    "temp_F": 86.53,
    "temp_K": 303.45
}
```

### Exemplo de CEP Inválido

```bash
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "123"}'
```
Resposta esperada:
```json
{
    "message": "invalid zipcode"
}
```

### Exemplo de CEP Não encontrado
```bash
curl -i http://localhost:8080/weather/99999999  # CEP com formato válido mas inexistente
```
Resposta:
```json
{
    "message": "can not find zipcode"
}
```

## Observabilidade

O sistema inclui os seguintes componentes:

### OpenTelemetry Collector
- Recebe traces de ambos os serviços
- Processa e exporta traces para o Zipkin
- Portas:
  - 4317: OTLP gRPC receiver
  - 4318: OTLP HTTP receiver
  - 8888: Endpoint de métricas

### Zipkin
- Visualização de rastreamento distribuído
- Acesse a interface em: http://localhost:9411

## Operações Rastreáveis
- Service A: Validação e encaminhamento de CEP
- Service B: Consulta de CEP e busca de clima
- Chamadas HTTP entre os serviços
- Chamadas externas (ViaCEP e WeatherAPI)
