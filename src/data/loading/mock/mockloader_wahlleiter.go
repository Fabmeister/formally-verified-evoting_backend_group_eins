package mock

import (
	"e-voting-service/data/dto"
	"sync"
)

var (
	wahlleiter       = map[int]dto.Wahlleiter{}
	mu_wahlleiter    sync.Mutex
	maxId_wahlleiter int = 1
)

type MockWahlleiterLoader struct {
}

func (this MockWahlleiterLoader) GetWahlleiter(id int) (dto.Wahlleiter, error) {
	mu_wahlleiter.Lock()
	defer mu_wahlleiter.Unlock()

	return wahlleiter[id], nil
}

func (this MockWahlleiterLoader) GetWahlleiterByUsername(name string) (dto.Wahlleiter, error) {
	mu_wahlleiter.Lock()
	defer mu_wahlleiter.Unlock()

	for _, wahlleiter := range wahlleiter {
		if wahlleiter.Username == name {
			return wahlleiter, nil
		}
	}

	return dto.Wahlleiter{}, nil // evtl error
}

func (this MockWahlleiterLoader) GetWahlleiterCountByName(name string) (int64, error) {
	mu_wahlleiter.Lock()
	defer mu_wahlleiter.Unlock()

	count := int64(0)
	for _, wahlleiter := range wahlleiter {
		if wahlleiter.Username == name {
			count++
		}
	}

	return count, nil
}

func (this MockWahlleiterLoader) InsertWahlleiter(user *dto.Wahlleiter) error {
	mu_wahlleiter.Lock()
	defer mu_wahlleiter.Unlock()

	user.ID = maxId_wahlleiter
	maxId_wahlleiter++
	wahlleiter[user.ID] = *user

	return nil
}
