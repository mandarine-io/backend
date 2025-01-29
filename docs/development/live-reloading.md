# Live reloading

Живая перезагрузка позволяет редактировать код и сразу видеть изменения, без перезапуска сервера.
Для этих целей мы используем [air](https://github.com/air-verse/air).

Для запуска сервера в режиме живой перезагрузки, можно использовать команду `Makefile`:

```bash
make start.dev
```

или

```bash
go install github.com/air-verse/air@latest
air
```