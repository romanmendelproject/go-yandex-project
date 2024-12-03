# go-yandex-project

Репозитория для системы хранения приватных данных

## Запуск программы
1. Клонируем репозиторий и переходим в него
2. Запускаем БД
- cd docker
- docker-compose up -d
3. Запуск сервера
- go run cmd/server/main.go

## Запуск тестов
1. Клонируем репозиторий и переходим в него
2. Запускаем БД
- cd docker
- docker-compose up -d
3. go test ./...  -coverprofile cover.out 
4. go tool cover -func cover.out
