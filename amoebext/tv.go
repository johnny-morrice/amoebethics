package amoebext

import (
    "encoding/json"
    "io"
    lib "github.com/johnny-morrice/amoebethics/libamoebethics"
)

type Tv struct {
    Speaker
}

var _ lib.Entity = (*Tv)(nil)

func NewTvNode(t lib.Torus) *lib.SimNode {
    tv := &Tv{}
    tv.R = 3
    tv.InvF = 3
    n := lib.EmptyNode()
    n.Entity = tv
    n.Name = "tv"
    n.P = randPlace(t)
    return n
}

func (tv *Tv) Handle(n *lib.SimNode, s *lib.Sim) {
    n.Expression.Clear()
    if tv.Speaking() {
        b := n.Beliefs.Rand()
        n.Expression.HoldBelief(b)
    }
}

func (tv *Tv) Greet(n *lib.SimNode, m *lib.SimNode, s *lib.Sim) {
    if tv.Heard(n, m, s.Torus) {
        m.AddNeigbour(n)
    }
}

func (tv *Tv) Serialize(w io.Writer) error {
    enc := json.NewEncoder(w)
    return enc.Encode(tv)
}

func (tv *Tv) Deserialize(r io.Reader) error {
    dec := json.NewDecoder(r)
    return dec.Decode(tv)
}

func (tv *Tv) Interpolate(n *lib.SimNode, time float64) {
    // No movement
}

type TvFactory struct {}

var _ lib.EntityFactory = TvFactory{}

func (sf TvFactory) Build(un *lib.UserNode, base lib.SimBase) (lib.Entity, error) {
    tv := &Tv{}
    err := decodeEntity(tv, un)
    return tv, err
}