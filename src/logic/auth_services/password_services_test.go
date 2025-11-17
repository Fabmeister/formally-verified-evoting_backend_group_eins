// Formally verified E-Voting using Dafny
// Copyright (C) 2025 Authors Gruppe EinS
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package authservices

import (
	"encoding/hex"
	"testing"
)

func Test_generateSalt_and_encodeSalt(t *testing.T) {
	salt := generateSalt()
	encoded := encodeSalt(salt)
	if len(encoded) != 24 {
		t.Errorf("encoded salt length = %d, want 24", len(encoded))
	}
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		t.Errorf("encodeSalt produced invalid hex string: %v", err)
	}
	if len(decoded) != 12 {
		t.Errorf("decoded salt length = %d, want 12", len(decoded))
	}
}

func Test_genPepper(t *testing.T) {
	for i := 0; i < 100; i++ {
		pepper := genPepper()
		if pepper >= PEPPER_MAX {
			t.Errorf("pepper >= PEPPER_MAX: %d", pepper)
		}
	}
}

func Test_hashPassword_and_encodePassword(t *testing.T) {
	salt := generateSalt()
	pepper := genPepper()
	password := "mypassword"
	hashed := hashPassword(password, pepper, salt)
	encoded := encodePassword(hashed)
	if len(encoded) != 64 {
		t.Errorf("encoded hash length = %d, want 64", len(encoded))
	}
}

func Test_HashForStorage_and_Decode(t *testing.T) {
	password := "myTestPassword123"
	hash, salt := HashForStorage(password)
	if len(hash) != 64 {
		t.Errorf("HashForStorage hash length = %d, want 64", len(hash))
	}
	if len(salt) != 24 {
		t.Errorf("HashForStorage salt length = %d, want 24", len(salt))
	}
	_, err := DecodeSalt(salt)
	if err != nil {
		t.Errorf("DecodeSalt error: %v", err)
	}
	decodedPwd, err := DecodePassword(hash)
	if err != nil {
		t.Errorf("DecodePassword error: %v", err)
	}
	// Hash should match original encoding
	if encodePassword(decodedPwd) != hash {
		t.Errorf("Decoded password does not re-encode to original hash")
	}
	// Check error for wrong salt length
	_, err = DecodeSalt("1234")
	if err == nil {
		t.Errorf("expected error for short salt, got nil")
	}
	// Check error for wrong password length
	_, err = DecodePassword("abcd")
	if err == nil {
		t.Errorf("expected error for short password hash, got nil")
	}
}

func Test_ComparePasswords(t *testing.T) {
	password := "securePassword"
	hash, saltStr := HashForStorage(password)
	salt, _ := DecodeSalt(saltStr)
	hashed, _ := DecodePassword(hash)
	// Should match with correct password
	if !ComparePasswords(password, salt, hashed) {
		t.Error("ComparePasswords failed for correct password")
	}
	// Should not match with incorrect password
	if ComparePasswords("wrongPassword", salt, hashed) {
		t.Error("ComparePasswords succeeded for wrong password")
	}
}
