# Сервис подписок

Это пример реализации REST API для управления подписками, включая создание, удаление, обновление, получение отдельных подписок, получение списка подписок и подсчет общей стоимости подписок за указанный период.

---

## Зависимости

Проект использует минимальный набор зависимостей:

- [Gin](https://github.com/gin-gonic/gin) — HTTP-роутер
- [Zerolog](https://github.com/rs/zerolog) — используемый логгер
- [Testify](https://github.com/stretchr/testify) — фреймворк для тестирования и утверждений
- [GoMock](https://github.com/golang/mock) — генератор моков для тестов
- [pgx/v5](https://github.com/jackc/pgx) — драйвер для PostgreSQL
- [Migrate](https://github.com/golang-migrate/migrate) — миграции базы данных

---

## Начало работы

### Установка

1. **Клонируйте репозиторий**:
   ```bash
   git clone https://github.com/prolgrammer/TaskManager
   cd task-service
   ```

2. **Установите зависимости при необходимости**:
   ```bash
   go mod tidy
   ```

3. **Запустите приложение с помощью Docker**:
   ```bash
   docker-compose -f docker-compose.yml up -d
   ```

   Сервис будет доступен по адресу: `http://localhost:8080`.


### Тестирование

Перейдите по адресу http://localhost:8080/swagger/index.html#/default для открытия документации и проведения тестов.

---


## Структура папок

```
subscription-service/
├── cmd/
│   ├── app/
│   │   └── app.go
│   └── main.go
├── config/
│   ├── pg/
│   │   ├── migrations/
│   │   │   ├── 000001_init.up.sql
│   │   │   └── 000001_init.down.sql
│   │   └── config.go
│   ├── config.go
│   └── config.yaml
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── infrastructure/
│   ├── postgres/
│   │   ├── commands/
│   │   │   ├── subscription/
│   │   │   │   ├── calculate_total_cost.go
│   │   │   │   ├── delete_by_id.go
│   │   │   │   ├── insert.go
│   │   │   │   ├── select_all.go
│   │   │   │   ├── select_by_id.go
│   │   │   │   ├── sub_repository.go
│   │   │   │   └── update_by_id.go
│   │   │   └── utils.go
│   │   └── client.go
├── internal/
│   ├── controllers/
│   │   ├── http/
│   │   │   ├── middleware/
│   │   │   │   ├── error_handler.go
│   │   │   │   ├── errors.go
│   │   │   │   └── middleware.go
│   │   │   ├── calculate_total_cost.go
│   │   │   ├── create_subscription.go
│   │   │   ├── delete_subscription.go
│   │   │   ├── get_all_subscriptions.go
│   │   │   ├── get_subscription.go
│   │   │   ├── router.go
│   │   │   └── update_subscription.go
│   │   ├── requests/
│   │   │   ├── calculate_total_cost.go
│   │   │   └── subscription.go
│   │   ├── responses/
│   │   │   ├── calculate_total_cost.go
│   │   │   └── subscription.go
│   │   └── errors.go
│   ├── entities/
│   │   └── subscription.go
│   ├── subscription/
│   │   └── subscription.go
│   ├── usecases/
│   │   ├── calculate_total_cost.go
│   │   ├── calculate_total_cost_test.go
│   │   ├── contracts.go
│   │   ├── create_subscription.go
│   │   ├── create_subscription_test.go
│   │   ├── delete_subscription.go
│   │   ├── delete_subscription_test.go
│   │   ├── errors.go
│   │   ├── get_all_subscriptions.go
│   │   ├── get_all_subscriptions_test.go
│   │   ├── get_subscription.go
│   │   ├── get_subscription_test.go
│   │   ├── mock_test.go
│   │   ├── update_subscription.go
│   │   └── update_subscription_test.go
├── pkg/
│   └── logger/
│       ├── logger.go
│       ├── mock_logger.go
│       └── zerolog.go
├── .env
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

### Описание структуры

- **cmd/**: Точка входа приложения.
- **config/**: Конфигурации приложения.
- **internal/controllers/http/**: HTTP-контроллеры для обработки входящих запросов.
- **internal/controllers/requests/**: Структуры для входных данных запросов.
- **internal/controllers/responses/**: Структуры для ответов.
- **internal/entities/**: Определения сущностей.
- **internal/usecases/**: Бизнес-логика приложения.
- **internal/infrastructure/postgres/commands/**: Константы для имен таблиц и полей базы данных.
- **pkg/**: Вспомогательные пакеты (например, логгер).

---

## API

### Создание подписки
- **Метод**: `POST /subscriptions`
- **Тело запроса**:
  ```json
  {
    "service_name": "строка",
    "price": "целое число",
    "user_id": "uuid",
    "start_date": "MM-YYYY",
    "end_date": "MM-YYYY"
  }
  ```
- **Ответ** (200 OK):
  ```json
  {
    "id": "uuid",
    "service_name": "строка",
    "price": "целое число",
    "user_id": "uuid",
    "start_date": "MM-YYYY",
    "end_date": "MM-YYYY"
  }
  ```
- **Ошибки**:
    - `400 Bad Request`: Некорректный формат запроса (например, неверный UUID или формат даты).
    - `500 Internal Server Error`: Внутренняя ошибка сервера (например, сбой базы данных).
- **Пример**:
  ```bash
  curl -X POST http://localhost:8080/subscriptions -H "Content-Type: application/json" -d '{"service_name":"Yandex Plus","price":400,"user_id":"6060ffee-2bf1-4721-ae6f-7636e979a0cb","start_date":"07-2025","end_date":"12-2025"}'
  ```

### Удаление подписки
- **Метод**: `DELETE /subscriptions/{sub_id}`
- **Ответ** (200 OK):
  ```json
  {
    "message": "подписка успешно удалена"
  }
  ```
- **Ошибки**:
    - `400 Bad Request`: Неверный формат `sub_id` (должен быть валидным UUID).
    - `404 Not Found`: Подписка с указанным `sub_id` не найдена.
    - `500 Internal Server Error`: Внутренняя ошибка сервера.
- **Пример**:
  ```bash
  curl -X DELETE http://localhost:8080/subscriptions/6060ffee-2bf1-4721-ae6f-7636e979a0cb
  ```

### Обновление подписки
- **Метод**: `PUT /subscriptions/{sub_id}`
- **Тело запроса**:
  ```json
  {
    "service_name": "строка",
    "price": "целое число",
    "user_id": "uuid",
    "start_date": "MM-YYYY",
    "end_date": "MM-YYYY"
  }
  ```
- **Ответ** (200 OK):
  ```json
  {
    "id": "uuid",
    "service_name": "строка",
    "price": "целое число",
    "user_id": "uuid",
    "start_date": "MM-YYYY",
    "end_date": "MM-YYYY"
  }
  ```
- **Ошибки**:
    - `400 Bad Request`: Некорректный формат запроса (например, неверный UUID или формат даты).
    - `404 Not Found`: Подписка с указанным `sub_id` не найдена.
    - `500 Internal Server Error`: Внутренняя ошибка сервера.
- **Пример**:
  ```bash
  curl -X PUT http://localhost:8080/subscriptions/6060ffee-2bf1-4721-ae6f-7636e979a0cb -H "Content-Type: application/json" -d '{"service_name":"Yandex Plus Updated","price":500,"user_id":"6060ffee-2bf1-4721-ae6f-7636e979a0cb","start_date":"07-2025","end_date":"12-2026"}'
  ```

### Получение подписки
- **Метод**: `GET /subscriptions/{sub_id}`
- **Ответ** (200 OK):
  ```json
  {
    "id": "uuid",
    "service_name": "строка",
    "price": "целое число",
    "user_id": "uuid",
    "start_date": "MM-YYYY",
    "end_date": "MM-YYYY"
  }
  ```
- **Ошибки**:
    - `400 Bad Request`: Неверный формат `sub_id` (должен быть валидным UUID).
    - `404 Not Found`: Подписка с указанным `sub_id` не найдена.
    - `500 Internal Server Error`: Внутренняя ошибка сервера.
- **Пример**:
  ```bash
  curl http://localhost:8080/subscriptions/6060ffee-2bf1-4721-ae6f-7636e979a0cb
  ```

### Получение списка подписок
- **Метод**: `GET /subscriptions`
- **Параметры запроса**:
    - `limit` (целое число): Количество подписок на странице (по умолчанию: 10).
    - `offset` (целое число): Смещение для пагинации (по умолчанию: 0).
- **Ответ** (200 OK):
  ```json
  [
    {
      "id": "uuid",
      "service_name": "строка",
      "price": "целое число",
      "user_id": "uuid",
      "start_date": "MM-YYYY",
      "end_date": "MM-YYYY"
    }
  ]
  ```
- **Ошибки**:
    - `400 Bad Request`: Некорректные параметры пагинации (например, отрицательный `limit`).
    - `500 Internal Server Error`: Внутренняя ошибка сервера.
- **Пример**:
  ```bash
  curl http://localhost:8080/subscriptions?limit=10&offset=0
  ```

### Подсчет общей стоимости подписок
- **Метод**: `POST /subscriptions/total`
- **Тело запроса**:
  ```json
  {
    "start_period": "MM-YYYY",
    "end_period": "MM-YYYY",
    "user_id": "uuid",
    "service_name": "строка"
  }
  ```
- **Ответ** (200 OK):
  ```json
  {
    "total": "целое число"
  }
  ```
- **Ошибки**:
    - `400 Bad Request`: Некорректный формат запроса (например, неверный UUID или формат даты).
    - `500 Internal Server Error`: Внутренняя ошибка сервера.
- **Пример**:
  ```bash
  curl -X POST http://localhost:8080/subscriptions/total -H "Content-Type: application/json" -d '{"start_period":"07-2025","end_period":"12-2025","user_id":"6060ffee-2bf1-4721-ae6f-7636e979a0cb","service_name":"Yandex Plus"}'
  ```
