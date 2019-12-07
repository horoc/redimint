package core

import "github.com/emirpasic/gods/sets/hashset"

type Vote map[string]VoteCount

type VoteCount struct {
	N          int
	AddressSet *hashset.Set
}

func NewVote() Vote {
	return make(map[string]VoteCount)
}

func (v Vote) addVote(data string, address string) {
	if _, ok := v[data]; !ok {
		v[data] = VoteCount{
			N:          1,
			AddressSet: hashset.New(address),
		}
	} else {
		count := v[data]
		if !count.AddressSet.Contains(address) {
			count.AddressSet.Add(address)
			count.N++
		}
	}
}

func (v Vote) getVoteNum(data string) int {
	if _, ok := v[data]; !ok {
		return v[data].N
	} else {
		return 0
	}
}
