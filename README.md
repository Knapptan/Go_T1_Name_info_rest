# Структура проекта

```
├── Go_test_effective-mobile.pdf
├── README.md
├── cmd
│   └── main.go
├── go.mod
├── go.sum
├── ini.env                     # Конфигурационные переменные
├── internal
│   ├── clients
│   │   └── enrichment.go
│   ├── config
│   │   └── config.go           # Конфигурация
│   ├── handler
│   │   └── handlers.go         # HTTP-обработчики
│   ├── models
│   │   └── models.go           # Сущности и DTO
│   ├── repository              # Работа с БД
│   │   ├── person.go
│   │   └── postgres.go
│   └── service                 # Бизнес-логика
│       └── person.go
├── migration.bash
└── migrations                  # SQL-миграции
    ├── 000001_create_people_table.down.sql
    └── 000001_create_people_table.up.sql
```

Роутинг: chi

Конфигурация: godotenv

Миграции: golang-migrate

Логгер: zap

Swagger: swaggo
