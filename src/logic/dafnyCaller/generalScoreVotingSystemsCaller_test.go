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

package votingSystemsCaller

import (
	"reflect"
	"testing"
)

func TestDafnyCaller_ApprovalVoting(t *testing.T) {
	type args struct {
		candidates []int
		votes      [](map[int]int)
	}
	tests := []struct {
		name                  string
		args                  args
		wantWinners           []int
		wantScorePerCandidate map[int]int
		wantErr               bool
	}{
		{
			name: "Valid approval voting with clear winner",
			args: args{
				candidates: []int{1, 2, 3},
				votes: [](map[int]int){
					{1: 1, 2: 0, 3: 0},
					{1: 1, 2: 1, 3: 0},
					{1: 0, 2: 1, 3: 1},
				},
			},
			wantWinners:           []int{1, 2},
			wantScorePerCandidate: map[int]int{1: 2, 2: 2, 3: 1},
			wantErr:               false,
		},
		{
			name: "Nil candidates",
			args: args{
				candidates: nil,
				votes:      [](map[int]int){{1: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Nil votes",
			args: args{
				candidates: []int{1, 2},
				votes:      nil,
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Empty candidates",
			args: args{
				candidates: []int{},
				votes:      [](map[int]int){{1: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Invalid score (2 for approval voting)",
			args: args{
				candidates: []int{1, 2},
				votes:      [](map[int]int){{1: 2, 2: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Vote for nonexistent candidate",
			args: args{
				candidates: []int{1, 2},
				votes:      [](map[int]int){{1: 1, 3: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Single candidate",
			args: args{
				candidates: []int{5},
				votes:      [](map[int]int){{5: 1}, {5: 0}, {5: 1}},
			},
			wantWinners:           []int{5},
			wantScorePerCandidate: map[int]int{5: 2},
			wantErr:               false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWinners, gotScorePerCandidate, err := DafnyCaller_ApprovalVoting(tt.args.candidates, tt.args.votes)
			if (err != nil) != tt.wantErr {
				t.Errorf("DafnyCaller_ApprovalVoting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(gotWinners, tt.wantWinners) {
					t.Errorf("DafnyCaller_ApprovalVoting() gotWinners = %v, want %v", gotWinners, tt.wantWinners)
				}
				if !reflect.DeepEqual(gotScorePerCandidate, tt.wantScorePerCandidate) {
					t.Errorf("DafnyCaller_ApprovalVoting() gotScorePerCandidate = %v, want %v", gotScorePerCandidate, tt.wantScorePerCandidate)
				}
			}
		})
	}
}

func TestDafnyCaller_CombinedApprovalVoting(t *testing.T) {
	type args struct {
		candidates []int
		votes      [](map[int]int)
	}
	tests := []struct {
		name                  string
		args                  args
		wantWinners           []int
		wantScorePerCandidate map[int]int
		wantErr               bool
	}{
		{
			name: "Valid combined approval voting",
			args: args{
				candidates: []int{1, 2, 3},
				votes: [](map[int]int){
					{1: 1, 2: -1, 3: 0},
					{1: 1, 2: 1, 3: -1},
					{1: -1, 2: 1, 3: 1},
				},
			},
			wantWinners:           []int{1, 2},
			wantScorePerCandidate: map[int]int{1: 1, 2: 1, 3: 0},
			wantErr:               false,
		},
		{
			name: "Invalid score (2 for combined approval voting)",
			args: args{
				candidates: []int{1, 2},
				votes:      [](map[int]int){{1: 2, 2: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Invalid score (-2 for combined approval voting)",
			args: args{
				candidates: []int{1, 2},
				votes:      [](map[int]int){{1: -2, 2: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "All negative votes",
			args: args{
				candidates: []int{1, 2, 3},
				votes: [](map[int]int){
					{1: -1, 2: -1, 3: -1},
					{1: -1, 2: -1, 3: -1},
				},
			},
			wantWinners:           []int{1, 2, 3},
			wantScorePerCandidate: map[int]int{1: -2, 2: -2, 3: -2},
			wantErr:               false,
		},
		{
			name: "Nil candidates",
			args: args{
				candidates: nil,
				votes:      [](map[int]int){{1: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWinners, gotScorePerCandidate, err := DafnyCaller_CombinedApprovalVoting(tt.args.candidates, tt.args.votes)
			if (err != nil) != tt.wantErr {
				t.Errorf("DafnyCaller_CombinedApprovalVoting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(gotWinners, tt.wantWinners) {
					t.Errorf("DafnyCaller_CombinedApprovalVoting() gotWinners = %v, want %v", gotWinners, tt.wantWinners)
				}
				if !reflect.DeepEqual(gotScorePerCandidate, tt.wantScorePerCandidate) {
					t.Errorf("DafnyCaller_CombinedApprovalVoting() gotScorePerCandidate = %v, want %v", gotScorePerCandidate, tt.wantScorePerCandidate)
				}
			}
		})
	}
}

func TestDafnyCaller_ScoreVoting(t *testing.T) {
	type args struct {
		candidates []int
		votes      [](map[int]int)
	}
	tests := []struct {
		name                  string
		args                  args
		wantWinners           []int
		wantScorePerCandidate map[int]int
		wantErr               bool
	}{
		{
			name: "Valid score voting",
			args: args{
				candidates: []int{1, 2, 3},
				votes: [](map[int]int){
					{1: 5, 2: 3, 3: 1},
					{1: 4, 2: 5, 3: 2},
					{1: 2, 2: 4, 3: 5},
				},
			},
			wantWinners:           []int{2},
			wantScorePerCandidate: map[int]int{1: 11, 2: 12, 3: 8},
			wantErr:               false,
		},
		{
			name: "Invalid score (6 for score voting)",
			args: args{
				candidates: []int{1, 2},
				votes:      [](map[int]int){{1: 6, 2: 3}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Invalid score (-1 for score voting)",
			args: args{
				candidates: []int{1, 2},
				votes:      [](map[int]int){{1: -1, 2: 3}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "All minimum scores",
			args: args{
				candidates: []int{1, 2, 3},
				votes: [](map[int]int){
					{1: 0, 2: 0, 3: 0},
					{1: 0, 2: 0, 3: 0},
				},
			},
			wantWinners:           []int{1, 2, 3},
			wantScorePerCandidate: map[int]int{1: 0, 2: 0, 3: 0},
			wantErr:               false,
		},
		{
			name: "All maximum scores",
			args: args{
				candidates: []int{1, 2},
				votes: [](map[int]int){
					{1: 5, 2: 5},
					{1: 5, 2: 5},
				},
			},
			wantWinners:           []int{1, 2},
			wantScorePerCandidate: map[int]int{1: 10, 2: 10},
			wantErr:               false,
		},
		{
			name: "Negative candidate ID",
			args: args{
				candidates: []int{-1, 2},
				votes:      [](map[int]int){{-1: 3, 2: 4}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWinners, gotScorePerCandidate, err := DafnyCaller_ScoreVoting(tt.args.candidates, tt.args.votes)
			if (err != nil) != tt.wantErr {
				t.Errorf("DafnyCaller_ScoreVoting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(gotWinners, tt.wantWinners) {
					t.Errorf("DafnyCaller_ScoreVoting() gotWinners = %v, want %v", gotWinners, tt.wantWinners)
				}
				if !reflect.DeepEqual(gotScorePerCandidate, tt.wantScorePerCandidate) {
					t.Errorf("DafnyCaller_ScoreVoting() gotScorePerCandidate = %v, want %v", gotScorePerCandidate, tt.wantScorePerCandidate)
				}
			}
		})
	}
}

func TestDafnyCaller_MajorityVoting(t *testing.T) {
	type args struct {
		candidates []int
		votes      [](map[int]int)
	}
	tests := []struct {
		name                  string
		args                  args
		wantWinners           []int
		wantScorePerCandidate map[int]int
		wantErr               bool
	}{
		{
			name: "Valid approval voting with clear winner",
			args: args{
				candidates: []int{1, 2, 3},
				votes: [](map[int]int){
					{1: 1},
					{1: 1},
					{3: 1},
					{},
				},
			},
			wantWinners:           []int{1},
			wantScorePerCandidate: map[int]int{1: 2, 2: 0, 3: 1},
			wantErr:               false,
		},
		{
			name: "Zero Vote",
			args: args{
				candidates: []int{1, 2, 3},
				votes: [](map[int]int){
					{1: 1},
					{1: 1},
					{3: 1},
					{3: 0},
				},
			},
			wantWinners:           []int{1},
			wantScorePerCandidate: map[int]int{1: 2, 2: 0, 3: 1},
			wantErr:               true,
		},
		{
			name: "Nil candidates",
			args: args{
				candidates: nil,
				votes:      [](map[int]int){{1: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Nil votes",
			args: args{
				candidates: []int{1, 2},
				votes:      nil,
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Empty candidates",
			args: args{
				candidates: []int{},
				votes:      [](map[int]int){{1: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Invalid score (2 for approval voting)",
			args: args{
				candidates: []int{1, 2},
				votes:      [](map[int]int){{1: 2, 2: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Vote for nonexistent candidate",
			args: args{
				candidates: []int{1, 2},
				votes:      [](map[int]int){{1: 1, 3: 1}},
			},
			wantWinners:           []int{},
			wantScorePerCandidate: map[int]int{},
			wantErr:               true,
		},
		{
			name: "Single candidate",
			args: args{
				candidates: []int{5},
				votes:      [](map[int]int){{5: 1}, {5: 1}},
			},
			wantWinners:           []int{5},
			wantScorePerCandidate: map[int]int{5: 2},
			wantErr:               false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWinners, gotScorePerCandidate, err := DafnyCaller_MajorityVoting(tt.args.candidates, tt.args.votes)
			if (err != nil) != tt.wantErr {
				t.Errorf("DafnyCaller_MajorityVoting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(gotWinners, tt.wantWinners) {
					t.Errorf("DafnyCaller_MajorityVoting() gotWinners = %v, want %v", gotWinners, tt.wantWinners)
				}
				if !reflect.DeepEqual(gotScorePerCandidate, tt.wantScorePerCandidate) {
					t.Errorf("DafnyCaller_MajorityVoting() gotScorePerCandidate = %v, want %v", gotScorePerCandidate, tt.wantScorePerCandidate)
				}
			}
		})
	}
}
