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
	databaseconn "e-voting-service/data/dto/connection"
	"e-voting-service/data/loading/mock"
	"log"
	"strconv"
)

type ILoadWahlleiter interface {
	GetWahlleiter(id int) (dto.Wahlleiter, error)
	GetWahlleiterByUsername(name string) (dto.Wahlleiter, error)
	GetWahlleiterCountByName(name string) (int64, error)
	InsertWahlleiter(user *dto.Wahlleiter) error
}

type WahlleiterLoader struct {
}

func WahlleiterLoaderFactory() ILoadWahlleiter {
	var db_loader ILoadWahlleiter = WahlleiterLoader{}
	var mock_loader ILoadWahlleiter = mock.MockWahlleiterLoader{}
	return generalizedLoaderFactory(db_loader, mock_loader)
}

func (this WahlleiterLoader) GetWahlleiter(id int) (dto.Wahlleiter, error) {
	db, err := databaseconn.GetDB()
	var wahlleiter dto.Wahlleiter
	wahlleiter.ID = -1

	if err != nil {
		log.Printf("Fehler in loader Wahlleiter GetWahlleiter (db Pointer): %v", err.Error())
		return wahlleiter, err
	}

	res := db.First(&wahlleiter, id)

	if res.Error != nil {
		log.Printf("Fehler in loader Wahlleiter GetWahlleiter: %v", res.Error.Error())
		return wahlleiter, &databaseconn.LoaderError{Key: "Id", KeyValue: strconv.Itoa(id), Table: "Wahlleiter"}
	}

	return wahlleiter, nil
}

func (this WahlleiterLoader) GetWahlleiterByUsername(name string) (dto.Wahlleiter, error) {
	db, err := databaseconn.GetDB()
	var wahlleiter dto.Wahlleiter
	wahlleiter.ID = 0

	if err != nil {
		log.Printf("Fehler in loader Wahlleiter GetWahlleiter (db Pointer): %v", err.Error())
		return wahlleiter, err
	}

	res := db.Where("username = ?", name).First(&wahlleiter)

	if res.Error != nil {
		log.Printf("Fehler in loader Wahlleiter GetWahlleiter: %v", res.Error.Error())
		return wahlleiter, &databaseconn.LoaderError{Key: "username", KeyValue: name, Table: "Wahlleiter"}
	}

	return wahlleiter, nil
}

func (this WahlleiterLoader) GetWahlleiterCountByName(name string) (int64, error) {
	db, err := databaseconn.GetDB()
	var count int64

	if err != nil {
		log.Printf("Fehler in loader Wahlleiter GetWahlleiterCountByName (db Pointer): %v", err.Error())
		return -1, err
	}

	res := db.Model(&dto.Wahlleiter{}).Where("username = ?", name).Count(&count)

	if res.Error != nil {
		log.Printf("Fehler in loader Wahlleiter GetWahlleiterCountByName: %v", res.Error.Error())
		return -1, &databaseconn.LoaderError{Key: "username", KeyValue: name, Table: "Wahlleiter"}
	}

	return count, nil
}

func (this WahlleiterLoader) InsertWahlleiter(user *dto.Wahlleiter) error {
	db, err := databaseconn.GetDB()

	if err != nil {
		log.Printf("Fehler in loader Wahlleiter InsertWahlleiter (db Pointer): %v", err.Error())
		return err
	}

	var nextID int
	err = db.Raw("SELECT NEXT VALUE FOR seq_wahlleiter;").Scan(&nextID).Error
	if err != nil {
		log.Printf("failed to get next ID from sequence: %v", err)
		return err
	}

	user.ID = nextID
	res := db.Create(user)

	if res.Error != nil {
		log.Printf("Fehler in loader Wahlleiter InsertWahlleiter: %v", res.Error.Error())
		return &databaseconn.LoaderError{Key: "username", KeyValue: user.Username, Table: "Wahlleiter"}
	}

	return nil
}

func (this WahlleiterLoader) InsertVotesForOpenElection(candidateids []int32, electionid int32) error {
	// UNIMPLEMENTED
	return nil
}
