## Architecture Layers of the project

- Router
- Controller
- Usecase
- Repository
- Domain

![Go Backend Clean Architecture Diagram](https://gitlab.com/naonweh-studio/bubbme-backend/blob/main/assets/go-backend-arch-diagram.png?raw=true)

## Major Packages used in this project

- **gin**: Gin is an HTTP web framework written in Go (Golang). It features a Martini-like API with much better performance -- up to 40 times faster. If you need a smashing performance, get yourself some Gin.
- **mongo go driver**: The Official Golang driver for MongoDB.
- **mysql go driver**: The Official Golang driver for Mysql.
- **jwt**: JSON Web Tokens are an open, industry-standard RFC 7519 method for representing claims securely between two parties. Used for Access Token and Refresh Token.
- **viper**: For loading configuration from the `.env` file. Go configuration with fangs. Find, load, and unmarshal a configuration file in JSON, TOML, YAML, HCL, INI, envfile, or Java properties formats.
- **bcrypt**: Package bcrypt implements Provos and Mazières's bcrypt adaptive hashing algorithm.
- **testify**: A toolkit with common assertions and mocks that plays nicely with the standard library.
- **mockery**: A mock code autogenerator for Golang used in testing.
- Check more packages in `go.mod`.

### Public API Request Flow without JWT Authentication Middleware

![Public API Request Flow](https://gitlab.com/naonweh-studio/bubbme-backend/blob/main/assets/go-arch-public-api-request-flow.png?raw=true)

### Private API Request Flow with JWT Authentication Middleware

> JWT Authentication Middleware for Access Token Validation.

![Private API Request Flow](https://gitlab.com/naonweh-studio/bubbme-backend/blob/main/assets/go-arch-private-api-request-flow.png?raw=true)

### How to run this project?

We can run this Go Backend Clean Architecture project with or without Docker. Here, I am providing both ways to run this project.

- Clone this project

```bash
# Move to your workspace
cd your-workspace

# Clone this project into your workspace
git clone https://gitlab.com/naonweh-studio/bubbme-backend.git

# Move to the project root directory
cd go-backend-clean-architecture
```

#### Run without Docker

- Create a file `.env` similar to `.env.example` at the root directory with your configuration.
- Install `go` if not installed on your machine.
- Install `MongoDB` if not installed on your machine.
- Important: Change the `DB_HOST` to `localhost` (`DB_HOST=localhost`) in `.env` configuration file. `DB_HOST=mongodb` is needed only when you run with Docker.
- Run `go run cmd/main.go`.
- Access API using `http://localhost:8080`

#### Run with Docker

- Create a file `.env` similar to `.env.example` at the root directory with your configuration.
- Install Docker and Docker Compose.
- Run `docker-compose up -d`.
- Access API using `http://localhost:8080`

### How to run the test?

```bash
# Run all tests
go test ./...
```

### How to generate the mock code?

In this project, to test, we need to generate mock code for the use-case, repository, and database.

```bash
# Generate mock code for the usecase and repository
mockery --dir=domain --output=domain/mocks --outpkg=mocks --all

# Generate mock code for the database
mockery --dir=mongo --output=mongo/mocks --outpkg=mocks --all
```

Whenever you make changes in the interfaces of these use-cases, repositories, or databases, you need to run the corresponding command to regenerate the mock code for testing.

### The Complete Project Folder Structure

```
.
├── Dockerfile
├── api
│   ├── controller
│   │   ├── login_controller.go
│   │   ├── profile_controller.go
│   │   ├── profile_controller_test.go
│   │   ├── refresh_token_controller.go
│   │   ├── signup_controller.go
│   │   └── task_controller.go
│   ├── middleware
│   │   └── jwt_auth_middleware.go
│   └── route
│       ├── login_route.go
│       ├── profile_route.go
│       ├── refresh_token_route.go
│       ├── route.go
│       ├── signup_route.go
│       └── task_route.go
├── bootstrap
│   ├── app.go
│   ├── database.go
│   └── env.go
├── cmd
│   └── main.go
├── docker-compose.yaml
├── domain
│   ├── error_response.go
│   ├── jwt_custom.go
│   ├── login.go
│   ├── profile.go
│   ├── refresh_token.go
│   ├── signup.go
│   ├── success_response.go
│   ├── task.go
│   └── user.go
├── go.mod
├── go.sum
├── internal
│   └── tokenutil
│       └── tokenutil.go
├── mongo
│   └── mongo.go
├── repository
│   ├── task_repository.go
│   ├── user_repository.go
│   └── user_repository_test.go
└── usecase
    ├── login_usecase.go
    ├── profile_usecase.go
    ├── refresh_token_usecase.go
    ├── signup_usecase.go
    ├── task_usecase.go
    └── task_usecase_test.go
```

### API documentation of Go Backend Clean Architecture

<a href="https://documenter.getpostman.com/view/391588/2s8Z75S9xy" target="_blank">
    <img alt="View API Doc Button" src="https://gitlab.com/naonweh-studio/bubbme-backend/blob/main/assets/button-view-api-docs.png?raw=true" width="200" height="60"/>
</a>

### Example API Request and Response

- signup

  - Request

  ```
  curl --location --request POST 'http://localhost:8080/signup' \
  --data-urlencode 'email=test@gmail.com' \
  --data-urlencode 'password=test' \
  --data-urlencode 'name=Test Name'
  ```

  - Response

  ```json
  {
    "accessToken": "access_token",
    "refreshToken": "refresh_token"
  }
  ```

- login

  - Request

  ```
  curl --location --request POST 'http://localhost:8080/login' \
  --data-urlencode 'email=test@gmail.com' \
  --data-urlencode 'password=test'
  ```

  - Response

  ```json
  {
    "accessToken": "access_token",
    "refreshToken": "refresh_token"
  }
  ```

- profile

  - Request

  ```
  curl --location --request GET 'http://localhost:8080/profile' \
  --header 'Authorization: Bearer access_token'
  ```

  - Response

  ```json
  {
    "name": "Test Name",
    "email": "test@gmail.com"
  }
  ```

- task create

  - Request

  ```
  curl --location --request POST 'http://localhost:8080/task' \
  --header 'Authorization: Bearer access_token' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'title=Test Task'
  ```

  - Response

  ```json
  {
    "message": "Task created successfully"
  }
  ```

- task fetch

  - Request

  ```
  curl --location --request GET 'http://localhost:8080/task' \
  --header 'Authorization: Bearer access_token'
  ```

  - Response

  ```json
  [
    {
      "title": "Test Task"
    },
    {
      "title": "Test Another Task"
    }
  ]
  ```

- refresh token

  - Request

  ```
  curl --location --request POST 'http://localhost:8080/refresh' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'refreshToken=refresh_token'
  ```

  - Response

  ```json
  {
    "accessToken": "access_token",
    "refreshToken": "refresh_token"
  }
  ```
  

install golang migrate 

  ```json
go install -tags "postgres,mysql" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate create -ext sql -dir db/namefile

1. migrate create -ext sql -dir db/migrations -seq create_table_category
2. create databases
3. migrate -database 'mysql://root:ZXCasdqwe123!@tcp(127.0.0.1:3306)/bubbme_db' -path db/migrations up
  ```

Noted: 
```
#show user granted
select * from information_schema.user_privileges;
```

proxySQL 
```
mysql -u admin -p admin -h 127.0.0.1 -P6032 --prompt='Admin> '

mysql -h mysqldb -P 3306 -u root -p -e 'select * from bubbme_db.user_admin'

SELECT * FROM mysql_users;

select @@admin-hash_passwords;

```



Kong Api Gateway 
```
https://docs.konghq.com/gateway/latest/get-started/
```
  
