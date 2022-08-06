# fizzbuzz-server
Home task for leboncoin - fizzbuzz server implementation  
  
## Table of Contents  
* [fizzbuzz-server](#fizzbuzz-server)
   * [Project structure](#project-structure)
   * [Env vars](#env-vars)
   * [How to run](#how-to-run)  
  
## Project structure  
```shell
.
├── config # load configuration from env vars
│   └── config.go
├── Dockerfile
├── fizzbuzz # fizzbuzz algorithm implementation
│   ├── fizzbuzz.go
│   └── fizzbuzz_test.go
├── go.mod
├── go.sum
├── http # manages the API routes
│   ├── http.go
│   └── http_test.go
├── main.go
└── README.md
```  
  
## Env vars  
  
| Env var | Mandatory | Default | Description                             |  
| ------- | --------- | ------- | --------------------------------------- |  
| PORT    | no        | 8080    | Port on which the API will be listening |  

## How to run  
The simplest way is to use docker:  
 - Build the image: `docker build --tag fizzbuzz-server .`  
 - Run the container (and publish a port) `docker run --publish 8080:8080 -d fizzbuzz-server`  
 - You can now access the API on the port you published
