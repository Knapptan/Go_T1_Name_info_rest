# Сервис управления персональными данными

REST API сервис для работы с персональными данными с обогащением информации из внешних источников

## Технологический стек

- **Язык**: Go 1.19+
- **Роутинг**: [chi](https://github.com/go-chi/chi)
- **База данных**: PostgreSQL
- **Конфигурация**: [godotenv](https://github.com/joho/godotenv)
- **Миграции**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Логирование**: [zap](https://github.com/uber-go/zap)
- **Документация**: [swaggo](https://github.com/swaggo/swag)

## Функционал

✅ Реализованные требования:

- Добавление новых людей с автоматическим обогащением данных
- Фильтрация и пагинация при получении данных
- Обновление и удаление записей
- Полная интеграция с внешними API:
  - Определение возраста (agify.io)
  - Определение пола (genderize.io)
  - Определение национальности (nationalize.io)
- Автоматическая генерация Swagger-документации
- Подробное логирование операций
- Конфигурация через .env файл
- Миграции базы данных

## Установка и запуск

### Предварительные требования

- Go 1.19+
- PostgreSQL 12+
- Утилита migrate ([инструкции по установке](https://github.com/golang-migrate/migrate))

### Шаги установки

1. Клонировать репозиторий:

git clone https://github.com/Knapptan/Go_T1_Name_info_rest
cd Go_T1_Name_info_rest

2. Установить зависимости:

go mod download

3. Настроить базу данных:

createdb people_db

4. Выполнить миграции:

migrate -path migrations -database postgres://user:password@localhost:5432/people_db?sslmode=disable up

5. Создать файл конфигурации .env или использовать готовый:

```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=people_db
```

6. Сгенерировать Swagger-документацию или использовать готовую:

swag init -g cmd/main.go --output docs

7. Запустить сервис:

go run cmd/main.go

# Примеры запросов

## Создание новой записи

curl -X POST http://localhost:8080/api/v1/people \
 -H "Content-Type: application/json" \
 -d '{"name": "Ivan", "surname": "Petrov"}'

## Получение списка людей

curl "http://localhost:8080/api/v1/people?page=1&limit=10&gender=male"

## Обновление данных

curl -X PATCH http://localhost:8080/api/v1/people/1 \
 -H "Content-Type: application/json" \
 -d '{"age": 35}'

## Удаление записи

curl -X DELETE http://localhost:8080/api/v1/people/1
Документация API

### После запуска сервиса документация доступна по адресу:

http://localhost:8080/swagger/index.html

![swagger](swagg.png)

## Структура проекта

```
├── Go_test_effective-mobile.pdf
├── README.md
├── cmd
│   └── main.go
├── docs                        # Swagger документация
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
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

## Особенности реализации

Архитектура: Чистая архитектура с разделением на слои

Логирование: Подробные логи операций с контекстом

Обработка ошибок: Единый формат ошибок API

Безопасность: Валидация входных данных
