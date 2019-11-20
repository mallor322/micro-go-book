# Micro-Go-Pracrise
Micro-Go-Pracrise

### setup

```sh
go mod tidy
```
### 安装 protobuf

```
protoc --version
go get github.com/golang/protobuf
go install github.com/golang/protobuf/protoc-gen-go/
```

```
protoc string.proto --go_out=plugins=grpc:.
```