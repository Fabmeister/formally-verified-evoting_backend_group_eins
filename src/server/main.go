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
