// Formally verified E-Voting using Dafny
// Copyright (C) 2025 Authors Gruppe EinS
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"e-voting-service/api"
	"e-voting-service/data/configuration"
	databaseconn "e-voting-service/data/dto/connection"
	authservices "e-voting-service/logic/auth_services"
	pb "e-voting-service/proto/proto"

	"google.golang.org/grpc"
)

func ScheduledJobs() {
	for {
		log.Println("Beginning Jobs...")

		log.Println("Starting Cleanup Auth Tokens...")
		authservices.CleanBearerTokens()

		// Hier können weitere Jobs hinzugefügt werden, z.B. zum Wahlenbeenden nach Deadline

		log.Println("Ended Jobs!")
		time.Sleep(time.Minute * 10)
	}
}

func main() {
	flag.Parse()

	// Config laden
	conf := configuration.Read_config("./config.json")

	// Test Datenbank
	err := databaseconn.TestConnection(conf)

	if err != nil {
		log.Printf("Datenbank fehler: %v", err)
	} else {
		log.Printf("Datenbank erfolgreich verbunden")
	}

	// Jobs schedulen (z.B. Auth Token Cleanup)
	go ScheduledJobs()

	// ab hier gRPC Server Connection und starten
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterWahlServicesServer(s, &api.WahlServices_Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
