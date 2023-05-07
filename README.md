## CRUD API unit-testing example

Example API built with:
- labstack/echo as web framework
- go-playground/validator for validation
- stretchr/testify for testing

## Installation
Install go dependencies:
```
go mod download
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