[ref_source](https://github.com/learning-cloud-native-go/myapp/tree/main)


# Rivall Backend
This backend was constructed for the Rivall mobile app.  This project is part of a undergrad capstone project for Nathaniel Reeves at Utah Tech University in Jan 2025.

## Notes

run tests command
```bash
go test ./api/resources... ./api/router...
```
run tests
-count n
    Run each test, benchmark, and fuzz seed n times (default 1).
    If -cpu is set, run n times for each GOMAXPROCS value.
    Examples are always run once. -count does not apply to
    fuzz tests matched by -fuzz.

-cover
    Enable coverage analysis.
    Note that because coverage works by annotating the source
    code before compilation, compilation and test failures with
    coverage enabled may report line numbers that don't correspond
    to the original sources.

-failfast
    Do not start new tests after the first test failure.

-fullpath
    Show full file names in the error messages.

-json
    Log verbose output and test results in JSON. This presents the
    same information as the -v flag in a machine-readable format.

-v
    Verbose output: log all tests as they are run. Also print all
    text from Log and Logf calls even if the test succeeds.

-coverprofile cover.out
    Write a coverage profile to the file after all tests have passed.
    Sets -cover.

-trace trace.out
    Write an execution trace to the specified file before exiting.

```bash
go test -v -cover ./api/resources... ./api/router...
```

go test help docs
```bash
go help testflag
```

## 🔋 Batteries included

- The idiomatic structure based on the resource-oriented design.
- The usage of Docker, Docker compose, Alpine images, and linters on development.
- Healthcheck and CRUD API implementations with OpenAPI specifications.
- The usage of [Zerolog](https://github.com/rs/zerolog) as the centralized Syslog logger.
- The usage of [Validator.v10](https://github.com/go-playground/validator) as the form validator.
- The usage of GitHub actions to run tests and linters, generate OpenAPI specifications, and build and push production images to the Docker registry.

## 🚀 Endpoints

| Name        | HTTP Method | Route          |
|-------------|-------------|----------------|
| Health      | GET         | /livez         |
|             |             |                |
| List Books  | GET         | /v1/books      |
| Create Book | POST        | /v1/books      |
| Read Book   | GET         | /v1/books/{id} |
| Update Book | PUT         | /v1/books/{id} |
| Delete Book | DELETE      | /v1/books/{id} |

💡 [swaggo/swag](https://github.com/swaggo/swag) : `swag init -g cmd/api/main.go -o .swagger -ot yaml`

## 🗄️ Database design

| Column Name    | Datatype  | Not Null | Primary Key |
|----------------|-----------|----------|-------------|
| id             | UUID      | ✅        | ✅           |
| title          | TEXT      | ✅        |             |
| author         | TEXT      | ✅        |             |
| published_date | DATE      | ✅        |             |
| image_url      | TEXT      |          |             |
| description    | TEXT      |          |             |
| created_at     | TIMESTAMP | ✅        |             |
| updated_at     | TIMESTAMP | ✅        |             |
| deleted_at     | TIMESTAMP |          |             |

## 📦 Container image sizes

- DB: 241MB
- API
    - Development environment: 655MB
    - Production environment: 28MB ; 💡`docker build -f prod.Dockerfile . -t myapp_app`

## 📁 Project structure

```shell
myapp
├── cmd
│  ├── api
│  │  └── main.go
│  └── migrate
│     └── main.go
│
├── api
│  ├── resource
│  │  ├── book
│  │  │  ├── handler.go
│  │  │  ├── model.go
│  │  │  ├── repository.go
│  │  │  └── repository_test.go
│  │  ├── common
│  │  │  └── err
│  │  │     └── err.go
│  │  └── health
│  │     └── handler.go
│  │
│  └── router
│     ├── middleware
│     │  ├── request_id.go
│     │  ├── request_id_test.go
│     │  ├── requestlog
│     │  │  ├── handler.go
│     │  │  └── log_entry.go
│     │  ├── content_type.go
│     │  └── content_type_test.go
│     └── router.go
│
├── migrations
│  └── 00001_create_books_table.sql
│
├── config
│  └── config.go
│
├── util
│  ├── logger
│  │  └── logger.go
│  └── validator
│     └── validator.go
│
├── .env
│
├── go.mod
├── go.sum
│
├── docker-compose.yml
├── Dockerfile
│
├── prod.Dockerfile
└── k8s
   ├── app-configmap.yaml
   ├── app-secret.yaml
   ├── app-deployment.yaml
   └── app-service.yaml
```

## 📸 Form validations and logs
![Form validation](doc/assets/form_validation.png)

![Logs in app init](doc/assets/logs_app_init.png)
![Logs in crud](doc/assets/logs_crud.png)