<h1 align="center">Backend</h1>
<p align="center">
  <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/mandarine-io/backend?color=3fbc11&style=flat">
  <img alt="Go version" src="https://img.shields.io/github/go-mod/go-version/mandarine-io/backend?color=3fbc11&style=flat">
  <img alt="License" src="https://img.shields.io/github/license/mandarine-io/backend?color=3fbc11&style=flat">
  <img alt="Github issues" src="https://img.shields.io/github/issues/mandarine-io/backend?color=3fbc11&style=flat" />
  <img alt="Github forks" src="https://img.shields.io/github/forks/mandarine-io/backend?color=3fbc11&style=flat" />
  <img alt="Github stars" src="https://img.shields.io/github/stars/mandarine-io/backend?color=3fbc11&style=flat" />
</p>

**Mandarine** - это платформа для записи на услуги красоты и ухода, объединяющая клиентов и мастеров. Здесь мы изучим
один из ее компонентов - *сервер*.

Mandarine имеет клиент-серверную архитектуру, поэтому сервер инкапсулирует достаточно много функционала:

+ **Регистрация и авторизация**
+ **Управление аккаунтами**
+ **Профиль, услуги, портфолио мастеров**
+ **Формирование расписаний и запись на услуги**
+ **Поиск мастеров**
+ **Отзывы о мастерах и их рейтинг**
+ **Уведомления**
+ **И много другого**

## Быстрый старт

Для того чтобы запустить проект локально и ознакомиться с его основными функциями, вам нужно:

### Предварительные условия

Подготовить следующие инструменты:

- [Golang](https://go.dev/)
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/) (опционально)

### Клонирование проекта

Склонировать репозиторий сервера:

```bash
git clone https://github.com/mandarine-io/backend
```

### Среда выполнения

Запустить среду выполнения в Docker. Для этого выполните команду:

```bash
docker compose -f docker-compose.local.yml up -d
```

### Конфигурация

Используйте `config/config.default.yaml` в качестве шаблона для создания своего файла конфигурации:

```bash
cp config/config.default.yaml config/config.yaml
nano config/config.yaml
```

### Запуск

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

## Лицензия

Этот проект распространяется по [Лицензии Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0.html).