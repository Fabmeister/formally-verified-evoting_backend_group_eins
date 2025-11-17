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
	"e-voting-service/data/configuration"
	dto "e-voting-service/data/dto"
	"testing"
)

// Setup function for wahlleiter auth tests
func setupWahlleiterAuthMockEnvironment() {
	configuration.GlobalConfig.Use_mock_data = true
}

func TestPasswordRejectedError_Error(t *testing.T) {
	var err error = PasswordRejectedError{}
	if err.Error() != "Falsches Passwort!" {
		t.Errorf("unexpected error string: %s", err.Error())
	}
}

func TestCheckAnmeldung_InvalidUser(t *testing.T) {
	setupWahlleiterAuthMockEnvironment()

	_, err := CheckAnmeldung("invalid", "password")
	if err == nil {
		t.Error("expected error for invalid user")
	}
}

func TestRegisterWahlleiter_NewUser(t *testing.T) {
	setupWahlleiterAuthMockEnvironment()

	user := dto.Wahlleiter{Username: "newuser", Password: "testpassword", Email: "test@example.com"}
	token, err := RegisterWahlleiter(user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token for new user")
	}
}
