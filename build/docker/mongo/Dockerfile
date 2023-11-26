FROM mongo:7.0.3 

RUN apt update \
    && apt install wget curl -y \
    && wget https://github.com/wal-g/wal-g/releases/download/v2.0.1/wal-g-mongo-ubuntu-20.04-amd64.tar.gz \
    && tar -zxvf wal-g-mongo-ubuntu-20.04-amd64.tar.gz \
    && mv wal-g-mongo-ubuntu-20.04-amd64 /usr/local/bin/wal-g
