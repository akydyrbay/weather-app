# Weather App

REST API сервис на Go

Проект разделён на два сервиса:
- **api-service** - основной сервис: REST API, бизнес-логика, PostgreSQL.
- **gateway-service** - отдельный сервис, ходит во внешний публичный API (Open-Meteo) и проксирует данные.

## Запуск

```bash
docker compose up --build
```

## Переменные окружения

### api-service
| Переменная | Default | Описание |
|------------|---------|----------|
| `ADDR` | `:8080` | адрес HTTP-сервера |
| `DATABASE_URL` | — | строка подключения к PostgreSQL |
| `JWT_SECRET` | — | секрет для подписи JWT |
| `GATEWAY_URL` | — | базовый URL gateway-сервиса |

### gateway-service
| Переменная | Default | Описание |
|------------|---------|----------|
| `ADDR` | `:8081` | адрес HTTP-сервера |
| `OPEN_METEO_URL` | `https://api.open-meteo.com/v1/forecast` | endpoint погоды |
| `GEOCODING_URL` | `https://geocoding-api.open-meteo.com/v1/search` | endpoint геокодинга |
| `REQUEST_TIMEOUT` | `10s` | таймаут `http.Client` |

## Эндпоинты api-service (`:8080`)

### Healthcheck
```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
```

### Аутентификация
| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/auth/register` | Регистрация |
| POST | `/auth/login` | Вход, возвращает JWT |

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Ali","email":"ali@example.com","password":"secret"}'

curl -X POST http://localhost:8080/auth/login \
  -d '{"email":"ali@example.com","password":"secret"}'
# → {"access_token":"<jwt>"}
```

Защищённые маршруты требуют заголовок:
```
Authorization: Bearer <access_token>
```

### Погода (публично)
| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/weather?lat=&lon=` | Погода по координатам |
| GET | `/weather/{city}` | Погода в городе + рекомендация по одежде |
| GET | `/weather/country/{country}` | Погода по городам страны |
| GET | `/weather/country/{country}/top` | Топ-3 самых тёплых города |

### Пользователи (admin)
| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/users` | Список пользователей |
| GET | `/users/{id}` | Получить пользователя |
| DELETE | `/users/{id}` | Мягкое удаление |

### Текущий пользователь
```bash
curl http://localhost:8080/users/me -H "Authorization: Bearer <token>"
```

### Города пользователя
| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/cities` | Добавить город |
| GET | `/cities` | Список городов |
| DELETE | `/cities/{city_id}` | Удалить город |

### Погода пользователя + история
```bash
curl http://localhost:8080/users/weather -H "Authorization: Bearer <token>"
curl "http://localhost:8080/weather/history?city=Almaty&limit=10" \
  -H "Authorization: Bearer <token>"
```

## Эндпоинты gateway-service (`:8081`)

Используются api-service, но доступны и снаружи через `localhost:8081`:

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/health` | Health |
| GET | `/weather?lat=&lon=` | Погода по координатам |
| GET | `/weather/city/{city}` | Погода по городу |
| GET | `/weather/country/{country}` | Погода по списку городов страны |
| GET | `/weather/country/{country}/top` | Топ-3 тёплых |

```bash
curl "http://localhost:8081/weather?lat=43.25&lon=76.92"
curl http://localhost:8081/weather/city/Almaty
```

## Стек

- **Go**, **net/http**, **go-chi/chi** - сервер и роутинг
- **pgx/v5** - PostgreSQL (только api-service)
- **golang-jwt/jwt/v5** + **bcrypt** - auth
- **Open-Meteo** - погода и геокодинг (только gateway-service)
- **Uber Zap** - structured logging + logging middleware (method, path, status, duration, request_id)
- **testify** (assert / require / mock) - unit-тесты, мок-репозитории
- **Docker Compose** - оркестрация двух сервисов + PostgreSQL

## Тесты

```bash
cd api-service
go test ./...

TEST_DATABASE_URL=postgres://weather:weather@localhost:5432/weather_db?sslmode=disable \
  go test ./internal/repository/...
```

Мок репозитории живут в `api-service/internal/service/mocks/` - используются в `service` и `handler` юнит тестах, не ходя в БД.
