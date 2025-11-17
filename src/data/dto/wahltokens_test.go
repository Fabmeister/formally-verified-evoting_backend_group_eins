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

func TestWahltokenInitialization(t *testing.T) {
	token := Wahltoken{ID: 1, ElectionID: 2, Token: "abc", Voted: false}
	if token.ID != 1 || token.ElectionID != 2 || token.Token != "abc" || token.Voted != false {
		t.Errorf("Wahltoken struct not initialized correctly: %+v", token)
	}
}

func TestWahltokenVotedFlag(t *testing.T) {
	token := Wahltoken{Voted: false}
	if token.Voted {
		t.Error("Token should not be marked as voted")
	}
	token.Voted = true
	if !token.Voted {
		t.Error("Token should be marked as voted")
	}
}
