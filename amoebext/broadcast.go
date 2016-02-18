package amoebext

import (
    "math/rand"
    lib "teoma.io/amoebethics/libamoebethics"
)

type Speaker struct {
    R float64
    InvF uint32 // Frequency of speech = 1 / InvF
}

func (b Speaker) Speaking() bool {
    const max32 = ^uint32(0)
    cut := max32 / b.InvF
    return rand.Uint32() < cut
}

func (b Speaker) Heard(n *lib.SimNode, m *lib.SimNode, t lib.Torus) bool {
    return t.Explodes(b.R, n.P, m.P)
}