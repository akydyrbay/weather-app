# Weather App

REST API сервис на Go с PostgreSQL, который:
- управляет пользователями
- позволяет пользователю следить за несколькими городами
- получает погоду из внешнего API
- сохраняет историю погодных запросов
- поддерживает фильтрацию истории по городу

## Запуск

```bash
# Применить миграцию
psql $DATABASE_URL -f migrations/001_init.sql

# Запустить сервер
DATABASE_URL=postgres://akydyrbay@/weather_db go run ./cmd/app
```

Сервер стартует на `http://localhost:8080`.

## Эндпоинты

### Healthcheck

```bash
curl http://localhost:8080/health
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

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/users` | Создать пользователя |
| GET | `/users` | Список пользователей |
| GET | `/users/{id}` | Получить пользователя |
| PUT | `/users/{id}` | Обновить пользователя |
| DELETE | `/users/{id}` | Мягкое удаление |

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Алибек","email":"ali@example.com"}'

curl http://localhost:8080/users/1

curl -X PUT http://localhost:8080/users/1 \
  -d '{"name":"Алибек Б.","email":"ali@example.com"}'

curl -X DELETE http://localhost:8080/users/1
```

### Города пользователя

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/users/{id}/cities` | Добавить город |
| GET | `/users/{id}/cities` | Список городов |
| DELETE | `/users/{id}/cities/{city_id}` | Удалить город |

```bash
curl -X POST http://localhost:8080/users/1/cities \
  -d '{"city":"Almaty"}'

curl http://localhost:8080/users/1/cities

curl -X DELETE http://localhost:8080/users/1/cities/3
```

### Погода пользователя

Запрашивает погоду по всем городам пользователя параллельно и сохраняет результаты в историю.

```bash
curl http://localhost:8080/users/1/weather
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

```bash
# Только Алматы, последние 10 записей
curl "http://localhost:8080/users/1/weather/history?city=Almaty&limit=10"

# Вся история с пагинацией
curl "http://localhost:8080/users/1/weather/history?limit=20&offset=40"
```

Параметры запроса:

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
- **Open-Meteo** - погода и геокодинг (без ключа)
- Параллельные запросы к API через горутины
- Кэш погоды в памяти (TTL 5 минут)
- Мягкое удаление пользователей
- Индекс `(user_id, city)` для быстрой фильтрации истории