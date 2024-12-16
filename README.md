# useragent-value-store-api

This is a user-agent value store, a take on a key-value store where all keys are the user-agent of the HTTP client accessing the API. One string item is stored per user-agent. As many user-agents are allowed to store an item as will fit in memory. Storage is ephemeral and just for kicks :)

Build & Run:

```shell
docker build -t useragent-value-store-api . && docker run -dp 80:80 --name useragent-value-store-api useragent-value-store-api
```

Example use:

```shell
$ curl -X POST http://localhost/v1/value -d test1
{"userAgent":"curl/8.7.1","value":"test1"}

$ curl -X GET http://localhost/v1/value             
{"userAgent":"curl/8.7.1","value":"test1"}

$ curl -X GET http://localhost
{"numUserAgents:":1,"userAgent":"curl/8.7.1"}

$ curl -X DELETE http://localhost/v1/value
```
