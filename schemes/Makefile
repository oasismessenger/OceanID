all:generate

generate:
	rm -rf *.go types
	protoc --go_out=../ --go-grpc_out=../ *.proto

	mkdir types
	protoc --go_opt=paths=source_relative --go_out=../schemes/ proto_types/*.proto
	mv proto_types/*.go types/