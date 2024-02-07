postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres createdb  --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up 2>&1 | tee migrate.log

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down


execdb:
	docker exec -it postgres14 psql -U root -d simple_bank

sqlc:
	sqlc generate

cleandb:
	docker exec -it postgres14 psql -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" ${PSQL_URL}


sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go "git/db/sqlc" Store  

.PHONY:postgres createdb dropdb migrateup migratedown sqlc mock