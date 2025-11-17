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
	"errors"
	"log"

	dto "e-voting-service/data/dto"
	loading "e-voting-service/data/loading"
)

func GetElectionForVoter_usecase(electionid int, wahlToken dto.Wahltoken) (dto.Election, []dto.Candidate, error) {

	// Überpüfen ob Wahltoken existiert
	var wahltokenLoader = loading.WahltokenLoaderFactory()
	tokenExists, err := wahltokenLoader.CheckVotertokenExists(wahlToken)
	if err != nil {
		return dto.Election{}, nil, err
	}
	if !tokenExists {
		err = errors.New("given wahltoken not valid or already used")
		return dto.Election{}, nil, err
	}

	var wahlLoader loading.ILoadWahl = loading.WahlLoaderFactory()
	election, candidates, err := wahlLoader.GetElectionForVoter(electionid)
	if err != nil {
		return dto.Election{}, nil, err
	}
	return election, candidates, nil
}

func HandleVote_usecase(votes []dto.UnifiedVote, wahltoken dto.Wahltoken) error {
	// Sorgt für das einfügen des Votes in die Datenbank
	// Input:
	// 		wahltoken muss mit electionid und Token befüllt sein
	// 		votedCandidateids sind alle ids, die vom Wähler "approved" sind

	var loaderWahl loading.ILoadWahl = loading.WahlLoaderFactory()

	// get Election for voter
	election, _, err := loaderWahl.GetElectionForVoter(wahltoken.ElectionID)
	if err != nil {
		log.Printf("error in GetElection to figure out if election is open or closed: %v", err)
		return err
	}

	//check election hasn't already ended
	isActive, err := loaderWahl.IsElectionActive(election.Id)
	if err != nil {
		log.Printf("IsElectionActive failed: %v", err)
		return err
	}

	if !isActive {
		log.Printf("In HandleVote_usecase already finished election called")
		return dto.ElectionAlreadyEnded{}
	}

	// Already checked that not voted twice for a candidate in API layer

	handle, err := unifiedVotingHandleFactory(election)

	if err != nil {
		return err
	}

	err = handle.HandleVotePrecondition(&votes)

	if err != nil {
		return err
	}

	// Generalisiert
	if election.Open_wahl {
		err = loaderWahl.InsertVotesForOpenElection(votes, election)
		log.Printf("last error %v", err)
	} else {
		err = loaderWahl.InsertVoteAndInvalidateToken(votes, wahltoken)
	}
	return err
}

func GetVotertokenStatus_Usecase(wahltoken dto.Wahltoken) (tokenExists bool, tokenUnused bool, err error) {

	var wahltokenLoader = loading.WahltokenLoaderFactory()
	tokenExists, err = wahltokenLoader.CheckVotertokenExists(wahltoken)
	if err != nil {
		return false, false, err
	}
	if !tokenExists {
		return false, false, nil
	}

	tokenUnused, err = wahltokenLoader.CheckVotertokenNotYetVoted(wahltoken)
	if err != nil {
		return false, false, err
	}
	if tokenUnused {
		return true, true, nil
	} else {
		return true, false, nil
	}

}
