# URL Shortener

Простой сервис сокращения ссылок на Go — учебный проект для практики построения
HTTP-сервера "с нуля": без `net/http`, поверх сырого TCP (`net.Listen`), с ручным
парсингом запроса и сборкой ответа.

## Особенности

- **Ноль внешних зависимостей в базовой (in-memory) версии** — только стандартная
  библиотека Go. Единственная внешняя зависимость во всём проекте — драйвер `pgx`,
  и только если включён Postgres-бэкенд.
- **HTTP реализован вручную поверх TCP**: чтение request line и заголовков,
  сборка статус-строки и тела ответа — без `net/http` (см. `internal/web`).
- Каждое соединение обрабатывается в отдельной горутине (`go handleConnection`).
- In-memory хранилище со счётчиком кликов, защищённое `sync.RWMutex` от гонок
  при параллельном доступе.
- Хранилище переключается конфигом (`STORAGE_DRIVER`): in-memory по умолчанию или
  PostgreSQL для персистентности между перезапусками.

## Стек

Go (стандартная библиотека): `net`, `bufio`, `sync`, `encoding/json`.
Опционально, для Postgres-бэкенда: `database/sql` + драйвер
[`github.com/jackc/pgx/v5`](https://github.com/jackc/pgx).

## Как запустить

```bash
go run ./cmd/server
```

Сервер слушает `:8080`, хранилище — in-memory (данные не переживают перезапуск).

### С PostgreSQL

Чтобы данные переживали перезапуск, можно подключить PostgreSQL через Docker Compose.

1. Поднять Postgres:

```bash
docker compose up -d
```

2. Запустить сервер с указанием бэкенда:

```bash
STORAGE_DRIVER=postgres \
DATABASE_URL="postgres://shortener:shortener@localhost:5432/url_shortener?sslmode=disable" \
go run ./cmd/server
```

Таблица создаётся автоматически при старте (`CREATE TABLE IF NOT EXISTS`), отдельно
накатывать схему не нужно.

## Примеры запросов

### Создать короткую ссылку

```bash
curl -X POST http://localhost:8080/links -d 'https://example.com'
```

Ответ: `201 Created`, тело — сгенерированный ID, например `aB3dEfGh`.

### Перейти по короткой ссылке

```bash
curl -i http://localhost:8080/aB3dEfGh
```

Ответ: `302 Found`, заголовок `Location: https://example.com`.

### Статистика по ссылке

```bash
curl http://localhost:8080/stats/aB3dEfGh
```

Ответ: `200 OK`, JSON:

```json
{"id":"aB3dEfGh","url":"https://example.com","clicks":1,"created_at":"..."}
```

## Тесты

```bash
go test ./... -race -cover
```

## Структура

- `cmd/server` — точка входа, TCP-листенер.
- `internal/web` — ручной парсинг HTTP-запроса и сборка ответа.
- `internal/router` — маршрутизация по методу и пути.
- `internal/handler` — обработчики (`Home`, `CreateLink`, `Redirect`, `Stats`, `NotFound`).
- `internal/storage` — интерфейс `Storage`, in-memory и PostgreSQL реализации.
- `internal/model` — доменная модель `Link`.
