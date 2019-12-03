package consensus

type Vote map[string]VoteCount

type VoteCount struct {
	N          int
	AddressMap map[string]bool
}

func NewVote() Vote {
	return make(map[string]VoteCount)
}

func (v Vote) addVote(data string, address string) {
	if _, ok := v[data]; !ok {
		addressMap := make(map[string]bool)
		addressMap[address] = true
		v[data] = VoteCount{
			N:          1,
			AddressMap: addressMap,
		}
	} else {
		count := v[data]
		if _, ok := count.AddressMap[address]; !ok {
			count.AddressMap[address] = true
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
