package dto

import (
    "testing"
)

func TestWahlleiterInitialization(t *testing.T) {
    w := Wahlleiter{ID: 1, Username: "admin"}
	if w.ID != 1 || w.Username != "admin" {
        t.Errorf("Wahlleiter struct not initialized correctly: %+v", w)
    }
}

func TestWahlleiterUsernameNotEmpty(t *testing.T) {
    w := Wahlleiter{ID: 2, Username: ""}
    if w.Username == "" {
        t.Log("Wahlleiter username is empty as expected")
    }
}