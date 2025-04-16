package app

import "fmt"

type VoteResponse struct {
	from        *Address
	to          *Address
	term        int
	voteGranted bool
}

func NewVoteResponse(from *Address, to *Address, term int, voteGranted bool) *VoteResponse {
	return &VoteResponse{from: from, to: to, term: term, voteGranted: voteGranted}
}

func (v *VoteResponse) String() string {
	return fmt.Sprintf("VoteResponse{from: %s, to: %s, term: %d, voteGranted: %v}", v.from, v.to, v.term, v.voteGranted)
}

func (v *VoteResponse) Equal(other *VoteResponse) bool {
	return v.from.Equal(other.from) && v.to.Equal(other.to) && v.term == other.term && v.voteGranted == other.voteGranted
}

func (v *VoteResponse) GetFrom() *Address {
	return v.from
}

func (v *VoteResponse) GetTo() *Address {
	return v.to
}

func (v *VoteResponse) GetTerm() int {
	return v.term
}

func (v *VoteResponse) GetVoteGranted() bool {
	return v.voteGranted
}

func (v *VoteResponse) SetTerm(term int) {
	v.term = term
}

func (v *VoteResponse) SetFrom(from *Address) {
	v.from = from
}

func (v *VoteResponse) SetTo(to *Address) {
	v.to = to
}

func (v *VoteResponse) SetVoteGranted(voteGranted bool) {
	v.voteGranted = voteGranted
}
