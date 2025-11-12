package dto

// VoteHeader represents the minimal fields for a vote header in a voting system.
type Candidate struct {
	Id         int `gorm:"primaryKey"`
	Name       string
	ElectionId int `gorm:"column:election_id"`
}
