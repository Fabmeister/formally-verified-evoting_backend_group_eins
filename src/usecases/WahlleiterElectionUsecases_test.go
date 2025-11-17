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
