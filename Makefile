##@ Build

db-migration:
	go run ./cmd/migrator --storage-path="./storage/aaa.db" --migrations-path="./migrations"  


test:
	go run ./cmd/migrator --storage-path="./tests/storage/test.db" --migrations-path="./migrations"  
	go run ./cmd/migrator --storage-path="./tests/storage/test.db" --migrations-path="./tests/migrations"  --migrations-table="migrations_test"
	go test ./...