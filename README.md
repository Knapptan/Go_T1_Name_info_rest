# Структура проекта

```
├── README.md
├── cmd
│   └── main.go
├── go.mod
├── go.sum
├── ini.env         # Конфигурационные переменные
├── internal
│   ├── config      # Конфигурация
│   ├── handler     # HTTP-обработчики
│   ├── models      # Сущности и DTO
│   │   └── models.go
│   ├── repository  # Работа с БД
│   └── service     # Бизнес-логика
├── migrations      # SQL-миграции
│   ├── 000001_create_people_table.down.sql
│   ├── 000001_create_people_table.up.sql
│   └── mig_start.bash
└── pkg
    ├── clients     # Клиенты внешних API
    └── logger      # Логирование
```
