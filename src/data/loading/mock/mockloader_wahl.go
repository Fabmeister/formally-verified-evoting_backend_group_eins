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

package mock

import (
	"e-voting-service/data/dto"
	"fmt"
	"maps"
	"slices"
	"sync"
	"time"
)

var (
	wahlen         = map[int]dto.Election{} // nur mit Mutex! evtl umbau auf sync.Map? Hat aber ned so schöne api
	mu_wahl        sync.Mutex
	maxId_wahl     int                         = 1
	vote_headers                               = map[int][]dto.VoteHeader{}  // index nach Wahlid, dann Liste der VoteHeader
	single_votes                               = map[int][]dto.UnifiedVote{} // index nach Headerid, einfach unified speichern, ist einfacer...
	vote_mutex     sync.Mutex                                                // mutex für vote header und einzelne votes
	candidates     = map[int][]dto.Candidate{}                               // index nach Wahlid, dann Liste der Kandidaten
	canidate_mutex sync.Mutex
)

type MockWahlLoader struct {
}

func (loader MockWahlLoader) InsertElection(wahl *dto.Election, candidateNames []string) error {
	mu_wahl.Lock()
	defer mu_wahl.Unlock()

	wahl.Id = maxId_wahl
	maxId_wahl += 1
	wahlen[wahl.Id] = *wahl // wird kopiert?

	// Kandidaten einfügen
	canidate_mutex.Lock()
	defer canidate_mutex.Unlock()
	for _, name := range candidateNames {
		candidate := dto.Candidate{
			Id:         len(candidates[wahl.Id]) + 1, // mock
			Name:       name,
			ElectionId: wahl.Id,
		}
		candidates[wahl.Id] = append(candidates[wahl.Id], candidate)
	}

	return nil
}

func (loader MockWahlLoader) GetElection(wahlid int) (dto.Election, error) {
	mu_wahl.Lock()
	defer mu_wahl.Unlock()

	return wahlen[wahlid], nil
}

func (loader MockWahlLoader) insertVoteHelper(votes []dto.UnifiedVote, electionid int) {
	// Header
	var nextVoteHeaderID int = len(vote_headers[electionid]) + 1
	voteHeader := dto.VoteHeader{Id: nextVoteHeaderID, ElectionId: electionid}

	vote_headers[electionid] = append(vote_headers[electionid], voteHeader)

	single_votes[voteHeader.Id] = make([]dto.UnifiedVote, len(votes))
	copy(single_votes[voteHeader.Id], votes)
}

func (loader MockWahlLoader) InsertVoteAndInvalidateToken(votes []dto.UnifiedVote, wahltoken dto.Wahltoken) error {
	vote_mutex.Lock()
	defer vote_mutex.Unlock()

	loader.insertVoteHelper(votes, wahltoken.ElectionID)

	// Token invalidieren
	MockWahltokenLoader{}.InvalidateVotertoken(nil, wahltoken) // mock hat keine transaktion

	return nil
}

// Ohne invalidieren
func (loader MockWahlLoader) InsertVotesForOpenElection(votes []dto.UnifiedVote, election dto.Election) error {
	vote_mutex.Lock()
	defer vote_mutex.Unlock()

	loader.insertVoteHelper(votes, election.Id)

	return nil
}

func (loader MockWahlLoader) GetElectionForVoter(wahlid int) (dto.Election, []dto.Candidate, error) {
	mu_wahl.Lock() // eigentlich nur read, aber ned schlecht
	defer mu_wahl.Unlock()

	election := wahlen[wahlid]

	canidate_mutex.Lock()
	defer canidate_mutex.Unlock()
	kandidaten := candidates[wahlid]

	// TODO: error handling falls nicht exisistiert, aber nur Mock, also geringe prio

	return election, kandidaten, nil
}

func (loader MockWahlLoader) GetCandidates(wahlid int) ([]int, map[int]dto.Candidate, error) {
	canidate_mutex.Lock()
	defer canidate_mutex.Unlock()

	var id_list []int = make([]int, 0)
	cand_map := make(map[int]dto.Candidate)

	for _, candidate := range candidates[wahlid] {
		id_list = append(id_list, candidate.Id)
		cand_map[candidate.Id] = candidate
	}

	return id_list, cand_map, nil
}

func (loader MockWahlLoader) GetVotesApproval(wahlid int) ([][]dto.Single_vote_approval, error) {
	vote_mutex.Lock()
	defer vote_mutex.Unlock()

	var votes [][]dto.Single_vote_approval = make([][]dto.Single_vote_approval, 0)

	for _, header := range vote_headers[wahlid] {
		var singleVote []dto.Single_vote_approval = make([]dto.Single_vote_approval, 0)
		for _, singleVoteApproval := range single_votes[header.Id] {
			singleVote = append(singleVote,
				dto.Single_vote_approval{
					Vote_id:      header.Id,
					Candidate_id: int(singleVoteApproval.CandidateID),
					Approved:     singleVoteApproval.WahlInfo == 1,
				},
			)
		}
		votes = append(votes, singleVote)
	}

	return votes, nil
}

func (loader MockWahlLoader) GetVotesType2(wahlid int) ([][]dto.Vote_Type2, error) {
	vote_mutex.Lock()
	defer vote_mutex.Unlock()

	var votes [][]dto.Vote_Type2 = make([][]dto.Vote_Type2, 0)

	for _, header := range vote_headers[wahlid] {
		var singleVote []dto.Vote_Type2 = make([]dto.Vote_Type2, 0)
		for _, singleVoteApproval := range single_votes[header.Id] {
			singleVote = append(singleVote,
				dto.Vote_Type2{
					Vote_id:      header.Id,
					Candidate_id: int(singleVoteApproval.CandidateID),
					Info:         int(singleVoteApproval.WahlInfo),
				},
			)
		}
		votes = append(votes, singleVote)
	}

	return votes, nil
}

func (loader MockWahlLoader) IsElectionActive(wahlid int) (bool, error) {
	mu_wahl.Lock()
	defer mu_wahl.Unlock()

	election, ok := wahlen[wahlid]
	if !ok {
		err := fmt.Errorf("In (MockWahlLoader) IsElectionActive Election with id %d doesn't exist", wahlid)
		return false, err
	}

	if !(election.Is_active) {
		return false, nil
	} else {
		if time.Now().After(election.End_time) {
			// election should have already ended, but field not yet updated
			election.Is_active = false
			wahlen[wahlid] = election
		}
		return election.Is_active, nil
	}
}

func (loader MockWahlLoader) GetElectionsOfWahlleiter(wahlleiterid int) ([]dto.Election, error) {
	mu_wahl.Lock()
	defer mu_wahl.Unlock()

	var listAllElections []dto.Election = slices.Collect(maps.Values(wahlen))
	electionsOfWahlleiter := make([]dto.Election, 0)

	for _, election := range listAllElections {
		if election.Wahlleiter_id == wahlleiterid {
			electionsOfWahlleiter = append(electionsOfWahlleiter, election)
		}
	}
	return electionsOfWahlleiter, nil

}
