CREATE DATABASE IF NOT EXISTS clickhouse;

USE clickhouse;

CREATE TABLE IF NOT EXISTS mytable (
    id Int64,
    name String
) ENGINE = MergeTree() ORDER BY id;

INSERT INTO mytable (id, name) VALUES (1, 'John'), (2, 'Alice');
