Инициализация зависимостей
```
go mod init github.com/EvgeniyBudaev/love-server
```

Сборка
```
go build -v ./cmd/
```

Удаление неиспользуемых зависимостей
```
go mod tidy -v
```

Библиотека для работы с переменными окружения ENV
https://github.com/joho/godotenv
```
go get -u github.com/joho/godotenv
```

ENV Config
https://github.com/kelseyhightower/envconfig
```
go get -u github.com/kelseyhightower/envconfig
```

Логирование
https://pkg.go.dev/go.uber.org/zap
```
go get -u go.uber.org/zap
```

Подключение к БД
Драйвер для Postgres
```
go get -u github.com/lib/pq
```

Миграции
https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md
https://www.appsloveworld.com/go/83/golang-migrate-installation-failing-on-ubuntu-22-04-with-the-following-gpg-error
```
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
sudo sh -c 'echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list'
sudo apt-get update
sudo apt-get install -y golang-migrate
```

Если ошибка E: Указаны конфликтующие значения параметра Signed-By из источника
https://packagecloud.io/golang-migrate/migrate/ubuntu/
jammy: /etc/apt/keyrings/golang-migrate_migrate-archive-keyring.gpg !=
```
cd /etc/apt/sources.list.d
ls
sudo rm migrate.list
```

Создание миграционного репозитория
```
migrate create -ext sql -dir migrations ProfilesCreationMigration
migrate create -ext sql -dir migrations ProfileImagesCreationMigration
```

Создание up sql файлов
```
migrate -path migrations -database "postgres://localhost:5432/tgbot?sslmode=disable&user=postgres&password=root" up
```

Создание down sql файлов
```
migrate -path migrations -database "postgres://localhost:5432/tgbot?sslmode=disable&user=postgres&password=root" down
```

Если ошибка Dirty database version 1. Fix and force version
```
migrate create -ext sql -dir migrations ProfilesCreationMigration force 20240127080026
```

Fiber
https://github.com/gofiber/fiber
```
go get -u github.com/gofiber/fiber/v2
```

CORS
https://github.com/gorilla/handlers
```
go get -u github.com/gorilla/handlers
```