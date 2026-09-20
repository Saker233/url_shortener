migrateup:
	migrate -path migrations/ -database "postgres://postgres:postgres@localhost:5433/url_shortener?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations/ -database "postgres://postgres:postgres@localhost:5433/url_shortener?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: migrateup migratedown sqlc