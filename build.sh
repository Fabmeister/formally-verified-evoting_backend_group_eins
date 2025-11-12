#!/bin/bash


pushd . > /dev/null
cd src/proto/proto
echo "building proto files..."
rm *.pb.go
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wahlServices.proto
popd > /dev/null

pushd . > /dev/null

cd src/logic/dafnyVoting
echo "Verifying Dafny and transpiling it into Go..."

rm -r votingSystems-go
dafny translate go votingSystems.dfy --include-runtime --verify-included-files
pushd . > /dev/null
cd votingSystems-go/src/dafny
go mod init dafny
popd > /dev/null
cd votingSystems-go/src/System_
go mod init System_

popd > /dev/null

echo "building Go..."
mkdir bin
cd src
go build -o ../bin/evotingServer server/main.go 
