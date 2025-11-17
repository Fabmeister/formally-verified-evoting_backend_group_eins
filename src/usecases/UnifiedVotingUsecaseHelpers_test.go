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

package usecases

import (
	dto "e-voting-service/data/dto"
	"reflect"
	"testing"
)

func Test_handleApproval_HandleVotePrecondition(t *testing.T) {
	type args struct {
		votes *[]dto.UnifiedVote
	}
	tests := []struct {
		name    string
		this    handleApprovalVote
		args    args
		wantErr bool
	}{
		{
			name: "Single Vote",
			this: handleApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Multi Vote",
			this: handleApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Including Empty Vote",
			this: handleApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 0},
					{CandidateID: 3, WahlInfo: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Including illegal Vote",
			this: handleApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 2},
					{CandidateID: 3, WahlInfo: 1},
				},
			},
			wantErr: true,
		},
		// Double vote test nicht sinnvoll, da an anderer Stelle geprüft wird
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := handleApprovalVote{}
			if err := this.HandleVotePrecondition(tt.args.votes); (err != nil) != tt.wantErr {
				t.Errorf("handleApproval.HandleVotePrecondition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_handleCombinedApproval_HandleVotePrecondition(t *testing.T) {
	type args struct {
		votes *[]dto.UnifiedVote
	}
	tests := []struct {
		name    string
		this    handleCombinedApprovalVote
		args    args
		wantErr bool
	}{
		{
			name: "Single Vote",
			this: handleCombinedApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Multi Vote",
			this: handleCombinedApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: -1},
				},
			},
			wantErr: false,
		},
		{
			name: "Including Empty Vote",
			this: handleCombinedApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 0},
					{CandidateID: 3, WahlInfo: -1},
				},
			},
			wantErr: false,
		},
		{
			name: "Including illegal Vote",
			this: handleCombinedApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 2},
					{CandidateID: 3, WahlInfo: 1},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := handleCombinedApprovalVote{}
			if err := this.HandleVotePrecondition(tt.args.votes); (err != nil) != tt.wantErr {
				t.Errorf("handleCombinedApproval.HandleVotePrecondition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_handleScoreVoting_HandleVotePrecondition(t *testing.T) {
	type args struct {
		votes *[]dto.UnifiedVote
	}
	tests := []struct {
		name    string
		this    handleScoreVotingVote
		args    args
		wantErr bool
	}{
		{
			name: "Single Vote",
			this: handleScoreVotingVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Multi Vote",
			this: handleScoreVotingVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 5},
				},
			},
			wantErr: false,
		},
		{
			name: "Including Empty Vote",
			this: handleScoreVotingVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 0},
					{CandidateID: 3, WahlInfo: 3},
				},
			},
			wantErr: false,
		},
		{
			name: "Including illegal Vote",
			this: handleScoreVotingVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: -1},
					{CandidateID: 3, WahlInfo: 5},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := handleScoreVotingVote{}
			if err := this.HandleVotePrecondition(tt.args.votes); (err != nil) != tt.wantErr {
				t.Errorf("handleScoreVoting.HandleVotePrecondition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_handleIRV_HandleVotePrecondition(t *testing.T) {
	type args struct {
		votes *[]dto.UnifiedVote
	}
	tests := []struct {
		name    string
		this    handleIRVVote
		args    args
		wantErr bool
	}{
		{
			name: "Single Vote",
			this: handleIRVVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Multi Vote",
			this: handleIRVVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 2},
					{CandidateID: 3, WahlInfo: 3},
				},
			},
			wantErr: false,
		},
		{
			name: "Multi Vote diff order",
			this: handleIRVVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 3, WahlInfo: 2},
					{CandidateID: 2, WahlInfo: 3},
				},
			},
			wantErr: false,
		},
		{
			name: "Skip a rank",
			this: handleIRVVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 2},
					{CandidateID: 3, WahlInfo: 4},
				},
			},
			wantErr: true,
		},
		{
			name: "Including Empty",
			this: handleIRVVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 0},
					{CandidateID: 3, WahlInfo: 2},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := handleIRVVote{}
			if err := this.HandleVotePrecondition(tt.args.votes); (err != nil) != tt.wantErr {
				t.Errorf("handleIRV.HandleVotePrecondition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_handleMajority_HandleVotePrecondition(t *testing.T) {
	type args struct {
		votes *[]dto.UnifiedVote
	}
	tests := []struct {
		name    string
		this    handleApprovalVote
		args    args
		wantErr bool
	}{
		{
			name: "Single Vote",
			this: handleApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Multi Vote",
			this: handleApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 1},
				},
			},
			wantErr: true,
		},
		{
			name: "Including Empty Vote",
			this: handleApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 1},
					{CandidateID: 2, WahlInfo: 0},
					{CandidateID: 3, WahlInfo: 0},
				},
			},
			wantErr: false,
		},
		{
			name: "Including illegal Vote",
			this: handleApprovalVote{},
			args: args{
				votes: &[]dto.UnifiedVote{
					{CandidateID: 1, WahlInfo: 0},
					{CandidateID: 2, WahlInfo: 2},
					{CandidateID: 3, WahlInfo: 0},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := handleMajorityVote{}
			if err := this.HandleVotePrecondition(tt.args.votes); (err != nil) != tt.wantErr {
				t.Errorf("handleMajorityVote.HandleVotePrecondition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_unifiedVotingHandleFactory(t *testing.T) {
	type args struct {
		wahl dto.Election
	}
	tests := []struct {
		name    string
		args    args
		want    iHandleVote
		wantErr bool
	}{
		{
			name:    "Approval",
			args:    args{dto.Election{Type: dto.APPROVAL_VOTING}},
			want:    handleApprovalVote{},
			wantErr: false,
		},
		{
			name:    "CombinedApproval",
			args:    args{dto.Election{Type: dto.COMBINED_APPROVAL_VOTING}},
			want:    handleCombinedApprovalVote{},
			wantErr: false,
		},
		{
			name:    "Score",
			args:    args{dto.Election{Type: dto.SCORE_VOTING}},
			want:    handleScoreVotingVote{},
			wantErr: false,
		},
		{
			name:    "IRV",
			args:    args{dto.Election{Type: dto.IRV}},
			want:    handleIRVVote{},
			wantErr: false,
		},
		{
			name:    "Unspecified",
			args:    args{dto.Election{Type: dto.ELECTION_TYPE_UNSPECIFIED}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unifiedVotingHandleFactory(tt.args.wahl)
			if (err != nil) != tt.wantErr {
				t.Errorf("unifiedVotingHandleFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unifiedVotingHandleFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}
