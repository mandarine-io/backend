# Конфигурация YAML и ENV

Конфигурирование сервера осуществляется с помощью YAML файла и переменных среды. Переменные среды имеют больший
приоритет, чем конфигурационный файл, поэтому ими можно переопределить значения из файла.

Указанные в следующих разделах значения по умолчанию будут использоваться при отсутствии их в конфигурационном файле и
переменных среды.

## Кэш

Настройки кэша Redis (Предоставлены значения по умолчанию).

```yaml
cache:
    address:
    dbindex: 0
    password:
    username: default
    ttl: 86400
```

```dotenv
APP_CACHE_ADDRESS=
APP_CACHE_DBINDEX=0
APP_CACHE_PASSWORD=
APP_CACHE_USERNAME=default
APP_CACHE_TTL=86400
```

## База данных

Настройки базы данных PostgreSQL (Предоставлены значения по умолчанию).

```yaml
database:
    address:
    dbname:
    password:
    username:
```

```dotenv
APP_DATABASE_ADDRESS=
APP_DATABASE_NAME=
APP_DATABASE_PASSWORD=
APP_DATABASE_USERNAME=
```

## Локализация

Настройки локализации (По умолчанию: язык — `ru`, путь — `locales`).

```yaml
locale:
    language: ru
    path: locales
```

```dotenv
APP_LOCALE_LANGUAGE=ru
APP_LOCALE_PATH=locales
```

## Логгирование

Настройки логгера (Предоставлены значения по умолчанию).

```yaml
logger:
    level: info
    console:
        enable: true
        encoding: text
    file:
        enable: false
        dirpath: logs
        maxage:
        maxsize:
```

```dotenv
APP_LOGGER_LEVEL=info

APP_LOGGER_CONSOLE_ENABLE=true
APP_LOGGER_CONSOLE_ENCODING=text

APP_LOGGER_FILE_ENABLE=false
APP_LOGGER_FILE_DIRPATH=logs
APP_LOGGER_FILE_MAXAGE=
APP_LOGGER_FILE_MAXSIZE=
```

## Миграции

Настройки миграции БД (По умолчанию путь до директории скриптов - `migrations`).

```yaml
migrations:
    path: migrations
```

```dotenv
APP_MIGRATIONS_PATH=migrations
```

## OAuth-клиенты

Настройки клиентов OAuth 2.0 (Предоставлены значения для примера. Поддерживаются `google`, `yandex`, `mailru`).

```yaml
oauthclients:
    -   name: mock
        clientid: mock
        clientsecret: mock
```

```dotenv
APP_OAUTH_CLIENTS_0_NAME=mock
APP_OAUTH_CLIENTS_0_CLIENTID=mock
APP_OAUTH_CLIENTS_0_CLIENTSECRET=mock
```

## Клиенты геокодирования

Настройки клиентов геокодирования (Предоставлены значения для примера. Поддерживаются `here`, `locationiq`,
`graphhopper`, `osm_nominatim` `yandex`).

```yaml
geocodingclients:
    -   name: mock
        apikey: mock
```

```dotenv
APP_GEOCODING_CLIENTS_0_NAME=mock
APP_GEOCODING_CLIENTS_0_APIKEY=mock
```

## Pub/Sub

Настройки Redis Pub/Sub (Представлены значения по умолчанию).

```yaml
pubsub:
    address:
    dbindex: 0
    password:
    username: default
```

```dotenv
APP_PUBSUB_ADDRESS=
APP_PUBSUB_DBINDEX=0
APP_PUBSUB_PASSWORD=
APP_PUBSUB_USERNAME=default
```

## S3

Настройки MinIO S3.

```yaml
s3:
    address:
    accesskey:
    bucket:
    secretkey:
```

```dotenv
APP_S3_ADDRESS=
APP_S3_ACCESSKEY=
APP_S3_BUCKET=
APP_S3_SECRETKEY=
```

## Безопасность

Настройки безопасности (Предоставлены значения по умолчанию).

```yaml
security:
    jwt:
        accesstokenttl: 3600
        refreshtokenttl: 86400
        secret:
    otp:
        length: 6
        ttl: 300
```

```dotenv
APP_SECURITY_JWT_ACCESSTOKENTTL=3600
APP_SECURITY_JWT_REFRESHTOKENTTL=86400
APP_SECURITY_JWT_SECRET=

APP_SECURITY_OTP_LENGTH=6
APP_SECURITY_OTP_TTL=300
```

## Сервер

Настройки сервера (Предоставлены значения по умолчанию).

```yaml
server:
    externalorigin: http://localhost:8000
    mode: local
    name: server
    port: 8080
    version: 0.0.0
```

```dotenv
APP_SERVER_EXTERNALORIGIN=http://localhost:8000
APP_SERVER_MODE=local
APP_SERVER_NAME=server
APP_SERVER_PORT=8080
APP_SERVER_VERSION=0.0.0
```

## SMTP

Настройки SMTP (Предоставлены значения для примера).

```yaml
smtp:
    from: 'Mandarine <mandarine.app@yandex.ru>'
    host:
    port:
    ssl:
    username:
    password:
```

```dotenv
APP_SMTP_FROM='Mandarine <mandarine.app@yandex.ru>'
APP_SMTP_HOST=
APP_SMTP_PORT=
APP_SMTP_SSL=
APP_SMTP_USERNAME=
APP_SMTP_PASSWORD=
```

## Шаблоны

Настройки шаблона (По умолчанию путь до директории с шаблонами - `templates`).

```yaml
template:
    path: templates
```

```dotenv
APP_TEMPLATE_PATH=templates
```

## WebSocket

Настройки WebSocket пула (Размер пула по умолчанию: 1024).

```yaml
websocket:
    poolsize: 1024
```

```dotenv
APP_WEBSOCKET_POOLSIZE=1024
```