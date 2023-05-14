## CRUD API unit-testing example

Example API built with:
- [Echo](https://echo.labstack.com/) as web framework
- [tesidy](https://github.com/stretchr/testify) for testing
- [Mockery](https://vektra.github.io/mockery) for mock generation
- [go-playground/validator](https://github.com/go-playground/validator) for input validation

## Installation
Install go dependencies:
```
go mod vendor
```
## Generate mocks
Make sure you have mockery installed:
```
go install github.com/vektra/mockery/v2@v2.20.0
```
or
```
brew install mockery
brew upgrade mockery
```
Generate mocks:
```
make mocks
```

## Environment variables:
Create .env:
```
cp .env.example .env
```
## Run:
```
make run
```
## Test:
```
make test
```