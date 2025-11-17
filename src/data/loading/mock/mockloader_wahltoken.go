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
	"sync"

	"gorm.io/gorm"
)

var (
	wahltokens      = []dto.Wahltoken{}
	mu_wahltokens   sync.Mutex
	maxId_wahltoken int = 1
)

type MockWahltokenLoader struct{}

func (MockWahltokenLoader) InsertVotertokens(tokens []dto.Wahltoken) error {
	mu_wahltokens.Lock()
	defer mu_wahltokens.Unlock()

	for _, token := range tokens {
		token.ID = maxId_wahltoken
		maxId_wahltoken++
		wahltokens = append(wahltokens, token)
	}

	return nil
}

func (MockWahltokenLoader) InsertSingleVotertoken(token dto.Wahltoken) error {
	mu_wahltokens.Lock()
	defer mu_wahltokens.Unlock()

	token.ID = maxId_wahltoken
	maxId_wahltoken++
	wahltokens = append(wahltokens, token)

	return nil
}

func (MockWahltokenLoader) GetVotertokensByElectionid(electionid int) ([]dto.Wahltoken, error) {
	mu_wahltokens.Lock()
	defer mu_wahltokens.Unlock()

	var tokens []dto.Wahltoken
	for _, token := range wahltokens {
		if token.ElectionID == electionid {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

func (MockWahltokenLoader) CheckVotertokenNotYetVoted(inputToken dto.Wahltoken) (bool, error) {
	mu_wahltokens.Lock()
	defer mu_wahltokens.Unlock()

	for _, token := range wahltokens {
		if token.Token == inputToken.Token && token.ElectionID == inputToken.ElectionID {
			return !token.Voted, nil
		}
	}
	return false, nil // Token not found
}

func (MockWahltokenLoader) InvalidateVotertoken(db *gorm.DB, token dto.Wahltoken) error {
	mu_wahltokens.Lock()
	defer mu_wahltokens.Unlock()

	for i, t := range wahltokens {
		if t.ID == token.ID {
			wahltokens[i].Voted = true
			return nil
		}
	}

	return nil // evtl error?
}

func (MockWahltokenLoader) CheckVotertokenExists(token dto.Wahltoken) (bool, error) {
	mu_wahltokens.Lock()
	defer mu_wahltokens.Unlock()

	for _, t := range wahltokens {
		if t.Token == token.Token && t.ElectionID == token.ElectionID {
			return true, nil
		}
	}

	return false, nil
}
