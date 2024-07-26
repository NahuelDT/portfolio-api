
# Portfolio API

Portfolio API es un sistema de gestión de portafolios financieros que permite a los usuarios realizar operaciones de compra y venta de instrumentos financieros, así como gestionar su saldo y ver su portafolio actual.

## Características

- Gestión de órdenes de compra y venta (MARKET y LIMIT)
- Cálculo dinámico de saldos de usuario
- Manejo de múltiples instrumentos financieros
- Validaciones para garantizar la integridad de las operaciones
- API RESTful para interactuar con el sistema

## Tecnologías Utilizadas

- Go (Golang)
- PostgreSQL
- GORM (ORM para Go)
- Gin (Framework web para Go)

## Configuración del Proyecto

### Prerrequisitos

- Go 1.15 o superior
- PostgreSQL 12 o superior

### Instalación

1. Clonar el repositorio:
```bash 
git clone https://github.com/NahuelDT/portfolio-api.git
```
2. Navegar al directorio del proyecto:
```bash 
cd portfolio-api
```
3. Instalar las dependencias:
```bash 
go mod tidy
```
4. Configurar las variables de entorno (crear un archivo `.env` en la raíz del proyecto):
```bash 
DB_HOST=localhost
DB_USER=yourusername
DB_PASSWORD=yourpassword
DB_NAME=dbname
DB_PORT=5432
```
5. Iniciar la aplicación:
```bash 
go run cmd/api/main.go
```

## Estructura del Proyecto

```bash 
portfolio-api/
├── cmd
│   └── api
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── handlers
│   │   │   ├── order.go
│   │   │   ├── portfolio.go
│   │   │   └── search.go
│   │   ├── middleware
│   │   │   └── error_handler.go
│   │   └── routes.go
│   ├── config
│   │   └── database.go
│   ├── mocks
│   │   ├── repository
│   │   │   ├── InstrumentRepositorer.go
│   │   │   ├── MarketDataRepositorer.go
│   │   │   ├── OrderRepositorer.go
│   │   │   └── UserRepositorer.go
│   │   └── service
│   │       ├── OrderServicer.go
│   │       ├── PortfolioServicer.go
│   │       └── SearchServicer.go
│   ├── models
│   │   ├── instrument.go
│   │   ├── marketdata.go
│   │   ├── order.go
│   │   ├── portfolio.go
│   │   └── user.go
│   ├── repository
│   │   ├── instrument_repository.go
│   │   ├── interfaces.go
│   │   ├── marketdata_repository.go
│   │   ├── order_repository.go
│   │   └── user_repository.go
│   └── service
│       ├── interfaces.go
│       ├── order_service.go
│       ├── order_service_test.go
│       ├── portfolio_service.go
│       ├── portfolio_service_test.go
│       ├── search_service.go
│       └── search_service_test.go
├── README.md
└── tests
    └── functional
        └── order_test.go
```

## Uso de la API

Ejemplos de endpoints:

- `POST /api/orders`: Crear una nueva orden
- `POST /orders/:orderID/cancel`: Cancelar una orden
- `GET /api/portfolio/{userID}`: Obtener el portafolio de un usuario
- `GET /api/instruments`: Listar instrumentos disponibles

## Pruebas

Para ejecutar las pruebas unitarias
```bash 
go test ./internal/...
```
Para ejecutar las pruebas funcionales:
```bash 
go test ./tests/functional/...
```

## Postman Collection

Para facilitar las pruebas y la interacción con la API, `Portfolio_API.postman_collection.json` es la colección de Postman. Esta colección incluye todos los endpoints disponibles con ejemplos de requests preconfigurados.
