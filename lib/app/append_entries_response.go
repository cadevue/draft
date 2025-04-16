package app

import (
	"fmt"
)

type AppendEntriesResp struct {
	from       *Address
	to         *Address
	term       int
	success    bool
	matchIndex int
}

func NewAppendEntriesResp(from *Address, to *Address, term int, success bool, matchIndex int) *AppendEntriesResp {
	return &AppendEntriesResp{from: from, to: to, term: term, success: success, matchIndex: matchIndex}
}

func (a *AppendEntriesResp) String() string {
	return fmt.Sprintf("AppendEntriesResp{from: %s, to: %s, term: %d, success: %v, matchIndex: %d}", a.from, a.to, a.term, a.success, a.matchIndex)
}

func (a *AppendEntriesResp) Equal(other *AppendEntriesResp) bool {
	return a.from.Equal(other.from) && a.to.Equal(other.to) && a.term == other.term && a.success == other.success && a.matchIndex == other.matchIndex
}

func (a *AppendEntriesResp) NotEqual(other *AppendEntriesResp) bool {
	return !a.Equal(other)
}

func (a *AppendEntriesResp) GetFrom() *Address {
	return a.from
}

func (a *AppendEntriesResp) GetTo() *Address {
	return a.to
}

func (a *AppendEntriesResp) GetTerm() int {
	return a.term
}

func (a *AppendEntriesResp) GetSuccess() bool {
	return a.success
}

func (a *AppendEntriesResp) GetMatchIndex() int {
	return a.matchIndex
}

func (a *AppendEntriesResp) SetMatchIndex(matchIndex int) {
	a.matchIndex = matchIndex
}

func (a *AppendEntriesResp) SetSuccess(success bool) {
	a.success = success
}

func (a *AppendEntriesResp) SetTerm(term int) {
	a.term = term
}

func (a *AppendEntriesResp) SetFrom(from *Address) {
	a.from = from
}

func (a *AppendEntriesResp) SetTo(to *Address) {
	a.to = to
}
