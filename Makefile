##@ Build
# .PHONY: db-migration
db-migration:
	go run ./cmd/migrator --storage-path="./storage/aaa.db" --migrations-path="./migrations"  