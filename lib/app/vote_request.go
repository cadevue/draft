package app

import "fmt"

type VoteRequest struct {
	from         *Address
	to           *Address
	term         int
	lastLogIndex int
	lastLogTerm  int
}

func NewVoteRequest(from *Address, to *Address, term int, lastLogIndex int, lastLogTerm int) *VoteRequest {
	return &VoteRequest{from: from, to: to, term: term, lastLogIndex: lastLogIndex, lastLogTerm: lastLogTerm}
}

func (v *VoteRequest) String() string {
	return fmt.Sprintf("VoteRequest{from: %s, to: %s, term: %d, lastLogIndex: %d, lastLogTerm: %d}", v.from, v.to, v.term, v.lastLogIndex, v.lastLogTerm)
}

func (v *VoteRequest) Equal(other *VoteRequest) bool {
	return v.from.Equal(other.from) && v.to.Equal(other.to) && v.term == other.term && v.lastLogIndex == other.lastLogIndex && v.lastLogTerm == other.lastLogTerm
}

func (v *VoteRequest) GetFrom() *Address {
	return v.from
}

func (v *VoteRequest) GetTo() *Address {
	return v.to
}

func (v *VoteRequest) GetTerm() int {
	return v.term
}

func (v *VoteRequest) GetLastLogIndex() int {
	return v.lastLogIndex
}

func (v *VoteRequest) GetLastLogTerm() int {
	return v.lastLogTerm
}

func (v *VoteRequest) SetTerm(term int) {
	v.term = term
}

func (v *VoteRequest) SetFrom(from *Address) {
	v.from = from
}

func (v *VoteRequest) SetTo(to *Address) {
	v.to = to
}

func (v *VoteRequest) SetLastLogIndex(lastLogIndex int) {
	v.lastLogIndex = lastLogIndex
}

func (v *VoteRequest) SetLastLogTerm(lastLogTerm int) {
	v.lastLogTerm = lastLogTerm
}
