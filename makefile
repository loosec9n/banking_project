postgres:
	sudo docker run --name simplebanksql -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14-alpine

indb:
	sudo docker exec -it simplebanksql psql -U root simple_bank

createdb:
	sudo docker exec -it simplebanksql createdb --username=root --owner=root simple_bank

dropdb:
	sudo docker exec -it simplebanksql dropdb simple_bank

dbup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

dbdown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run .

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: postgres indb createdb dropdb dbup dbdown sqlc test proto server