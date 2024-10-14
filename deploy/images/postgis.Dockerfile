FROM postgres:17.0

RUN apt update && apt  install -y postgresql-17-postgis-3

CMD ["/usr/local/bin/docker-entrypoint.sh", "postgres"]