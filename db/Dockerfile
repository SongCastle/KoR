FROM mysql:8.4

COPY my.cnf /etc/mysql/conf.d/
COPY docker-entrypoint-initdb.d/ /docker-entrypoint-initdb.d/

RUN chmod -R +x /docker-entrypoint-initdb.d
