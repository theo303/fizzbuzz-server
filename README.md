# fizzbuzz-server
Home task for leboncoin - fizzbuzz server implementation  
  
## Table of Contents  
* [fizzbuzz-server](#fizzbuzz-server)
   * [Project structure](#project-structure)
   * [Env vars](#env-vars)
   * [How to run](#how-to-run)  
   * [Endpoints](#endpoints)
  
## Project structure  
```shell
.
├── api # manages the API routes
│   ├── api.go
│   ├── api_test.go # integration test
│   ├── clienterr # formatted error for client
│   │   ├── clienterr.go
│   │   └── clienterr_test.go
│   ├── fizzbuzzhandler # handler for fizzbuzz request
│   │   ├── fizzbuzzhandler.go
│   │   └── fizzbuzzhandler_test.go
│   └── mostfreqreqhandler # handler for mostfreqreq request
│       ├── mostfreqreqhander.go
│       └── mostfreqreqhander_test.go
├── config # load configuration from env vars
│   └── config.go
├── Dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── fizzbuzz # fizzbuzz algorithm implementation
│   │   ├── fizzbuzz.go
│   │   └── fizzbuzz_test.go
│   └── stats # request counter
│       ├── stats.go
│       └── stats_test.go
├── main.go
└── README.md
```  
  
## Env vars  
You can set env vars before starting the program in order to configure it. If you use Docker you can do this in the Dockerfile.  

| Env var   | Mandatory | Default | Description                             |  
| --------- | --------- | ------- | --------------------------------------- |  
| PORT      | no        | 8080    | Port on which the API will be listening |
| LOG_LEVEL | no        | info    | Level minimum for a log to be displayed |

## How to run  
The simplest way is to use docker:  
 - Build the image: `docker build --tag fizzbuzz-server .`  
 - Run the container (and publish a port) `docker run --publish 8080:8080 -d fizzbuzz-server`  
 - You can now access the API on the port you published
  
### Windows usage  
Fizzbuzz-server was not tested on Windows  
  
## Endpoints  
The API has 2 routes availables  
  
### FizzBuzz - /fizzbuzz (GET)
The Fizzbuzz endpoints allows the user to execute the fizzbuzz process on a set of parameters.  
The endpoint is `/fizzbuzz`. The only method accepted is GET.  
The parameters are sent to the API in the body, using JSON.  
request example:  
```json
{
    "int1":3,
    "int2":5,
    "limit":16,
    "str1":"fizz",
    "str2":"buzz"
}
```  
  
If all parameters are correct the API will return the processed list directly.  
response example:  
```json
["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz","16"]
```
  
### Most frequent request - /mostfreqreq (GET)
The most frequent request endpoints allows the user to retrieve the parameters of the most frequent request.  
It only counts requests to the fizzbuzz route with valid parameters.  
The endpoint is `/mostfreqreq`. The only method accepted is GET.  
No parameters are required.  
  
The response will be formatted in JSON and will give the set of parameters for the most frequent request and the number of times it was requested.  
response example:
```json
{
    "count": 1,
    "params": [
        {
            "int1": 3,
            "int2": 5,
            "limit": 16,
            "str1": "fizz",
            "str2": "buzz"
        }
    ]
}
```

## TODO / Improvements  
 - CI
 - Add swagger
 - Use prometheus instead of the stats.FizzbuzzCounter -> This would be usefull if monitoring other routes was needed
