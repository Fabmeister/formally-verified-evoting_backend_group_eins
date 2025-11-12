package dto

import (
    "testing"
)

func TestElectionInitialization(t *testing.T) {
    e := Election{Id: 1, Name: "Testwahl"}
    if e.Id != 1 || e.Name != "Testwahl" {
        t.Errorf("Election struct not initialized correctly: %+v", e)
    }
}

func TestElectionActiveFlag(t *testing.T) {
    e := Election{Is_active: true}
    if !e.Is_active {
        t.Error("Election should be active")
    }
}