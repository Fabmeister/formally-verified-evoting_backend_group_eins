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

import "fmt"

// Kein Datenbank Objekt!
type UnifiedVote struct {
	CandidateID int32
	WahlInfo    int32
}

type UnifiedVotePreconditionError struct {
	Type    ElectionType
	Message string
}

func (this UnifiedVotePreconditionError) Error() string {
	return fmt.Sprintf("Unified Vote Precondition Failed, type = %s, Message = %s", this.Type.String(), this.Message)
}
