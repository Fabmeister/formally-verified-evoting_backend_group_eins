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
	"math"
	"slices"

	votingSystems "e-voting-service/logic/dafnyVoting/votingSystems-go/src"
)

func fulfillsPreconditions_InstantRunoffVoting(candidates []int, votes [][]int) (err error) {
	if candidates == nil || votes == nil {
		return errors.New("candidates or votes are nil")
	}
	if len(votes) > math.MaxInt32 {
		return errors.New("unsupported: too many votes")
	}
	if len(candidates) > 0 && (slices.Max(candidates) > math.MaxInt32) {
		return errors.New("unsupported: too high candidateid")
	}

	candMap := make(map[int]int8)
	for i := 0; i < len(candidates); i++ {
		if candidates[i] < 0 {
			return errors.New("preconditions of instant runoff voting weren't met: negative candidate")
		}
		if candidates[i] == 0 {
			return errors.New("preconditions of instant runoff voting weren't met: candidate with id 0")
		}
		_, ok := candMap[candidates[i]]
		if ok {
			// candidate already in candMap
			return errors.New("duplicate candidate for instant runoff voting")
		}
		candMap[candidates[i]] = 1 // inidcate, that number is in candidates (1 is a random choice)
	}

	for i := 0; i < len(votes); i++ {

		// check not voted twice for the same canddiate
		voteCopy := slices.Clone(votes[i])
		slices.Sort(voteCopy)
		voteCopy = slices.Compact(voteCopy)
		if len(votes[i]) != len(voteCopy) {
			return errors.New("precondidtions of instant runoff voting weren't met:  multiple votes for the same candidate")
		}

		for j := 0; j < len(votes[i]); j++ {
			_, ok := candMap[votes[i][j]]
			if !ok {
				// votedCandidate not in candidates
				return errors.New("precondidtions of instant runoff voting weren't met: vote for candidate not in candidates")
			}
		}
	}
	return nil
}

func DafnyCaller_InstantRunoffVoting(candidates []int, votes [][]int) (winner int, err error) {
	// Input:
	// candidates: Kandidaten-ID's die alle > 0 und kein Eintrag kommt doppelt
	// votes: Liste an Votes. Ein Vote besteht ist eine (potentiell leere) List an Kandidaten.
	// 	Der Kandidat an erster Stelle findet der Wähler am besten, den zweiten am zweitbesten usw.
	//  Die Preconditions der Dafny Spezifikation müssen erfüllt sein, sonst wird ein error zurückgeben
	// Output:
	//   winner: 0 falls kein Gewinner gefunden worden konnte, sonst die Kandidaten-ID des/der Gewinner/in

	if candidates == nil || votes == nil {
		err := errors.New("candidates or votes are nil in DafnyCaller_InstantRunoffVoting")
		return 0, err
	}
	precondidtionsEnsuredErr := fulfillsPreconditions_InstantRunoffVoting(candidates, votes)
	if precondidtionsEnsuredErr != nil {
		return 0, precondidtionsEnsuredErr
	}

	// Übersetzen von candidates in <set<nat>>
	var candidatesSet dafny.Set = sliceInt_to_DafnySet(candidates)

	// Übersetzen von votes in seq<seq<nat>>
	var votesSeq dafny.Sequence = sliceSliceInt_to_dafnySeqSeqInt(votes)

	// Wahl auswerten lassen von Dafny-Funktion
	companionStruct := votingSystems.Companion_Default___
	winnerDafnyInt := companionStruct.CalculateElectionResult__IRV(candidatesSet, votesSeq)

	winner = int(winnerDafnyInt.Int32())
	// winner == 0 means no winner was found

	return winner, nil
}
