protoc --proto_path=./proto \
-I ./proto \
--go_out=proto \
--go_opt=paths=source_relative \
--validate_out "lang=go:./proto" \
--gofast_out=google/protobuf/any.proto=github.com/gogo/protobuf/types,google/protobuf/duration.proto=github.com/gogo/protobuf/types,google/protobuf/struct.proto=github.com/gogo/protobuf/types,google/protobuf/timestamp.proto=github.com/gogo/protobuf/types,google/protobuf/wrappers.proto=github.com/gogo/protobuf/types,paths=source_relative,plugins=grpc:./proto \
--openapiv2_out=doc \
--openapiv2_opt=logtostderr=true,json_names_for_fields=false,disable_default_errors=true \
./proto/app.proto
