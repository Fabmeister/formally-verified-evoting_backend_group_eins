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

package votingSystemsCaller

import (
	"dafny"
	"errors"
	"fmt"

	votingSystems "e-voting-service/logic/dafnyVoting/votingSystems-go/src"
)

type generalScoreVotingType int8

const (
	votingtype_ApprovalVoting generalScoreVotingType = iota
	votingtype_CombinedApprovalVoting
	votingtype_ScoreVoting
	votingtype_Majority
)

func fulfillsPreconditions_General_Score_Voting_Systems(candidates []int, votes [](map[int]int), lowestAllowedVotingScore int, highestAllowedVotingScore int) (err error) {
	// Use for Approval Voting, Score Voting and Combined Approval Voting.
	// e.g for Score Voting lowestAllowedVotingScore must be 0 and highest 5
	if candidates == nil || votes == nil {
		return errors.New("candidates or voters nil")
	}
	if len(candidates) == 0 {
		return errors.New("no candidate given")
	}
	// Check all votes are for existing candidate and have a valid score (in {0, 1, 2, 3, 4, 5})

	// for this create map to be able to check if votedCandidate is in candidates
	candidatesMap := make(map[int]int)
	for i := 0; i < len(candidates); i++ {
		if candidates[i] < 0 {
			return errors.New("given candidate is not a natural number")
		}
		candidatesMap[candidates[i]] = 1 // arbitrary constant to check existence in map
	}
	for _, singleVote := range votes {
		// handle a single vote
		for votedCandiate, givenScore := range singleVote {
			if !(lowestAllowedVotingScore <= givenScore && givenScore <= highestAllowedVotingScore) {
				return errors.New("gave invalid score for this voting system to a candidate")
			}
			_, ok := candidatesMap[votedCandiate]
			if !ok {
				return errors.New("vote for nonexistent candidate")
			}
		}
	}
	return nil
}

func DafnyCallerForGeneralScoreVotingSystems(candidates []int, votes [](map[int]int), votingSystemType generalScoreVotingType) (winners []int, scorePerCandidate map[int]int, err error) {
	// general function to use for Caller for Approval, Combined Approval and Score Voting

	// Panic catcher
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic in Dafny Caller: %s", r)
		}
	}()

	if candidates == nil || votes == nil {
		err := errors.New("candidates or votes are nil in DafnyCallerForGeneralScoreVotingSystems")
		return nil, nil, err
	}

	var precondidtionsEnsuredErr error
	switch votingSystemType {
	case votingtype_ApprovalVoting:
		precondidtionsEnsuredErr = fulfillsPreconditions_General_Score_Voting_Systems(candidates, votes, 0, 1)
	case votingtype_CombinedApprovalVoting:
		precondidtionsEnsuredErr = fulfillsPreconditions_General_Score_Voting_Systems(candidates, votes, -1, 1)
	case votingtype_ScoreVoting:
		precondidtionsEnsuredErr = fulfillsPreconditions_General_Score_Voting_Systems(candidates, votes, 0, 5)
	case votingtype_Majority:
		precondidtionsEnsuredErr = fulfillsPreconditions_General_Score_Voting_Systems(candidates, votes, 0, 1)
	default:
		return make([]int, 0), make(map[int]int), errors.New("no valid voting system in DafnyCallerForGeneralScoreVotingSystems")
	}

	if precondidtionsEnsuredErr != nil {
		return make([]int, 0), make(map[int]int), precondidtionsEnsuredErr
	}

	// Übersetzen von candidates in set<nat>
	candidatesSet := sliceInt_to_DafnySet(candidates)

	// Übersetzen von votes in seq<map<int,int>>
	votesSeq := sliceOfmapInt2Int_to_DafnySeqOfMapInt2Int(votes)

	// Wahl auswerten lassen von Dafny-Funktion
	companionStruct := votingSystems.Companion_Default___
	var winnersDafnySet dafny.Set
	var scorePerCandidateDafnyMap dafny.Map

	switch votingSystemType {
	case votingtype_ApprovalVoting:
		winnersDafnySet, scorePerCandidateDafnyMap = companionStruct.CalculateElectionResult__APPROVAL(candidatesSet, votesSeq)
	case votingtype_CombinedApprovalVoting:
		winnersDafnySet, scorePerCandidateDafnyMap = companionStruct.CalculateElectionResult__COMBINED__APPROVAL(candidatesSet, votesSeq)
	case votingtype_ScoreVoting:
		winnersDafnySet, scorePerCandidateDafnyMap = companionStruct.CalculateElectionResult__SCORE(candidatesSet, votesSeq)
	case votingtype_Majority:
		winnersDafnySet, scorePerCandidateDafnyMap = companionStruct.CalculateElectionResult__Majority(candidatesSet, votesSeq)
	}

	// Übersetzen votesPerCandidateDafnyMap in Golang-map
	scorePerCandidate, err = dafnyMapInt2Int_to_mapInt2Int(scorePerCandidateDafnyMap)
	if err != nil {
		return make([]int, 0), make(map[int]int), errors.New("error casting votespercandidate from dafny map to golang map")
	}

	// Übersetzen winnersDafnySet in Golang-slice
	winners, err = dafnySetInt_to_sliceInt(winnersDafnySet)
	if err != nil {
		return make([]int, 0), make(map[int]int), err
	}

	return winners, scorePerCandidate, nil
}

func DafnyCaller_ApprovalVoting(candidates []int, votes [](map[int]int)) (winners []int, scorePerCandidate map[int]int, err error) {
	// Input:
	// candidates: Kandidaten-ID's die alle >=0 und kein Eintrag kommt doppelt vor
	// votes: map von candidate auf 0 oder 1
	// Die Preconditions der Dafny Spezifikation müssen erfüllt sein, sonst wird ein error zurückgeben

	winners, scorePerCandidate, err = DafnyCallerForGeneralScoreVotingSystems(candidates, votes, votingtype_ApprovalVoting)

	if err != nil {
		return make([]int, 0), make(map[int]int), err
	}
	return winners, scorePerCandidate, nil
}

func DafnyCaller_MajorityVoting(candidates []int, votes [](map[int]int)) (winners []int, scorePerCandidate map[int]int, err error) {
	// Input:
	// candidates: Kandidaten-ID's die alle >=0 und kein Eintrag kommt doppelt vor
	// votes: map von candidate auf 0 oder 1
	// Die Preconditions der Dafny Spezifikation müssen erfüllt sein, sonst wird ein error zurückgeben

	winners, scorePerCandidate, err = DafnyCallerForGeneralScoreVotingSystems(candidates, votes, votingtype_Majority)

	if err != nil {
		return make([]int, 0), make(map[int]int), err
	}
	return winners, scorePerCandidate, nil
}

func DafnyCaller_CombinedApprovalVoting(candidates []int, votes [](map[int]int)) (winners []int, scorePerCandidate map[int]int, err error) {
	// Input:
	// candidates: Kandidaten-ID's die alle >=0 und kein Eintrag kommt doppelt vor
	// votes: map von candidate auf {-1, 0, 1}
	// Die Preconditions der Dafny Spezifikation müssen erfüllt sein, sonst wird ein error zurückgeben

	winners, scorePerCandidate, err = DafnyCallerForGeneralScoreVotingSystems(candidates, votes, votingtype_CombinedApprovalVoting)
	if err != nil {
		return make([]int, 0), make(map[int]int), err
	}
	return winners, scorePerCandidate, nil
}

func DafnyCaller_ScoreVoting(candidates []int, votes [](map[int]int)) (winners []int, scorePerCandidate map[int]int, err error) {
	// Input:
	// candidates: Kandidaten-ID's die alle >=0 und kein Eintrag kommt doppelt vor
	// votes: map von candidate auf {0, 1, 2, 3, 4, 5}
	// Die Preconditions der Dafny Spezifikation müssen erfüllt sein, sonst wird ein error zurückgeben

	winners, scorePerCandidate, err = DafnyCallerForGeneralScoreVotingSystems(candidates, votes, votingtype_ScoreVoting)
	if err != nil {
		return make([]int, 0), make(map[int]int), err
	}
	return winners, scorePerCandidate, nil
}
