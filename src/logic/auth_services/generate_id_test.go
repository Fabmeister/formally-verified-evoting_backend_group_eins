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
	"reflect"
	"testing"
	"time"
)

func TestDecodeToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    Token
		wantErr bool
	}{
		{
			name:    "Valid Token",
			args:    args{token: "00000000000000000000000000000000"},
			want:    Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "Invalid Token Length",
			args:    args{token: "000000000000000000000000000000"}, // nur 15 byte
			want:    Token{},
			wantErr: true,
		},
		{
			name:    "Invalid Hex Length",
			args:    args{token: "0000000000000000000000000000000"}, // nur 15,5 byte
			want:    Token{},
			wantErr: true,
		},
		{
			name:    "Valid Token Non-Zero",
			args:    args{token: "AA0000000000000000000000000000FF"},
			want:    Token{0xAA, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseWahlTokenString(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name            string
		args            args
		wantElection_id int
		wantRaw_token   string
		wantErr         bool
	}{
		{
			name:            "Valid Token",
			args:            args{token: "123-AA0000000000000000000000000000FF"},
			wantElection_id: 123,
			wantRaw_token:   "AA0000000000000000000000000000FF",
			wantErr:         false,
		},
		{
			name:            "Invalid Format: Token Hex length",
			args:            args{token: "123-AA0000000000000000000000000000F"}, // nur 15,5 byte
			wantElection_id: -1,
			wantRaw_token:   "",
			wantErr:         true,
		},
		{
			name:            "Invalid Format: Token length",
			args:            args{token: "123-AA0000000000000000000000000000"}, // nur 15 byte
			wantElection_id: -1,
			wantRaw_token:   "",
			wantErr:         true,
		},
		{
			name:            "Invalid Format: No Seperator",
			args:            args{token: "123AA0000000000000000000000000000"}, // nur 15 byte
			wantElection_id: -1,
			wantRaw_token:   "",
			wantErr:         true,
		},
		{
			name:            "Invalid Format: Wrong Seperator",
			args:            args{token: "123 - AA0000000000000000000000000000"}, // nur 15 byte
			wantElection_id: -1,
			wantRaw_token:   "",
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotElection_id, gotRaw_token, err := ParseVoterTokenString(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWahlTokenString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotElection_id != tt.wantElection_id {
				t.Errorf("ParseWahlTokenString() gotElection_id = %v, want %v", gotElection_id, tt.wantElection_id)
			}
			if gotRaw_token != tt.wantRaw_token {
				t.Errorf("ParseWahlTokenString() gotRaw_token = %v, want %v", gotRaw_token, tt.wantRaw_token)
			}
		})
	}
}

func TestCheckToken(t *testing.T) {

	// Token einfügen
	token_list = []BearerTokenHolder{}
	InsertBearerToken(Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1)
	InsertBearerToken(Token{0xAA, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 2)
	InsertBearerToken(Token{0xAB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 3)
	token_list[2].EndValid = token_list[2].EndValid.Add(-25 * time.Hour) // Token ablaufen lassen

	type args struct {
		token Token
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Check Valid Token",
			args: args{token: Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			want: 1,
		},
		{
			name: "Check Valid Non-Zero Token",
			args: args{token: Token{0xAA, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}},
			want: 2,
		},
		{
			name: "Check Expired Token",
			args: args{token: Token{0xAB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}},
			want: -1,
		},
		{
			name: "Check Non-Existant Token",
			args: args{token: Token{0xAB, 0, 0, 0, 0, 0, 0x34, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}},
			want: -1,
		},
		{
			name: "Check Short Token with Zero",
			args: args{token: Token{0xFF, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetWahlleiteridFromBearerToken(tt.args.token); got != tt.want {
				t.Errorf("CheckToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertToken(t *testing.T) {
	token_list = []BearerTokenHolder{}

	type args struct {
		token Token
		id    int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Insert Valid Token",
			args: args{token: Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, id: 1},
			want: true,
		},
		{
			name: "Insert Second Token",
			args: args{token: Token{0xAA, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, id: 3},
			want: true,
		},
		{
			name: "Insert Duplicate Token",
			args: args{token: Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, id: 2},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InsertBearerToken(tt.args.token, tt.args.id); got != tt.want {
				t.Errorf("InsertToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanTokens(t *testing.T) {

	// Token einfügen
	token_list = []BearerTokenHolder{}
	InsertBearerToken(Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1)
	InsertBearerToken(Token{0xAA, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 2)
	InsertBearerToken(Token{0xAB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 3) // abelaufen
	InsertBearerToken(Token{0xAB, 0xBB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 4)
	token_list[2].EndValid = token_list[2].EndValid.Add(-25 * time.Hour) // Token ablaufen lassen

	CleanBearerTokens()

	type args struct {
		token Token
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Check Valid Token 1",
			args: args{token: Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			want: 1,
		},
		{
			name: "Check Valid Token 2",
			args: args{token: Token{0xAA, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}},
			want: 2,
		},
		{
			name: "Check Valid Token 3",
			args: args{token: Token{0xAB, 0xBB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}},
			want: 4,
		},
		{
			name: "Check Removed Token",
			args: args{token: Token{0xAB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetWahlleiteridFromBearerToken(tt.args.token); got != tt.want {
				t.Errorf("CheckToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {

	// Token einfügen
	token_list = []BearerTokenHolder{}
	InsertBearerToken(Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1)
	InsertBearerToken(Token{0xAA, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 2)
	InsertBearerToken(Token{0xAB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 3) // abelaufen
	InsertBearerToken(Token{0xAB, 0xBB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 4)
	token_list[2].EndValid = token_list[2].EndValid.Add(-25 * time.Hour) // Token ablaufen lassen

	type args struct {
		token Token
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Refresh Valid Token",
			args: args{token: Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			want: true,
		},
		{
			name: "Refresh Expired Token",
			args: args{token: Token{0xAB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}},
			want: true, // Token wird erneuert
		},
		{
			name: "Refresh Non-Existant Token",
			args: args{token: Token{0xAB, 0, 0, 0, 0, 0, 0x34, 0, 0, 0, 0, 0, 0, 0, 0xFF}},
			want: false, // Token existiert nicht
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RefreshToken(tt.args.token); got != tt.want {
				t.Errorf("RefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}

	// Prüfen ob Token erneuert wurde (angenommen CleanTokens() funktioniert, aber anderer Test)
	t.Run("Check Refreshed Token pre clean", func(t *testing.T) {
		if got := GetWahlleiteridFromBearerToken(Token{0xAB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}); got != 3 {
			t.Errorf("CheckToken() = %v, want %v", got, 3)
		}
	})

	CleanBearerTokens()
	t.Run("Check Refreshed Token after clean", func(t *testing.T) {
		if got := GetWahlleiteridFromBearerToken(Token{0xAB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}); got != 3 {
			t.Errorf("CheckToken() = %v, want %v", got, 3)
		}
	})
}

func TestRemoveToken(t *testing.T) {

	token_list = []BearerTokenHolder{}
	InsertBearerToken(Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1)
	InsertBearerToken(Token{0xAA, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 2)
	InsertBearerToken(Token{0xAB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 3)
	InsertBearerToken(Token{0xAB, 0xBB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}, 4)

	type args struct {
		token Token
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Remove Valid Token",
			args: args{token: Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			want: true,
		},
		{
			name: "Remove Non-Existant Token",
			args: args{token: Token{0xAB, 0, 0, 0, 0, 0, 0x34, 0, 0, 0, 0, 0, 0, 0xFF}},
			want: false, // Token existiert nicht
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveBearerToken(tt.args.token); got != tt.want {
				t.Errorf("RemoveToken() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("Check Removed Token still there", func(t *testing.T) {
		if got := GetWahlleiteridFromBearerToken(Token{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}); got != -1 {
			t.Errorf("CheckToken() = %v, want %v", got, -1)
		}
	})

	t.Run("Check Some Token still there", func(t *testing.T) {
		if got := GetWahlleiteridFromBearerToken(Token{0xAB, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF}); got != 3 {
			t.Errorf("CheckToken() = %v, want %v", got, 3)
		}
	})
}
