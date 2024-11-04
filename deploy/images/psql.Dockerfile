FROM debian:stable-20241016-slim
RUN apt update && apt install -y postgresql-client