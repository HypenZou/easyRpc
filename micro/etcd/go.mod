module etcd

go 1.13

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/wubbalubbaaa/easyRpc v0.0.0-20201116144448-f320fdb04c6b
	go.uber.org/zap v1.16.0 // indirect
	google.golang.org/genproto v0.0.0-20201116144945-7adebfbe6a3f // indirect
	google.golang.org/grpc v1.33.2 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
