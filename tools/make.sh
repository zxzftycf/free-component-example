cd ..
rm -rf bin
mkdir -p bin
cd bin
go build ../server/main.go
cp ../server/layout.json layout.json
cp main account_server 
cp main websocket_server
cp main data_server