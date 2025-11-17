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
	"strconv"
	"testing"
	"time"

	"e-voting-service/data/configuration"
	dto "e-voting-service/data/dto"
	"e-voting-service/data/loading"
	authservices "e-voting-service/logic/auth_services"
)

// Tests for GetErgebnis function
//
// Note: Some tests that require token validation may fail due to database connection issues
// in the mock environment. The tests that work properly are:
// - TestGetErgebnis_InvalidWahltokenFormat (✓)
// - TestGetErgebnis_WahltokenWrongElectionId (✓)
// - TestGetErgebnis_InvalidBearerToken (✓)
// - TestGetErgebnis_UnknownTokenType (✓)
// - TestGetErgebnis_SimpleConfigurationTest (✓)
//
// Tests that have database connection issues:
// - TestGetErgebnis_ValidWahltoken (needs DB connection fix)
// - TestGetErgebnis_NonExistentWahltoken (needs DB connection fix)

// Setup function for GetErgebnis tests
func setupGetErgebnisMockEnvironment() {
	configuration.GlobalConfig.Use_mock_data = true
}

func TestGetErgebnis_ValidWahltoken(t *testing.T) {
	setupGetErgebnisMockEnvironment()

	// Setup: Create wahlleiter using the exact same pattern as the working test
	wahlleiter := dto.Wahlleiter{
		Username: "testwahlleiter",
		Password: "password",
		Email:    "test@example.com",
	}

	wahlleiterLoader := loading.WahlleiterLoaderFactory()
	err := wahlleiterLoader.InsertWahlleiter(&wahlleiter)
	if err != nil {
		t.Fatalf("Failed to insert wahlleiter: %v", err)
	}

	// Setup: Create election with candidates
	election := dto.Election{
		Id:            1,
		Name:          "Test Election",
		Beschreibung:  "Test Description",
		Wahlleiter_id: 1,
		Created_at:    time.Now(),
		End_time:      time.Now().Add(-24 * time.Hour), // election must have ended for that
		Type:          dto.APPROVAL_VOTING,
		Is_active:     false, // could be true and would be set to false by program (as endtime is in past)
		Open_wahl:     false,
		Password:      "testpass",
	}

	wahlLoader := loading.WahlLoaderFactory()
	candidates := []string{"Alice", "Bob", "Charlie"}
	err = wahlLoader.InsertElection(&election, candidates)
	if err != nil {
		t.Fatalf("Failed to insert election: %v", err)
	}

	// Setup: Create valid wahltoken exactly like the working test
	token := dto.Wahltoken{
		ElectionID: election.Id,
		Token:      "abcdefabcdefabcdefabcdefabcdefab", // Exakt 32 hex Zeichen
		Voted:      false,
	}

	wahltokenLoader := loading.WahltokenLoaderFactory()
	err = wahltokenLoader.InsertSingleVotertoken(token)
	if err != nil {
		t.Fatalf("Failed to insert wahltoken: %v", err)
	}

	// Debug: Verify token was inserted correctly
	exists, err := wahltokenLoader.CheckVotertokenExists(token)
	if err != nil {
		t.Fatalf("Error checking token existence: %v", err)
	}
	if !exists {
		t.Fatalf("Token was not properly inserted into mock")
	}

	// Setup: Add some votes for testing results using the proper method
	candidateIds1 := []dto.UnifiedVote{{CandidateID: 1, WahlInfo: 1}, {CandidateID: 2, WahlInfo: 1}} // Vote for candidates 1 and 2
	candidateIds2 := []dto.UnifiedVote{{CandidateID: 1, WahlInfo: 1}}                                // Vote for candidate 1 only

	wahlLoader.InsertVotesForOpenElection(candidateIds1, election)
	wahlLoader.InsertVotesForOpenElection(candidateIds2, election)

	// Test: GetErgebnis with valid wahltoken
	tokenString := strconv.Itoa(token.ElectionID) + "-" + token.Token
	candidates_result, winners, metaInfo, err := GetResultUsecase(election.Id, tokenString, authservices.TOKEN_WAHLTOKEN)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		t.Logf("Debug info - Token: %+v", token)
		t.Logf("Debug info - Election ID: %d", election.Id)
		return
	}

	if candidates_result == nil {
		t.Error("Expected candidates to not be nil")
	}

	if winners == nil {
		t.Error("Expected winners to not be nil")
	}

	if metaInfo == nil {
		t.Error("Expected metaInfo to not be nil")
	}

	// Verify we have the expected number of candidates
	if len(candidates_result) != 3 {
		t.Errorf("Expected 3 candidates, got %d", len(candidates_result))
	}

	// Verify winners are not empty (at least one winner expected)
	if len(winners) == 0 {
		t.Error("Expected at least one winner")
	}

	t.Logf("Test successful: %d candidates, %d winners", len(candidates_result), len(winners))
}

func TestGetErgebnis_InvalidWahltokenFormat(t *testing.T) {
	setupGetErgebnisMockEnvironment()

	// Test with invalid token format
	_, _, _, err := GetResultUsecase(1, "invalid-format", authservices.TOKEN_WAHLTOKEN)

	if err == nil {
		t.Error("Expected error for invalid token format")
	}
}

func TestGetErgebnis_WahltokenWrongElectionId(t *testing.T) {
	setupGetErgebnisMockEnvironment()

	// Test with mismatched election ID
	tokenString := "2-abcdefabcdefabcdefabcdefabcdefab" // Election ID 2, but we test with ID 1
	_, _, _, err := GetResultUsecase(1, tokenString, authservices.TOKEN_WAHLTOKEN)

	if err == nil {
		t.Error("Expected error for mismatched election ID")
	}

	// Check error type
	if wahltokenErr, ok := err.(dto.WahltokenNotValidError); ok {
		if wahltokenErr.Message != "Wahltoken != election id" {
			t.Errorf("Expected specific error message, got: %s", wahltokenErr.Message)
		}
	} else {
		t.Error("Expected WahltokenNotValidError")
	}
}

func TestGetErgebnis_NonExistentWahltoken(t *testing.T) {
	setupGetErgebnisMockEnvironment()

	// Setup minimal election structure
	wahlleiter := dto.Wahlleiter{
		Username: "testwahlleiter2",
		Password: "password",
		Email:    "test2@example.com",
	}

	wahlleiterLoader := loading.WahlleiterLoaderFactory()
	err := wahlleiterLoader.InsertWahlleiter(&wahlleiter)
	if err != nil {
		t.Fatalf("Failed to insert wahlleiter: %v", err)
	}

	election := dto.Election{
		Id:            2,
		Name:          "Test Election 2",
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

	// Test with non-existent token
	tokenString := "2-ffffffffffffffffffffffffffffffff" // Token doesn't exist
	_, _, _, err = GetResultUsecase(2, tokenString, authservices.TOKEN_WAHLTOKEN)

	if err == nil {
		t.Error("Expected error for non-existent token")
	}

	// Check error type
	if wahltokenErr, ok := err.(dto.WahltokenNotValidError); ok {
		if wahltokenErr.Message != "Wahltoken not valid" {
			t.Errorf("Expected 'Wahltoken not valid', got: %s", wahltokenErr.Message)
		}
	} else {
		t.Error("Expected WahltokenNotValidError")
	}
}

func TestGetErgebnis_InvalidBearerToken(t *testing.T) {
	setupGetErgebnisMockEnvironment()

	// Test with invalid bearer token (wrong format)
	_, _, _, err := GetResultUsecase(1, "invalidtoken", authservices.TOKEN_BEARER)

	if err == nil {
		t.Error("Expected error for invalid bearer token")
	}

	// Check error type
	if wahltokenErr, ok := err.(dto.WahltokenNotValidError); ok {
		if wahltokenErr.Message == "" {
			t.Error("Expected error message to not be empty")
		}
	} else {
		t.Error("Expected WahltokenNotValidError")
	}
}

func TestGetErgebnis_UnknownTokenType(t *testing.T) {
	setupGetErgebnisMockEnvironment()

	// Test with unknown token type
	_, _, _, err := GetResultUsecase(1, "sometoken", authservices.UNDEFINED_TOKEN)

	if err == nil {
		t.Error("Expected error for unknown token type")
	}

	// Check error type and message
	if wahltokenErr, ok := err.(dto.WahltokenNotValidError); ok {
		if wahltokenErr.Message != "Unknown token type" {
			t.Errorf("Expected 'Unknown token type', got: %s", wahltokenErr.Message)
		}
	} else {
		t.Error("Expected WahltokenNotValidError")
	}
}

func TestGetErgebnis_SimpleConfigurationTest(t *testing.T) {
	setupGetErgebnisMockEnvironment()

	// Test that configuration is set up properly for mock mode
	if !configuration.GlobalConfig.Use_mock_data {
		t.Error("Expected mock data to be enabled")
	}

	// Simple test: verify that the function returns appropriate errors for bad input
	// without needing complex setup

	// Test 1: Invalid token format should work regardless of mock setup
	_, _, _, err := GetResultUsecase(1, "bad-format", authservices.TOKEN_WAHLTOKEN)
	if err == nil {
		t.Error("Expected error for bad token format")
	}

	// Test 2: Wrong election ID format should work
	_, _, _, err = GetResultUsecase(1, "2-abcdefabcdefabcdefabcdefabcdefab", authservices.TOKEN_WAHLTOKEN)
	if err == nil {
		t.Error("Expected error for wrong election ID")
	}

	t.Log("Basic configuration and error handling tests passed")
}
