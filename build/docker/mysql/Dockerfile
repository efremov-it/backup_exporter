FROM mysql:8.0.35-debian

RUN apt update \
    && apt install wget lsb-release curl -y  \
    && wget https://repo.percona.com/apt/percona-release_latest.$(lsb_release -sc)_all.deb \
    && dpkg -i percona-release_latest.$(lsb_release -cs)_all.deb \
    && apt update && apt install percona-xtrabackup-80 -y --fix-missing \
    && wget https://github.com/wal-g/wal-g/releases/download/v2.0.1/wal-g-mysql-ubuntu-20.04-amd64.tar.gz \
    && tar -zxvf wal-g-mysql-ubuntu-20.04-amd64.tar.gz \
    && mv wal-g-mysql-ubuntu-20.04-amd64 /usr/local/bin/wal-g
    
## Change default mysql home directory and allow mysql write to /xtrabackup_backupfiles
RUN usermod -d /var/lib/mysql mysql \
    && mkdir /xtrabackup_backupfiles \
    && chown mysql:mysql /xtrabackup_backupfiles
