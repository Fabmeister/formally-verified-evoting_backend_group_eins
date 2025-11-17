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
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Token [16]byte // Für Wahltoken und Bearer Token, 128-Bit

type TokenTypes int

const (
	UNDEFINED_TOKEN TokenTypes = iota // Unbekannt/Default falls nix
	TOKEN_BEARER                      // Wahlleiter Auth Token
	TOKEN_WAHLTOKEN                   // Wahltoken für Wähler Format: <wahl_id>-<16 Byte Hex>
)

func DecodeToken(tokenstring string) (Token, error) {
	decoded, err := hex.DecodeString(tokenstring)
	if err != nil {
		return Token{}, err // ungültiges Token
	}

	if len(decoded) != 16 {
		return Token{}, &HexLengthError{len(decoded), 16, "Falsche Länge für Decode (Token)"}
	}

	var token Token
	copy(token[:], decoded)
	return token, nil
}

func GenerateRandom128BitToken() Token {
	// generate random 16 bytes
	var random = [16]byte{}
	rand.Read(random[:])
	return random
}

func GenerateMultipleDistinct128BitTokens(amountOfTokens int) []Token {
	// generate amountOfTokens many Tokens (16bytes each) with no duplicate Tokens
	// amountOfTokens must be >0

	tokens := make([]Token, amountOfTokens)
	tokensMap := make(map[Token]bool, amountOfTokens)
	i := 0
	for i < amountOfTokens {
		tokenToCheckUniqueness := GenerateRandom128BitToken()
		// Check that this token wasn't already generated
		_, okay := tokensMap[tokenToCheckUniqueness]
		if !okay {
			// Token doesn't yet exist
			tokensMap[tokenToCheckUniqueness] = true
			tokens[i] = tokenToCheckUniqueness
			i = i + 1
		}
		// Else Try with another Token
	}
	return tokens
}

func GenerateUniqueVoterTokens(wahlId int32, amountOfTokens int) (full_tokens []string, raw_stringtokens []string, err error) {
	// generate amountOfTokens many Votertoken:
	// generates raw_token: a 128 Bit UniueId
	// full_tokens: each of form: wahlid-rawToken
	// raw_token: each of form: 128 bit random string

	if amountOfTokens <= 0 {
		err := errors.New("GenerateUniqueVoterTokens called with amountOfTokens<=0")
		return nil, nil, err
	}
	raw_tokenBytes := GenerateMultipleDistinct128BitTokens(amountOfTokens)

	var fulltoken_prefix = strconv.Itoa(int(wahlId)) + "-"
	raw_stringtokens = make([]string, amountOfTokens)
	full_tokens = make([]string, amountOfTokens)
	for i := 0; i < amountOfTokens; i++ {
		raw_stringtoken := hex.EncodeToString(raw_tokenBytes[i][:])
		raw_stringtokens[i] = raw_stringtoken
		full_tokens[i] = fulltoken_prefix + raw_stringtoken
	}
	return full_tokens, raw_stringtokens, nil
}

func ParseVoterTokenString(token string) (election_id int, raw_token string, err error) {
	// parse a Votertoken (form electionid-rawtoken) into it's two parts
	splited := strings.Split(token, "-")

	if len(splited) != 2 {
		return -1, "", &HexLengthError{len(splited), 2, "Falsches Format für Wahltoken"} // eigentlich kein "HexError", aber Parameter dazu passen trotzdem
	}
	if len(splited[1]) != 32 { // 32 zeichen weil 16 Byte
		return -1, "", &HexLengthError{len(splited[1]), 32, "Falsche Länge für Hexteil von Wahltoken"}
	}

	election_id, err = strconv.Atoi(splited[0])

	if err != nil {
		return -1, "", err
	}

	return election_id, splited[1], nil
}

// ----------------------------------------------------
// Bearer Token verwaltung
// ----------------------------------------------------
type BearerTokenHolder struct {
	Token    Token
	Id       int // Wahlleiter id
	EndValid time.Time
}

var (
	token_list []BearerTokenHolder // nur mit Mutex!
	mu         sync.Mutex
)

// Absichtlich ned als String, sollte eine höhrere Schicht das Error handling machen
// Return: Wahlleiter id falls valid, sonst -1
func GetWahlleiteridFromBearerToken(token Token) int {
	mu.Lock()
	defer mu.Unlock()

	// Prüfen ob Token existiert und gültig ist
	var zw_token *BearerTokenHolder = nil
	for _, k := range token_list {
		if k.Token == token {
			zw_token = &k
			break
		}
	}

	if zw_token == nil {
		return -1
	}

	if time.Now().Before(zw_token.EndValid) {
		return zw_token.Id
	}
	return -1
}

func InsertBearerToken(token Token, id int) bool {
	// Prüfen ob Token bereits existiert
	if GetWahlleiteridFromBearerToken(token) != -1 {
		// Unwahrscheinlicher Fall, dass Token bereits existiert
		// Fall, dass Token bereits existiert, aber abgelaufen ist, passiert hoffentlich nicht
		return false
	}

	mu.Lock() // WICHTIG, MUTEX ERST NACH CHECK SPERREN, SONST DEADLOCK
	defer mu.Unlock()
	zw_token := BearerTokenHolder{
		Token:    token,
		Id:       id,
		EndValid: time.Now().Add(24 * time.Hour), // Token ist 24 Stunden gültig
	}
	token_list = append(token_list, zw_token)
	return true
}

func GenerateAndInsertAuthBearerToken(id int) string {
	token := GenerateRandom128BitToken()
	ret := hex.EncodeToString(token[:])

	if !InsertBearerToken(token, id) {
		// Falls Token schon da, einfach neuen machen
		ret = GenerateAndInsertAuthBearerToken(id)
	}

	return ret
}

// Sollte irggenwie von dem extra thread hin und wieder aufgerufen werden
func CleanBearerTokens() {
	mu.Lock()
	defer mu.Unlock()

	// i --> 0
	for i := len(token_list) - 1; i >= 0; i-- {
		if time.Now().After(token_list[i].EndValid) {
			// Token ist abgelaufen, entfernen
			// Entfernen: mit letzten Element tauschen und dann das letzte weglassen
			// Hoffe das ist schneller als das mit append und slices
			zw_token := token_list[i]
			token_list[i] = token_list[len(token_list)-1]
			token_list[len(token_list)-1] = zw_token
			token_list = token_list[:len(token_list)-1]
		}
	}
}

func RefreshToken(token Token) bool {
	mu.Lock()
	defer mu.Unlock()

	// Prüfen ob Token existiert
	var zw_token *BearerTokenHolder = nil
	for i := range token_list {
		if token_list[i].Token == token {
			zw_token = &token_list[i]
			break
		}
	}

	if zw_token == nil {
		return false // Token nicht gefunden
	}

	// Token verlängern
	zw_token.EndValid = time.Now().Add(24 * time.Hour)
	return true
}

func RemoveBearerToken(token Token) bool {
	mu.Lock()
	defer mu.Unlock()

	// Prüfen ob Token existiert
	var found bool = false
	for i, k := range token_list {
		if k.Token == token {
			found = true
			token_list[i] = token_list[len(token_list)-1]
			token_list = token_list[:len(token_list)-1]
			break
		}
	}

	return found // Token wurde gefunden und entfernt
}
