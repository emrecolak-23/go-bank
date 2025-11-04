DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres18 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine
createdb:
	docker exec -it postgres18 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres18 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v ./...
server:
	go run main.go
mockdb:
	mockgen -package mockdb -destination db/mock/store.go github.com/emrecolak-23/go-bank/db/sqlc Store
db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown db_docs db_schema migratedown1 sqlc test server mockdb