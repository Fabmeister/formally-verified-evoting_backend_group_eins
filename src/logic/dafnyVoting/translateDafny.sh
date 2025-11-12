#!/bin/bash

rm -r votingSystems-go
dafny translate go votingSystems.dfy --include-runtime --verify-included-files
cd votingSystems-go/src/dafny
go mod init dafny
cd ../System_
go mod init System_
