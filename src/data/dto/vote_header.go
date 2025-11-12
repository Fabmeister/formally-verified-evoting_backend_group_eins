package dto

import ()

// VoteHeader represents the minimal fields for a vote header in a voting system.
type VoteHeader struct {
	Id    int 
	ElectionId int
}

func (VoteHeader) TableName() string {
	return "vote_header"
}

