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
	"e-voting-service/usecases"

	pb "e-voting-service/proto/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *WahlServices_Server) GetResult(ctx context.Context, in *pb.GetResultRequest) (*pb.GetResultResponse, error) {
	// usage requires a voter-Token OR Bearer Token of Wahlleiter (who must have created the requested election)

	electionid := in.GetElectionId()
	log.Printf("Neuer GetResult-Request, Electionid: %d", electionid)

	// Gibt Wahlergebnis zurück, um "Wahlzettel" für Wähler bereitzustellen
	md, _ := metadata.FromIncomingContext(ctx)

	// Auth laden wahltoken vs wahlleiter auth, prüfen tut der usecase
	token, token_type, err := GetVoteORBearerTokenStringFromHeaderMetadata(md)

	if err != nil {
		log.Printf("GetResult error mit Token: %v", err)
		return nil, err
	}

	candidates, winners, metaInfo, err := usecases.GetResultUsecase(int(electionid), token, token_type)
	if err != nil {
		log.Printf("GetResult error vom usecase: %v", err)

		if reflect.TypeOf(err) == reflect.TypeOf(&dto.WahltokenNotValidError{}) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		return nil, err
	}

	// Konvertieren von int zu int32...
	metaInfoInt32 := make(map[int32]int32)
	for k, v := range metaInfo {
		metaInfoInt32[int32(k)] = int32(v)
	}

	var proto_candidates []*pb.Candidate = make([]*pb.Candidate, len(candidates))
	for i, cand := range candidates {
		proto_candidates[i] = dtoCandidateIntoProtoCandidatePointer(cand)
	}

	var proto_winners []*pb.Candidate = make([]*pb.Candidate, len(winners))
	for i, cand := range winners {
		proto_winners[i] = dtoCandidateIntoProtoCandidatePointer(cand)
	}

	log.Printf("GetResult success!")
	return &pb.GetResultResponse{Candidates: proto_candidates, MetaInfo: metaInfoInt32, Winner: proto_winners}, nil
}
