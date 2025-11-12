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