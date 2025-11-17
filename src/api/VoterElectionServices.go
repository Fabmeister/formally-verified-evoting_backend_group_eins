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
	"reflect"

	"e-voting-service/data/dto"
	authservices "e-voting-service/logic/auth_services"
	"e-voting-service/usecases"

	pb "e-voting-service/proto/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *WahlServices_Server) GetElectionForVoter(ctx context.Context, _ *emptypb.Empty) (*pb.GetElectionResponse, error) {
	// Gibt Wahldaten zurück, um "Wahlzettel" für Wähler bereitzustellen

	md, _ := metadata.FromIncomingContext(ctx)

	wahltokenString, _, err := GetVoteORBearerTokenStringFromHeaderMetadata(md) // Tokentype ignorieren, da später abgefangen wird
	if err != nil {
		log.Printf("error getting VoteToken from Header: %v", err)
		return nil, err
	}

	electionid, rawToken, err := authservices.ParseVoterTokenString(wahltokenString)
	if err != nil {
		log.Printf("error in wahltoken: %v", err)
		return nil, err
	}

	log.Printf(
		"Received GetElectionForVoter-Request:\n"+
			"\tElection id: %d\n",
		electionid,
	)

	var wahltoken dto.Wahltoken = dto.Wahltoken{ElectionID: electionid, Token: rawToken}

	var dtoCandidates []dto.Candidate
	var dtoElection dto.Election
	dtoElection, dtoCandidates, err = usecases.GetElectionForVoter_usecase(int(electionid), wahltoken)
	if err != nil {
		log.Printf("error in GetElection_usecase: %v", err)
		return nil, err
	}

	// Convert data transfer objects into correct type  for grpc Response

	// Convert dtocandidates into []*pb.Candidate
	var protoCandidates []*pb.Candidate = make([]*pb.Candidate, len(dtoCandidates))
	for i, cand := range dtoCandidates {
		protoCandidates[i] = dtoCandidateIntoProtoCandidatePointer(cand)
	}

	// convert election into *pb.Election
	var protoElection *pb.Election = dtoElectionIntoProtoElectionPointer(dtoElection)

	log.Printf("GetElectionForVoter success!")
	return &pb.GetElectionResponse{Election: protoElection, Candidates: protoCandidates}, nil
}

func (s *WahlServices_Server) SendVote(ctx context.Context, in *pb.SendVoteRequest) (*emptypb.Empty, error) {
	// Überprüfen ob Voter wahlberechtigt und einfügen des Votes in die Datenbank
	md, _ := metadata.FromIncomingContext(ctx)

	wahltokenString, tokentype, err := GetVoteORBearerTokenStringFromHeaderMetadata(md)
	if err != nil {
		log.Printf("error getting VoteToken from Header: %v", err)
		return nil, err
	}
	if tokentype != authservices.TOKEN_WAHLTOKEN {
		log.Printf("Non Votertoken in Sendvote")
		err := status.Errorf(codes.FailedPrecondition, "Non Votertoken in Metadata")
		return nil, err
	}

	electionid, rawToken, err := authservices.ParseVoterTokenString(wahltokenString)
	if err != nil {
		log.Printf("error in wahltoken: %v", err)
		return nil, err
	}
	var wahltoken dto.Wahltoken = dto.Wahltoken{ElectionID: electionid, Token: rawToken}

	votes := in.GetVotes()

	log.Printf(
		"Received sendVote-Request:\n"+
			"\telectionid id: %d\n"+
			"\tvoted Candidates ids: %v\n",
		electionid,
		votes,
	)
	if electionid == 0 {
		err := status.Errorf(codes.FailedPrecondition, "no electionid given")
		log.Printf("error in SendVote: %v", err)
		return nil, err
	}

	if len(votes) != 0 {
		// Überprüfen, ob für einen Kandidat doppelt gewählt wurde
		cand_map := map[int32]bool{}
		for _, vote := range votes {
			cand_map[vote.CandidateId] = true
		}

		if len(cand_map) != len(votes) {
			err := status.Errorf(codes.FailedPrecondition, "voted (at least) twice for the same candidate")
			log.Printf("error in SendVote: %v", err)
			return nil, err
		}
	}

	err = usecases.HandleVote_usecase(protoVoteListIntoDtoVoteList(votes), wahltoken)
	if err != nil {
		log.Printf("error in handleVote_usecase %v:", err)

		if reflect.TypeOf(err) == reflect.TypeOf(&dto.WahltokenNotValidError{}) { // soll anscheinend isInstanceof sein
			return nil, status.Error(codes.PermissionDenied, err.Error())
		} else if reflect.TypeOf(err) == reflect.TypeOf(&dto.UnifiedVotePreconditionError{}) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		return nil, err
	}
	log.Printf("SendVote success")
	return &emptypb.Empty{}, nil
}

func (s *WahlServices_Server) GetVotertokenStatus(ctx context.Context, in *emptypb.Empty) (*pb.GetVotertokenStatusResponse, error) {
	// Überprüfen ob Voter wahlberechtigt und einfügen des Votes in die Datenbank
	md, _ := metadata.FromIncomingContext(ctx)

	wahltokenString, tokentype, err := GetVoteORBearerTokenStringFromHeaderMetadata(md)
	if err != nil {
		log.Printf("error getting VoteToken from Header: %v", err)
		return nil, err
	}
	log.Printf(
		"Received GetVotertokenStatus-Request:\n"+
			"\tVotertoken: %v\n",
		wahltokenString,
	)
	if tokentype != authservices.TOKEN_WAHLTOKEN {
		log.Printf("Non Votertoken in Sendvote")
		err := status.Errorf(codes.FailedPrecondition, "Non Votertoken in Metadata")
		return nil, err
	}

	electionid, rawToken, err := authservices.ParseVoterTokenString(wahltokenString)
	if err != nil {
		log.Printf("error in wahltoken: %v", err)
		return nil, err
	}

	var wahltoken dto.Wahltoken = dto.Wahltoken{ElectionID: electionid, Token: rawToken}
	tokenExists, tokenUnused, err := usecases.GetVotertokenStatus_Usecase(wahltoken)
	if err != nil {
		log.Printf("error in GetVotertokenStatus_Usecase: %v", err)
		return nil, err
	}
	log.Printf("DEBUG exists %v unused %v", tokenExists, tokenUnused)

	return &pb.GetVotertokenStatusResponse{TokenExists: tokenExists, TokenUnused: tokenUnused}, nil
}
