package libamoebethics

import (
    "math/rand"
)

type Belief struct {
    Opp Opinion
    Name string
    Id int
}

type Opinion uint8

const (
    DontKnow = Opinion(iota)
    IsFalse
    IsTrue
    IsBoth
)

type BeliefSet []Opinion

func MakeBeliefSet(belm map[string]uint, held []Belief) BeliefSet {
    cnt := len(belm)
    sl := make([]Opinion, cnt)
    for _, b := range held {
        sl[b.Id] = b.Opp
    }
    return BeliefSet(sl)
}

func (bs BeliefSet) HoldIrratBelief(b Belief) {
    old := bs[b.Id]
    switch old {
    case DontKnow:
        bs[b.Id] = b.Opp
    case IsTrue:
        if b.Opp == IsFalse {
            bs[b.Id] = IsBoth
        }
    case IsFalse:
        if b.Opp == IsTrue {
            bs[b.Id] = IsBoth
        }
    }
}

func (bs BeliefSet) Rand() Belief {
    cands := make([]Belief, 0)
    for i, opp := range ([]Opinion)(bs) {
        if opp != DontKnow {
            b := Belief{}
            b.Id = i
            b.Opp = opp
            cands = append(cands, b)
        }
    }
    r := rand.Intn(len(cands))
    return cands[r]
}

func (bs BeliefSet) Slice() []Belief {
    bels := make([]Belief, len(bs))
    for i, opp := range ([]Opinion)(bs) {
        if opp != DontKnow {
            b := Belief{}
            b.Id = i
            b.Opp = opp
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