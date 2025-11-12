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
