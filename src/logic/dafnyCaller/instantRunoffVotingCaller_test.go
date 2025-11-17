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
	"testing"
)

func Test_fulfillsPreconditions_InstantRunoffVoting(t *testing.T) {
	type args struct {
		candidates []int
		votes      [][]int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid instant runoff voting",
			args: args{
				candidates: []int{1, 2, 3},
				votes: [][]int{
					{1, 2, 3},
					{2, 3, 1},
					{3, 1, 2},
				},
			},
			wantErr: false,
		},
		{
			name: "Empty candidates",
			args: args{
				candidates: []int{},
				votes:      [][]int{{1}},
			},
			wantErr: true,
		},
		{
			name: "Negative candidate",
			args: args{
				candidates: []int{-1, 2, 3},
				votes:      [][]int{{-1, 2, 3}},
			},
			wantErr: true,
		},
		{
			name: "Candidate with ID 0",
			args: args{
				candidates: []int{0, 1, 2},
				votes:      [][]int{{0, 1, 2}},
			},
			wantErr: true,
		},
		{
			name: "Duplicate candidate",
			args: args{
				candidates: []int{1, 2, 2},
				votes:      [][]int{{1, 2}},
			},
			wantErr: true,
		},
		{
			name: "Vote for nonexistent candidate",
			args: args{
				candidates: []int{1, 2, 3},
				votes:      [][]int{{1, 2, 4}},
			},
			wantErr: true,
		},
		{
			name: "Multiple votes for same candidate in single vote",
			args: args{
				candidates: []int{1, 2, 3},
				votes:      [][]int{{1, 2, 2}},
			},
			wantErr: true,
		},
		{
			name: "Empty vote (valid)",
			args: args{
				candidates: []int{1, 2, 3},
				votes:      [][]int{{}, {1, 2}, {3}},
			},
			wantErr: false,
		},
		{
			name: "Single candidate single vote",
			args: args{
				candidates: []int{5},
				votes:      [][]int{{5}},
			},
			wantErr: false,
		},
		{
			name: "Partial voting (not all candidates ranked)",
			args: args{
				candidates: []int{1, 2, 3, 4},
				votes:      [][]int{{1, 3}, {2, 4, 1}, {4}},
			},
			wantErr: false,
		},
		{
			name: "High candidate ID within limit",
			args: args{
				candidates: []int{1, 1000, 2000000},
				votes:      [][]int{{1, 1000, 2000000}},
			},
			wantErr: false,
		},
		{
			name: "Nil candidates",
			args: args{
				candidates: nil,
				votes:      [][]int{{1}},
			},
			wantErr: true,
		},
		{
			name: "Nil votes",
			args: args{
				candidates: []int{1, 2},
				votes:      nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fulfillsPreconditions_InstantRunoffVoting(tt.args.candidates, tt.args.votes); (err != nil) != tt.wantErr {
				t.Errorf("fulfillsPreconditions_InstantRunoffVoting() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDafnyCaller_InstantRunoffVoting(t *testing.T) {
	type args struct {
		candidates []int
		votes      [][]int
	}
	tests := []struct {
		name       string
		args       args
		wantWinner int
		wantErr    bool
	}{
		{
			name: "Valid instant runoff voting with clear winner",
			args: args{
				candidates: []int{1, 2, 3},
				votes: [][]int{
					{1, 2, 3},
					{1, 3, 2},
					{2, 1, 3},
				},
			},
			wantWinner: 1,
			wantErr:    false,
		},
		{
			name: "Single candidate",
			args: args{
				candidates: []int{5},
				votes:      [][]int{{5}, {5}, {}},
			},
			wantWinner: 5,
			wantErr:    false,
		},
		{
			name: "Nil candidates",
			args: args{
				candidates: nil,
				votes:      [][]int{{1}},
			},
			wantWinner: 0,
			wantErr:    true,
		},
		{
			name: "Nil votes",
			args: args{
				candidates: []int{1, 2},
				votes:      nil,
			},
			wantWinner: 0,
			wantErr:    true,
		},
		{
			name: "Invalid preconditions - negative candidate",
			args: args{
				candidates: []int{-1, 2},
				votes:      [][]int{{-1, 2}},
			},
			wantWinner: 0,
			wantErr:    true,
		},
		{
			name: "Invalid preconditions - candidate ID 0",
			args: args{
				candidates: []int{0, 1},
				votes:      [][]int{{0, 1}},
			},
			wantWinner: 0,
			wantErr:    true,
		},
		{
			name: "Invalid preconditions - duplicate candidate",
			args: args{
				candidates: []int{1, 1, 2},
				votes:      [][]int{{1, 2}},
			},
			wantWinner: 0,
			wantErr:    true,
		},
		{
			name: "Invalid preconditions - vote for nonexistent candidate",
			args: args{
				candidates: []int{1, 2},
				votes:      [][]int{{1, 2, 3}},
			},
			wantWinner: 0,
			wantErr:    true,
		},
		{
			name: "Invalid preconditions - multiple votes for same candidate",
			args: args{
				candidates: []int{1, 2, 3},
				votes:      [][]int{{1, 2, 2}},
			},
			wantWinner: 0,
			wantErr:    true,
		},
		{
			name: "Empty votes (valid but may result in no winner)",
			args: args{
				candidates: []int{1, 2, 3},
				votes:      [][]int{{}, {}, {}},
			},
			wantWinner: 0,
			wantErr:    false,
		},
		{
			name: "Partial preferences",
			args: args{
				candidates: []int{1, 2, 3, 4},
				votes: [][]int{
					{1, 2},
					{2, 3},
					{3, 4},
					{1},
				},
			},
			wantWinner: 1,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWinner, err := DafnyCaller_InstantRunoffVoting(tt.args.candidates, tt.args.votes)
			if (err != nil) != tt.wantErr {
				t.Errorf("DafnyCaller_InstantRunoffVoting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWinner != tt.wantWinner {
				t.Errorf("DafnyCaller_InstantRunoffVoting() = %v, want %v", gotWinner, tt.wantWinner)
			}
		})
	}
}
