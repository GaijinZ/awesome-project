FROM postgres:latest

COPY ./pkg/database/migrations /docker-entrypoint-initdb.d/
RUN chmod -R 755 /docker-entrypoint-initdb.d/

CMD ["docker-entrypoint.sh", "postgres"]
