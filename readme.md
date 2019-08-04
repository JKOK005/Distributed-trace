### Distributed trace

### To compile protobuf
Execute the command: 
```go
protoc -I api/proto/v1/ --go_out=plugins=grpc:pkg/api/proto/ api/proto/v1/messages.proto
```
