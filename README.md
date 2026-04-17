# Weather API (Go + chi)

HTTP-сервис для получения текущей погоды через внешний API (Open-Meteo).

---

## Запуск

```bash
go mod tidy
go run ./cmd/app
```

Сервер стартует на:

```
http://localhost:8080
```

---

## Проверка

### 🔹 Healthcheck

```bash
curl http://localhost:8080/health
```

Ответ:

```json
{"status":"ok"}
```

---

### 🔹 Получение погоды

```bash
curl "http://localhost:8080/api/weather?lat=43.2389&lon=76.8897"
```

Пример ответа:

```json
{
  "latitude": 43.2389,
  "longitude": 76.8897,
  "temperature": 18.4,
  "wind_speed": 7.2,
  "weather_code": 1,
  "time": "2026-04-14T14:00",
  "description": "Переменная облачность"
}
```

---

### 🔹 Погода в городе
 
```
GET /weather/{city}
```
 
```bash
curl http://localhost:8080/weather/Almaty
```
 
Ответ:
 
```json
{
  "city": "Almaty",
  "latitude": 43.25,
  "longitude": 76.9167,
  "temperature": 12.3,
  "wind_speed": 5.1,
  "weather_code": 0,
  "time": "2026-04-14T14:00",
  "description": "Ясно",
  "clothing": "куртка"
}
```
---

Поле `clothing` формируется на основе температуры:
- холодно — тёплая одежда
- прохладно — куртка
- тепло — лёгкая одежда

 
--- 

 
### 🔹 Погода по стране
 
```
GET /weather/country/{country}
```
 
Возвращает список городов страны с текущей погодой.
 
```bash
curl http://localhost:8080/weather/country/Kazakhstan
```
 
Ответ:
 
```json
[
  {
    "city": "Almaty",
    "latitude": 43.25,
    "longitude": 76.9167,
    "temperature": 12.3,
    "wind_speed": 5.1,
    "weather_code": 0,
    "time": "2026-04-14T14:00",
    "description": "Ясно"
  },
  {
    "city": "Astana",
    "latitude": 51.1801,
    "longitude": 71.446,
    "temperature": 4.7,
    "wind_speed": 11.2,
    "weather_code": 3,
    "time": "2026-04-14T14:00",
    "description": "Переменная облачность"
  }
]
```
 
---
 
### 🔹 Топ-3 самых тёплых города страны
 
```
GET /weather/country/{country}/top
```
 
Возвращает три города с наибольшей температурой, отсортированные по убыванию.
 
```bash
curl http://localhost:8080/weather/country/Kazakhstan/top
```
 
Ответ:
 
```json
[
  {
    "city": "Shymkent",
    "latitude": 42.3,
    "longitude": 69.6,
    "temperature": 10.1,
    "wind_speed": 23.5,
    "weather_code": 61,
    "time": "2026-04-16T23:45",
    "description": "Дождь"
  },
  {
    "city": "Atyrau",
    "latitude": 47.1167,
    "longitude": 51.8833,
    "temperature": 8.3,
    "wind_speed": 16.2,
    "weather_code": 1,
    "time": "2026-04-16T23:45",
    "description": "Переменная облачность"
  },
  {
    "city": "Almaty",
    "latitude": 43.25,
    "longitude": 76.9167,
    "temperature": 8.2,
    "wind_speed": 1.5,
    "weather_code": 3,
    "time": "2026-04-16T23:45",
    "description": "Переменная облачность"
  }
]
```
 
---
 
## Поддерживаемые страны
 
| Название   | Параметр запроса |
|------------|-----------------|
| Казахстан  | `Kazakhstan`    |
| Россия     | `Russia`        |
| США        | `USA`           |
| Германия   | `Germany`       |
| Франция    | `France`        |
| Китай      | `China`         |
| Япония     | `Japan`         |
| Индия      | `India`         |
| Бразилия   | `Brazil`        |
| Канада     | `Canada`        |
 
---

## Где взять координаты

Проще всего через Google Maps — клик по карте → копировать координаты.

---

## Структура

```
cmd/app         — точка входа
internal/handler — HTTP слой
internal/service — бизнес-логика
internal/client  — внешний API
```

---

## Стек

* Go
* net/http
* go-chi
* JSON
* Open-Meteo API
* Open-Meteo Geocoding

---