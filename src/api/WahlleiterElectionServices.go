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

package api

import (
	"context"
	"log"
	"slices"
	"time"

	//"errors"

	"e-voting-service/data/dto"
	"e-voting-service/usecases"

	pb "e-voting-service/proto/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// in der API Schicht die gRPC Funktionen implementieren
// vor allem Datenkonvertierung in und aus Protobuf Typen
// Richtige Logik in der Businesslayer Schicht (Use cases)
// oder Low level Funktionsschicht (logik)

func (s *WahlServices_Server) CreateElection(ctx context.Context, in *pb.CreateElectionRequest) (*pb.WahlTokens, error) {
	var electionCreationData *pb.ElectionCreationData = in.GetElectionCreationData()
	var endTimeAsTime time.Time = electionCreationData.EndTime.AsTime()
	md, _ := metadata.FromIncomingContext(ctx)

	log.Printf(
		"Received createElection-Request:\n"+

			"\tName: %s\n"+
			"\tBeschreibung: %s\n"+
			"\tCandidates: %s\n"+
			"\tEnd Time: %s\n"+
			"\tElection Type: %d\n"+
			"\tOpen Wahl: %t\n"+
			"\tVoter-mails: %v\n"+
			"\tMetadata: %v\n",
		electionCreationData.Name,
		electionCreationData.Beschreibung,
		electionCreationData.Candidates,
		endTimeAsTime.Format(time.DateTime),
		electionCreationData.ElectionType,
		electionCreationData.IsOpen,
		electionCreationData.VoterEmails,
		md["authorization"],
	)
	// TODO MAILS einbinden

	// Überprüfen ob Wahlleiter authentifiziert ist
	isAuth, err := IsAuthenticatedFromHeaderMetadata(md)
	if err != nil {
		log.Printf("Error in IsAuthenticatedFromHeaderMetadata: %v", err)
		return nil, err
	}

	err = nil
	if !isAuth {
		err = status.Errorf(codes.Unauthenticated, "Authtoken nicht angenommen")

		log.Printf("%v", err)
		return nil, err
	}

	// Zugehörige Wahlleiterid aus Bearer Token holen
	wahlleiterid, err := GetWahlleiterIdFromHeaderMetadata(md)
	if err != nil {
		log.Printf("Error in GetWahlleiterIdFromHeaderMetadata: %v", err)
		return nil, err
	}
	log.Printf(
		"found wahlleiterid of request: %d\n",
		wahlleiterid)

	// Überprüfen, ob Wahldaten valide
	if electionCreationData == nil {
		err := status.Errorf(codes.FailedPrecondition, "Election is nil")
		log.Printf("error in CreateElection: %v", err)
		return nil, err
	}

	if electionCreationData.Name == "" {
		err := status.Errorf(codes.FailedPrecondition, "Empty or no Election name")
		log.Printf("error in CreateElection: %v", err)
		return nil, err
	}

	if len(electionCreationData.Candidates) == 0 {
		err := status.Errorf(codes.FailedPrecondition, "No candidates provided")
		log.Printf("error in CreateElection: %v", err)
		return nil, err
	}

	if len(electionCreationData.Password) == 0 {
		err := status.Errorf(codes.FailedPrecondition, "No password provided")
		log.Printf("error in CreateElection: %v", err)
		return nil, err
	}

	if electionCreationData.ElectionType == pb.ElectionType_ELECTION_TYPE_UNSPECIFIED {
		err := status.Errorf(codes.FailedPrecondition, "ElectionType not set")
		log.Printf("error in CreateElection: %v", err)
		return nil, err
	}

	// Überprüfen logischer Preconditions für Wahlerstellung

	// Gültigkeit der EndTime
	if time.Now().After(endTimeAsTime) {
		err := status.Errorf(codes.FailedPrecondition, "Given EndTime of election is in the past")
		log.Printf("error in CreateElection: %v", err)
		return nil, err
	}

	// Falls Election nicht offen ist, muss es votermails geben
	if len(electionCreationData.VoterEmails) == 0 && !electionCreationData.IsOpen {
		return nil, status.Errorf(codes.FailedPrecondition, "no voter emails given")
	}

	// Candidaten-namen nicht doppelt
	candNamesCopy := slices.Clone(electionCreationData.Candidates)
	slices.Sort(candNamesCopy)
	candNamesCopy = slices.Compact(candNamesCopy)
	if len(candNamesCopy) != len(electionCreationData.Candidates) {
		err := status.Errorf(codes.FailedPrecondition, "At least two candidates have the same name")
		log.Printf("error in CreateElection: %v", err)
		return nil, err
	}

	// VoterEmails nicht doppelt
	voterEmailsCopy := slices.Clone(electionCreationData.VoterEmails)
	slices.Sort(voterEmailsCopy)
	voterEmailsCopy = slices.Compact(voterEmailsCopy)
	if len(voterEmailsCopy) != len(electionCreationData.VoterEmails) {
		err := status.Errorf(codes.FailedPrecondition, "Votermail given (at least) twice")
		log.Printf("error in CreateElection: %v", err)
		return nil, err
	}

	election := dto.Election{
		Name:          electionCreationData.Name,
		Beschreibung:  electionCreationData.Beschreibung,
		Wahlleiter_id: wahlleiterid,
		End_time:      endTimeAsTime,
		Type:          dto.ElectionType(electionCreationData.ElectionType),
		Password:      electionCreationData.Password,
		Open_wahl:     electionCreationData.IsOpen,
	}

	tokens, err := usecases.CreateElection_Usecase(election, electionCreationData.Candidates, electionCreationData.VoterEmails)
	if err != nil {
		log.Printf("error in CreateElection_Usecase: %v", err)
		return nil, err
	}

	var mesTok pb.WahlTokens
	mesTok.Tokens = tokens
	log.Printf("CreateElection success!")
	return &mesTok, nil
}

func (s *WahlServices_Server) GetElectionsOfWahlleiter(ctx context.Context, _ *emptypb.Empty) (*pb.GetElectionsOfWahlleiterResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	log.Printf(
		"Received GetElectionsOfWahlleiter-Request:\n"+
			"\tMetadata: %v\n",
		md["authorization"],
	)

	// Überprüfen ob Wahlleiter authentifiziert ist
	isAuth, err := IsAuthenticatedFromHeaderMetadata(md)
	if err != nil {
		log.Printf("Error in IsAuthenticatedFromHeaderMetadata: %v", err)
		return nil, err
	}

	err = nil
	if !isAuth {
		err = status.Errorf(codes.Unauthenticated, "Authtoken nicht angenommen")

		log.Printf("%v", err)
		return nil, err
	}

	// Zugehörige Wahlleiterid aus Bearer Token holen
	wahlleiterid, err := GetWahlleiterIdFromHeaderMetadata(md)
	if err != nil {
		log.Printf("Error in GetWahlleiterIdFromHeaderMetadata: %v", err)
		return nil, err
	}
	log.Printf(
		"found wahlleiterid of request: %d\n",
		wahlleiterid)

	var dtoElections []dto.Election
	dtoElections, err = usecases.GetElectionsOfWahlleiter_Usecase(wahlleiterid)
	if err != nil {
		log.Printf("error calling GetElectionsOfWahlleiter_Usecase: %v", err)
		return nil, err
	}

	// convert []dto.Election into []*pb.Election
	var protoElections []*pb.Election = make([]*pb.Election, len(dtoElections))
	for i, dtoElection := range dtoElections {
		var protoElection *pb.Election = dtoElectionIntoProtoElectionPointer(dtoElection)
		protoElections[i] = protoElection
	}
	return &pb.GetElectionsOfWahlleiterResponse{Elections: protoElections}, nil
}
