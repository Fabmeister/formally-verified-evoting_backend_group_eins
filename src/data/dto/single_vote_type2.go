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

/*
vote_id int,
candidate_id int,
info int,
primary key(vote_id, candidate_id),
constraint fk_single_vote_type2_vote foreign key (vote_id) references vote_header(id),
constraint fk_single_vote_type2_candidate foreign key (candidate_id) references candidates(id)
*/
package dto

/*
Datenbankobjekt für die Tabelle "single_vote_type2"
Nicht verwechseln mit unified_vote! Ähnliche Felder aber dieser hier nur Datenbank!
*/
type Vote_Type2 struct {
	Vote_id      int `gorm:"primaryKey"`
	Candidate_id int
	Info         int
}

func (Vote_Type2) TableName() string {
	return "single_vote_type2"
}
