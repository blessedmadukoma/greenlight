migrateup:
	migrate -path=./migrations -database "postgres://greenlight:greenlight@localhost/greenlight?sslmode=disable" -verbose up

migratedown:
	migrate -path=./migrations -database "postgres://greenlight:greenlight@localhost/greenlight?sslmode=disable" -verbose down

.PHONY : migrateup	migratedown