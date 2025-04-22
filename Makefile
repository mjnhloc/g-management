# MIGRATE_CMD=migrate -path db/migrations -database "mysql://gadmin:gpassword@tcp(127.0.0.1:3306)/gdb"

# migrate-up:
# 	$(MIGRATE_CMD) up

# migrate-down:
# 	$(MIGRATE_CMD) down

run-main:
	docker-compose up --build -d
