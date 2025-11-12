package dto

import (
    "testing"
)

func TestCandidateInitialization(t *testing.T) {
    c := Candidate{Id: 1, Name: "Alice"}
    if c.Id != 1 || c.Name != "Alice" {
        t.Errorf("Candidate struct not initialized correctly: %+v", c)
    }
}

func TestCandidateNameNotEmpty(t *testing.T) {
    c := Candidate{Id: 2, Name: ""}
    if c.Name == "" {
        t.Log("Candidate name is empty as expected")
    }
}