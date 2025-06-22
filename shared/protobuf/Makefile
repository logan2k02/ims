gen: inventory_protobuf products_protobuf

products_protobuf:
	@protoc --proto_path=./ --go_out=./ --go_opt=paths=source_relative \
		--go-grpc_out=./ --go-grpc_opt=paths=source_relative \
		products.proto

inventory_protobuf:
	@protoc --proto_path=./ --go_out=./ --go_opt=paths=source_relative \
		--go-grpc_out=./ --go-grpc_opt=paths=source_relative \
		inventory.proto