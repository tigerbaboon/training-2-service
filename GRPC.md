# gRPC Get start
## Install Protocol Buffers v3
Install the protoc compiler that is used to generate gRPC service code. The simplest way to do this is to download pre-compiled binaries for your platform(protoc-<version>-<platform>.zip) from here: https://github.com/google/protobuf/releases

- Unzip this file.
- Update the environment variable PATH to include the path to the protoc binary file.

Next, install the protoc plugin for Go
```bash
$ go get -u github.com/golang/protobuf/protoc-gen-go
```
The compiler plugin, protoc-gen-go, will be installed in $GOBIN, defaulting to $GOPATH/bin. It must be in your $PATH for the protocol compiler, protoc, to find it.
```bash
$ export PATH=$PATH:$GOPATH/bin
```

## Test Install
```bash
$ protoc --version
# libprotoc 3.8.0
$ which protoc-gen-go
# .../bin//protoc-gen-go
```

## Use
Compile proto file to Golang
```bash
$ protoc --go_out=plugins=grpc:. app/grpc/pb/*.proto
```