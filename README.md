# Бекапилка не несет ответственности за то в каком виде конфиг файлы.
Установка wal-g, clickhouse-backup осуществляется отдельно

Пример команды для запуска 
```
main --project projectName --backup_type <psql|mysql|mariadb|mongodb|clickhouse> --cron "* * * * *" 
Support only for postgresql --delete_cron "*/1 * * * *" --delete_retain 5
```

# Вся логика по созданию бекапа в create/create.go
Удаление старых бекапов поддерживается долько для postgresql
Проверка проходит на стадии задания переменных в config.go
Также здесь создаю `config.DeleteEnable` которая в main.go определяет состояние удаления старых бекапов

# Переменные создаются в config/config.go

# Метрики exporter/collector.go

# Крон запускает выполнение тасок cron/cron.go

# Локальная проверка
В build/docker лежат Докер файлы и конфиги для запуска и теста
Бекап предпологает наличие конфиг файлов в дирах по умолчанию (кроме mysql там установил home dir /var/lib/mysql руками)
Если нужно изменить место конфига, добовляем --config флаг.
Команды запускаются от имени пользователя бекапилки.
Обязательные поля для заполнения (в конфигах бекапилок)
- "AWS_ACCESS_KEY_ID": "",
- "AWS_SECRET_ACCESS_KEY": "",
- "AWS_REGION": "eu-central-1",
- "AWS_ENDPOINT": "",
- "WALE_S3_PREFIX": "",

