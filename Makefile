.DEFAULT_GOAL := up
LIST := psql mysql mongo mariadb

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

all: down config up 
	@make -s create
	@docker compose logs postgres
	@echo your developer Environment is ready.


config:
	@for i in ${LIST};do [ -f build/docker/$$i/config/.walg.json ] ||\
	cp build/docker/$$i/config/.walg.json.example build/docker/$$i/config/.walg.json;done
	@[ -f build/docker/clickhouse/config/config.yml ] ||\
	cp build/docker/clickhouse/config/config.yml.example build/docker/clickhouse/config/config.yml

test:
	go build -o app/backup_exporter cmd/main.go
	docker compose exec -u postgres postgres /app/backup_exporter --project test --backup_type psql --backup_cron "*/1 * * * *" --delete_cron "*/1 * * * *" --delete_retain 5 /usr/local/bin/wal-g 

tc:
	go build -o app/backup_exporter cmd/main.go
	docker compose exec -u postgres postgres /app/backup_exporter --project test --backup_type psqlq

tm:
	go build -o app/backup_exporter cmd/main.go
	docker compose exec -u mysql mariadb /app/backup_exporter --project test --backup_type mariadb --backup_cron "*/1 * * * *"
tmy:
	go build -o app/backup_exporter cmd/main.go
	docker compose exec -u mysql mysql /app/backup_exporter --project test --backup_type mysql --backup_cron "*/1 * * * *"
tmd:
	go build -o app/backup_exporter cmd/main.go
	docker compose exec -u mongodb mongo /app/backup_exporter --project test --backup_type mongodb --backup_cron "*/1 * * * *" --delete_cron "*/1 * * * *" --delete_retain 5
	


tch:
	go build -o app/backup_exporter cmd/main.go
	docker compose exec -u clickhouse clickhouse /app/backup_exporter --project test --backup_type clickhouse --backup_cron "*/1 * * * *" /usr/local/bin/clickhouse-backup

curl:
	docker compose exec postgres curl localhost:9023/metrics|grep backup


.PHONY: up down build exec create all
