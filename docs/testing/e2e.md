# E2E тесты

**E2E (End-to-End), или сквозное тестирование** – это тип тестирования, целью которого является моделирование работы
пользователя от начала до конца.

## Запуск

Запустите среды тестирование в Docker с помощью команды Makefile:

```bash
make deploy.e2e-test
```

Запустить E2E тесты можно с помощью Makefile команды:

```bash
make test.e2e
```