FROM clickhouse/clickhouse-server:23.10.3.5-alpine

RUN apk add wget curl \
    && wget https://github.com/AlexAkulov/clickhouse-backup/releases/download/v2.4.5/clickhouse-backup-linux-amd64.tar.gz \
    && tar -xf clickhouse-backup-linux-amd64.tar.gz \
    && cp build/linux/amd64/clickhouse-backup /usr/local/bin/
