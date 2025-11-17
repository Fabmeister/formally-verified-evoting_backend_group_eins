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
	"crypto/rand"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
)

type salttype [12]byte
type passwordtype [32]byte

type HexLengthError struct {
	hexLength      int
	expectedLength int
	message        string
}

func (e *HexLengthError) Error() string {
	return fmt.Sprintf("Länge: %d, Erwartet: %d - %s", e.hexLength, e.expectedLength, e.message)
}

const PEPPER_MAX = 12

func generateSalt() salttype {
	var salt [12]byte
	rand.Read(salt[:])
	return salt
}

func encodeSalt(salt salttype) string {
	return hex.EncodeToString(salt[:])
}

func genPepper() byte {
	return byte(mrand.Int31n(PEPPER_MAX))
}

func hashPassword(password string, pepper byte, salt salttype) passwordtype {
	var pwdBytes = []byte(password)
	var concat = append(pwdBytes, salt[:]...)
	concat = append(concat, pepper)
	return sha3.Sum256(concat)
}

func HashForStorage(password string) (string, string) {
	var salt = generateSalt()
	var pepper = genPepper()
	var hashed = hashPassword(password, pepper, salt)
	return encodePassword(hashed), encodeSalt(salt)
}

func encodePassword(password passwordtype) string {
	return hex.EncodeToString(password[:])
}

func DecodeSalt(salt string) (salttype, error) {
	decoded, err := hex.DecodeString(salt)
	if err != nil {
		return salttype{}, err
	}
	if len(decoded) != 12 {
		return salttype{}, &HexLengthError{len(decoded), 12, "Falsche Länge für Decode (salt)"}
	}

	var saltlen salttype
	copy(saltlen[:], decoded[0:12])
	return salttype(saltlen), nil
}

func DecodePassword(pwd string) (passwordtype, error) {
	decoded, err := hex.DecodeString(pwd)
	if err != nil {
		return passwordtype{}, err
	}
	if len(decoded) != 32 {
		return passwordtype{}, &HexLengthError{len(decoded), 32, "Falsche Länge für Decode (pwd)"}
	}

	var pwdlen passwordtype
	copy(pwdlen[:], decoded[0:32])
	return pwdlen, nil
}

func ComparePasswords(pwdInput string, salt salttype, pwdServer passwordtype) bool {
	for i := 0; i < PEPPER_MAX; i++ {
		var hashInput = hashPassword(pwdInput, byte(i), salt)
		if hashInput == pwdServer {
			return true
		}
	}
	return false
}
