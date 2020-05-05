## http_auth_apikey

this example authorize requests depending
 on the passed in apikey header and allows traffic to http://httpbin.org/uuid. Only one value of apikey is supported: 'mashery'. Any other value will result in access forbidden.

### build && run
```
tinygo build -o apikey.wasm -wasm-abi=generic -target wasm ./main.go && docker-compose up
```

now you can make requests authorized randomly:  

```bash
$ curl -H "apikey: mashery" localhost:18000/uuid -v
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 18000 (#0)
> GET /uuid HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< date: Wed, 25 Mar 2020 09:06:33 GMT
< content-type: application/json
< content-length: 53
< server: envoy
< access-control-allow-origin: *
< access-control-allow-credentials: true
< x-envoy-upstream-service-time: 1056
<
{
  "uuid": "e1020f65-f97a-47cd-9b31-368ba2063b6a"
}


# curl localhost:18000/uuid -v
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 18000 (#0)
> GET /uuid HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 403 Forbidden
< content-length: 16
< content-type: text/plain
< powered-by: proxy-wasm-go!!
< date: Wed, 25 Mar 2020 09:07:36 GMT
< server: envoy
<
* Connection #0 to host localhost left intact
access forbidden

```
