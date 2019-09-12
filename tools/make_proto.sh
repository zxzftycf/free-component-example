cd ../server/grpc
protoc --go_out=plugins=grpc:. *.proto
protoc --js_out=import_style=commonjs,binary:../../client/test/tools/ *.proto
cd ../../client/test/tools
sh change.sh