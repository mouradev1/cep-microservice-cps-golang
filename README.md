# API CEP - Golang

API de consulta de CEPs com cache Redis, MongoDB e múltiplos provedores (ViaCEP, BrasilAPI, ApiCEP).

## Funcionalidades

- Consulta de CEPs usando ViaCEP, BrasilAPI e ApiCEP
- Cache com Redis (1h para encontrados, 4h para não encontrados)  
- Registro de requisições por IP no MongoDB
- Rotação automática de provedores

## Como rodar

### Local
```bash
# Configure .env
MONGODB_URI=mongodb://root:teste@localhost:27017/cepdb?authSource=admin
REDIS_URL=redis://localhost:6379
PORT=3000

# Execute
go run cmd/main.go
```

### Docker (2 instâncias + Load Balancer)
```bash
docker-compose up -d --build
```

## Endpoint

### `GET /cep/:cep`
```bash
curl http://localhost:3000/cep/01001000
```

**Resposta:**
```json
{
  "cep": "01001000",
  "logradouro": "Praça da Sé",
  "bairro": "Sé", 
  "localidade": "São Paulo",
  "uf": "SP",
  "ibge": "3550308",
  "ddd": "11"
}
```

## Portas Docker
- Load Balancer: `http://localhost:80`
- API 1: `http://localhost:8887`  
- API 2: `http://localhost:8888`

## Referências
- [ViaCEP](https://viacep.com.br/)
- [BrasilAPI - CEP](https://brasilapi.com.br/docs#tag/CEP)
- [ApiCEP](https://apicep.com/)
- [Fiber Framework](https://fiber.gofiber.io/)
- [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [Redis Go Client](https://redis.uptrace.dev/)