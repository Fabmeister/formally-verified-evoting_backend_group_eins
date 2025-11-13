# Formally-Verified-E-Voting-Backend

Have you ever wanted to vote, but it beeing formally verified? Now you can!

Backend for a formally verified voting system using Dafny.
This Repository contains a gRPC-Server which was written using Golang.

The API is defined using protobufs, which can be found here: <https://github.com/Fabmeister/formally-verified-evoting_protobufs_group_eins/>

## Requirements
Dafny, Go, protoc, mariadb

## Setup

### Protobuf Repo

git submodule update --init

### Database

Create database by sequentialy running all sql commands in DatabaseSchema.sql

### Build server

1. Check for correct version of wahlServices.proto in src/proto/proto (git pull or update submodule)
2. execute build.sh script, evotingServer file in bin/ is beeing created.
   The script generates Go code from the .proto and transpiles Dafny for use in Go

## Host backend

1. In example_external_config/Envoy-gRPC/: sudo docker compose up
2. In bin/: ./evotingServer 

   On first startup, a config.json will be created in the current folder, which has to be filled with proper information.
   Start again and the server runs now.

   Also make sure backend port in the frontend matches that of envoy.

## About and credits

This is a student project made as part of the computer science major at the University Regensburg.
During the duration of the project gitlab was used but now ported to github in way that email addresses wouldn't be leaked. Unfortunatly the git history was lost during this process.

It is currently hosted on: <eldarica.org/eins>

The project group (called Gruppe EinS) consisted of (with the respected area of responsobility):
   1. <https://github.com/Vok321> Backend dev and Dafny
   2. <https://github.com/Fabmeister> Backend dev, architecture and technical planning
   3. <https://github.com/4lekse1> DevOps (in the original gitlab) and most of the unit tests
   4. <https://github.com/JBTastic> Frontend/Web dev
   5. <https://github.com/Tobox-xD> Frontend/Web dev