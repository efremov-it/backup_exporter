.DEFAULT_GOAL := up

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

exec:
	docker compose exec postgres bash

create:
	docker compose exec -u postgres postgres wal-g backup-push /var/lib/postgresql/data

all: down up
	@make -s create
	docker compose logs postgres

test:
	go build -o app/main cmd/main.go
	docker compose exec -u postgres postgres /app/main --project test --backup_type psql --backup_cron "*/1 * * * *" --delete_cron "*/1 * * * *" --delete_retain 5

tc:
	go build -o app/main cmd/main.go
	docker compose exec -u postgres postgres /app/main --project test --backup_type psql

tm:
	go build -o app/main cmd/main.go
	docker compose exec -u mysql mariadb /app/main --project test --backup_type mariadb --backup_cron "*/1 * * * *"
tmy:
	go build -o app/main cmd/main.go
	docker compose exec -u mysql mysql /app/main --project test --backup_type mysql --backup_cron "*/1 * * * *"
tmd:
	go build -o app/main cmd/main.go
	docker compose exec -u mongodb mongo /app/main --project test --backup_type mongodb --backup_cron "*/1 * * * *" --delete_cron "*/1 * * * *" --delete_retain 5
	


tch:
	go build -o app/main cmd/main.go
	docker compose exec -u clickhouse clickhouse /app/main --project test --backup_type clickhouse


.PHONY: up down build exec create all
