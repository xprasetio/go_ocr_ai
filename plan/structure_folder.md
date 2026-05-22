├── cmd
│   └── api
│       └── main.go
├── config
│   └── config.go
├── docker-compose.yml
├── Dockerfile
├── env.example
├── go.mod
├── go.sum
├── internal
│   ├── domain
│   │   ├── ocr
│   │   │   └── service.go
│   │   └── ocr
│   │       ├── entity.go
│   │       ├── repository.go
│   │       └── usecase.go
│   ├── infrastructure
│   │   ├── container
│   │   │   └── container.go
│   │   ├── database
│   │   │   └── sqlite.go
│   │   └── repository
│   │       └── ocr_repository.go
│   └── interfaces
│       └── http
│           ├── handler
│           │   └── ocr_handler.go
│           └── router.go
├── README.md