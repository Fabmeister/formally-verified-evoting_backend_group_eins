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
