
# Overview
`riggs` tool allows you to get geoip data from [Geolocation](https://geolocation-db.com) db, it can encode geoip data to Morse code.

Also it can encode to morse code any string data you pass it.
It is client/server application which can be easily deployed to K8s (deployment files are in `deploy/` folder).

## Configuration

Configuration of both server and client can be done via cli flags or env. vars.

### Server
* `RIGGS_LISTEN` - listen `address:port`, default `127.0.0.1:8080`
* `RIGGS_GEO_DB` - geo db API fqdn, default `geolocation-db.com`
  
### Client (Cli)
* `RIGGS_ADDR` - riggs server `address:port`, default `127.0.0.1:8080`
* `RIGGS_CONN_TIMEOUT` - gRPC connection timeout (in seconds), Default 5
* `RIGGS_OP_TIMEOUT` - operation (call) timeout (in seconds), Default 5


## Build/deploy
### K8s

_Note: Tested on local K8s provided by "Docker for Mac_

By default Riggs service will listen on `0.0.0.0:5050`. K8s Service will expose 8080 port with `LoadBalancer` type.

* Build docker image
  ```
  docker build -t riggs -f deploy/Dockerfile .
  ```
* Deploy server
  ```
  kubectl apply -f deploy/k8s/
  ```  

### K8s Update

Due to lack of external Docker image registry for app updates we need delete currently running pod. So Update process will include following steps:

* Build docker image
  ```
  docker build -t riggs -f deploy/Dockerfile .
  ```
* Restart pod with newly build image
  ```
  kubectl delete pod $(kubectl get pod -l app=riggs -o custom-columns=:metadata.name)
  ```  

  In future this issue is easily solved by templating `image` in k8s deployment spec., storing images in external registry, tagging by `git_commit`

### Local

* Run 
  ```
  CGO_ENABLED=0 go build
  ```

## CLI Usage
After you've built binary locally you can list available commands with
```
riggs cli -h
```

### Geo
* Get geo ip data
  ```
  ./riggs cli geo -ip 8.8.8.8
  ```
* Get geo ip of self IP
  ```
  ./riggs cli geo
  ```
* Get geo ip data with Morse encoding
  ```
  ./riggs cli geo -morsify -ip 8.8.8.8
  ```
* Get geo ip data with Morse encoding with operation timeout (in seconds)
  ```
  ./riggs cli geo -morsify -ip 8.8.8.8 -t 3
  ```

### Morsify
* Get geo ip data
  ```
  ./riggs cli morsify some text goes here foo bar
  ```

## Development
### Protobuf
  Generate PB go code
  ```
  protoc -I=./ --go_out=plugins=grpc+retag:./ ./pb/riggs.proto
  ```
