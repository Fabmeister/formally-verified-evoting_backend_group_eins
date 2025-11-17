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
	"time"

	"e-voting-service/data/dto"
	databaseconn "e-voting-service/data/dto/connection"
	"e-voting-service/data/loading/mock"
)

type ILoadWahl interface {
	GetElectionForVoter(wahlid int) (dto.Election, []dto.Candidate, error)
	InsertElection(wahl *dto.Election, candidateNames []string) error
	InsertVoteAndInvalidateToken(votes []dto.UnifiedVote, wahltoken dto.Wahltoken) error
	InsertVotesForOpenElection(votes []dto.UnifiedVote, election dto.Election) error
	GetElection(wahlid int) (dto.Election, error)
	GetCandidates(wahlid int) ([]int, map[int]dto.Candidate, error) // Liste der ids mit map für dto
	GetVotesApproval(wahlid int) ([][]dto.Single_vote_approval, error)
	GetVotesType2(wahlid int) ([][]dto.Vote_Type2, error)
	IsElectionActive(wahlid int) (bool, error)
	GetElectionsOfWahlleiter(wahlleiterid int) ([]dto.Election, error)
}

func WahlLoaderFactory() ILoadWahl {
	var db_loader ILoadWahl = WahlLoader{}
	var mock_loader ILoadWahl = mock.MockWahlLoader{}
	return generalizedLoaderFactory(db_loader, mock_loader)
}

type WahlLoader struct{}

func (loader WahlLoader) InsertElection(wahl *dto.Election, candidateNames []string) error {
	// election muss alle Felder, außer Electionid gefüllt haben

	db, err := databaseconn.GetDB()
	if err != nil {
		return err
	}

	// Datenbank Transaction Anfang, dass entweder alles oder nichts eingefügt
	transaction := db.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}

	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	var nextElectionID int
	err = transaction.Raw("SELECT NEXT VALUE FOR seq_election;").Scan(&nextElectionID).Error
	if err != nil {
		transaction.Rollback()
		return err
	}

	wahl.Id = nextElectionID
	// Insert election into Database
	err = transaction.Create(wahl).Error
	if err != nil {
		log.Printf("Insertion of election failed: %v", err)
		transaction.Rollback()
		return err
	}

	// insert Candidates

	for _, candidateName := range candidateNames {
		var nextCandidateID int
		err = transaction.Raw("SELECT NEXT VALUE FOR seq_candidates;").Scan(&nextCandidateID).Error
		if err != nil {
			log.Printf("failed to get next ID from candidate sequence: %v", err)
			transaction.Rollback()
			return err
		}

		cand := dto.Candidate{Id: int(nextCandidateID), ElectionId: wahl.Id, Name: candidateName}
		err = transaction.Create(&cand).Error
		if err != nil {
			log.Printf("Insertion of candidate for election failed: %v", err)
			transaction.Rollback()
			return err
		}
	}
	return transaction.Commit().Error
}

func (loader WahlLoader) GetElectionForVoter(electionid int) (dto.Election, []dto.Candidate, error) {
	// Erwartet, dass die Authentifizierung des Wählers bereits stattgefunden hat und returnt alle Informationen,
	// die für das Erstellen des HTML-Wahlzettel-Form benötigt sind.

	db, err := databaseconn.GetDB()
	if err != nil {
		return dto.Election{}, nil, err
	}

	// Datenbank Transaction Anfang, dass entweder alles oder nichts eingefügt
	transaction := db.Begin()
	if transaction.Error != nil {
		return dto.Election{}, nil, transaction.Error
	}

	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	var election dto.Election
	err = transaction.First(&election, electionid).Error
	if err != nil {
		log.Printf("Failed to get Election %v", err)
		transaction.Rollback()
		return dto.Election{}, nil, err
	}

	var candidates []dto.Candidate
	err = transaction.Where("election_id = ?", electionid).Find(&candidates).Error
	if err != nil {
		log.Printf("Failed to get Candidates for Election %v", err)
		transaction.Rollback()
		return dto.Election{}, nil, err
	}
	return election, candidates, transaction.Commit().Error
}

func (loader WahlLoader) InsertVoteAndInvalidateToken(votes []dto.UnifiedVote, wahlToken dto.Wahltoken) error {
	// Behandelt einen Wahlzettel eines einzelnen Wählers.
	// Einfügen der einzelnen Votes in die Datenbank
	// Input:
	// 		wahltoken muss mit electionid und Token befüllt sein
	// 		votes: ein vote muss je nach Wahlsystem entsprechende WahlInfo zu Kandidaten haben.
	//      z.B:Bei Approval voting 0/1 für Enthaltung/Zustimmung

	db, err := databaseconn.GetDB()
	if err != nil {
		return err
	}
	// Datenbank Transaction Anfang, dass entweder alles oder nichts eingefügt
	transaction := db.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}

	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	// Check validity of VoteToken
	var wahltokenLoader WahltokenLoader = WahltokenLoader{}
	useLock := true
	isValid, err := wahltokenLoader.CheckTokenValidWithLockingOption(transaction, wahlToken, useLock)
	if err != nil {
		log.Printf("error in CheckTokenValidWithLockingOption: %v", err)
		transaction.Rollback()
		return err
	}
	if !isValid {
		err = &dto.WahltokenNotValidError{Message: "Given Wahltoken not valid or already used"}
		log.Printf("%v", err)
		transaction.Rollback()
		return err
	}

	// Insert Votes:
	// Create Vote Header:
	var nextVoteHeaderID int
	err = transaction.Raw("SELECT NEXT VALUE FOR seq_vote_header;").Scan(&nextVoteHeaderID).Error
	if err != nil {
		log.Printf("failed to get next ID from seq_vote header_sequence: %v", err)
		transaction.Rollback()
		return err
	}

	// Vote Header in db einfügen
	voteHeader := dto.VoteHeader{Id: nextVoteHeaderID, ElectionId: wahlToken.ElectionID}

	err = transaction.Create(&voteHeader).Error
	if err != nil {
		log.Printf("Insertion of voteHeader failed: %v", err)
		transaction.Rollback()
		return err
	}

	election, err := loader.GetElection(wahlToken.ElectionID)
	if err != nil {
		log.Printf("GetWahl failed: %v", err)
		transaction.Rollback()
		return err
	}

	// AB HIER GENERALISIERT

	single_vote_inserter, err := singleVoteInserterFactory(election)
	if err != nil {
		log.Printf("error Insert Single Factory! Wahlid (token) %d, wahlid (election) %d", wahlToken.ID, election.Id)
		transaction.Rollback()
		return err
	}
	err = single_vote_inserter.insert_votes(votes, voteHeader, transaction)

	if err != nil {
		transaction.Rollback()
		return err
	}

	// Generalisierung ende

	// Invalidate VoteToken so voter can't vote twice
	err = wahltokenLoader.InvalidateVotertoken(transaction, wahlToken)
	if err != nil {
		log.Printf("error invalidaing Wahltoken: %v", err)
		transaction.Rollback()
		return err
	}
	return transaction.Commit().Error
}

func (loader WahlLoader) InsertVotesForOpenElection(votes []dto.UnifiedVote, election dto.Election) error {
	// Behandelt einen Wahlzettel eines einzelnen Wählers.
	// Offene Wahl bedeutet, dass kein Votetoken invalidiert werden muss.
	// Einfügen der einzelnen Approval-Votes in die Datenbank
	// Input:
	// 		wahltoken muss mit electionid und Token befüllt sein
	// 		votedCandidateids sind alle ids, die vom Wähler "approved" sind

	db, err := databaseconn.GetDB()
	if err != nil {
		return err
	}
	// Datenbank Transaction Anfang, dass entweder alles oder nichts eingefügt
	transaction := db.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}

	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	// Insert Votes:
	// Create Vote Header:
	var nextVoteHeaderID int
	err = transaction.Raw("SELECT NEXT VALUE FOR seq_vote_header;").Scan(&nextVoteHeaderID).Error
	if err != nil {
		log.Printf("failed to get next ID from seq_vote header_sequence: %v", err)
		transaction.Rollback()
		return err
	}

	// Vote Header in db einfügen
	voteHeader := dto.VoteHeader{Id: nextVoteHeaderID, ElectionId: int(election.Id)}

	err = transaction.Create(&voteHeader).Error
	if err != nil {
		log.Printf("Insertion of voteHeader failed: %v", err)
		transaction.Rollback()
		return err
	}

	// AB HIER GENERALISIERT

	single_vote_inserter, err := singleVoteInserterFactory(election)
	if err != nil {
		log.Printf("error creating Insert Single Factory! wahlid (election) %d, wahltype %v", election.Id, election.Type)
		transaction.Rollback()
		return err
	}
	err = single_vote_inserter.insert_votes(votes, voteHeader, transaction)

	if err != nil {
		transaction.Rollback()
		return err
	}

	// generalisierung ende

	return transaction.Commit().Error
}

func (loader WahlLoader) GetElection(wahlid int) (dto.Election, error) {
	db, err := databaseconn.GetDB()
	if err != nil {
		log.Printf("Fehler in loader Wahltoken GetWahl (db Pointer): %v", err.Error())
		return dto.Election{}, err
	}
	var election dto.Election
	err = db.First(&election, wahlid).Error
	if err != nil {
		log.Printf("Fehler in loader Wahltoken GetWahl: %v", err.Error())
		return dto.Election{}, err
	}
	return election, nil
}

func (loader WahlLoader) GetCandidates(wahlid int) (id_list []int, cand_map map[int]dto.Candidate, err error) {
	db, err := databaseconn.GetDB()
	if err != nil {
		log.Printf("Fehler in loader Wahltoken GetCandidates (db Pointer): %v", err.Error())
		return nil, nil, err
	}

	rows, err := db.Table("candidates").Where("election_id = ?", wahlid).Rows()
	if err != nil {
		log.Printf("Fehler in loader Wahltoken GetCandidates: %v", err.Error())
		return nil, nil, err
	}
	defer rows.Close()

	cand_map = make(map[int]dto.Candidate)
	id_list = make([]int, 0)
	for rows.Next() {
		var candidate dto.Candidate
		err = db.ScanRows(rows, &candidate)
		if err != nil {
			log.Printf("Fehler in loader Wahltoken GetCandidates ScanRows: %v", err.Error())
			return nil, nil, err
		}
		cand_map[candidate.Id] = candidate
		id_list = append(id_list, candidate.Id)
	}

	return id_list, cand_map, nil
}

func (loader WahlLoader) GetVotesApproval(wahlid int) ([][]dto.Single_vote_approval, error) {
	// get Votes for Approval Voting
	db, err := databaseconn.GetDB()
	if err != nil {
		return nil, err
	}
	query := "select sva.* from vote_header vh inner join single_vote_approval sva on sva.vote_id = vh.id where vh.election_id = ? and approved != 0 order by vh.id"
	rows, err := db.Raw(query, wahlid).Rows()

	if err != nil {
		log.Printf("Fehler in loader Wahltoken GetVotesApproval: %v", err.Error())
		return nil, err
	}

	defer rows.Close()

	cur_header_id := -1
	var approvals [][]dto.Single_vote_approval
	for rows.Next() {
		var approv dto.Single_vote_approval
		err := rows.Scan(&approv.Vote_id, &approv.Candidate_id, &approv.Approved)
		if err != nil {
			log.Printf("in (WahlLoader)GetVotesType2 error scanning db row into dto.Vote_Type2: %v", err)
			return nil, err
		}

		// Sind nach voteheader sortiert, darum prüfen wann headerid anders ist
		if cur_header_id != approv.Vote_id {
			cur_header_id = approv.Vote_id
			approvals = append(approvals, []dto.Single_vote_approval{})
		}
		approvals[len(approvals)-1] = append(approvals[len(approvals)-1], approv)
	}

	return approvals, nil
}

func (loader WahlLoader) GetVotesType2(wahlid int) ([][]dto.Vote_Type2, error) {
	// get Votes for Combined Approval, Score and Instant Runoff Voting

	db, err := databaseconn.GetDB()
	if err != nil {
		return nil, err
	}
	query := "select t2.* from vote_header vh inner join single_vote_type2 t2 on t2.vote_id = vh.id where vh.election_id = ? order by vh.id"
	rows, err := db.Raw(query, wahlid).Rows()

	if err != nil {
		log.Printf("Fehler in loader Wahltoken GetVotesType2: %v", err.Error())
		return nil, err
	}

	defer rows.Close()

	cur_header_id := -1
	var votes [][]dto.Vote_Type2
	for rows.Next() {
		var vote dto.Vote_Type2
		err := rows.Scan(&vote.Vote_id, &vote.Candidate_id, &vote.Info)
		if err != nil {
			log.Printf("in (WahlLoader)GetVotesType2 error scanning db row into dto.Vote_Type2: %v", err)
			return nil, err
		}

		// Sind nach voteheader sortiert, darum prüfen wann headerid anders ist
		if cur_header_id != vote.Vote_id {
			cur_header_id = vote.Vote_id
			votes = append(votes, []dto.Vote_Type2{})
		}
		votes[len(votes)-1] = append(votes[len(votes)-1], vote)
	}

	return votes, nil
}

func (loader WahlLoader) IsElectionActive(wahlid int) (bool, error) {
	// returns field is_active from database AND updates it from false to true,
	// if the election endtime is now already in the past
	db, err := databaseconn.GetDB()
	if err != nil {
		return false, err
	}

	transaction := db.Begin()
	if transaction.Error != nil {
		return false, transaction.Error
	}

	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	// get election from db
	var election dto.Election
	err = transaction.First(&election, wahlid).Error
	if err != nil {
		log.Printf("Failed to get Election %v", err)
		transaction.Rollback()
		return false, err
	}

	// check is_active status
	if !(election.Is_active) {
		return false, nil
	} else {
		if time.Now().After(election.End_time) {
			// election should have already ended, but field not yet updated
			election.Is_active = false
			// update in database
			return false, db.Save(&election).Error
		}
		return election.Is_active, nil
	}
}

func (loader WahlLoader) GetElectionsOfWahlleiter(wahlleiterid int) ([]dto.Election, error) {
	// get all Elections the wahleiter with the specified id created

	db, err := databaseconn.GetDB()
	if err != nil {
		return nil, err
	}
	var elections []dto.Election
	err = db.Where("wahlleiter_id = ?", wahlleiterid).Find(&elections).Error
	if err != nil {
		return nil, err
	}
	return elections, nil
}
