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
	"e-voting-service/data/dto"
	"e-voting-service/data/loading"
	votingSystemsCaller "e-voting-service/logic/dafnyCaller"
	"sort"

	"log"
)

type iGet_ergebnis_handler interface {
	handle_load_votes(wahl_id int, loader loading.ILoadWahl) error
	send_to_danfy(candid_list []int) (winners []int, scorePerCandidate map[int]int, err error)
}

type handle_get_ergeb_approval struct {
	votes [](map[int]int)
}

func (this *handle_get_ergeb_approval) handle_load_votes(wahl_id int, loader loading.ILoadWahl) error {
	approvals, err := loader.GetVotesApproval(wahl_id)
	if err != nil {
		return err
	}
	this.votes = make([]map[int]int, len(approvals))

	for i, wahlzettel := range approvals {
		var wahlzettel_map map[int]int = make(map[int]int)
		for _, v := range wahlzettel {
			if v.Approved {
				wahlzettel_map[v.Candidate_id] = 1
			}
		}
		this.votes[i] = wahlzettel_map
	}
	return nil
}

func (this *handle_get_ergeb_approval) send_to_danfy(candid_list []int) (winners []int, scorePerCandidate map[int]int, err error) {
	return votingSystemsCaller.DafnyCaller_ApprovalVoting(candid_list, this.votes)
}

func helper_4_Type2_votes(wahl_id int, loader loading.ILoadWahl) ([]map[int]int, error) {
	votes, err := loader.GetVotesType2(wahl_id)
	if err != nil {
		return nil, err
	}
	alle_wahlzettel := make([]map[int]int, len(votes))

	for i, wahlzettel := range votes {
		var wahlzettel_map map[int]int = make(map[int]int)
		for _, v := range wahlzettel {
			wahlzettel_map[v.Candidate_id] = v.Info
		}
		alle_wahlzettel[i] = wahlzettel_map
	}
	return alle_wahlzettel, nil
}

type handle_get_ergeb_combined_approval struct {
	votes [](map[int]int)
}

func (this *handle_get_ergeb_combined_approval) handle_load_votes(wahl_id int, loader loading.ILoadWahl) error {
	votes, err := helper_4_Type2_votes(wahl_id, loader)
	if err != nil {
		return err
	}
	this.votes = votes
	return nil
}

func (this *handle_get_ergeb_combined_approval) send_to_danfy(candid_list []int) (winners []int, scorePerCandidate map[int]int, err error) {
	return votingSystemsCaller.DafnyCaller_CombinedApprovalVoting(candid_list, this.votes)
}

type handle_get_ergeb_scored struct {
	votes [](map[int]int)
}

func (this *handle_get_ergeb_scored) handle_load_votes(wahl_id int, loader loading.ILoadWahl) error {
	votes, err := helper_4_Type2_votes(wahl_id, loader)
	if err != nil {
		return err
	}
	this.votes = votes
	return nil
}

func (this *handle_get_ergeb_scored) send_to_danfy(candid_list []int) (winners []int, scorePerCandidate map[int]int, err error) {
	return votingSystemsCaller.DafnyCaller_ScoreVoting(candid_list, this.votes)
}

type handle_get_ergeb_IRV struct {
	votes [][]int
}

func (this *handle_get_ergeb_IRV) handle_load_votes(wahl_id int, loader loading.ILoadWahl) error {
	votes_db, err := loader.GetVotesType2(wahl_id)
	if err != nil {
		log.Printf("error when callling GetVotesType2: %v", err)
		return err
	}

	this.votes = make([][]int, len(votes_db))

	for i, wahlzettel := range votes_db {
		// sort the votes in ascending order according to wahlinfo (so the candidate with wahlinfo 1 is first)
		sort.Slice(wahlzettel, func(i, j int) bool { return (wahlzettel)[i].Info < (wahlzettel)[j].Info })

		wahlzettel_ints := make([]int, len(wahlzettel))

		indexIntoVote := 0
		for _, vote := range wahlzettel {
			// Sollte schon bei einfügen in DB gefiltered worden sein, aber sicherheitshalber nochmal
			// sichergehen, dass Kandidaten mit Wahlinfo 0 ignoriert werden
			if vote.Info != 0 {
				wahlzettel_ints[indexIntoVote] = vote.Candidate_id
				indexIntoVote++
			}
		}

		if indexIntoVote < len(wahlzettel) {
			wahlzettel = wahlzettel[:indexIntoVote]
		}

		this.votes[i] = wahlzettel_ints
	}

	return nil
}

func (this *handle_get_ergeb_IRV) send_to_danfy(candid_list []int) (winners []int, scorePerCandidate map[int]int, err error) {

	singular_winner, err := votingSystemsCaller.DafnyCaller_InstantRunoffVoting(candid_list, this.votes)

	// if singular_winner is 0 that means no winner could be found:
	if singular_winner == 0 {
		winners = make([]int, 0)
	} else {
		winners = make([]int, 1)
		winners[0] = singular_winner

	}

	// score per Candidate ist leer, gibt keine meta info bei IRV
	return winners, scorePerCandidate, err
}

type handle_get_ergeb_majority struct {
	votes [](map[int]int)
}

func (this *handle_get_ergeb_majority) handle_load_votes(wahl_id int, loader loading.ILoadWahl) error {
	approvals, err := loader.GetVotesApproval(wahl_id)
	if err != nil {
		return err
	}
	this.votes = make([]map[int]int, len(approvals))

	for i, wahlzettel := range approvals {
		var wahlzettel_map map[int]int = make(map[int]int)
		for _, v := range wahlzettel {
			if v.Approved {
				wahlzettel_map[v.Candidate_id] = 1
			}
		}
		this.votes[i] = wahlzettel_map
	}
	return nil
}

func (this *handle_get_ergeb_majority) send_to_danfy(candid_list []int) (winners []int, scorePerCandidate map[int]int, err error) {
	return votingSystemsCaller.DafnyCaller_MajorityVoting(candid_list, this.votes)
}

func getErgebnisHandlerFactory(wahl dto.Election) (iGet_ergebnis_handler, error) {
	switch wahl.Type {
	case dto.APPROVAL_VOTING:
		return &handle_get_ergeb_approval{}, nil
	case dto.COMBINED_APPROVAL_VOTING:
		return &handle_get_ergeb_combined_approval{}, nil
	case dto.SCORE_VOTING:
		return &handle_get_ergeb_scored{}, nil
	case dto.IRV:
		return &handle_get_ergeb_IRV{}, nil
	case dto.MAJORITY:
		return &handle_get_ergeb_majority{}, nil
	}

	return nil, dto.UnifiedVotePreconditionError{Type: dto.ELECTION_TYPE_UNSPECIFIED, Message: "Factory: Kein Typ"}
}
