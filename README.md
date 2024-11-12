<h1 id="home" align="center" style="font-weight: bold;">Mandarine (Backend)</h1>

<h2 id="technologies">Technologies</h2>

- [Golang](https://go.dev/)
- [Gin Gonic](https://gin-gonic.com/)
- [GORM](https://gorm.io/index.html)
- [gocron](https://github.com/go-co-op/gocron)
- [gorilla/websocket](https://github.com/gorilla/websocket)
- OAuth2
  providers ([Google](https://developers.google.com/identity/protocols/oauth2?hl=ru), [Yandex](https://yandex.ru/dev/id/doc/ru/concepts/ya-oauth-intro), [Mail.ru](https://help.mail.ru/developers/oauth))
- Geocoding
  providers ([Graphhoper](https://docs.graphhopper.com/), [Here](https://developer.here.com/develop/rest-apis), [LocationIQ](https://docs.locationiq.com/reference/search), [OpenStreetMap (Nominatim)](https://nominatim.org/release-docs/latest/api/Overview/), [Yandex](https://yandex.ru/dev/geocode/doc/ru/))
- [PostgreSQL (PostGIS, gin)](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [MinIO](https://min.io/)
- [WebSockets](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API)
- [Nginx](https://nginx.org/)
- [Mailhog](https://github.com/mailhog/MailHog)
- [Testcontainers](https://testcontainers.com/)
- [K6](https://k6.io/)
- [Docker](https://www.docker.com/)
- [Git](https://git-scm.com/)
- [Make](https://www.gnu.org/software/make/)
- [EditorConfig](https://editorconfig.org/)

<h2 id="getting-started">Getting started</h2>

Here you describe how to run project locally

<h3 id="prerequisites">Prerequisites</h3>

To launch a project, you need:

- [Golang](https://go.dev/)
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Make](https://www.gnu.org/software/make/)
- [NPM](https://www.npmjs.com/)

<h3 id="cloning">Cloning</h3>

Сlone this project:

```bash
git clone https://github.com/mandarine-io/Backend
```

<h3 id="config">Configuration</h3>

<h4 id="yaml-file">YAML file</h4>

YAML configuration file contains all base application settings.
Use the `config/config.example.yaml` as reference to create your configuration file:

```bash
cp config/config.example.yaml config/config.yaml
nano config/config.yaml
```

<h4 id="envs">Environment variables</h4>

To overwrite some properties from YAML file, you can use environment variables.
Use the `config/.env.example` as reference to create your env file `.env`:

```bash
cp config/.env.example .env
nano .env
```

<h3 id="launch">Launch</h3>

To start server, you can run Makefile command:

```bash
make start
```

or you can run manually to use custom YAML config file and environment variables file:

```bash
go mod tidy
go build -o build/server cmd/api
MANDARINE_CONFIG_FILE=config/config.yaml ./build/server
```

<h2 id="dev">Development</h2>

To start server with hot reload (development mode), you can run Makefile command:

```bash
make start.dev
```

<h3 id="format">Formatting</h3>

To format code, you can run Makefile command:

```bash
make format
```

With fixing found issues:

```bash
make format.fix
```

<h3 id="lint">Linting</h3>

All linters and its settings describes file `golangcli.yaml`. To run linters, you can execute Makefile
command:

```bash
make lint
```

With fixing found issues:

```bash
make lint.fix
```

<h2 id="testing">Testing</h2>

The system is covered with various types of tests.

<h3 id="unit-testing">Unit tests</h3>

Created unit tests for services, various custom managers and clients, and util functions:

```bash
make test.unit
```

After finishing, you can see the results in the `logs/unit-test` folder (logs and coverage reports).

<h3 id="e2e-testing">E2E tests</h3>

The main business scenarios are covered with e2e tests, and for them a test environment is deployed in Docker
containers:

```bash
make test.e2e
```

After finishing, you can see the results in the `logs/e2e-test` folder (logs and coverage reports).

<h3 id="load-testing">Load tests</h3>

To test the system under load and identify bottlenecks, load tests are written:

```bash
make test.load LOAD_TEST_NAME=<test-file-name>
```

After finishing, you can see the results in the `logs/load-test` folder (logs and performance reports).

<h2 id="license">License</h2>

This project is licensed under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0.html).