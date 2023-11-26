FROM postgres:14.8-bookworm

RUN apt update \
    && apt install wget curl -y \
    && wget https://github.com/wal-g/wal-g/releases/download/v2.0.0/wal-g-pg-ubuntu-20.04-amd64.tar.gz \
    && tar -zxvf wal-g-pg-ubuntu-20.04-amd64.tar.gz \
    && mv wal-g-pg-ubuntu-20.04-amd64 /usr/local/bin/wal-g 
