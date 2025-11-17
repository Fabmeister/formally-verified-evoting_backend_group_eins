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
	"log"

	"e-voting-service/data/dto"
	databaseconn "e-voting-service/data/dto/connection"
	"e-voting-service/data/loading/mock"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Allgemein: Loopkup hier nicht über ID, sondern über Wahlid und Token wert (beide im dto.Wahltoken)
type ILoadWahltokens interface {
	InsertVotertokens(tokens []dto.Wahltoken) error                     // mit Transaktion
	InsertSingleVotertoken(token dto.Wahltoken) error                   // kann weniger schief gehen
	GetVotertokensByElectionid(electionid int) ([]dto.Wahltoken, error) // aller Token zu einer Wahl
	CheckVotertokenNotYetVoted(dto.Wahltoken) (bool, error)             // prüfen ob: Token existiert und ob valid (schon gewählt)
	InvalidateVotertoken(db *gorm.DB, token dto.Wahltoken) error        // valid auf falsch setzen
	CheckVotertokenExists(token dto.Wahltoken) (bool, error)            // prüfen ob Token existiert
}

type WahltokenLoader struct{}

func WahltokenLoaderFactory() ILoadWahltokens {
	var db_loader ILoadWahltokens = WahltokenLoader{}
	var mock_loader ILoadWahltokens = mock.MockWahltokenLoader{}
	return generalizedLoaderFactory(db_loader, mock_loader)
}

func getNextValue_seq_wahltoken() (int, error) {
	db, err := databaseconn.GetDB()
	if err != nil {
		log.Printf("error getting Database in getNextValue_seq_wahltoken: %v", err.Error())
		return -1, err
	}
	var nextID int
	err = db.Raw("SELECT NEXT VALUE FOR seq_wahltoken;").Scan(&nextID).Error
	if err != nil {
		log.Printf("failed to get next ID from sequence: %v", err)
		return -1, err
	}
	return nextID, nil
}

func (this WahltokenLoader) InsertVotertokens(tokens []dto.Wahltoken) error {
	db, err := databaseconn.GetDB()
	if err != nil {
		log.Printf("Fehler in loader Wahltoken InsertWahltokens (db Pointer): %v", err.Error())
		return err
	}

	tx := db.Begin()
	for _, token := range tokens {
		token.ID, err = getNextValue_seq_wahltoken()

		if err == nil {
			err = tx.Create(&token).Error
		}

		if err != nil {
			tx.Rollback()
			log.Printf("Fehler in loader Wahltoken InsertWahltokens: %v", err.Error())
			return err
		}
	}
	return tx.Commit().Error
}

func (this WahltokenLoader) InsertSingleVotertoken(token dto.Wahltoken) error {
	db, err := databaseconn.GetDB()
	if err != nil {
		log.Printf("Fehler in loader Wahltoken InsertSingleWahltoken (db Pointer): %v", err.Error())
		return err
	}

	token.ID, err = getNextValue_seq_wahltoken()
	if err != nil {
		return err
	}

	if err := db.Create(&token).Error; err != nil {
		log.Printf("Fehler in loader Wahltoken InsertSingleWahltoken: %v", err.Error())
		return err
	}
	return nil
}

func (this WahltokenLoader) GetVotertokensByElectionid(electionid int) ([]dto.Wahltoken, error) {
	db, err := databaseconn.GetDB()
	if err != nil {
		log.Printf("Fehler in loader Wahltoken GetWahltokensByWahlid (db Pointer): %v", err.Error())
		return nil, err
	}

	var tokens []dto.Wahltoken
	res := db.Where("election_id = ?", electionid).Find(&tokens)

	if res.Error != nil {
		log.Printf("Fehler in loader Wahltoken GetWahltokensByWahlid: %v", res.Error.Error())
		return nil, res.Error
	}

	return tokens, nil
}

func (this WahltokenLoader) CheckVotertokenNotYetVoted(token dto.Wahltoken) (bool, error) {
	// Überprüfen ob ein Wahltoken gültig ist (noch nicht gewählt), um Informationen zur Wahl zurückzugeben.
	// Für Abgabe eines Votes die Funktion CheckTokenValidWithLockingOption direkt benutzten
	db, err := databaseconn.GetDB()
	if err != nil {
		log.Printf("Fehler in loader Wahltoken CheckTokenValid (db Pointer): %v", err.Error())
		return false, err
	}

	useLock := false
	isValidToken, err := this.CheckTokenValidWithLockingOption(db, token, useLock)
	if err != nil {
		log.Printf("error: %v", err)
		return false, err
	}
	return isValidToken, nil
}

func (this WahltokenLoader) CheckTokenValidWithLockingOption(db *gorm.DB, token dto.Wahltoken, useLock bool) (bool, error) {
	// Überprüfen ob ein Wahltoken gültig ist
	// Für Abgabe eines Votes useLock=true und db ist in einer Transaktion.
	// Dies lockt die Token-Reihe in der Datenbank, bis die Transaktion beendet wird.
	// Um nur Gültigkeit zu prüfen useLock=false und db eine normale db-connection
	// (wie berteis in CheckTokenValidWithLockingOption getan wird)

	query := db.Model(&dto.Wahltoken{}).
		Where("token = ? AND election_id = ? AND voted = ?", token.Token, token.ElectionID, false)

	if useLock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var count int64
	res := query.Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	return count > 0, nil
}

func (this WahltokenLoader) InvalidateVotertoken(transaction *gorm.DB, token dto.Wahltoken) error {
	// setze voted auf true in Datenbank für Token.
	// db muss in einer Transaktion sein, um eine Racecondition zu verhindern

	res := transaction.Model(&dto.Wahltoken{}).Where("token = ? AND election_id = ?", token.Token, token.ElectionID).Update("voted", true)

	if res.Error != nil {
		log.Printf("Fehler in loader Wahltoken InvalidateVotertoken: %v", res.Error.Error())
		return res.Error
	}

	return nil
}

func (this WahltokenLoader) CheckVotertokenExists(token dto.Wahltoken) (bool, error) {
	// Überprüfen ob Token existiert
	db, err := databaseconn.GetDB()
	if err != nil {
		log.Printf("Fehler in loader Wahltoken CheckVotertokenExists (db Pointer): %v", err.Error())
		return false, err
	}

	var count int64
	res := db.Model(&dto.Wahltoken{}).Where("token = ? AND election_id = ?", token.Token, token.ElectionID).Count(&count)

	if res.Error != nil {
		log.Printf("Fehler in loader Wahltoken CheckVotertokenExists: %v", res.Error.Error())
		return false, res.Error
	}

	return count > 0, nil
}
