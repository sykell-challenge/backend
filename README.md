# Sykell Challenge - Backend API

This is the backend git submodule for the sykell challenge. 

For details on setup, see the [setup](https://github.com/sykell-challenge/setup) repo.

## Tech stack

- Go (Golang) as the main programming language
- Gin for HTTP web framework
- GORM for ORM/database access (MySQL)
- Colly for web crawling
- Socket IO for realtime messages
- TaskQ for background task queue
- JWT for authentication
- Docker for containerization
- Air for hot reloading in development


## APIs

Import the Postman collection in the `docs` folder to test out the APIs.

## ER Diagram

![alt text](<docs/ER Diagram.png>)

## Next Steps

- Microservice architecture
- Automatic cleanup of long running jobs or very old jobs
- Seed database for better testing
- Writing tests
- Automatic generation of Open API Specs / Swagger UI Docs
- Admin APIs
- Refresh tokens