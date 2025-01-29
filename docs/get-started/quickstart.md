# Быстрый старт

Для того, чтобы запустить проект локально и ознакомиться с его основными функциями, вам нужно:

## Предварительные условия

Подготовить следующие инструменты:

- [Golang](https://go.dev/)
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/) (опционально)

## Клонирование проекта

Клонируйте репозиторий сервера:

```bash
git clone https://github.com/mandarine-io/backend
```

## Среда выполнения

Запустите среду выполнения в Docker. Для этого выполните команду:

```bash
docker compose -f docker-compose.local.yml up -d
```

## Конфигурация

Используйте `config/config.default.yaml` в качестве шаблона для создания своего файла конфигурации:

```bash
cp config/config.default.yaml config/config.yaml
nano config/config.yaml
```

## Запуск

Чтобы запустить сервер, вы можете запустить команду Makefile:

```bash
make start
```

или вы можете запустить вручную:

```bash
go mod tidy
go build -o build/server cmd/api
./build/server
```

## Swagger

Переходите по ссылке [Swagger](http://localhost:8080/swagger/index.html) для просмотра документации

![swagger-screen](./_assets/swagger-screen.png 'Скриншот Swagger')

## Остановка

Сервер можно остановить при помощи сигналов `SIGINT`, `SIGTERM`, `SIGQUIT`