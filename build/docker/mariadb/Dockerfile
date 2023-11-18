FROM mariadb:11.0.4-jammy

RUN apt update \
    && apt install wget vim lsb-release curl mariadb-backup -y  \
    && apt --fix-broken install \
    && wget https://github.com/wal-g/wal-g/releases/download/v2.0.1/wal-g-mysql-ubuntu-20.04-amd64.tar.gz \
    && tar -zxvf wal-g-mysql-ubuntu-20.04-amd64.tar.gz \
    && mv wal-g-mysql-ubuntu-20.04-amd64 /usr/local/bin/wal-g
