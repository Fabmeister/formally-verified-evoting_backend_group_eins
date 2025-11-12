method RuntimeRequires2(cond: bool)
   ensures cond
{
   expect cond, "precondition failure";
}

// Explanation Instant Runoff Voting (IRV):
// Voting:
// - Each voter ranks the candidates according to their preference.
//    Not every candidate has to be ranked
// Finding the winner:
// - For each candidate: count how many times they are the highest priority in a vote
// - If one candidate has a total majority (only counting NONEMPTY votes!)
//    e.g. if votes are [[3,1,5],[],[3,5], [6]] then 3 has a total majority with 2 votes out of 3
// - Otherwise: Remove all candidates who have the lowest votecount. 
//    Adjust all votes by removing them from the priority lists of the voter:
//      e.g [[2], [3], [4], [4,3]] -> remove 2 and 3 -> new votes: [[], [], [4], [4]]
// - Repeat this until a winner is found OR all candidates are removed => then output is 0


// Input:
//  candidates: set of Candidate-Ids
//  votes: set of votes, where one vote consists of all Candidate-Ids, in descending order based on preference 
// Output:
//  winner: winning candidate OR 0 if no winner can be found
method CalculateElectionResult_IRV_Internal(candidates: set<nat>, votes: seq<seq<nat>>) returns (winner: nat)
  requires 0 !in candidates
  // all votes are for existing candidates:
  requires forall i,j:: (0<=i<|votes| && 0<=j<|votes[i]|) ==> votes[i][j] in candidates
  // make sure a vote doesn't vote for the same candidate multiple times
  requires forall i,j,k:: (0<=i<|votes| && 0<=k<|votes[i]| && 0<=j<|votes[i]| && k!=j) ==> votes[i][j] != votes[i][k]


  // if no votes were given, no one can be a winner
  ensures |votes| == 0 ==> winner==0
  // No candidates given => no candidate can be a winner
  ensures candidates == {} ==> winner == 0
  // No winner ==> no candidate can be a winner (neither in the first round, nor next vote counting rounds (where lowest candidates get removed))
  ensures winner==0 && |candidates|!=0 ==> forall c:: (c in candidates) ==> !isWinner_IRV(c, candidates, votes)

  ensures winner!=0 ==> winner in candidates
  // winner either has majority, or will get it in a later round (when removing the lowest candidates)
  ensures winner!=0 ==> isWinner_IRV(winner, candidates, votes)
{
  winner := 0;
  if |candidates|==0 {
    return 0;
  }
  assert candidates == {} ==> winner == 0;
  var votesPerCandidate: map<nat,nat> := CountVotesPerCandidate_IRV(candidates, votes);

  var nonEmptyVoteCount := CountNonEmptyVotes_IRV(votes);
  var requiredVoteCountForMajority := (nonEmptyVoteCount/2)+1;
  var candidatesNotYetCheckedForMajority:= candidates;

  // Iterate over all candidates to check if someone got the votes needed for a majority:
  while (candidatesNotYetCheckedForMajority != {})
  decreases candidatesNotYetCheckedForMajority
  invariant forall c:: (c in candidates-candidatesNotYetCheckedForMajority) ==> votesPerCandidate[c] < requiredVoteCountForMajority
  {
    var cand:|(cand in candidatesNotYetCheckedForMajority);
    if (votesPerCandidate[cand] >= requiredVoteCountForMajority){
      // candidate got majority
      assert cand!=0;
      assert cand in candidates;
      winner := cand;
      return winner;
    }
    candidatesNotYetCheckedForMajority := candidatesNotYetCheckedForMajority - {cand};
  }



  // Calculate the candidates who got the lowest votes
  var candidatesNotYetCheckedForLowestVotes := candidates;

  //handle first candidate separately outside of loop to have invariant |lowestVotedCandidates|>=1
  var cand:|(cand in candidatesNotYetCheckedForLowestVotes);
  var lowestVoteCount:int:=votesPerCandidate[cand];
  var lowestVotedCandidates:= {cand};

  while (candidatesNotYetCheckedForLowestVotes != {})
  decreases candidatesNotYetCheckedForLowestVotes
  invariant |lowestVotedCandidates|>=1
  invariant lowestVotedCandidates<=candidates
  invariant forall c:: (c in lowestVotedCandidates) ==> votesPerCandidate[c] == lowestVoteCount
  // if cand not in lowestVotedCandidates, but already handled, it has more votes than lowestVoteCount:
  invariant forall c:: ((c !in lowestVotedCandidates)&&(c in candidates-candidatesNotYetCheckedForLowestVotes))==> votesPerCandidate[c] > lowestVoteCount
  {
    var cand:|(cand in candidatesNotYetCheckedForLowestVotes);
    
    var count := votesPerCandidate[cand];
    if (count <= lowestVoteCount){
      if (count < lowestVoteCount){
        // lower votecount than all previously checked candidates
        lowestVoteCount := count;
        lowestVotedCandidates := {cand}; 
      }
      else{
        // same votecount as the lowest vote count encountered until now
        lowestVotedCandidates := lowestVotedCandidates + {cand};
      }
    }
    candidatesNotYetCheckedForLowestVotes := candidatesNotYetCheckedForLowestVotes - {cand};
  }

  var newCandidates := candidates - lowestVotedCandidates;
  assert lowestVotedCandidates == getLowestCandidatesInCurrentRound_IRV(candidates, votes);
  var newVotes := VotesWithRemovedCandidates_IRV(candidates,votes, lowestVotedCandidates);
  winner := CalculateElectionResult_IRV_Internal(newCandidates, newVotes);
}

method CalculateElectionResult_IRV(candidates: set<nat>, votes: seq<seq<nat>>) returns (winner: nat)
{
  RuntimeRequires2(0 !in candidates);
  RuntimeRequires2(forall i,j:: (0<=i<|votes| && 0<=j<|votes[i]|) ==> votes[i][j] in candidates);
  RuntimeRequires2(forall i,j,k:: (0<=i<|votes| && 0<=k<|votes[i]| && 0<=j<|votes[i]| && k!=j) ==> votes[i][j] != votes[i][k]);

  winner := CalculateElectionResult_IRV_Internal(candidates, votes);
}

ghost function countNonEmptyVotes_IRV(votes: seq<seq<nat>>) : (count: nat)
ensures |votes|==0 ==> count == 0
ensures |votes|>0 ==> count == countNonEmptyVotes_IRV(votes[..|votes|-1]) + if votes[|votes|-1] == [] then 0 else 1
{
  if votes == [] then 0 else
  countNonEmptyVotes_IRV(votes[..|votes|-1]) + if votes[|votes|-1]==[] then 0 else 1
}

// count how many votes are [], to be able to know how many votes are required for a total majority
method CountNonEmptyVotes_IRV(votes: seq<seq<nat>>) returns (count: nat)
ensures count == countNonEmptyVotes_IRV(votes)
{
  count := 0;
  var i := 0;
  while (i<|votes|)
  invariant 0 <= i <= |votes|
  invariant count == countNonEmptyVotes_IRV(votes[..i])
  {
    if votes[i] != [] {
      count := count + 1;
    }
    assert votes[..i+1]== votes[..i] + [votes[i]];
    i := i+1;
  }
  assert votes[..i] == votes;
}


// return all votes with the candidates in the candidatesToRemove set cut out
ghost function singleVoteWithRemovedCandidates_IRV(fullCandidates: set<nat>, vote: seq<nat>, candidatesToRemove: set<nat>): (newVote:seq<nat>)

  // all votes are for existing candidates
  requires forall i:: (0<=i<|vote|) ==> vote[i] in fullCandidates
  // only listed every candidate at most once
  requires forall i,j:: (0<=i<|vote| && 0<=j<|vote| && i!=j) ==> vote[i] != vote[j]

  // single entry in newVote comes from somewhere in original vote (mostly to help with the unique/duplicate vote postcondidtion)
  // IMPORTANT: without this assertion the verification times out
  // AND if the position of this assertion is switched to the bottom, then it still verifies for me locally,
  // but the gitlab pipeline times out in verifying this function!!!
  ensures forall i:: (0<=i<|newVote|) ==> exists j:: (0<=j<|vote|)&& vote[j] == newVote[i]

  ensures vote == [] ==> newVote == []
  // ensure removal of candidate, if in candidatesToRemove
  ensures |vote|!=0 ==> newVote == singleVoteWithRemovedCandidates_IRV(fullCandidates, vote[..|vote|-1], candidatesToRemove) + if vote[|vote|-1] in candidatesToRemove then [] else [vote[|vote|-1]]


  // no duplicate votes for a candidate
  ensures forall i,j:: (0<=i<|newVote| && 0<=j<|newVote| && i!=j) ==> newVote[i] != newVote[j]

  // remaining Candidates original Candidates, but not in candidatesToRemove
  ensures forall i ::(0<=i<|newVote|) ==> newVote[i] in fullCandidates - candidatesToRemove
{
  if vote == [] then [] else
  ghost var removedVoteWithoutLastEntry := singleVoteWithRemovedCandidates_IRV(fullCandidates, vote[..|vote|-1], candidatesToRemove);
  removedVoteWithoutLastEntry + 
  if vote[|vote|-1] in candidatesToRemove then [] 
  else 
  [vote[|vote|-1]]
}


ghost function votesWithRemovedCandidates_IRV(fullCandidates: set<nat>, votes: seq<seq<nat>>, candidatesToRemove: set<nat>) : (newVotes:seq<seq<nat>>)
  requires forall i,j:: (0<=i<|votes| && 0<=j<|votes[i]|) ==> votes[i][j] in fullCandidates
  requires forall i,j,k:: (0<=i<|votes| && 0<=k<|votes[i]| && 0<=j<|votes[i]| && k!=j) ==> votes[i][j] != votes[i][k]
  requires candidatesToRemove<=fullCandidates

  ensures |votes|==|newVotes|
  ensures forall i,j,k:: (0<=i<|newVotes| && 0<=k<|newVotes[i]| && 0<=j<|newVotes[i]| && k!=j) ==> newVotes[i][j] != newVotes[i][k]
  ensures forall i:: (0<=i<|votes|) ==> (newVotes[i] == singleVoteWithRemovedCandidates_IRV(fullCandidates, votes[i], candidatesToRemove))
  ensures forall i,j :: (0<=i<|newVotes|) && (0<=j<|newVotes[i]|) ==> newVotes[i][j] in fullCandidates - candidatesToRemove
{
  if |votes| == 0 then []
  else
  var firstVote:seq<nat> := votes[0];
  [singleVoteWithRemovedCandidates_IRV(fullCandidates,firstVote, candidatesToRemove)] + 
  if |votes|==1 then [] else votesWithRemovedCandidates_IRV(fullCandidates, votes[1..], candidatesToRemove) 
}

// remove the candidates in candidatesToRemove out of votes. To calculate new votes for the next simulated voting round
method VotesWithRemovedCandidates_IRV(fullCandidates: set<nat>, votes: seq<seq<nat>>, candidatesToRemove: set<nat>) returns (newVotes:seq<seq<nat>>)
  requires forall i,j:: (0<=i<|votes| && 0<=j<|votes[i]|) ==> votes[i][j] in fullCandidates
  requires forall i,j,k:: (0<=i<|votes| && 0<=k<|votes[i]| && 0<=j<|votes[i]| && k!=j) ==> votes[i][j] != votes[i][k]
  requires candidatesToRemove<=fullCandidates

  ensures newVotes==votesWithRemovedCandidates_IRV(fullCandidates, votes, candidatesToRemove)
{
  newVotes := [];
  var i:=0;
  while (i<|votes|)
  invariant 0<=i<=|votes|
  invariant |newVotes| == i
  invariant forall j:: (0<=j<i) ==> newVotes[j] == singleVoteWithRemovedCandidates_IRV(fullCandidates, votes[j], candidatesToRemove)
  {
    
    var j := 0;
    var newVote := [];
    while (j<|votes[i]|)
    invariant 0 <= j <= |votes[i]|
    invariant newVote == singleVoteWithRemovedCandidates_IRV(fullCandidates, votes[i][..j], candidatesToRemove)
    {
      // build new Vote, which only takes elements of the old vote votes[i], if the candidate is not in candidatesToRemove
      if (votes[i][j] !in candidatesToRemove) {
        newVote := newVote + [votes[i][j]];
      }
      assert votes[i][..j+1] == votes[i][..j] + [votes[i][j]];
      j := j+1;
    }
    assert votes[i][..j] == votes[i];
    
    newVotes := newVotes + [newVote];
    i := i+1;
  }
}



// count how often c is the first entry/preference in the votes
ghost function voteCountForSingleCandidate_IRV(c:nat, votes:seq<seq<nat>>) : (r:nat)
  ensures |votes|==0 ==> r==0
  ensures |votes|!=0 ==> r == voteCountForSingleCandidate_IRV(c, votes[..|votes|-1]) + if ((|votes[|votes|-1]|>0) && (c == votes[|votes|-1][0])) then 1 else 0
{
  if |votes|==0 then 0
  else voteCountForSingleCandidate_IRV(c, votes[..|votes|-1]) + if (|votes[|votes|-1]|>0) && (c == votes[|votes|-1][0]) then 1 else 0
}


// Return a map which maps every candidate to their vote count (number of first-preference-votes) they got
method CountVotesPerCandidate_IRV(candidates: set<nat>, votes: seq<seq<nat>>) returns (votesPerCandidate:map<nat,nat>)

  // make sure all votes are for existing candidates:
  requires forall i,j:: (0<=i<|votes| && 0<=j<|votes[i]|) ==> votes[i][j] in candidates
  // make sure theres a count for all candidates
  ensures forall cand:nat :: (cand in candidates) ==> (cand in votesPerCandidate.Keys)
  // make sure each count is correct
  ensures forall cand:nat :: (cand in candidates) ==> (voteCountForSingleCandidate_IRV(cand, votes) == votesPerCandidate[cand])
{
  //Create Candidate-to-Vote map
  votesPerCandidate := map cand| cand in candidates:: 0;

  //Iterate over all Votes and count them up
  var i:=0;
  while (i<|votes|)
  invariant 0<=i<=|votes|
  invariant candidates == votesPerCandidate.Keys
  invariant forall cand :: (cand in candidates) ==> voteCountForSingleCandidate_IRV(cand, votes[..i]) == votesPerCandidate[cand]
  {
    if |votes[i]|>0 {
      var votedCandidate := votes[i][0];
      var oldVotesOfCandidate := votesPerCandidate[votedCandidate];
      votesPerCandidate := votesPerCandidate[votedCandidate:= oldVotesOfCandidate+1];
    }

    assert votes[..i+1] == votes[..i] + [votes[i]];
    i := i+1;  
  }
  assert votes[..i]==votes;
  return votesPerCandidate;
}


// lowest number of first-preference-votes some candidate got
ghost function lowestVoteCountInCurrentRound_IRV(candidates: set<nat>, votes: seq<seq<nat>>) : (r: nat)
  requires candidates != {}
  ensures forall c:: (c in candidates) ==> r<=voteCountForSingleCandidate_IRV(c, votes)
  ensures exists lowestc:: (lowestc in candidates) && r==voteCountForSingleCandidate_IRV(lowestc, votes)
{
  ghost var cand :| (cand in candidates);
  ghost var candCount := voteCountForSingleCandidate_IRV(cand, votes);

  if |candidates|==1 then 
  assert (|candidates|==1 && cand in candidates) ==> |candidates - {cand}| == 0;
  candCount
  else 
  assert |candidates-{cand}|>=1;
  ghost var lowestCountInRemainingCandidates := lowestVoteCountInCurrentRound_IRV(candidates-{cand}, votes);
  if candCount < lowestCountInRemainingCandidates then candCount
  else lowestCountInRemainingCandidates

}

// highest number of first-preference-votes some candidate got
ghost function highestVoteCountInCurrentRound_IRV(candidates: set<nat>, votes: seq<seq<nat>>) : (r: nat)
  requires candidates != {}
  ensures forall c:: (c in candidates) ==> r>=voteCountForSingleCandidate_IRV(c, votes)
  ensures exists highestc:: (highestc in candidates) && r==voteCountForSingleCandidate_IRV(highestc, votes)
{
  ghost var cand :| (cand in candidates);
  ghost var candCount := voteCountForSingleCandidate_IRV(cand, votes);

  if |candidates|==1 then 
  assert (|candidates|==1 && cand in candidates) ==> |candidates - {cand}| == 0;
  candCount
  else 
  assert |candidates-{cand}|>=1;
  ghost var highestCountInRemainingCandidates := highestVoteCountInCurrentRound_IRV(candidates-{cand}, votes);
  if candCount > highestCountInRemainingCandidates then candCount
  else highestCountInRemainingCandidates

}


ghost function helper_getLowestCandidatesInCurrentRound_IRV(candidates: set<nat>, remainingCandidates: set<nat>, votes: seq<seq<nat>>) : (lowestRemainingCandidates:set<nat>)
  requires candidates != {}
  requires remainingCandidates <= candidates

  // if it is the first function call (no recursive call has been made, then the lowestRemainingCandidates are non-Empty)
  ensures candidates == remainingCandidates ==> lowestRemainingCandidates != {}
  ensures lowestRemainingCandidates <= remainingCandidates
  ensures remainingCandidates=={} ==> lowestRemainingCandidates=={}
  ensures forall lowestc, c :: (c in candidates && lowestc in lowestRemainingCandidates) ==> voteCountForSingleCandidate_IRV(lowestc, votes) <= voteCountForSingleCandidate_IRV(c, votes)
  ensures forall lowestc:: (lowestc in lowestRemainingCandidates) <==> lowestVoteCountInCurrentRound_IRV(candidates, votes) == voteCountForSingleCandidate_IRV(lowestc, votes) && (lowestc in remainingCandidates)
  decreases remainingCandidates
{
  if |remainingCandidates|==0 then {}
  else 
  ghost var cand :| (cand in remainingCandidates);
  ghost var lowestCount := lowestVoteCountInCurrentRound_IRV(candidates, votes);
  ghost var candCount := voteCountForSingleCandidate_IRV(cand, votes);

  if lowestCount == candCount then {cand} + helper_getLowestCandidatesInCurrentRound_IRV(candidates, remainingCandidates-{cand}, votes)
  else helper_getLowestCandidatesInCurrentRound_IRV(candidates, remainingCandidates-{cand}, votes)
}

// return the candidates who have the lowest vote count in votes
ghost function getLowestCandidatesInCurrentRound_IRV(candidates: set<nat>, votes: seq<seq<nat>>) : (lowestCandidates:set<nat>)
  requires candidates != {}
  ensures lowestCandidates != {}
  ensures lowestCandidates <= candidates
  ensures forall lowestc, c :: (c in candidates && lowestc in lowestCandidates) ==> voteCountForSingleCandidate_IRV(lowestc, votes) <= voteCountForSingleCandidate_IRV(c, votes)
  ensures forall lowestc:: (lowestc in lowestCandidates) <==> lowestVoteCountInCurrentRound_IRV(candidates, votes) == voteCountForSingleCandidate_IRV(lowestc, votes) && (lowestc in candidates)
{
  helper_getLowestCandidatesInCurrentRound_IRV(candidates, candidates, votes)
}


// just a help function for shorter pre/postconditions:
// Returns votes with the lowest Candidates removed
// This is just shorthand for votesWithRemovedCandidates_IRV(candidates, votes, getLowestCandidatesInCurrentRound_IRV(candidates, votes))
ghost function votesForNextRound_IRV(candidates: set<nat>, votes: seq<seq<nat>>) : (newVotes: seq<seq<nat>>)
  requires candidates != {}
  requires 0 !in candidates
  // all votes are for existing candidates:
  requires forall i,j:: (0<=i<|votes| && 0<=j<|votes[i]|) ==> votes[i][j] in candidates
  // make sure a vote doesn't vote for the same candidate multiple times
  requires forall i,j,k:: (0<=i<|votes| && 0<=k<|votes[i]| && 0<=j<|votes[i]| && k!=j) ==> votes[i][j] != votes[i][k]

  ensures newVotes == votesWithRemovedCandidates_IRV(candidates, votes, getLowestCandidatesInCurrentRound_IRV(candidates, votes))
{
  votesWithRemovedCandidates_IRV(candidates, votes, getLowestCandidatesInCurrentRound_IRV(candidates, votes))
}


// return true if and only if cand has a total majority in the (non-empty) votes 
ghost function hasMajorityInCurrentRound_IRV(cand: nat, candidates: set<nat>, votes: seq<seq<nat>>) : (hasMajority:bool)
  requires cand in candidates
  requires 0 !in candidates
  // all votes are for existing candidates:
  requires forall i,j:: (0<=i<|votes| && 0<=j<|votes[i]|) ==> votes[i][j] in candidates

  ensures !hasMajority ==> voteCountForSingleCandidate_IRV(cand, votes) < countNonEmptyVotes_IRV(votes)/2 +1
  ensures hasMajority ==> voteCountForSingleCandidate_IRV(cand, votes) >= countNonEmptyVotes_IRV(votes)/2 +1
{

  ghost var requiredVoteCountForMajority := (countNonEmptyVotes_IRV(votes)/2) + 1;
  if voteCountForSingleCandidate_IRV(cand, votes) >= requiredVoteCountForMajority 
  then true
  else 
  false
}


ghost function hasWinnerInCurrentRound_IRV(candidates: set<nat>, votes: seq<seq<nat>>) : (r:bool)
  requires candidates != {}
  requires 0 !in candidates
  // all votes are for existing candidates:
  requires forall i,j:: (0<=i<|votes| && 0<=j<|votes[i]|) ==> votes[i][j] in candidates

  ensures r ==> exists c:: c in candidates && hasMajorityInCurrentRound_IRV(c, candidates, votes)
  ensures !r ==> forall c:: c in candidates ==> !hasMajorityInCurrentRound_IRV(c, candidates, votes)
{

  ghost var requiredVoteCountForMajority := (countNonEmptyVotes_IRV(votes)/2) + 1;
  ghost var highestCount := highestVoteCountInCurrentRound_IRV(candidates, votes);

  if highestCount >= requiredVoteCountForMajority then true
  else false
}

//  isWinner: true if c wins the election false otherwise
ghost function isWinner_IRV(candidateToCheck: nat, candidates: set<nat>, votes: seq<seq<nat>>) : (isWinner: bool)
  requires candidateToCheck in candidates
  requires 0 !in candidates
  // all votes are for existing candidates:
  requires forall i,j:: (0<=i<|votes| && 0<=j<|votes[i]|) ==> votes[i][j] in candidates
  // make sure a vote doesn't vote for the same candidate multiple times
  requires forall i,j,k:: (0<=i<|votes| && 0<=k<|votes[i]| && 0<=j<|votes[i]| && k!=j) ==> votes[i][j] != votes[i][k]

  ensures |votes| == 0 ==> !isWinner
  // Can't be winner if candidates are empty
  ensures candidates == {} ==> !isWinner
  // Can't have majority if no Winner
  ensures !isWinner && candidates != {} ==> !hasMajorityInCurrentRound_IRV(candidateToCheck, candidates, votes)

  // If no winner, then either 
  // - someone else wins in the first round
  // - candidateToCheck gets eliminated in the first Round (as lowest Voted Candidate)
  // - candidateToCheck doesn't win in the later vote-counting rounds
  ensures !isWinner && candidates != {} ==> hasWinnerInCurrentRound_IRV(candidates, votes) || candidateToCheck in getLowestCandidatesInCurrentRound_IRV(candidates, votes) || !isWinner_IRV(candidateToCheck, candidates-getLowestCandidatesInCurrentRound_IRV(candidates, votes), votesWithRemovedCandidates_IRV(candidates, votes, getLowestCandidatesInCurrentRound_IRV(candidates, votes)))

  // If candidateToCheck is a winner, then either:
  // - has a majority in the first Round
  // - or wins in the later vote-counting rounds (and doesn't get eliminated in the first round)
  // (and it can be said that there must be at least on vote)
  ensures isWinner ==> |votes|!=0 && (hasMajorityInCurrentRound_IRV(candidateToCheck, candidates, votes) || (candidateToCheck !in getLowestCandidatesInCurrentRound_IRV(candidates, votes) && isWinner_IRV(candidateToCheck, candidates-getLowestCandidatesInCurrentRound_IRV(candidates, votes), votesForNextRound_IRV(candidates, votes))))
{
  if candidates == {} then 
   //no winner exists
   false
  else
    if hasMajorityInCurrentRound_IRV(candidateToCheck, candidates, votes) then 
      // candidate got majority in first round
      true
    else
      ghost var requiredVoteCountForMajority := (countNonEmptyVotes_IRV(votes)/2) + 1;
      ghost var highestCount := highestVoteCountInCurrentRound_IRV(candidates, votes);

      if highestCount >= requiredVoteCountForMajority then
        //someone else won in the first round
        false
      else 
        // calculate next Voting iteration by removing lowest candidates:
        ghost var newVotes := votesWithRemovedCandidates_IRV(candidates, votes, getLowestCandidatesInCurrentRound_IRV(candidates, votes));

        ghost var newCandidates := candidates - getLowestCandidatesInCurrentRound_IRV(candidates, votes);
          if candidateToCheck !in newCandidates then 
            // The candidate got eliminated in this round, so isn't a winner
            false
          else
            // assert candidateToCheck in newCandidates;
            isWinner_IRV(candidateToCheck, newCandidates, newVotes)
}


