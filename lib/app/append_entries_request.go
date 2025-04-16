package app

import (
	"fmt"
)

type AppendEntriesReq struct {
	from        *Address
	to          *Address
	term        int
	prevIndex   int
	prevTerm    int
	entries     []LogEntry
	commitIndex int
}

func NewAppendEntriesReq(from *Address, to *Address, term int, prevIndex int, prevTerm int, entries []LogEntry, commitIndex int) *AppendEntriesReq {
	return &AppendEntriesReq{
		from:        from,
		to:          to,
		term:        term,
		prevIndex:   prevIndex,
		prevTerm:    prevTerm,
		entries:     entries,
		commitIndex: commitIndex,
	}
}

func (a *AppendEntriesReq) String() string {
	return fmt.Sprintf("AppendEntriesReq{from: %s, to: %s, term: %d, prevIndex: %d, prevTerm: %d, entries: %v, commitIndex: %d}", a.from, a.to, a.term, a.prevIndex, a.prevTerm, a.entries, a.commitIndex)
}

func (a *AppendEntriesReq) Equal(other *AppendEntriesReq) bool {
	return a.from.Equal(other.from) && a.to.Equal(other.to) && a.term == other.term && a.prevIndex == other.prevIndex && a.prevTerm == other.prevTerm && a.commitIndex == other.commitIndex
}

func (a *AppendEntriesReq) NotEqual(other *AppendEntriesReq) bool {
	return !a.Equal(other)
}

func (a *AppendEntriesReq) GetFrom() *Address {
	return a.from
}

func (a *AppendEntriesReq) GetTo() *Address {
	return a.to
}

func (a *AppendEntriesReq) GetTerm() int {
	return a.term
}

func (a *AppendEntriesReq) GetPrevIndex() int {
	return a.prevIndex
}

func (a *AppendEntriesReq) GetPrevTerm() int {
	return a.prevTerm
}

func (a *AppendEntriesReq) GetEntries() []LogEntry {
	return a.entries
}

func (a *AppendEntriesReq) GetCommitIndex() int {
	return a.commitIndex
}

func (a *AppendEntriesReq) SetFrom(from *Address) {
	a.from = from
}

func (a *AppendEntriesReq) SetTo(to *Address) {
	a.to = to
}

func (a *AppendEntriesReq) SetTerm(term int) {
	a.term = term
}

func (a *AppendEntriesReq) SetPrevIndex(prevIndex int) {
	a.prevIndex = prevIndex
}

func (a *AppendEntriesReq) SetPrevTerm(prevTerm int) {
	a.prevTerm = prevTerm
}

func (a *AppendEntriesReq) SetEntries(entries []LogEntry) {
	a.entries = entries
}
