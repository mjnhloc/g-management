.PHONY: migrate-up
migrate-up:
	migrate -path db/migrations -database "mysql://gadmin:gpassword@tcp(127.0.0.1:3306)/gdb" up

.PHONY: migrate-down
migrate-up:
	migrate -path db/migrations -database "mysql://gadmin:gpassword@tcp(127.0.0.1:3306)/gdb" down

.PHONY: run-main
run-main:
	docker-compose up -d
