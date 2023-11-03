.PHONY : protos

protos:
	protoc -I=protos/ --go_out=protos/files --go-grpc_out=protos/files protos/app.proto

