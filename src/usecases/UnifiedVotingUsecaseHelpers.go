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
	dto "e-voting-service/data/dto"
	"fmt"
	"sort"
)

type iHandleVote interface {
	HandleVotePrecondition(votes *[]dto.UnifiedVote) error
}

type handleApprovalVote struct{}

func (this handleApprovalVote) HandleVotePrecondition(votes *[]dto.UnifiedVote) error {
	voted4Trunc := make([]dto.UnifiedVote, len(*votes))
	vote_len := 0
	for _, vote := range *votes {
		// vote Element in der Menge {0, 1}
		if vote.WahlInfo == 1 {
			voted4Trunc[vote_len] = vote
			vote_len++
		} else if vote.WahlInfo == 0 {
			continue
		} else {
			return dto.UnifiedVotePreconditionError{Type: dto.APPROVAL_VOTING, Message: fmt.Sprintf("Vote für Cand %d mit Element %d", vote.CandidateID, vote.WahlInfo)}
		}
	}
	// Slice kürzen, falls einige Votes keinen Approval haben
	if vote_len < len(voted4Trunc) {
		voted4Trunc = voted4Trunc[:vote_len]
		*votes = voted4Trunc
	}

	return nil
}

type handleCombinedApprovalVote struct{}

func (this handleCombinedApprovalVote) HandleVotePrecondition(votes *[]dto.UnifiedVote) error {
	voted4Trunc := make([]dto.UnifiedVote, len(*votes))
	vote_len := 0
	for _, vote := range *votes {
		// vote Element in der Menge {-1, 0, 1}
		if vote.WahlInfo == 1 || vote.WahlInfo == -1 {
			voted4Trunc[vote_len] = vote
			vote_len++
		} else if vote.WahlInfo == 0 {
			continue
		} else {
			return dto.UnifiedVotePreconditionError{Type: dto.COMBINED_APPROVAL_VOTING, Message: fmt.Sprintf("Vote für Cand %d mit Element %d", vote.CandidateID, vote.WahlInfo)}
		}
	}
	// Slice kürzen, falls einige Votes einen score von 0 bekommen haben (dann haben sie für das Wahlergebnis keinen Einfluss)
	if vote_len < len(voted4Trunc) {
		voted4Trunc = voted4Trunc[:vote_len]
		*votes = voted4Trunc
	}

	return nil
}

type handleScoreVotingVote struct{}

func (this handleScoreVotingVote) HandleVotePrecondition(votes *[]dto.UnifiedVote) error {
	voted4Trunc := make([]dto.UnifiedVote, len(*votes))
	vote_len := 0
	for _, vote := range *votes {
		// vote Element in der Menge {0, 1, 2, 3, 4, 5}
		if vote.WahlInfo >= 1 && vote.WahlInfo <= 5 {
			voted4Trunc[vote_len] = vote
			vote_len++
		} else if vote.WahlInfo == 0 {
			continue
		} else {
			return dto.UnifiedVotePreconditionError{Type: dto.SCORE_VOTING, Message: fmt.Sprintf("Vote für Cand %d mit Element %d", vote.CandidateID, vote.WahlInfo)}
		}
	}
	// Slice kürzen, falls einige Votes keinen Approval haben
	if vote_len < len(voted4Trunc) {
		voted4Trunc = voted4Trunc[:vote_len]
		*votes = voted4Trunc
	}

	return nil
}

type handleIRVVote struct{}

func (this handleIRVVote) HandleVotePrecondition(votes *[]dto.UnifiedVote) error {
	sort.Slice(*votes, func(i, j int) bool { return (*votes)[i].WahlInfo < (*votes)[j].WahlInfo })

	voted4Trunc := make([]dto.UnifiedVote, len(*votes))

	vote_len := 0
	var prev_rank int32 = 0
	for _, vote := range *votes {
		if vote.WahlInfo == 0 {
			continue
		} else if vote.WahlInfo == prev_rank+1 { // sequentiel immer eins mehr
			voted4Trunc[vote_len] = vote
			vote_len++
			prev_rank = vote.WahlInfo
		} else {
			return dto.UnifiedVotePreconditionError{Type: dto.IRV, Message: fmt.Sprintf("Vote für Cand %d mit Element %d, prev = %d", vote.CandidateID, vote.WahlInfo, prev_rank)}
		}
	}

	if vote_len < len(voted4Trunc) {
		voted4Trunc = voted4Trunc[:vote_len]
		*votes = voted4Trunc
	}

	return nil
}

type handleMajorityVote struct{}

func (this handleMajorityVote) HandleVotePrecondition(votes *[]dto.UnifiedVote) error {
	err := handleApprovalVote{}.HandleVotePrecondition(votes)

	if err == nil {
		if len(*votes) > 1 {
			err = dto.UnifiedVotePreconditionError{Type: dto.MAJORITY, Message: fmt.Sprintf("Zu viele Votes: %d", len(*votes))}
		}
	} else {
		// Casten um Vote type zu ändern
		if cast_err, ok := err.(dto.UnifiedVotePreconditionError); ok {
			err = dto.UnifiedVotePreconditionError{Type: dto.MAJORITY, Message: cast_err.Message}
		}
	}

	return err
}

func unifiedVotingHandleFactory(wahl dto.Election) (iHandleVote, error) {
	switch wahl.Type {
	case dto.APPROVAL_VOTING:
		return handleApprovalVote{}, nil
	case dto.COMBINED_APPROVAL_VOTING:
		return handleCombinedApprovalVote{}, nil
	case dto.SCORE_VOTING:
		return handleScoreVotingVote{}, nil
	case dto.IRV:
		return handleIRVVote{}, nil
	case dto.MAJORITY:
		return handleMajorityVote{}, nil
	}

	return nil, dto.UnifiedVotePreconditionError{Type: dto.ELECTION_TYPE_UNSPECIFIED, Message: "Factory: Kein Typ"}
}
