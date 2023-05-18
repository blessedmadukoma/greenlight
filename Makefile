migrateup:
	migrate -path=./migrations -database "postgres://greenlight:greenlight@localhost/greenlight?sslmode=disable" -verbose up

.PHONY : migrate