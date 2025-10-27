# Go CRUD API

Мой первый API на Go 🚀  
Стек:

- Chi (роутер)
- pgx (PostgreSQL драйвер)
- Viper (конфигурация)
- PostgreSQL

## Инструкция для запуска

### Настройка конфигурации

#### Первая настройка:

1. Создайте `config.yaml` из примера:

```bash
cp config.yaml.example config.yaml
```

2. Создайте `.env` из примера:

```bash
cp .env.example .env
```

3. Отредактируйте созданные файлы своими значениями:

**config.yaml** - замените:

- `your_password_here` на ваш пароль БД
- `your_jwt_secret_key_here` на ваш JWT секрет

**.env** - замените:

- `your_password_here` на ваш пароль БД

### Генерация JWT секрета:

Если нужен новый JWT секрет, выполните:

```bash
go run cmd/generate-jwt/main.go
```

### Запуск:

```bash
docker-compose up --build
```

### Применение миграции:

```bash
docker exec go-api-app goose -dir internal/db/migrations up

# или через две команды

docker exec -it go-api-app sh
goose -dir internal/db/migrations up
```
