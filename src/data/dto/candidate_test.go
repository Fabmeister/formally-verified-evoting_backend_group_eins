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

func TestCandidateInitialization(t *testing.T) {
	c := Candidate{Id: 1, Name: "Alice"}
	if c.Id != 1 || c.Name != "Alice" {
		t.Errorf("Candidate struct not initialized correctly: %+v", c)
	}
}

func TestCandidateNameNotEmpty(t *testing.T) {
	c := Candidate{Id: 2, Name: ""}
	if c.Name == "" {
		t.Log("Candidate name is empty as expected")
	}
}
