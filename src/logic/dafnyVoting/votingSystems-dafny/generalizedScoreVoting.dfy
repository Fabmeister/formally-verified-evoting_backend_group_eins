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

type Vote = map<nat,int>

method RuntimeRequires(cond: bool)
   ensures cond
{
   expect cond, "precondition failure";
}

// This file contains 3 voting systems: Approval Voting, Combined Approval Voting and Score Voting.
// Explanation of the systems:
//  Approval Voting: Every voter can approve or not approve every candidate individually
//  Combined Approval Voting: Every voter can approve/abstain/dissaprove every candiate individually
//  Score Voting: Every voter can give every candidate a score from 0-5


// A approval is counted as a score of 1, abstain as 0 and dissaprove as -1
// In all 3 systems the winner/winners is/are the ones who got the highest score
// (and there is always at least one winner)
// The obtained score of each candidate is also returned as a map

// The 3 systems can be generalized into one function CalculateElectionResult_SCORED_GENERAL,
// which has as Input allowedScores, which is a set of all scores voters are allowed to assign candidates.
// It returns the winner/winners, who have/has the highest score
// allowedScores for:
//  Approval: {0, 1}
//  Combined Approval: {-1, 0, 1}
//  Score: {0, 1, 2, 3, 4, 5}

// The split in wrapper and _Internal methods is important, because the wrapper function makes the runtime check and the internal does the static check 

// Input:
//  candidates: set of Candidate-Ids
//  votes: set of votes, where one vote consists of all Candidate-Ids, the voter approved 
// Output:
// winners: all candidates with the most approvals
// votesPerCandidate: Map from Candidates to their amount of votes they got
method CalculateElectionResult_APPROVAL_Internal(candidates: set<nat>, votes: seq<Vote>) returns (winners:set<nat>, scorePerCandidate:map<nat,int>)

  requires 0 !in candidates
  requires candidates!={}
  // all votes are for existing candidates and the given score is in allowedScores {0, 1}
  requires forall i :: (0<=i<|votes|) ==> (votes[i].Keys <= candidates && votes[i].Values <= {0, 1})


  // ensure winners are correct:
  ensures winners<=candidates
  ensures |winners|>=1
  // the Score of each candidate (number of approvals) is correct for every candidate
  ensures forall c:nat :: (c in candidates) ==> (c in scorePerCandidate.Keys)&&(scorePerCandidate[c]==scoreForSingleCandidate_SCORED_GENERAL(c,votes))

  // winner has higher (or same) score as anyone else
  ensures forall c1:: (c1 in winners) ==> forall c2:: (c2 in candidates) ==> scorePerCandidate[c1]>=scorePerCandidate[c2]
  // someone with higher score than everyone else is a winner
  ensures forall c1:: ((c1 in candidates)&&(forall c2:: (c2 in candidates) && scorePerCandidate[c1]>=scorePerCandidate[c2])) ==> c1 in winners
  // for someone not a winner there must be someone else with a higher score
  ensures forall c1:: (c1 in candidates)&&(c1 !in winners) ==> exists c2:: (c2 in candidates) && scorePerCandidate[c1]<scorePerCandidate[c2]
  {
    winners, scorePerCandidate := CalculateElectionResult_SCORED_GENERAL(candidates, votes, {0, 1});
  }

method CalculateElectionResult_APPROVAL(candidates: set<nat>, votes: seq<Vote>) returns (winners:set<nat>, scorePerCandidate:map<nat,int>)
{
  RuntimeRequires(0 !in candidates);
  RuntimeRequires(candidates!={});
  RuntimeRequires(forall i :: (0<=i<|votes|) ==> (votes[i].Keys <= candidates && votes[i].Values <= {0, 1}));

  winners, scorePerCandidate := CalculateElectionResult_APPROVAL_Internal(candidates, votes);
}

// Reduction of Majority Vote to Approval Voting with only one canididate per vote
method CalculateElectionResult_Majority_Internal(candidates: set<nat>, votes: seq<Vote>) returns (winners:set<nat>, scorePerCandidate:map<nat,int>)
  requires 0 !in candidates
  requires candidates!={}
  // all votes are for existing candidates and the given score is in allowedScores {0, 1}
  requires forall i :: (0<=i<|votes|) ==> (votes[i].Keys <= candidates && votes[i].Values <= {0, 1})
  // max one approval/canidate per vote 
  requires forall i :: (0<=i<|votes|) ==> (|votes[i].Keys| <= 1 && ((|votes[i].Keys| == 1) ==> (votes[i].Values == {1})))

  // ensure winners are correct:
  ensures winners<=candidates
  ensures |winners|>=1
  // the Score of each candidate (number of approvals) is correct for every candidate
  ensures forall c:nat :: (c in candidates) ==> (c in scorePerCandidate.Keys)&&(scorePerCandidate[c]==scoreForSingleCandidate_SCORED_GENERAL(c,votes))

  // winner has higher (or same) score as anyone else
  ensures forall c1:: (c1 in winners) ==> forall c2:: (c2 in candidates) ==> scorePerCandidate[c1]>=scorePerCandidate[c2]
  // someone with higher score than everyone else is a winner
  ensures forall c1:: ((c1 in candidates)&&(forall c2:: (c2 in candidates) && scorePerCandidate[c1]>=scorePerCandidate[c2])) ==> c1 in winners
  // for someone not a winner there must be someone else with a higher score
  ensures forall c1:: (c1 in candidates)&&(c1 !in winners) ==> exists c2:: (c2 in candidates) && scorePerCandidate[c1]<scorePerCandidate[c2]
  {
    winners, scorePerCandidate := CalculateElectionResult_APPROVAL_Internal(candidates, votes);
  }

method CalculateElectionResult_Majority(candidates: set<nat>, votes: seq<Vote>) returns (winners:set<nat>, scorePerCandidate:map<nat,int>)
{
  RuntimeRequires(0 !in candidates);
  RuntimeRequires(candidates!={});
  RuntimeRequires(forall i :: (0<=i<|votes|) ==> (votes[i].Keys <= candidates && votes[i].Values <= {0, 1}));
  RuntimeRequires(forall i :: (0<=i<|votes|) ==> (|votes[i].Keys| <= 1 && ((|votes[i].Keys| == 1) ==> (votes[i].Values == {1}))));
  
  winners, scorePerCandidate := CalculateElectionResult_Majority_Internal(candidates, votes);
}

method CalculateElectionResult_COMBINED_APPROVAL_Internal(candidates: set<nat>, votes: seq<Vote>) returns (winners:set<nat>, scorePerCandidate:map<nat,int>)

requires 0 !in candidates
requires candidates!={}
// all votes are for existing candidates and the given score is either -1,0 or 1 (for reject, abstain, approve)
requires forall i :: (0<=i<|votes|) ==> (votes[i].Keys <= candidates && votes[i].Values <= {-1, 0, 1})


// ensure winners are correct:
ensures winners<=candidates
ensures |winners|>=1
// the Score of each candidate is correct for every candidate
ensures forall c:nat :: (c in candidates) ==> (c in scorePerCandidate.Keys)&&(scorePerCandidate[c]==scoreForSingleCandidate_SCORED_GENERAL(c,votes))

// winner has higher (or same) score as anyone else
ensures forall c1:: (c1 in winners) ==> forall c2:: (c2 in candidates) ==> scorePerCandidate[c1]>=scorePerCandidate[c2]
// someone with higher score than everyone else is a winner
ensures forall c1:: ((c1 in candidates)&&(forall c2:: (c2 in candidates) && scorePerCandidate[c1]>=scorePerCandidate[c2])) ==> c1 in winners
  // for someone not a winner there must be someone else with a higher score
ensures forall c1:: (c1 in candidates)&&(c1 !in winners) ==> exists c2:: (c2 in candidates) && scorePerCandidate[c1]<scorePerCandidate[c2]
{
  winners, scorePerCandidate := CalculateElectionResult_SCORED_GENERAL(candidates, votes, {-1, 0, 1});
}

method CalculateElectionResult_COMBINED_APPROVAL(candidates: set<nat>, votes: seq<Vote>) returns (winners:set<nat>, scorePerCandidate:map<nat,int>)
{
  RuntimeRequires(0 !in candidates);
  RuntimeRequires(candidates!={});
  RuntimeRequires(forall i :: (0<=i<|votes|) ==> (votes[i].Keys <= candidates && votes[i].Values <= {-1, 0, 1}));
  winners, scorePerCandidate := CalculateElectionResult_COMBINED_APPROVAL_Internal(candidates, votes);
}
  
method CalculateElectionResult_SCORE_Internal(candidates: set<nat>, votes: seq<Vote>) returns (winners:set<nat>, scorePerCandidate:map<nat,int>)
  requires 0 !in candidates
  requires candidates!={}
  // all votes are for existing candidates and the given score is in {0,1,2,3,4,5}
  requires forall i :: (0<=i<|votes|) ==> (votes[i].Keys <= candidates && votes[i].Values <= {0, 1, 2, 3, 4, 5})


  // ensure winners are correct:
  ensures winners<=candidates
  ensures |winners|>=1
  // the Score of each candidate is correct for every candidate
  ensures forall c:nat :: (c in candidates) ==> (c in scorePerCandidate.Keys)&&(scorePerCandidate[c]==scoreForSingleCandidate_SCORED_GENERAL(c,votes))

  // winner has higher (or same) score as anyone else
  ensures forall c1:: (c1 in winners) ==> forall c2:: (c2 in candidates) ==> scorePerCandidate[c1]>=scorePerCandidate[c2]
  // someone with higher score than everyone else is a winner
  ensures forall c1:: ((c1 in candidates)&&(forall c2:: (c2 in candidates) && scorePerCandidate[c1]>=scorePerCandidate[c2])) ==> c1 in winners
  // for someone not a winner there must be someone else with a higher score
  ensures forall c1:: (c1 in candidates)&&(c1 !in winners) ==> exists c2:: (c2 in candidates) && scorePerCandidate[c1]<scorePerCandidate[c2]
  {
    winners, scorePerCandidate := CalculateElectionResult_SCORED_GENERAL(candidates, votes, {0, 1, 2, 3, 4, 5});
  }

method CalculateElectionResult_SCORE(candidates: set<nat>, votes: seq<Vote>) returns (winners:set<nat>, scorePerCandidate:map<nat,int>)
{
  RuntimeRequires(0 !in candidates);
  RuntimeRequires(candidates!={});
  RuntimeRequires(forall i :: (0<=i<|votes|) ==> (votes[i].Keys <= candidates && votes[i].Values <= {0, 1, 2, 3, 4, 5}));

  winners, scorePerCandidate := CalculateElectionResult_SCORE_Internal(candidates, votes);
}


//------------------------------------------------------------------------------------------------------------------------------
// general function for a voting system where every voter can give every candidate a score in allowedScores.
// winner/winners are the candidate/candidates with the highest obtained score
method CalculateElectionResult_SCORED_GENERAL(candidates: set<nat>, votes: seq<Vote>, allowedScores: set<int>) returns (winners: set<nat>, scorePerCandidate:map<nat,int>)
  requires 0 !in candidates
  requires 0 in allowedScores
  requires candidates!={}
  // all votes are for existing candidates and the given score is allowedScores
  requires forall i :: (0<=i<|votes|) ==> (votes[i].Keys <= candidates && votes[i].Values <= allowedScores)



  // ensure winners are correct:
  ensures winners<=candidates
  ensures |winners|>=1
  // the Score of each candidate is correct for every candidate
  ensures forall c:nat :: (c in candidates) ==> (c in scorePerCandidate.Keys)&&(scorePerCandidate[c]==scoreForSingleCandidate_SCORED_GENERAL(c,votes))

  // winner has higher (or same) score as anyone else
  ensures forall c1:: (c1 in winners) ==> forall c2:: (c2 in candidates) ==> scorePerCandidate[c1]>=scorePerCandidate[c2]
  // someone with higher score than everyone else is a winner
  ensures forall c1:: ((c1 in candidates)&&(forall c2:: (c2 in candidates) && scorePerCandidate[c1]>=scorePerCandidate[c2])) ==> c1 in winners
  ensures forall c1:: (c1 in candidates)&&(c1 !in winners) ==> exists c2:: (c2 in candidates) && scorePerCandidate[c1]<scorePerCandidate[c2]
{
  scorePerCandidate := getScorePerCandidate_SCORED_GENERAL(candidates, votes);


  // to track which candidate was already handled 
  // (has more(or the same) amount of votes as all handled candidates)
  var candidatesNotYetHandled := candidates; 

  //handle first candidate separately outside of loop to have invariant |winners|>=1
  var cand:|(cand in candidatesNotYetHandled);
  
  var maxVotes:int:=scorePerCandidate[cand];
  winners:= {cand};

  // Iterate over remaining candidates:
  while (candidatesNotYetHandled != {})
  decreases candidatesNotYetHandled
  invariant |winners|>=1
  invariant winners<=candidates
  invariant forall c:: (c in winners) ==> scorePerCandidate[c] == maxVotes
  // if cand not in winner, but already handled, it has less votes than maxVotes:
  invariant forall c:: ((c !in winners)&&(c in candidates-candidatesNotYetHandled)) ==> scorePerCandidate[c] < maxVotes
  {
    var cand:|(cand in candidatesNotYetHandled);
    
    var count := scorePerCandidate[cand];
    if (count>=maxVotes){
      if (count>maxVotes){
       maxVotes := count;
       winners := {cand}; 
      }
      else{
        winners := winners + {cand};
      }
    }
    candidatesNotYetHandled := candidatesNotYetHandled - {cand};
  }
  return winners, scorePerCandidate;
}

ghost function scoreForCandidateFromSingleVote_SCORED_GENERAL(c:nat, vote:Vote): (r:int)
ensures c in vote ==> r== vote[c]
ensures c !in vote ==> r==0
{
  if c in vote 
  then vote[c]
  else 0
}

// count the total score c got in votes
ghost function scoreForSingleCandidate_SCORED_GENERAL(c:nat, votes:seq<Vote>): (r:int)
ensures |votes|==0 ==> r==0
// remove all votes once and repeat this step by step, until no votes are left
ensures |votes|!=0 ==> r == scoreForSingleCandidate_SCORED_GENERAL(c, votes[..|votes|-1]) + scoreForCandidateFromSingleVote_SCORED_GENERAL(c, votes[|votes|-1])
{
  if |votes|==0 then 0
  else scoreForSingleCandidate_SCORED_GENERAL(c, votes[..|votes|-1]) + scoreForCandidateFromSingleVote_SCORED_GENERAL(c, votes[|votes|-1])
}

// calculate mapping from candidates to their respective scores
method getScorePerCandidate_SCORED_GENERAL(candidates: set<nat>, votes: seq<Vote>) returns (scorePerCandidate:map<nat, int>)
// make sure all votes are for existing candidates:
requires forall i:: (0<=i<|votes|) ==> votes[i].Keys <= candidates
// make sure theres a count for all candidates
ensures forall cand:nat :: (cand in candidates) ==> (cand in scorePerCandidate.Keys)
// make sure each count is correct
ensures forall cand:nat :: (cand in candidates) ==> (scoreForSingleCandidate_SCORED_GENERAL(cand,votes) == scorePerCandidate[cand])
{
  //Create Candidate-to-score map
  scorePerCandidate := map cand| cand in candidates:: 0;
  assert (candidates == scorePerCandidate.Keys   <==>   (forall c :: c in candidates <==> c in scorePerCandidate.Keys));
  assert(candidates==scorePerCandidate.Keys);

  //Iterate over all Votes and count up all scores
  var i:=0;
  while (i<|votes|)
  invariant 0<=i<=|votes|
  invariant candidates == scorePerCandidate.Keys
  invariant forall cand :: (cand in candidates) ==> scoreForSingleCandidate_SCORED_GENERAL(cand, votes[..i]) == scorePerCandidate[cand]
  {
    var singleVote:Vote:= votes[i];
    var singleVoteFull := singleVote;

    assert (singleVoteFull-singleVote.Keys) == map[];

    //handle one single Vote (iterate over all scores given for candidates):
    while(singleVote.Keys!={})
    decreases singleVote
    invariant candidates == scorePerCandidate.Keys
    invariant forall unhandledCand:: (unhandledCand in singleVoteFull.Keys && unhandledCand in singleVote.Keys) ==> singleVote[unhandledCand]==singleVoteFull[unhandledCand]
    invariant forall cand :: (cand in candidates) ==> scoreForSingleCandidate_SCORED_GENERAL(cand, votes[..i]+[singleVoteFull-singleVote.Keys]) == scorePerCandidate[cand] 
    {
      assert forall unhandledCand:: unhandledCand in singleVote.Keys ==> scoreForCandidateFromSingleVote_SCORED_GENERAL(unhandledCand, singleVoteFull-singleVote.Keys) == 0;
      assert forall handledCand:: handledCand in (singleVoteFull - singleVote.Keys).Keys ==> scoreForCandidateFromSingleVote_SCORED_GENERAL(handledCand, singleVoteFull-singleVote.Keys) == singleVoteFull[handledCand];
      // Update scores one vote gave for different candidates one by one
      var votedCandidateInScoreMap:nat :| (votedCandidateInScoreMap in singleVote.Keys);
      var oldScoreOfCandidate := scorePerCandidate[votedCandidateInScoreMap];
      var newScoreOfCandidate := oldScoreOfCandidate + singleVote[votedCandidateInScoreMap];
      scorePerCandidate := scorePerCandidate[votedCandidateInScoreMap:= newScoreOfCandidate];

      singleVote := singleVote - {votedCandidateInScoreMap};
    }
    assert votes[i] == singleVoteFull-singleVote.Keys;
    assert votes[..i]+[singleVoteFull-singleVote.Keys]==votes[..i+1];

    i := i+1;  
  }
  assert votes[..i]==votes;

  return scorePerCandidate;
}