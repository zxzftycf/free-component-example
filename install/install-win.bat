 xcopy vscode-package\*.* %GOPATH%\src\ /s/d
 go install github.com/mdempsky/gocode
 go install github.com/ramya-rao-a/go-outline
 go install github.com/uudashr/gopkgs
 go install github.com/sqs/goreturns
 go install github.com/go-delve/delve/cmd/dlv

echo vscode env ok

 xcopy project-package\*.* %GOPATH%\src\ /s/d
 go install github.com/golang/protobuf/protoc-gen-go

echo project package ok

cd ./tools
cd m4-1.4.13
./configure --prefix=/usr/local

make
 make install
cd ..
cd autoconf-2.65
./configure --prefix=/usr/local # ironic, isn't it?
make
 make install
cd ..
cd automake-1.11
./configure --prefix=/usr/local
make
 make install
cd ..
cd libtool-2.2.6b
./configure --prefix=/usr/local
make
 make install
cd ..

echo pre protoc ok

cd ..
cd ./tools/protobuf
./autogen.sh
./configure
make
#make check
 make install

echo protoc ok

pause