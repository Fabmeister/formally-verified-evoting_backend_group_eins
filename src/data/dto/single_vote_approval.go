package dto

type Single_vote_approval struct {
	Vote_id      int `gorm:"primaryKey"`
	Candidate_id int
	Approved     bool
}

func (Single_vote_approval) TableName() string {
	return "single_vote_approval"
}
