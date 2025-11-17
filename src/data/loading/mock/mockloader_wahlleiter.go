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
