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
	"time"

	dto "e-voting-service/data/dto"
	loading "e-voting-service/data/loading"
	authservices "e-voting-service/logic/auth_services"
	"e-voting-service/usecases/email"
)

func CreateElection_Usecase(election dto.Election, candidateNames []string, voterEmails []string) ([]string, error) {
	// erwartet folgende Felder von election korrekt befüllt: Name, Beschreibung, Wahlleiter_id, End_time, Type, open_wahl
	// and already checked endtime not in past

	election.Created_at = time.Now()
	// Da das geschickte Enddatum beim receiven des grpc-requests auf gültigkeit (datum nicht in Vergangenheit) gepüft wird
	election.Is_active = true

	// Passwort hashen
	hashed_pwd, salt := authservices.HashForStorage(election.Password)
	election.Password = hashed_pwd
	election.Salt = salt

	// Wahl in Datenbank einfügen
	var loader loading.ILoadWahl = loading.WahlLoaderFactory()
	err := loader.InsertElection(&election, candidateNames)
	if err != nil {
		return nil, err
	}

	// erstellen der/des Wahltoken(s)
	var neededVoterTokenAmount int
	if election.Open_wahl {
		neededVoterTokenAmount = 1
	} else {
		neededVoterTokenAmount = len(voterEmails)
	}
	tokens, raw_stringtokens, err := authservices.GenerateUniqueVoterTokens(int32(election.Id), neededVoterTokenAmount)
	if err != nil {
		return nil, err
	}
	db_tokens := make([]dto.Wahltoken, neededVoterTokenAmount)
	for i, raw_stringtoken := range raw_stringtokens {
		db_tokens[i] = dto.Wahltoken{ID: -1, ElectionID: election.Id, Token: raw_stringtoken, Voted: false}
	}

	// Insert Votertoken(s) into Database
	var loaderTokens loading.ILoadWahltokens = loading.WahltokenLoaderFactory()
	err = loaderTokens.InsertVotertokens(db_tokens)
	if err != nil {
		return nil, err
	}

	// Wähler benachrichtigen
	// TODO Test: len(voterEmails) == len(tokens)
	if !election.Open_wahl {
		// logged mur, Wahl wurde erstellt
		loadAndSendNotificationMails(voterEmails, tokens)
	}

	// Erweiterung, falls Frontend das will: Rückgabe Zuordnung tokens <-> voterEmails
	// (dann kann man bei mails auch tupel (mail, token) übergeben)
	return tokens, nil
}

func loadAndSendNotificationMails(voterEmails []string, tokens []string) {
	// Precondition: len(voterEmails) == len(tokens)
	for i := range voterEmails {
		mail := &email.VoteInviteMail{}
		mail.SetMailAddresses([]string{voterEmails[i]})
		mail.Token = tokens[i]

		email.SendMail(mail)
	}
}

func GetElectionsOfWahlleiter_Usecase(wahlleiterid int) ([]dto.Election, error) {
	// Get All Elections where wahlleiter is the one specified by wahlleiterid

	//Could in theory directly call the database function out of the API layer, skipping the Usecase Layer

	var loader loading.ILoadWahl = loading.WahlLoaderFactory()
	elections, err := loader.GetElectionsOfWahlleiter(wahlleiterid)
	if err != nil {
		return nil, err
	}

	return elections, nil
}
