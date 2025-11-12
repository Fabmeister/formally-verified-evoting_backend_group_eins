package dto

import (
	"time"
)

type ElectionType int

const (
	ELECTION_TYPE_UNSPECIFIED ElectionType = iota
	APPROVAL_VOTING
	COMBINED_APPROVAL_VOTING
	SCORE_VOTING
	IRV
	MAJORITY
)

func (e ElectionType) String() string {
	switch e {
	case ELECTION_TYPE_UNSPECIFIED:
		return "Unspecified"
	case APPROVAL_VOTING:
		return "Approval Voting"
	case COMBINED_APPROVAL_VOTING:
		return "Combined Approval Voting"
	case SCORE_VOTING:
		return "Score Voting"
	case IRV:
		return "Instant Runoff Voting"
	case MAJORITY:
		return "Majority Voting"
	default:
		return "UNKNOWN"
	}
}

// struct with all values present in the election table in database
type Election struct {
	Id            int `gorm:"primaryKey"`
	Name          string
	Beschreibung  string
	Wahlleiter_id int
	Created_at    time.Time
	End_time      time.Time
	Type          ElectionType
	Is_active     bool
	Password      string
	Salt          string
	Hmac          string
	Open_wahl     bool
}

// Override table name, as GORM automaticaly assumes the tablename Elections
func (Election) TableName() string {
	return "election"
}

type ElectionAlreadyEnded struct {
	Message string
}

func (e ElectionAlreadyEnded) Error() string {
	return "Election has already ended"
}

type ElectionStillActive struct {
	Message string
}

func (e ElectionStillActive) Error() string {
	return "Election is still active"
}
