# Бэкапилка не несет ответственности за то в каком виде конфиг файлы.
Установка [wal-g](https://github.com/wal-g/wal-g), [clickhouse-backup](https://github.com/Altinity/clickhouse-backup) осуществляется отдельно

Пример команды для запуска 
```
main --project projectName --backup_type <postgres|mysql|mariadb|mongodb|clickhouse> --cron "* * * * *" 
Support only for postgresql --delete_cron "*/1 * * * *" --delete_retain 5
```

в app/backup_exporter.service описан пример запуска команды
в зависимости от типа бекапа, нужно изменить пользователя, от чьего имени будет запускаться бекап.
User=postgres
Group=postgres
для clickhouse запускать от root.

## Флаги
  --host
  --port
	--project
	--instance
	--backup_type
	--config_file
	--backup_cron

# Вся логика по созданию бэкапа в create/create.go
Удаление старых бэкапов поддерживается только для postgresql
Проверка проходит на стадии задания переменных в config.go
Также здесь создаю `config.DeleteEnable` которая в main.go определяет состояние удаления старых бэкапов

# Переменные создаются в config/config.go

# Метрики exporter/collector.go

# Крон запускает выполнение тасок cron/cron.go

# Локальная проверка
В build/docker лежат Докер файлы и конфиги для запуска и теста
Бэкап предпологает наличие конфиг файлов в дирах по умолчанию (кроме mysql там установил home dir /var/lib/mysql руками)
Если нужно изменить место конфига, добовляем --config флаг.
Команды запускаются от имени пользователя бэкапилки.
Обязательные поля для заполнения (в конфигах бекапилок по умолчанию используется минио, поднятый через копоз)
- "AWS_ACCESS_KEY_ID": "",
- "AWS_SECRET_ACCESS_KEY": "",
- "AWS_REGION": "eu-central-1",
- "AWS_ENDPOINT": "",
- "WALE_S3_PREFIX": "",

`make all` --> Поднимет все окружение и проведет ручной запукс создания бекапа postgres
Дальше в Makefile указаны примеры использования для разных баз данных

Пример systemd файла в /app дире.
Программа должна запускаться от имени пользователя в зависимости от базы данных которую используем.

install

