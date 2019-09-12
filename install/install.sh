sudo cp -R vscode-package/* ${GOPATH}/src/
sudo go install github.com/mdempsky/gocode
sudo go install github.com/ramya-rao-a/go-outline
sudo go install github.com/uudashr/gopkgs
sudo go install github.com/sqs/goreturns
sudo go install github.com/go-delve/delve/cmd/dlv

echo vscode env ok

sudo cp -R project-package/* ${GOPATH}/src/
sudo go install github.com/golang/protobuf/protoc-gen-go

echo project package ok

cd ./tools
cd m4-1.4.13
./configure --prefix=/usr/local
make
sudo make install
cd ..
cd autoconf-2.65
./configure --prefix=/usr/local # ironic, isn't it?
make
sudo make install
cd ..
cd automake-1.11
./configure --prefix=/usr/local
make
sudo make install
cd ..
cd libtool-2.2.6b
./configure --prefix=/usr/local
make
sudo make install
cd ..

echo pre protoc ok

cd ..
cd ./tools/protobuf
./autogen.sh
./configure
make
#make check
sudo make install

echo protoc ok