# Weather App

REST API сервис на Go с PostgreSQL, Auth, Security который:
- управляет пользователями
- позволяет пользователю следить за несколькими городами
- получает погоду из внешнего API
- сохраняет историю погодных запросов
- поддерживает фильтрацию истории по городу

## Запуск

```bash
# Применить миграции
psql $DATABASE_URL -f migrations/001_init.sql

# Запустить сервер
DATABASE_URL=postgres://akydyrbay@/weather_db JWT_SECRET=your_secret go run ./cmd/app
```

Сервер стартует на `http://localhost:8080`.

## Эндпоинты

### Healthcheck

```bash
curl http://localhost:8080/health
```

### Аутентификация

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/auth/register` | Регистрация |
| POST | `/auth/login` | Вход, возвращает JWT |

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Алибек","email":"ali@example.com","password":"secret"}'

curl -X POST http://localhost:8080/auth/login \
  -d '{"email":"ali@example.com","password":"secret"}'
# → {"access_token":"<jwt>"}
```

Защищённые маршруты требуют заголовок:
```
Authorization: Bearer <access_token>
```

### Погода

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/weather?lat=&lon=` | Погода по координатам |
| GET | `/weather/{city}` | Погода в городе + рекомендация по одежде |
| GET | `/weather/country/{country}` | Погода по городам страны |
| GET | `/weather/country/{country}/top` | Топ-3 самых тёплых города |

```bash
curl "http://localhost:8080/weather?lat=43.25&lon=76.92"

curl http://localhost:8080/weather/Almaty

curl http://localhost:8080/weather/country/Kazakhstan

curl http://localhost:8080/weather/country/Kazakhstan/top
```

### Пользователи

> Только для `admin`

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/users` | Список пользователей |
| GET | `/users/{id}` | Получить пользователя |
| DELETE | `/users/{id}` | Мягкое удаление |

```bash
curl http://localhost:8080/users -H "Authorization: Bearer <token>"

curl http://localhost:8080/users/1 -H "Authorization: Bearer <token>"

curl -X DELETE http://localhost:8080/users/1 -H "Authorization: Bearer <token>"
```

### Текущий пользователь

> Требует авторизацию

```bash
curl http://localhost:8080/users/me -H "Authorization: Bearer <token>"
```

### Города пользователя

> Требует авторизацию. Пользователь берётся из JWT.

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/cities` | Добавить город |
| GET | `/cities` | Список городов |
| DELETE | `/cities/{city_id}` | Удалить город |

```bash
curl -X POST http://localhost:8080/cities \
  -H "Authorization: Bearer <token>" \
  -d '{"city":"Almaty"}'

curl http://localhost:8080/cities -H "Authorization: Bearer <token>"

curl -X DELETE http://localhost:8080/cities/3 -H "Authorization: Bearer <token>"
```

### Погода пользователя

> Требует авторизацию. Запрашивает погоду по всем городам параллельно и сохраняет в историю.

```bash
curl http://localhost:8080/users/weather -H "Authorization: Bearer <token>"
```

```json
{
  "user_id": 1,
  "results": [
    {
      "city": "Almaty",
      "temperature": 12.3,
      "description": "Ясно",
      "clothing": "куртка"
    }
  ]
}
```

### История погоды

> Требует авторизацию.

```bash
curl "http://localhost:8080/weather/history?city=Almaty&limit=10" \
  -H "Authorization: Bearer <token>"

curl "http://localhost:8080/weather/history?limit=20&offset=40" \
  -H "Authorization: Bearer <token>"
```

| Параметр | Обязательный | Описание |
|----------|-------------|----------|
| `city` | нет | Фильтр по городу |
| `limit` | нет | Максимум записей |
| `offset` | нет | Смещение (пагинация) |

```json
{
  "user_id": 1,
  "city": "Almaty",
  "history": [
    {
      "temperature": 18,
      "description": "Ясно",
      "requested_at": "2026-04-20T10:00:00Z"
    }
  ]
}
```

## Стек

- **Go**, **net/http**, **go-chi/chi** - сервер и роутинг
- **pgx/v5** - PostgreSQL
- **golang-jwt/jwt/v5** - JWT аутентификация
- **bcrypt** - хэширование паролей
- **Open Meteo** - погода и геокодинг (без ключа)
- **Uber Zap** - structured logging + logging middleware (method, path, status, duration, request_id)
- **testify** (assert / require / mock) - unit-тесты + mock-репозитории
- Параллельные запросы к API через горутины
- Кэш погоды в памяти (TTL 5 минут)
- Мягкое удаление пользователей
- Индекс `(user_id, city)` для быстрой фильтрации истории

## Тесты

```bash
go test -v ./...

go test -cover ./...

TEST_DATABASE_URL=postgres://akydyrbay@/weather_db \
  go test -tags integration -v ./internal/repository/...
```

Мок репозитории живут в `internal/service/mocks/` - используются в `service` и `handler` юнит тестах, не ходя в БД.