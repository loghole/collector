# BUILD DATABASE
FROM yandex/clickhouse-server:20.5.3.27

COPY ./init.sql ./docker-entrypoint-initdb.d/init-db.sql

ENTRYPOINT ["/entrypoint.sh"]
