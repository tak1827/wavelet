// Copyright (c) 2019 Perlin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package wavelet

import (
	"github.com/perlin-network/noise/skademlia"
	"github.com/perlin-network/wavelet/conf"
	"github.com/perlin-network/wavelet/sys"
	"sync"
	// "github.com/davecgh/go-spew/spew"
)

type syncVote struct {
	voter     *skademlia.ID
	outOfSync bool
}

type finalizationVote struct {
	voter *skademlia.ID
	round *Round
}

func CollectVotesForSync(
	accounts *Accounts,
	snowball *Snowball,
	voteChan <-chan syncVote,
	wg *sync.WaitGroup,
	snowballK int,
) {
	votes := make([]syncVote, 0, snowballK)
	voters := make(map[AccountID]struct{}, snowballK)

	for vote := range voteChan {
		if _, recorded := voters[vote.voter.PublicKey()]; recorded {
			continue // To make sure the sampling process is fair, only allow one vote per peer.
		}

		voters[vote.voter.PublicKey()] = struct{}{}
		votes = append(votes, vote)

		if len(votes) == cap(votes) {
			snapshot := accounts.Snapshot()

			stakes := make(map[AccountID]float64, len(votes))
			maxStake := float64(0)

			for _, vote := range votes {
				s, _ := ReadAccountStake(snapshot, vote.voter.PublicKey())

				if s < sys.MinimumStake {
					s = sys.MinimumStake
				}

				stake := float64(s)
				stakes[vote.voter.PublicKey()] = stake

				if maxStake < stake {
					maxStake = stake
				}
			}

			votesStakesPercentages := make(map[bool]float64, len(votes))
			totalStakePercentages := float64(0)

			for _, vote := range votes {
				percent := stakes[vote.voter.PublicKey()] / maxStake
				votesStakesPercentages[vote.outOfSync] += percent
				totalStakePercentages += percent
			}

			var majority Identifiable
			for _, vote := range votes {
				if votesStakesPercentages[vote.outOfSync]/totalStakePercentages >= conf.GetSyncVoteThreshold() {
					majority = &outOfSyncVote{outOfSync: vote.outOfSync}
					break
				}
			}

			snowball.Tick(majority)

			voters = make(map[AccountID]struct{}, snowballK)
			votes = votes[:0]
		}
	}

	if wg != nil {
		wg.Done()
	}
}

func CollectVotesForFinalization(
	accounts *Accounts,
	snowball *Snowball,
	voteChan <-chan finalizationVote,
	wg *sync.WaitGroup,
	snowballK int,
) {
	votes := make([]finalizationVote, 0, snowballK)
	voters := make(map[AccountID]struct{}, snowballK)
	tallies := make(map[[32]byte]uint32)
	var majorityTally uint32 = 0

	for vote := range voteChan {
		if _, recorded := voters[vote.voter.PublicKey()]; recorded {
			continue // To make sure the sampling process is fair, only allow one vote per peer.
		}

		voters[vote.voter.PublicKey()] = struct{}{}
		votes = append(votes, vote)

		tallies[vote.round.ID] += vote.round.Transactions
		// spew.Dump(vote.round.Transactions)

		if len(votes) == cap(votes) {
			// snapshot := accounts.Snapshot()

			// stakes := make(map[AccountID]float64, len(votes))
			// maxStake := float64(0)

			// for _, vote := range votes {
			// 	s, _ := ReadAccountStake(snapshot, vote.voter.PublicKey())

			// 	if s < sys.MinimumStake {
			// 		s = sys.MinimumStake
			// 	}

			// 	stake := float64(s)
			// 	stakes[vote.voter.PublicKey()] = stake

			// 	if maxStake < stake {
			// 		maxStake = stake
			// 	}
			// }

			// votesStakesPercentages := make(map[AccountID]float64, len(votes))
			// var totalStakePercentages float64

			// votesTransactionsNums := make(map[AccountID]uint32, len(votes))
			// var maxTransactionsNum uint32

			// votesEndDepths := make(map[AccountID]uint64, len(votes))
			// var minEndDepth uint64
			// minEndDepth-- // to have default value for minimal variable as max possible

			// for _, vote := range votes {
			// 	percent := stakes[vote.voter.PublicKey()] / maxStake

			// 	votesStakesPercentages[vote.round.ID] += percent
			// 	totalStakePercentages += percent

			// 	votesTransactionsNums[vote.round.ID] = vote.round.Transactions
			// 	if vote.round.Transactions > maxTransactionsNum {
			// 		maxTransactionsNum = vote.round.Transactions
			// 	}

			// 	depth := vote.round.End.Depth - vote.round.Start.Depth
			// 	votesEndDepths[vote.round.ID] = depth
			// 	if depth < minEndDepth {
			// 		minEndDepth = depth
			// 	}
			// }

			var majority *Round
			for _, vote := range votes {
				if tallies[vote.round.ID] > majorityTally {
					majority = vote.round
				}
				// stake := (votesStakesPercentages[vote.round.ID] / totalStakePercentages) * conf.GetStakeMajorityWeight()

				// var transactions float64
				// if maxTransactionsNum > 0 {
				// 	transactions = float64(votesTransactionsNums[vote.round.ID]/maxTransactionsNum) * conf.GetTransactionsNumMajorityWeight()
				// }

				// var depth float64
				// if votesEndDepths[vote.round.ID] > 0 {
				// 	depth = float64(minEndDepth/votesEndDepths[vote.round.ID]) * conf.GetRoundDepthMajorityWeight()
				// }

				// if stake+transactions+depth >= conf.GetFinalizationVoteThreshold() {
				// 	majority = vote.round
				// 	break
				// }
			}


			// spew.Dump(majority)
			snowball.Tick(majority)

			voters = make(map[AccountID]struct{}, snowballK)
			votes = votes[:0]
		}
	}

	if wg != nil {
		wg.Done()
	}
}
