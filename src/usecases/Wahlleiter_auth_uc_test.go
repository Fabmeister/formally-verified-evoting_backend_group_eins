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
