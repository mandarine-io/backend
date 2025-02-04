# Развертывание в Docker

Мы используем [Docker Compose](https://docs.docker.com/compose/) для развертывания проекта в Docker. В репозитории есть
конфигурационные файлы, которые нужно использовать при развертывании в различных окружениях (`local`, `dev`,
`integration-test`, `e2e-test`).

## `local`

В `local` окружении используется `docker-compose.local.yml` файл. Он содержит контейнеры:

+ `postgres` - база данных с установленным расширением [*PostGIS*](https://postgis.net/)
+ `redis` - кэш и Pub/Sub с встроенным в образ UI для просмотра содержимого кеша и логов публикаций
+ `minio` - хранилище медиафайлов
+ `mailhog` - почтовый сервис с встроенным в образ UI для просмотра содержимого почтовых сообщений

> [!NOTE]
> Особенность этого окружения в том, что все порты контейнером проброшены в `localhost`.

> [!NOTE]
> Все учетные данные хранятся и задаются в самом конфигурационном файле.

## `dev`

В `dev` окружении используется `docker-compose.dev.yml` файл. Он содержит контейнеры:

+ `postgres` - база данных с установленным расширением [*PostGIS*](https://postgis.net/)
+ `redis` - кэш и Pub/Sub с встроенным в образ UI для просмотра содержимого кеша и логов публикаций
+ `minio` - хранилище медиафайлов
+ `prometheus` - база данных для хранения метрик
+ `grafana` - веб-сервер для визуализации метрик
+ `traefik` - API gateway с встроенным в образ UI для балансировки между монолитом и UI сервисов

> [!NOTE]
> В этом окружении все порты контейнером не проброшены, исключением являются порты Traefik, доступ к UI сервисов и
> монолиту происходят через него.

> [!NOTE]
> Все учетные данные задаются в конфигурационном файле через переменные окружения.

## `integration-test`

В `integration-test` окружении используется `docker-compose.integration-test.yml` файл. Он содержит контейнеры:

+ `postgres` - база данных с установленным расширением [*PostGIS*](https://postgis.net/)
+ `redis` - кэш и Pub/Sub с встроенным в образ UI для просмотра содержимого кеша и логов публикаций
+ `minio` - хранилище медиафайлов
+ `mailhog` - почтовый сервис с встроенным в образ UI для просмотра содержимого почтовых сообщений

> [!NOTE]
> В этом окружении все порты контейнером проброшены в `localhost`.

> [!NOTE]
> Все учетные данные хранятся и задаются в конфигурационном файле через переменные окружения.

## `e2e-test`

В `e2e-test` окружении используется `docker-compose.e2e-test.yml` файл. Он содержит контейнеры:

+ `postgres` - база данных с установленным расширением [*PostGIS*](https://postgis.net/)
+ `redis` - кэш и Pub/Sub с встроенным в образ UI для просмотра содержимого кеша и логов публикаций
+ `minio` - хранилище медиафайлов
+ `mailhog` - почтовый сервис с встроенным в образ UI для просмотра содержимого почтовых сообщений
+ `server` - собранный и запущенный сервер системы

> [!NOTE]
> В этом окружении все порты контейнером проброшены в `localhost`.

> [!NOTE]
> Все учетные данные хранятся и задаются в конфигурационном файле через переменные окружения.

> [!NOTE]
> Конфигурация сервера задается через переменные окружения в `docker-compose.e2e-test.yml`.