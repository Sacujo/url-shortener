# URL Shortener

Простой сервис сокращения ссылок на Go — учебный проект для практики построения
HTTP-сервера "с нуля": без `net/http`, поверх сырого TCP (`net.Listen`), с ручным
парсингом запроса и сборкой ответа.

## Особенности

- **Ноль внешних зависимостей** — только стандартная библиотека Go.
- **HTTP реализован вручную поверх TCP**: чтение request line и заголовков,
  сборка статус-строки и тела ответа — без `net/http` (см. `internal/web`).
- Каждое соединение обрабатывается в отдельной горутине (`go handleConnection`).
- In-memory хранилище со счётчиком кликов, защищённое `sync.RWMutex` от гонок
  при параллельном доступе.

## Стек

Go (стандартная библиотека): `net`, `bufio`, `sync`, `encoding/json`.

## Как запустить

```bash
go run ./cmd/server
```

Сервер слушает `:8080`.

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
- `internal/storage` — интерфейс `Storage` и in-memory реализация.
- `internal/model` — доменная модель `Link`.
