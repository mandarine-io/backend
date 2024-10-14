# Configuration

<a id="server"></a>
## Server settings

### Name

Set server name (default: mandarine-server)

```yaml
server:
    name: mandarine-server
```
```dotenv
MANDARINE_SERVER__NAME=mandarine-server
```

### Port

Set server port (default: 8080)

```yaml
server:
    port: 8080
```
```dotenv
MANDARINE_SERVER__PORT=8080
```

### External origin

Set server external origin (default: http://localhost:8080)

```yaml
server:
    external_origin: http://localhost:8080
```
```dotenv
MANDARINE_SERVER__EXTERNAL_ORIGIN=http://localhost:8080
```

### Mode

Set server mode (development, production, local) (default: local)

```yaml
server:
    mode: local
```
```dotenv
MANDARINE_SERVER__MODE=development
```

### Version

Set server version (default: 0.0.0)

```yaml
server:
    version: 0.0.0
```
```dotenv
MANDARINE_SERVER__VERSION=0.0.0
```

<a id="database-provider"></a>
## Database provider settings (PostgreSQL)

### Host

Set database host (default: localhost)

```yaml
postgres:
    host: localhost
```
```dotenv
MANDARINE_POSTGRES__HOST=localhost
```

### Port

Set database port (default: 5432)

```yaml
postgres:
    port: 5432
```
```dotenv
MANDARINE_POSTGRES__PORT=5432
```

### Username

Set database username

```yaml
postgres:
    username: admin
```
```dotenv
MANDARINE_POSTGRES__USERNAME=admin
```

### Password

Set database password

```yaml
postgres:
    password: password
```
```dotenv
MANDARINE_POSTGRES__PASSWORD=password
```

### Password file

Set database password file

```yaml
postgres:
    password_file: /run/secrets/postgres-password
```
```dotenv
MANDARINE_POSTGRES__PASSWORD_FILE=/run/secrets/postgres-password
```

### Database name

Set database name

```yaml
postgres:
    db_name: mandarine
```
```dotenv
MANDARINE_POSTGRES__DB_NAME=mandarine
```

<a id="cache-provider"></a>
## Cache provider settings (Redis)

### Host

Set cache host (default: localhost)

```yaml
redis:
    host: localhost
```
```dotenv
MANDARINE_REDIS__HOST=localhost
```

### Port

Set cache port (default: 6379)

```yaml
redis:
    port: 6379
```
```dotenv
MANDARINE_REDIS__PORT=6379
```

### Username

Set cache username

```yaml
redis:
    username: default
```
```dotenv
MANDARINE_REDIS__USERNAME=default
```

### Password

Set cache password

```yaml
redis:
    password: password
```
```dotenv
MANDARINE_REDIS__PASSWORD=password
```

### Password file

Set cache password file

```yaml
redis:
    password_file: /run/secrets/redis-password
```
```dotenv
MANDARINE_REDIS__PASSWORD_FILE=/run/secrets/redis-password
```

### Database index

Set cache database index (default: 0)

```yaml
redis:
    db_index: 0
```
```dotenv
MANDARINE_REDIS__DB_INDEX=0
```

<a id="s3-provider"></a>
## S3 provider settings (MinIO)

### Host

Set S3 host (default: localhost)

```yaml
minio:
    host: localhost
```
```dotenv
MANDARINE_MINIO__HOST=localhost
```

### Port

Set S3 port (default: 9000)

```yaml
minio:
    port: 9000
```
```dotenv
MANDARINE_MINIO__PORT=9000
```

### Access key

Set S3 access key

```yaml
minio:
    access_key: admin
```
```dotenv
MANDARINE_MINIO__ACCESS_KEY=admin
```

### Secret key

Set S3 secret key

```yaml
minio:
    secret_key: secret_key
```
```dotenv
MANDARINE_MINIO__SECRET_KEY=secret_key
```

### Secret key file

Set S3 secret key file

```yaml
minio:
    secret_key_file: /run/secrets/minio-secret-key
```
```dotenv
MANDARINE_MINIO__SECRET_KEY_FILE=/run/secrets/minio-secret-key
```

### Bucket name

Set S3 bucket name

```yaml
minio:
    bucket_name: mandarine
```
```dotenv
MANDARINE_MINIO__BUCKET_NAME=mandarine
```

<a id="smtp"></a>
## SMTP settings

### Host

Set SMTP host (default: smtp.yandex.ru)

```yaml
smtp:
    host: smtp.yandex.ru
```
```dotenv
MANDARINE_SMTP__HOST=smtp.yandex.ru
```

### Port

Set SMTP port (default: 465)

```yaml
smtp:
    port: 465
```
```dotenv
MANDARINE_SMTP__PORT=465
```

### Username

Set SMTP username

```yaml
smtp:
    username: example@yandex.ru
```
```dotenv
MANDARINE_SMTP__USERNAME=example@yandex.ru
```

### Password

Set SMTP password

```yaml
smtp:
    password: password
```
```dotenv
MANDARINE_SMTP__PASSWORD=password
```

### Password file

Set SMTP password file

```yaml
smtp:
    password_file: /run/secrets/smtp-password
```
```dotenv
MANDARINE_SMTP__PASSWORD_FILE=/run/secrets/smtp-password
```

### SSL

Set SMTP SSL mode (default: true)

```yaml
smtp:
    ssl: true
```
```dotenv
MANDARINE_SMTP__SSL=true
```

### From

Set SMTP from

```yaml
smtp:
    from: 'Mandarine <example@yandex.ru>'
```
```dotenv
MANDARINE_SMTP__FROM='Mandarine <example@yandex.ru>'
```

<a id="oauth-provider"></a>
## OAuth provider settings (Google, Yandex, Mail.ru)

### Client ID

Set OAuth client ID

```yaml
oauth:
    google:
        client_id: client_id
```
```dotenv
MANDARINE_GOOGLE_OAUTH_CLIENT__CLIENT_ID=client_id
```

### Client secret

Set OAuth client secret

```yaml
oauth:
    google:
        client_secret: client_secret
```
```dotenv
MANDARINE_GOOGLE_OAUTH_CLIENT__CLIENT_SECRET=client_secret
```

### Client secret file

Set OAuth client secret file

```yaml
oauth:
    google:
        client_secret_file: /run/secrets/google_client_secret
```
```dotenv
MANDARINE_GOOGLE_OAUTH_CLIENT__CLIENT_SECRET_FILE=/run/secrets/google_client_secret
```

<a id="cache"></a>
## Cache settings

### TTL

Set cache TTL (default: 120)

```yaml
cache:
    ttl: 120
```
```dotenv
MANDARINE_CACHE__TTL=120
```

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
MANDARINE_JWT__SECRET=secret
```

#### Secret file

Set JWT secret file

```yaml
security:
    jwt:
        secret_file: /run/secrets/jwt-secret
```
```dotenv
MANDARINE_JWT__SECRET_FILE=/run/secrets/jwt-secret
```

#### Access token TTL

Set JWT access token TTL (default: 3600)

```yaml
security:
    jwt:
        access_token_ttl: 3600
```
```dotenv
MANDARINE_JWT__ACCESS_TOKEN_TTL=3600
```

#### Refresh token TTL

Set JWT refresh token TTL (default: 86400)

```yaml
security:
    jwt:
        refresh_token_ttl: 86400
```
```dotenv
MANDARINE_JWT__REFRESH_TOKEN_TTL=86400
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
MANDARINE_OTP__LENGTH=6
```

#### TTL

Set OTP TTL (default: 600)

```yaml
security:
    otp:
        ttl: 600
```
```dotenv
MANDARINE_OTP__TTL=600
```

<a id="security-rate-limit"></a>
### Rate limit

#### RPS

Set rate limit RPS (default: 100)

```yaml
security:
    rate_limit:
        rps: 100
```
```dotenv
MANDARINE_RATE_LIMIT__RPS=100
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
MANDARINE_LOCALE__PATH=locales
```

### Language

Set locale default language (default: ru)

```yaml
locale:
    language: ru
```
```dotenv
MANDARINE_LOCALE__LANGUAGE=ru
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
MANDARINE_TEMPLATE__PATH=templates
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
MANDARINE_MIGRATIONS__PATH=migrations
```

<a id="logger"></a>
## Logger settings

<a id="logger-console"></a>
### Console

#### Level

Set console logger level (default: info)

```yaml
logger:
    console:
        level: info
```
```dotenv
MANDARINE_LOGGER__CONSOLE_LEVEL=info
```

#### Encoding

Set console logger encoding (text, json) (default: text)

```yaml
logger:
    console:
        encoding: text
```
```dotenv
MANDARINE_LOGGER__CONSOLE_ENCODING=text
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
MANDARINE_LOGGER__FILE_ENABLE=false
```

#### Level

Set file logger level (default: info) 

```yaml
logger:
    file:
        level: info
```
```dotenv
MANDARINE_LOGGER__FILE_LEVEL=info
```

#### Dir path

Set file logger directory path (default: logs)

```yaml
logger:
    file:
        dir_path: logs
```
```dotenv
MANDARINE_LOGGER__FILE_DIR_PATH=logs
```

#### Max size

Set file logger max size (default: 1)

```yaml
logger:
    file:
        max_size: 1
```
```dotenv
MANDARINE_LOGGER__FILE_MAX_SIZE=1
```

#### Max age

Set file logger max age

```yaml
logger:
    file:
        max_age: 30
```
```dotenv
MANDARINE_LOGGER__FILE_MAX_AGE=30
```