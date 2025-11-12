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
