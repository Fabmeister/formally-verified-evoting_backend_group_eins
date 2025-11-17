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
	"e-voting-service/data/dto"
	"e-voting-service/data/loading"
	authservices "e-voting-service/logic/auth_services"
	"errors"
	"log"
)

type PasswordRejectedError struct{}

func (PasswordRejectedError) Error() string {
	return "Falsches Passwort!"
}

func CheckAnmeldung(username string, password string) (string, error) {

	// Passenden user aus db laden
	var wahl_load loading.ILoadWahlleiter = loading.WahlleiterLoaderFactory()

	wahlleiter, err := wahl_load.GetWahlleiterByUsername(username)

	if err != nil {
		return "", err
	}

	salt, err := authservices.DecodeSalt(wahlleiter.Salt)

	if err != nil {
		return "", err
	}

	srv_pwd, err := authservices.DecodePassword(wahlleiter.Password)
	if err != nil {
		return "", err
	}

	// eingebenes Passwort hashen (pepper beachten) und überprüfen
	success := authservices.ComparePasswords(password, salt, srv_pwd)

	// Auth bearer token generieren und rückgeben
	if success {
		return authservices.GenerateAndInsertAuthBearerToken(wahlleiter.ID), nil
	}

	return "", PasswordRejectedError{}
}

func RegisterWahlleiter(user dto.Wahlleiter) (string, error) {
	var wahl_load loading.ILoadWahlleiter = loading.WahlleiterLoaderFactory()

	// Prüfen ob User schon vorhanden
	count, err := wahl_load.GetWahlleiterCountByName(user.Username)
	if err != nil {
		return "", err
	}

	if count != 0 {
		// User existiert bereits
		return "", errors.New("wahlleiter mit diesem namen existiert bereits")
	}

	// Passwort hashen
	hashed_pwd, salt := authservices.HashForStorage(user.Password)

	// dto bauen und in db ballern
	user.Password = hashed_pwd
	user.Salt = salt
	user.ID = -1 // wird von db gesetzt
	err = wahl_load.InsertWahlleiter(&user)
	if err != nil {
		return "", err
	}
	return authservices.GenerateAndInsertAuthBearerToken(user.ID), nil
}

// returned wahlleiterid zum bearertoken
// kann auch zum Authstatus prüfen benutzt werden
// falls bearertoken ungültig ist, dann id = -1
func TokenToWahlleiter(bearertoken string) (int, error) {
	tok_str, err := authservices.DecodeToken(bearertoken)
	if err != nil {
		return -1, err
	}

	return authservices.GetWahlleiteridFromBearerToken(tok_str), nil
}

func Abmelden(bearertoken string) error {
	tok_str, err := authservices.DecodeToken(bearertoken)
	if err != nil {
		return err
	}

	if !authservices.RemoveBearerToken(tok_str) {
		// Token nicht gefunden, nicht umbedingt ein Fehler, aber sollte geloged werden
		log.Printf("Abmelden: Token %s nicht gefunden", bearertoken)
	}

	return nil
}
