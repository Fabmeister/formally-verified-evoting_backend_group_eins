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
	dto "e-voting-service/data/dto"
	authservices "e-voting-service/logic/auth_services"
	pb "e-voting-service/proto/proto"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WahlServices_Server struct {
	pb.UnimplementedWahlServicesServer
}

func dtoElectionIntoProtoElectionPointer(dtoElection dto.Election) *pb.Election {
	pbElection := pb.Election{
		Id:           int32(dtoElection.Id),
		Name:         dtoElection.Name,
		Beschreibung: dtoElection.Beschreibung,
		Wahlleiterid: int32(dtoElection.Wahlleiter_id),
		CreatedAt:    timestamppb.New(dtoElection.Created_at),
		EndTime:      timestamppb.New(dtoElection.End_time),
		Type:         pb.ElectionType(dtoElection.Type),
		IsActive:     dtoElection.Is_active,
	}

	return &pbElection
}

func dtoCandidateIntoProtoCandidatePointer(dtoCandidate dto.Candidate) *pb.Candidate {
	pbCandidate := pb.Candidate{
		Id:         int32(dtoCandidate.Id),
		ElectionId: int32(dtoCandidate.ElectionId),
		Name:       dtoCandidate.Name,
	}

	return &pbCandidate
}

func protoVoteIntoDtoVote(pbVote *pb.UnifiedVotingVote) dto.UnifiedVote {
	return dto.UnifiedVote{
		CandidateID: int32(pbVote.CandidateId),
		WahlInfo:    int32(pbVote.WahlInfo),
	}
}

func protoVoteListIntoDtoVoteList(pbVotes []*pb.UnifiedVotingVote) []dto.UnifiedVote {
	dtoVotes := make([]dto.UnifiedVote, len(pbVotes))
	for i, pbVote := range pbVotes {
		dtoVotes[i] = protoVoteIntoDtoVote(pbVote)
	}
	return dtoVotes
}

// erstmal nur für case Bearer Token
func checkHeadersForAuthToken(md metadata.MD) (string, error) {
	// Geht davon aus, dass "authorization": nur Bearer token hat, nichts anderes
	if len(md["authorization"]) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "Kein Auth Header!")
	}

	authToken := md["authorization"][0]
	if authToken == "Bearer null" {
		return "", status.Errorf(codes.Unauthenticated, "Auth Header leer!")
	}
	if authToken == "" {
		return "", status.Errorf(codes.Unauthenticated, "Auth Header leer!")
	}

	return authToken, nil
}

// sollte nach checkHeadersForAuthToken aufgerufen werden
// Precondition: checkHeadersForAuthToken ist Fehlerfrei
func parseAuthHeaderGetId(authToken string) (int, error) {
	authTokenParts := strings.Split(authToken, " ")

	// wenn man will, handel nach Auth art, aber wir haben nur Bearer
	if len(authTokenParts) != 2 {
		return -1, status.Errorf(codes.Unauthenticated, "Auth Header falsches Format! Erwartet: Bearer <token>")
	}

	bin_token, err := authservices.DecodeToken(authTokenParts[1])

	if err != nil {
		return -1, status.Errorf(codes.Unauthenticated, "Auth Token Konvertierungsfehler: %v", err)
	}

	id := authservices.GetWahlleiteridFromBearerToken(bin_token)

	return id, nil
}

func parseAndCheckAuthHeader(authToken string) (bool, error) {
	id, err := parseAuthHeaderGetId(authToken)
	return id != -1, err
}

func GetWahlleiterIdFromHeaderMetadata(md metadata.MD) (id int, err error) {
	authToken, err := checkHeadersForAuthToken(md)
	if err != nil {
		return -1, err
	}

	id, err = parseAuthHeaderGetId(authToken)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Erwartet erstmal nur Bearer Token
func IsAuthenticatedFromHeaderMetadata(md metadata.MD) (is_auth bool, err error) {
	authToken, err := checkHeadersForAuthToken(md)
	// authToken should now be of form "bearer xxxxxxxxxxx..."
	if err != nil {
		return false, err
	}

	isValid, err := parseAndCheckAuthHeader(authToken)
	if err != nil {
		return false, err
	}

	return isValid, nil
}

// Return either:
//   - voteToken from metadata (of form electionid-token)
//     or bearerToken
func GetVoteORBearerTokenStringFromHeaderMetadata(md metadata.MD) (token string, tokentype authservices.TokenTypes, err error) {

	// VoteToken or Bearer Zeile aus metadata parsen
	if len(md["authorization"]) == 0 {
		return "", authservices.UNDEFINED_TOKEN, status.Errorf(codes.Unauthenticated, "Kein VoteToken!")
	}

	token = md["authorization"][0]
	if token == "" {
		return "", authservices.UNDEFINED_TOKEN, status.Errorf(codes.Unauthenticated, "VoteToken leer!")
	}

	voteTokenParts := strings.Split(token, " ")

	// handel nach Auth art, aber wir haben nur Bearer und VoterToken
	if len(voteTokenParts) != 2 || !(voteTokenParts[0] == "VoterToken" || voteTokenParts[0] == "Bearer") {
		return "", authservices.UNDEFINED_TOKEN, status.Errorf(codes.Unauthenticated, "Auth Header falsches Format! Erwartet: Bearer/VoterToken <token>")
	}

	switch voteTokenParts[0] {
	case "VoterToken":
		tokentype = authservices.TOKEN_WAHLTOKEN
	case "Bearer":
		tokentype = authservices.TOKEN_BEARER
	default:
		tokentype = authservices.UNDEFINED_TOKEN // case wurde schon überprüft
	}

	return voteTokenParts[1], tokentype, nil

}
