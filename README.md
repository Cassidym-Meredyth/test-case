# Weather Test Case (Go + Docker + Prometheus)

Данный репозиторий содержит простое веб приложение на **Golang**, разработанное в рамках тестового задания.
<br>Приложение отображает текузую температуру в выбранном городе с использованием сервиса **WeatherAPI** и упаковано в **Docker-контейнер**.
<br>Доступность приложения контроллируется через HTTP-chec, а также собираются простые метрики в Prometheus.

## Функционал

- HTTP-сервер на Go (`net/http`, `html/template`)
- Получение текущей температуры из [WeatherAPI](https://www.weatherapi.com/my/)
- Передача города через query-параметр `?city=` (По умолчаниюЖ `Moscow`)
- Эндпоинт для проверки доступности (HTTP-check): `/healthcheck`
- Эндпоинт метрик Prometheus: `/metrics`
    - счетчик HTTP-запросов к главной странице - параметр `reqTotal`
    - гистограмма длительности запросов к WeatherAPI - параметр `weatherDuration`
    - счетчик ошибок при запросах к WeatherAPI - параметр `weatherErrors`
- Dockerfile для сборки образа приложения
- `docker-compose.yml` для запуска стека:
    - приложение на Golang
    - Prometheus

## Требования 

- Само веб-приложение было написано на Golang 1.25, но сервис может работать на версиях 1.22+
- Docker и Docker Compose, но есть возможность локального запуска (инструкция по локальному запуску [ниже](#локальный-запуск-без-docker))
- Аккаунт WeatherAPI и API-ключ

## Настройка

В корне проекта создайте файл `.env`:

```
WEATHER_API_KEY=YOUR_WEATHER_API_KEY
```

**PS:** для того, чтобы проверка проекта проходила быстро, файл .env так же был добавлен в репозиторий **в любом другом случае данный файл не был бы добавлен**.

## Локальный запуск (без Docker)

В терминале проекта необходимо ввести команды:

```golang
go mod tidy
export WEATHER_API_KEY=YOUR_WEATHER_API_KEY
go run main.go
```

Доступные эндпоинты:
- `http://localhost:8080/` - главная страница, текущая температура в Москве
    Есть возможность просмотра температуры в других городах:
    `http://localhost:8080/?city=London` - температура в Лондоне
- `http://localhost:8080/healthcheck` - HTTP-check (возвращает `statuscode - OK`)
- `http://localhost:8080/metrics` - метрики Prometheus

## Запуск в Docker

1. Сборка образа:
```docker
docker build -t weather-app:latest
```
2. Запуск контейнера:
```docker
docker run -p 8080:8080 -e WEATHER_API_KEY=YOUR_WEATHER_API_KEY weather-app:latest
```

## Запуск полного стека (веб-приложение + Prometheus)
```docker
docker build -t weather-app:latest
docker compose up -d
```

Порты по умолчанию:
- Приложение: `http://localhost/` (порт 80 на хосте, проброшен на 8080 в контейнере)
- Prometheus: `http://localhost:9090`

Prometheus собирает:
- метрики приложения с `app:8080/metrics`

## Структура проекта

```
test-case
├───html
│   ├───style
│   │   └───style.css
│   └───index.html
├───prometheus
│   └───prometheus.yml
├───.env
├───.gitignore
├───docker-compose.yml
├───Dockerfile
├───go.mod
├───go.sum
├───main.go
└───README.md
```