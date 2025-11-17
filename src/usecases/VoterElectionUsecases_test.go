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
	"testing"
	"time"

	"e-voting-service/data/configuration"
	dto "e-voting-service/data/dto"
	"e-voting-service/data/loading"
)

// Setup function for voter election tests
func setupVoterElectionMockEnvironment() {
	configuration.GlobalConfig.Use_mock_data = true
}

// Test für GetElectionForVoter_usecase mit Mock-Loader
func TestGetElectionForVoter_usecase_ValidToken(t *testing.T) {
	setupVoterElectionMockEnvironment()

	// Setup: Wahlleiter erstellen - KORREKTE Factory verwenden
	wahlleiter := dto.Wahlleiter{
		Username: "testwahlleiter",
		Password: "password",
		Email:    "test@example.com",
	}

	// WICHTIG: WahlleiterLoaderFactory für Wahlleiter-Operationen
	wahlleiterLoader := loading.WahlleiterLoaderFactory()
	err := wahlleiterLoader.InsertWahlleiter(&wahlleiter)
	if err != nil {
		t.Fatalf("Failed to insert wahlleiter: %v", err)
	}

	// Setup: Wahl und Kandidaten erstellen - KORREKTE Factory verwenden
	election := dto.Election{
		Id:            1,
		Name:          "Testwahl",
		Beschreibung:  "Test Description",
		Wahlleiter_id: 1,
		Created_at:    time.Now(),
		End_time:      time.Now().Add(24 * time.Hour), //important it is active!
		Type:          dto.APPROVAL_VOTING,
		Is_active:     true,
		Open_wahl:     false,
		Password:      "testpass",
	}

	// WICHTIG: WahlLoaderFactory für Wahl-Operationen
	wahlLoader := loading.WahlLoaderFactory()
	candidates := []string{"Alice", "Bob"}
	err = wahlLoader.InsertElection(&election, candidates)
	if err != nil {
		t.Fatalf("Failed to insert election: %v", err)
	}

	// Setup: Gültigen Wahltoken erstellen - KORREKTE Factory verwenden
	token := dto.Wahltoken{
		ElectionID: election.Id,
		Token:      "abcdefabcdefabcdefabcdefabcdefab", // Exakt 32 hex Zeichen
		Voted:      false,
	}

	// WICHTIG: WahltokenLoaderFactory für Token-Operationen
	wahltokenLoader := loading.WahltokenLoaderFactory()
	err = wahltokenLoader.InsertSingleVotertoken(token)
	if err != nil {
		t.Fatalf("Failed to insert wahltoken: %v", err)
	}

	// Debug: Überprüfen ob Token korrekt eingefügt wurde
	exists, err := wahltokenLoader.CheckVotertokenExists(token)
	if err != nil {
		t.Fatalf("Error checking token existence: %v", err)
	}
	if !exists {
		t.Fatalf("Token was not properly inserted into mock")
	}

	// Debug: Überprüfen ob Token als gültig erkannt wird
	isValid, err := wahltokenLoader.CheckVotertokenNotYetVoted(token)
	if err != nil {
		t.Fatalf("Error checking token validity: %v", err)
	}
	if !isValid {
		t.Fatalf("Token should be valid but is not recognized as such")
	}

	// Test: GetElectionForVoter_usecase
	gotElection, gotCandidates, err := GetElectionForVoter_usecase(1, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		t.Logf("Debug info - Token: %+v", token)
		t.Logf("Debug info - Election ID: %d", election.Id)
		return
	}

	// Validierungen
	if gotElection.Id != 1 {
		t.Errorf("Expected election ID 1, got %d", gotElection.Id)
	}

	if gotElection.Name != "Testwahl" {
		t.Errorf("Expected election name 'Testwahl', got '%s'", gotElection.Name)
	}

	if len(gotCandidates) != 2 {
		t.Errorf("Expected 2 candidates, got %d", len(gotCandidates))
	}

	t.Logf("Test successful: Election '%s' with %d candidates", gotElection.Name, len(gotCandidates))
}

// Test für GetElectionForVoter_usecase mit ungültigem Token
func TestGetElectionForVoter_usecase_InvalidToken(t *testing.T) {
	setupVoterElectionMockEnvironment()

	// Setup: Wahl erstellen (ohne Token)
	wahlleiter := dto.Wahlleiter{
		Username: "invalidtestwahlleiter",
		Password: "password",
		Email:    "invalid@example.com",
	}

	wahlleiterLoader := loading.WahlleiterLoaderFactory()
	err := wahlleiterLoader.InsertWahlleiter(&wahlleiter)
	if err != nil {
		t.Fatalf("Failed to insert wahlleiter: %v", err)
	}

	election := dto.Election{
		Id:            2,
		Name:          "Invalid Token Test",
		Wahlleiter_id: 1,
		End_time:      time.Now().Add(24 * time.Hour),
		Type:          dto.APPROVAL_VOTING,
		Is_active:     true,
		Open_wahl:     false,
	}

	wahlLoader := loading.WahlLoaderFactory()
	err = wahlLoader.InsertElection(&election, []string{"Candidate"})
	if err != nil {
		t.Fatalf("Failed to insert election: %v", err)
	}

	// Test mit ungültigem Token (nicht in Mock vorhanden)
	invalidToken := dto.Wahltoken{
		ElectionID: 2,
		Token:      "ffffffffffffffffffffffffffffffff", // Token existiert nicht im Mock
		Voted:      false,
	}

	_, _, err = GetElectionForVoter_usecase(2, invalidToken)
	if err == nil {
		t.Error("Expected error for invalid token")
	}

	expectedError := "given wahltoken not valid or already used"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}

	t.Logf("Correctly rejected invalid token: %v", err)
}

// Test für HandleVote_usecase mit offener Wahl
func TestHandleVote_usecase_OpenElection(t *testing.T) {
	setupVoterElectionMockEnvironment()

	// Setup: Wahlleiter und offene Wahl
	wahlleiter := dto.Wahlleiter{
		Username: "openwahlleiter",
		Password: "password",
		Email:    "open@example.com",
	}

	wahlleiterLoader := loading.WahlleiterLoaderFactory()
	err := wahlleiterLoader.InsertWahlleiter(&wahlleiter)
	if err != nil {
		t.Fatalf("Failed to insert wahlleiter: %v", err)
	}

	election := dto.Election{
		Id:            3,
		Name:          "Open Election Test",
		Wahlleiter_id: 1,
		End_time:      time.Now().Add(24 * time.Hour),
		Type:          dto.APPROVAL_VOTING,
		Is_active:     true,
		Open_wahl:     true, // Offene Wahl
	}

	wahlLoader := loading.WahlLoaderFactory()
	err = wahlLoader.InsertElection(&election, []string{"Open Candidate"})
	if err != nil {
		t.Fatalf("Failed to insert election: %v", err)
	}

	// Test: Vote in offener Wahl
	token := dto.Wahltoken{
		ElectionID: 3,
		Token:      "12345678901234567890123456789012", // 32 chars
		Voted:      false,
	}

	votes := []dto.UnifiedVote{{CandidateID: 1, WahlInfo: 1}}
	err = HandleVote_usecase(votes, token)
	if err != nil {
		t.Errorf("Unexpected error for open election vote: %v", err)
	}

	t.Log("Successfully handled vote in open election")
}

// Test für HandleVote_usecase mit beendeter Wahl
func TestHandleVote_usecase_ClosedElection(t *testing.T) {
	setupVoterElectionMockEnvironment()

	// Setup: Wahlleiter und beendete Wahl
	wahlleiter := dto.Wahlleiter{
		Username: "closedwahlleiter",
		Password: "password",
		Email:    "closed@example.com",
	}

	wahlleiterLoader := loading.WahlleiterLoaderFactory()
	err := wahlleiterLoader.InsertWahlleiter(&wahlleiter)
	if err != nil {
		t.Fatalf("Failed to insert wahlleiter: %v", err)
	}

	election := dto.Election{
		Id:            4,
		Name:          "Closed Election Test",
		Wahlleiter_id: 1,
		End_time:      time.Now().Add(-24 * time.Hour),
		Type:          dto.APPROVAL_VOTING,
		Is_active:     false, // Wahl beendet
		Open_wahl:     false,
	}

	wahlLoader := loading.WahlLoaderFactory()
	err = wahlLoader.InsertElection(&election, []string{"Closed Candidate"})
	if err != nil {
		t.Fatalf("Failed to insert election: %v", err)
	}

	// Test: Vote in beendeter Wahl
	token := dto.Wahltoken{
		ElectionID: 4,
		Token:      "87654321876543218765432187654321", // 32 chars
		Voted:      false,
	}

	votes := []dto.UnifiedVote{{CandidateID: 1, WahlInfo: 1}}
	err = HandleVote_usecase(votes, token)
	if err == nil {
		t.Error("Expected error for closed election")
	}

	var expectedError error = dto.ElectionAlreadyEnded{}
	if err != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}

	t.Logf("Correctly rejected vote in closed election: %v", err)
}
