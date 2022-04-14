.PHONY: protos

protos:
	protoc -I protos/ protos/currency.proto --go_out=protos --go_opt=paths=source_relative
	protoc -I protos/ protos/currency.proto --go-grpc_out=protos --go-grpc_opt=paths=source_relative