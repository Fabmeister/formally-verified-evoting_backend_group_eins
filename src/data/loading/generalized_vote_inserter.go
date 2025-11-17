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

package loading

import (
	"e-voting-service/data/dto"
	"errors"
	"log"

	"gorm.io/gorm"
)

type igeneralized_vote_inserter interface {
	insert_votes(votes []dto.UnifiedVote, voteHeader dto.VoteHeader, transaction *gorm.DB) error
}

func singleVoteInserterFactory(wahl dto.Election) (igeneralized_vote_inserter, error) {

	switch wahl.Type {
	case dto.MAJORITY:
		fallthrough
	case dto.APPROVAL_VOTING:
		return insertVoteApprovalVoting{}, nil

	case dto.COMBINED_APPROVAL_VOTING:
		fallthrough
	case dto.SCORE_VOTING:
		fallthrough
	case dto.IRV:
		return insertVoteType2{}, nil
	}

	err := errors.New("inserVoteForSingleVoterFactory called with unsupported election Type")
	return nil, err
}

type insertVoteApprovalVoting struct{}

func (this insertVoteApprovalVoting) insert_votes(votes []dto.UnifiedVote, voteHeader dto.VoteHeader, transaction *gorm.DB) error {
	// Insert a vote of one voter
	for _, vote := range votes {
		// insert single Approvals in Db
		singleVoteApproval := dto.Single_vote_approval{Vote_id: voteHeader.Id, Candidate_id: int(vote.CandidateID), Approved: (vote.WahlInfo > 0)}

		err := transaction.Create(&singleVoteApproval).Error
		if err != nil {
			log.Printf("Insertion of a approval for a singular candidate failed: %v", err)
			return err
		}
	}

	return nil
}

type insertVoteType2 struct{} //use type2 for Combined Approval, Score and Instant Runoff Voting

func (this insertVoteType2) insert_votes(votes []dto.UnifiedVote, voteHeader dto.VoteHeader, transaction *gorm.DB) error {
	for _, vote := range votes {
		singleVote := dto.Vote_Type2{Vote_id: voteHeader.Id, Candidate_id: int(vote.CandidateID), Info: int(vote.WahlInfo)}

		err := transaction.Create(&singleVote).Error
		if err != nil {
			log.Printf("Insertion of a Type 2 Vote for a singular candidate failed: %v", err)
			return err
		}
	}

	return nil
}
