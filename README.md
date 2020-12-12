# Beer Api
This is a project for a challenge in Golang. This api is responsible for saving the information of the beers and making currency exchange at the time of payment.

The following technologies were used in this project:
- [Golang 1.15](https://golang.org/dl/)
- [go-chi](https://github.com/go-chi/chi)
- [go-chi/cors](https://github.com/go-chi/cors)
- [pq](https://github.com/lib/pq)
- [godotenv](https://github.com/joho/godotenv)
- [yalm](https://github.com/go-yaml/yaml)
- [errors](https://github.com/pkg/errors)
- [errwrap](https://github.com/hashicorp/errwrap)
- [testify](https://github.com/stretchr/testify)
- [go-cmp](https://github.com/google/go-cmp)
- [mockery](https://github.com/vektra/mockery)
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- [PostgreSQL](https://www.postgresql.org/download/)
- [golang-migrate](https://github.com/golang-migrate/migrate/)

## Requirements
- Golang 14+
- Docker
- Docker Compose
- go-swagger
    * It is necessary to have this plugin installed, this is the installation process:
    ```bash
        go get -u github.com/go-swagger/go-swagger
        go install ./cmd/swagger
    ```
## Generate Mock Interface
This is an automatic mock generator using mockery, the first thing we must do is go to the path of the file that we want to autogenerate:

Download the library
```
go get -u github.com/vektra/mockery
```

We enter the route where you are
```
cd path
```

After entering the route we must execute the following command, Repository this is name the interface
```
mockery -name Repository
```

## Documentation
This is the command that runs the swagger autogenerated documentation
* the document can be generated in the following formats in:
    * yalm
    * json
```bash
swagger generate spec -o ./swagger.json --scan-models
```

This is the command to start the server with the documentation
* You can generate the documentation in _swagger_ and _redoc_ that is changed in the variable `-F=`
```bash
swagger serve -F=redoc swagger.json
```

Raise the server without automatic start, with specific port and path
````
swagger serve -F=redoc --host=0.0.0.0 --port=8082 --no-open swagger.json
````

You can get more information in the Swagger documentation:
```
localhost:8082/docs
```

## Test commands for the project
These are the commands to run the unit and integration tests of the project

#### Execute All Test
This is the command to run the white box tests, and the test report command
```
go test -v -coverprofile=coverage.out -coverpkg=./domain/... ./test/...

go tool cover -html=coverage.out
```
This command gets the total coverage of the project
```
go tool cover -func coverage.out
```

#### Execute All Test Integration
This is the command to run the black box tests, and the test report command
```
go test -v -coverprofile=coverage_integration.out -coverpkg=./domain/... ./test/integration

go tool cover -html=coverage_integration.out
```
This command gets the total coverage of the project
```
go tool cover -func coverage_integration.out
```

#### User Test Handler
The command to run the handler tests, and the command to generate the report
````
go test -v -coverprofile=coverage.out -coverpkg=./domain/beer/application/v1 ./test/handler/beer/v1

go tool cover -html=coverage.out
````

#### User Test Repository
The command to run the handler tests, and the command to generate the report
````
go test -v -coverprofile=coverage.out -coverpkg=./domain/beer/infrastructure/persistence ./test/repository/beer

go tool cover -html=coverage.out
````

#### User Test Integration
The command to run the handler tests, and the command to generate the report for these tests redis must be above
````
go test -v -coverprofile=coverage_integration.out -coverpkg=./domain/beer/application/v1/ ./test/integration

go tool cover -html=coverage_integration.out
````

## Quick Run Project
First clone the repo then go to beer-api folder. After that build your image and run by docker. Make sure you have docker in your machine.

```
git clone https://github.com/samuskitchen/beer-api.git

cd beer-api
```

#### Start Api
```
docker-compose up -d --build
```

#### Down Api
```
docker-compose down --remove-orphans --volumes
```
