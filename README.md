# go-uk-holiday-endpoint

Holiday endpoint is a Go application that provides information about bank holidays in the United Kingdom.

Features:

Get a list of bank holidays for a specific year.
Retrieve bank holiday data for England and Wales.
Obtain bank holiday titles and dates for a specific year.

Prerequisites:

Before running the application, ensure you have the following installed:

> Go (Golang) v1.20.1 or later

Deployment:

Download the source from git and run the command from source location.

Eg: D:\go\src\holiday>go run main.go

Endpoint details:

1. To get a list of bank holidays for a specific year:

GET http://localhost:8080/holidays/{year}

2. To retrieve bank holiday data for England and Wales:

GET http://localhost:8080/england-and-wales

3. To obtain bank holiday titles and dates for a specific year:

GET http://localhost:8080/holidays-title-date/{year}
