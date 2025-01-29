# Качество кода

В процессе разработки проекта вы можете использовать [gfmt](https://pkg.go.dev/cmd/gofmt)
и [golangci-lint](https://golangci-lint.run/) для проверки кода.

## Форматирование

Чтобы отформатировать код, вы можете запустить команду `Makefile`:

```bash
make format
```

или:

```bash
gofmt -w .
```

С исправлением найденных проблем:

```bash
make format.fix
```

или:

```bash
gofmt -s -w .
```

## Линтеры

Все линтеры и их настройки описаны в файле `golangcli.yaml`. Для запуска линтеров можно выполнить `Makefile`
команда:

```bash
make lint
```

или:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run -c golangci.yaml
```

С исправлением найденных проблем:

```bash
make lint.fix
```

или:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run -c golangci.yaml --fix
```