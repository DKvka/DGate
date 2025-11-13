# DGate

A simple and lightweight API Gateway/reverse proxy made for a self hosted server collection and written in Go.

## Features:

- JSON based configuration for gateway and servers behind it
- Reverse proxying
- Simple, extensible and easy to modify structure
- Websocket passthrough support

## How-to:

### Build:
```shell
go build ./cmd/gateway
```
### Run:
(This config.json path is the default and not required)
```shell
./gateway -cpath ./config.json
```
