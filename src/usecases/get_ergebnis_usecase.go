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

package usecases

import (
	"log"

	"e-voting-service/data/dto"
	"e-voting-service/data/loading"
	authservices "e-voting-service/logic/auth_services"
)

func GetResultUsecase(wahlid int, token string, token_type authservices.TokenTypes) (candidates []dto.Candidate, winners []dto.Candidate, metaInfoPerCandidate map[int]int, err error) {
	// token is either Voter-token or Bearer token according to token_type
	// Output: candidates in election
	//				 winners: winners of election
	//				 metaInfoPerCandidate: maps candidate to their respective score in Approval, Combined Approval, Score Voting
	// 															 empty for Instant Runoff Voting

	loaderWahl := loading.WahlLoaderFactory()

	election, err := loaderWahl.GetElection(wahlid)
	if err != nil {
		log.Printf("GetWahl failed: %v", err)
		return nil, nil, nil, err
	}

	switch token_type {
	case authservices.TOKEN_WAHLTOKEN:
		token_election_id, raw_wahltoken, err := authservices.ParseVoterTokenString(token)
		if err != nil {
			return nil, nil, nil, err
		}
		if token_election_id != wahlid {
			return nil, nil, nil, dto.WahltokenNotValidError{Message: "Wahltoken != election id"}
		}

		loaderToken := loading.WahltokenLoaderFactory()
		token_exists, err := loaderToken.CheckVotertokenExists(dto.Wahltoken{ElectionID: token_election_id, Token: raw_wahltoken})
		if err != nil {
			return nil, nil, nil, err
		}
		if !token_exists {
			return nil, nil, nil, dto.WahltokenNotValidError{Message: "Wahltoken not valid"}
		}

	case authservices.TOKEN_BEARER:
		token_bytes, err := authservices.DecodeToken(token)
		if err != nil {
			return nil, nil, nil, dto.WahltokenNotValidError{Message: err.Error()} // auch wenn auth token, konvertierung damit api layer leichter richtigen code zurückgeben kann
		}
		wahlleiter_id := authservices.GetWahlleiteridFromBearerToken(token_bytes)
		if wahlleiter_id == -1 {
			return nil, nil, nil, dto.WahltokenNotValidError{Message: "Wahlleiter token not valid"}
		}

		if election.Wahlleiter_id != wahlleiter_id {
			return nil, nil, nil, dto.WahltokenNotValidError{Message: "Wahlleiter from authtoken != wahlleiter from requested election"}
		}
	default:
		return nil, nil, nil, dto.WahltokenNotValidError{Message: "Unknown token type"}
	}

	// check election is still ongoing
	isActive, err := loaderWahl.IsElectionActive(wahlid)
	if err != nil {
		log.Printf("GetWahl failed: %v", err)
		return nil, nil, nil, err
	}
	if isActive {
		log.Printf("Get Ergebnis for ongoing election called")
		return nil, nil, nil, dto.ElectionStillActive{}
	}

	candidates_id, cand_map, err := loaderWahl.GetCandidates(wahlid)
	if err != nil {
		log.Printf("error in GetResultUsecase calling GetCandidates: %v", err)
		return nil, nil, nil, err
	}
	var winners_ids []int

	getErgebnisHandle, err := getErgebnisHandlerFactory(election)

	if err != nil {
		return nil, nil, nil, err
	}

	getErgebnisHandle.handle_load_votes(election.Id, loaderWahl)
	winners_ids, metaInfoPerCandidate, err = getErgebnisHandle.send_to_danfy(candidates_id)

	// z.B. preconditions verletzt vorm Dafny Code, Fehler beim laden aus der DB sollten auch vom Dafny code abgefangen werden
	if err != nil {
		return nil, nil, nil, err
	}

	winners = make([]dto.Candidate, len(winners_ids))
	for i, id := range winners_ids {
		winners[i] = cand_map[id]
		// ist garantiert enthalten (oder bei IRV, wenn 0 bedeutet kein Gewinner eh mit Defaultwerten gefüllt)
	}

	candidates = make([]dto.Candidate, len(candidates_id))
	for i, id := range candidates_id {
		candidates[i] = cand_map[id]
	}

	return candidates, winners, metaInfoPerCandidate, nil
}
