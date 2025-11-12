package usecases

import (
	"e-voting-service/data/configuration"
	dto "e-voting-service/data/dto"
	"testing"
)

func TestCreateElection_Usecase_OpenWahl(t *testing.T) {
	configuration.GlobalConfig.Use_mock_data = true

	election := dto.Election{
		Name:      "Testwahl",
		Password:  "secret",
		Open_wahl: true,
	}
	candidates := []string{"Alice", "Bob"}
	votermails := []string{}
	tokens, err := CreateElection_Usecase(election, candidates, votermails)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) == 0 {
		t.Error("expected at least one token")
	}
}
