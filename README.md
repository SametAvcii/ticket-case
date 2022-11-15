
# Gowit-Backend-Case
This repo is prepared for the gowit backend position case.




## Acknowledgements

 - [API Docs](http://localhost:8080/swagger/index.html)



## Run Locally

Clone the project

```bash
  git clone https://github.com/SametAvcii/gowit-case.git
```

Go to the project directory

```bash
  cd gowit-case
```

Install dependencies

```bash
  docker-compose up
```

Start the server

Server is running localhost:8080

## Running Tests

To run tests, run the following command

```bash
  go test -coverprofile=coverage.out ./...

  go tool cover -func=coverage.out
  
  or show browse
  
  go tool cover -html=coverage.out
```

