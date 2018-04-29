# gIntercept
Proxy-tool for gRPC/protobuf requests.

This is a part of a master thesis project and should not be used in production.

## Installing

To use this proxy tool for intercepting protobuf: clone the repo and generate certificates.

```
git clone https://github.com/Diddern/gRPC-simpleGCDService && cd gRPC-simpleGCDService/
openssl req -x509 -newkey rsa:4096 -keyout certs/server-key.pem -out certs/server-cert.pem -days 365 -nodes -subj '/CN=localhost'

```

Then start the server, the proxy and then test with the client.

```
go run gcd/main.go
```

```
cd ../gIntercept/
go run ../gIntercept/cmd/main.go
```

```
go run ../gRPC-simpleGCDService/client/main.go 294 462
```
