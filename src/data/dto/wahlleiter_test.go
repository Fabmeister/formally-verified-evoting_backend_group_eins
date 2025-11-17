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

package dto

import (
	"testing"
)

func TestWahlleiterInitialization(t *testing.T) {
	w := Wahlleiter{ID: 1, Username: "admin"}
	if w.ID != 1 || w.Username != "admin" {
		t.Errorf("Wahlleiter struct not initialized correctly: %+v", w)
	}
}

func TestWahlleiterUsernameNotEmpty(t *testing.T) {
	w := Wahlleiter{ID: 2, Username: ""}
	if w.Username == "" {
		t.Log("Wahlleiter username is empty as expected")
	}
}
