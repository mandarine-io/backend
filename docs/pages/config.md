# Configuration

<a id="server"></a>

## Server settings

### Name

Set server name (default: server)

```yaml
server:
    name: server
```

```dotenv
MANDARINE_SERVER_NAME=server
```

### Port

Set server port (default: 8080)

```yaml
server:
    port: 8080
```

```dotenv
MANDARINE_SERVER_PORT=8080
```

### External origin

Set server external origin (default: http://localhost:8080)

```yaml
server:
    externalorigin: http://localhost:8080
```

```dotenv
MANDARINE_SERVER_EXTERNALORIGIN=http://localhost:8080
```

### Mode

Set server mode (development, production, local) (default: local)

```yaml
server:
    mode: local
```

```dotenv
MANDARINE_SERVER_MODE=development
```

### Version

Set server version (default: 0.0.0)

```yaml
server:
    version: 0.0.0
```

```dotenv
MANDARINE_SERVER_VERSION=0.0.0
```

### Request per second (RPS) (default: 100)

```yaml
server:
    rps: 100
```

```dotenv
MANDARINE_SERVER_RPS=100
```

### Max request body size (bytes) (default: 524288000 (500MB))

```yaml
server:
    maxrequestsize: 524288000
```

```dotenv
MANDARINE_SERVER_MAXREQUESTSIZE=524288000
```

<a id="database-provider"></a>

## Database provider settings

### Type

Set database type (postgres) (default: postgres)

```yaml
database:
    type: postgres
```

```dotenv
MANDARINE_DATABASE_TYPE=postgres
```

<a id="postgres-database-provider"></a>

### PostgreSQL

#### Address

Set postgres database address

```yaml
database:
    postgres:
        address: localhost:5432
```

```dotenv
MANDARINE_DATABASE_POSTGRES_ADDRESS=localhost:5432
```

#### Username

Set postgres database username

```yaml
database:
    postgres:
        username: admin
```

```dotenv
MANDARINE_DATABASE_POSTGRES_USERNAME=admin
```

#### Password

Set postgres database password

```yaml
database:
    postgres:
        password: password
```

```dotenv
MANDARINE_DATABASE_POSTGRES_PASSWORD=password
```

#### Database name

Set postgres database name

```yaml
database:
    postgres:
        dbname: mandarine
```

```dotenv
MANDARINE_DATABASE_POSTGRES_DBNAME=mandarine
```

<a id="cache-provider"></a>

## Cache provider settings

### TTL

Set cache TTL (default: 120)

```yaml
cache:
    ttl: 120
```

```dotenv
MANDARINE_CACHE_TTL=120
```

### Type

Set cache type (memory, redis) (default: memory)

```yaml
cache:
    type: memory
```

```dotenv
MANDARINE_CACHE_TYPE=memory
```

<a id="redis-cache-provider"></a>

### Redis

#### Address

Set redis cache address

```yaml
cache:
    redis:
        address: localhost:6379
```

```dotenv
MANDARINE_CACHE_REDIS_ADDRESS=localhost:6379
```

#### Username

Set redis cache username

```yaml
cache:
    redis:
        username: default
```

```dotenv
MANDARINE_CACHE_REDIS_USERNAME=default
```

#### Password

Set redis cache password

```yaml
cache:
    redis:
        password: password
```

```dotenv
MANDARINE_CACHE_REDIS_PASSWORD=password
```

#### Database index

Set redis cache database index (default: 0)

```yaml
cache:
    redis:
        dbindex: 0
```

```dotenv
MANDARINE_CACHE_REDIS_DBINDEX=0
```

<a id="s3-provider"></a>

## S3 provider settings

### Type

Set S3 type (minio) (default: minio)

```yaml
s3:
    type: minio
```

```dotenv
MANDARINE_S3_TYPE=minio
```

<a id="minio-provider"></a>

### Minio

#### Address

Set Minio S3 address

```yaml
s3:
    minio:
        address: localhost:9000
```

```dotenv
MANDARINE_S3_MINIO_ADDRESS=localhost:9000
```

#### Access key

Set Minio S3 access key

```yaml
s3:
    minio:
        accesskey: admin
```

```dotenv
MANDARINE_S3_MINIO_ACCESSKEY=admin
```

#### Secret key

Set Minio S3 secret key

```yaml
s3:
    minio:
        secretkey: secret_key
```

```dotenv
MANDARINE_S3_MINIO_SECRETKEY=secret_key
```

#### Bucket name

Set Minio S3 bucket name

```yaml
s3:
    minio:
        bucket: mandarine
```

```dotenv
MANDARINE_S3_MINIO_BUCKET=mandarine
```

<a id="pubsub-provider"></a>

## Pub/Sub provider settings

### Type

Set pubsub type (memory, redis) (default: memory)

```yaml
pubsub:
    type: memory
```

```dotenv
MANDARINE_PUBSUB_TYPE=memory
```

<a id="redis-pubsub-provider"></a>

### Redis

#### Address

Set redis pubsub address

```yaml
pubsub:
    redis:
        address: localhost:6379
```

```dotenv
MANDARINE_PUBSUB_REDIS_ADDRESS=localhost:6379
```

#### Username

Set redis pubsub username (default: default)

```yaml
pubsub:
    redis:
        username: default
```

```dotenv
MANDARINE_PUBSUB_REDIS_USERNAME=default
```

#### Password

Set redis pubsub password (default: password)

```yaml
pubsub:
    redis:
        password: password
```

```dotenv
MANDARINE_PUBSUB_REDIS_PASSWORD=password
```

<a id="smtp"></a>

## SMTP settings

### Host

Set SMTP host

```yaml
smtp:
    host: smtp.yandex.ru
```

```dotenv
MANDARINE_SMTP_HOST=smtp.yandex.ru
```

### Port

Set SMTP port

```yaml
smtp:
    port: 465
```

```dotenv
MANDARINE_SMTP_PORT=465
```

### Username

Set SMTP username

```yaml
smtp:
    username: example@yandex.ru
```

```dotenv
MANDARINE_SMTP_USERNAME=example@yandex.ru
```

### Password

Set SMTP password

```yaml
smtp:
    password: password
```

```dotenv
MANDARINE_SMTP_PASSWORD=password
```

### SSL

Set SMTP SSL mode (default: false)

```yaml
smtp:
    ssl: true
```

```dotenv
MANDARINE_SMTP_SSL=true
```

### From

Set SMTP from

```yaml
smtp:
    from: 'Mandarine <example@yandex.ru>'
```

```dotenv
MANDARINE_SMTP_FROM='Mandarine <example@yandex.ru>'
```

<a id="websocket"></a>

## Websocket settings

### Pool size

Set websocket pool size (default: 1024)

```yaml
websocket:
    poolsize: 1024
```

```dotenv
MANDARINE_WEBSOCKET_POOLSIZE=1024
```

<a id="oauth-provider"></a>

## OAuth provider settings

### Client ID

Set OAuth client ID

```yaml
oauthclients:
    <provider>:
        clientid: client_id
```

```dotenv
MANDARINE_OAUTHCLIENTS_<PROVIDER>_CLIENTID=client_id
```

### Client secret

Set OAuth client secret

```yaml
oauthclients:
    <provider>:
        clientsecret: client_secret
```

```dotenv
MANDARINE_OAUTHCLIENTS_<PROVIDER>_CLIENTSECRET=client_secret
```

> **_NOTE:_** If you want to set provider API key using only environment variable and not showing it in the
> configuration file, set it to a random and not empty string in the configuration file.

<a id="geocoding-provider"></a>

## Geocoding provider settings

### API key

Set geocoding provider API key

```yaml
geocodingclients:
    <provider>:
       apikey: api_key
```

```dotenv
MANDARINE_GEOCODINGCLIENTS_<PROVIDER>_APIKEY=api_key
```

> **_NOTE:_** If the provider doesn't need an API key, just set it to a random and not empty string. This is important 
> because otherwise the parser will not recognize the provider with such a key.

> **_NOTE:_** If you want to set provider API key using only environment variable and not showing it in the 
> configuration file, set it to a random and not empty string in the configuration file.

<a id="security"></a>

## Security settings

<a id="security-jwt"></a>

### JWT

#### Secret

Set JWT secret

```yaml
security:
    jwt:
        secret: secret
```

```dotenv
MANDARINE_JWT_SECRET=secret
```

#### Access token TTL

Set JWT access token TTL (default: 3600)

```yaml
security:
    jwt:
        accesstokenttl: 3600
```

```dotenv
MANDARINE_JWT_ACCESSTOKENTTL=3600
```

#### Refresh token TTL

Set JWT refresh token TTL (default: 86400)

```yaml
security:
    jwt:
        refreshtokenttl: 86400
```

```dotenv
MANDARINE_JWT_REFRESHTOKENTTL=86400
```

<a id="security-otp"></a>

### OTP

#### Length

Set OTP length (default: 6)

```yaml
security:
    otp:
        length: 6
```

```dotenv
MANDARINE_OTP_LENGTH=6
```

#### TTL

Set OTP TTL (default: 600)

```yaml
security:
    otp:
        ttl: 600
```

```dotenv
MANDARINE_OTP_TTL=600
```

<a id="locale"></a>

## Locale settings

### Path

Set locale path (default: locales)

```yaml
locale:
    path: locales
```

```dotenv
MANDARINE_LOCALE_PATH=locales
```

### Language

Set locale default language (default: ru)

```yaml
locale:
    language: ru
```

```dotenv
MANDARINE_LOCALE_LANGUAGE=ru
```

<a id="template"></a>

## Template settings

### Path

Set template path (default: templates)

```yaml
template:
    path: templates
```

```dotenv
MANDARINE_TEMPLATE_PATH=templates
```

<a id="migrations"></a>

## Migrations settings

### Path

Set migrations path (default: migrations)

```yaml
migrations:
    path: migrations
```

```dotenv
MANDARINE_MIGRATIONS_PATH=migrations
```

<a id="logger"></a>

## Logger settings

<a id="logger-console"></a>

### Level

Set logger level (default: info)

```yaml
logger:
    level: info
```

```dotenv
MANDARINE_LOGGER_LEVEL=info
```

### Console

#### Enable

Set file logger enable (default: false)

```yaml
logger:
    console:
        enable: false
```

```dotenv
MANDARINE_LOGGER_CONSOLE_ENABLE=false
```

#### Encoding

Set console logger encoding (text, json) (default: text)

```yaml
logger:
    console:
        encoding: text
```

```dotenv
MANDARINE_LOGGER_CONSOLE_ENCODING=text
```

<a id="logger-file"></a>

### File

#### Enable

Set file logger enable (default: false)

```yaml
logger:
    file:
        enable: false
```

```dotenv
MANDARINE_LOGGER_FILE_ENABLE=false
```

#### Dir path

Set file logger directory path (default: logs)

```yaml
logger:
    file:
        dirpath: logs
```

```dotenv
MANDARINE_LOGGER_FILE_DIRPATH=logs
```

#### Max size

Set file logger max size (default: 1)

```yaml
logger:
    file:
        maxsize: 1
```

```dotenv
MANDARINE_LOGGER_FILE_MAXSIZE=1
```

#### Max age

Set file logger max age

```yaml
logger:
    file:
        maxage: 30
```

```dotenv
MANDARINE_LOGGER_FILE_MAXAGE=30
```