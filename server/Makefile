postgresinit:
	docker run --name postgres_cont -p 5444:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine

postgres:
	docker exec -it postgres_cont psql

createdb:
	docker exec -it postgres_cont createdb --username=root --owner=root go-chat

dropdb:
	docker exec -it postgres_cont dropdb go-chat

migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5444/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5444/go-chat?sslmode=disable" -verbose down

.PHONY: postgresinit postgres createdb dropdb migrateup migratedown