package libamoebethics

import (
    "math/rand"
)

type Belief struct {
    Op Opinion
    Name string
    Id int
}

type Opinion uint8

const (
    DontKnow = Opinion(iota)
    IsFalse
    IsTrue
)

type BeliefSet []Opinion

func MakeBeliefSet(bels []string, held []Belief) BeliefSet {
    sl := make([]Opinion, len(bels))
    for _, b := range held {
        sl[b.Id] = b.Op
    }
    return BeliefSet(sl)
}

func (bs BeliefSet) HoldBelief(b Belief) {
    bs[b.Id] = b.Op
}

func (bs BeliefSet) Rand() Belief {
    cands := bs.Slice()
    r := rand.Intn(len(cands))
    return cands[r]
}

func (bs BeliefSet) Slice() []Belief {
    bels := make([]Belief, 0, len(bs))
    for i, opp := range ([]Opinion)(bs) {
        if opp != DontKnow {
            b := Belief{}
            b.Id = i
            b.Op = opp
            bels[i] = b
        }
    }
    return bels
}

func (bs BeliefSet) Clear() {
    opps := ([]Opinion)(bs)
    for i := 0; i < len(opps); i++ {
        opps[i] = DontKnow
    }
}