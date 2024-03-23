[Unit]
Description=wal-g/clickhouse-backup Wrapper-exporter
After=network.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/srv/backup_exporter
# Command to start your application. You need to change for your needs.
ExecStart=/srv/backup_exporter/backup_exporter --project ch --backup_type clickhouse --port 9032 --backup_cron "32 1 * * *"
Restart=on-failure
StandardOutput=append:/var/log/backup_exporter.log
StandardError=append:/var/log/backup_exporter_err.log

[Install]
WantedBy=multi-user.target
