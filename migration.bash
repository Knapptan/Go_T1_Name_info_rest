migrate -path ./migrations -database "postgres://postgres:admin@localhost:5432/people_db?sslmode=disable" up

# migrate -path ./migrations -database "postgres://postgres:admin@localhost:5432/people_db?sslmode=disable" down