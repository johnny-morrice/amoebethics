package amoebext

import (
    "io"
    lib "teoma.io/amoebethics/libamoebethics"
)

type Tv struct {

}

var _ lib.Entity = (*Tv)(nil)

func (tv *Tv) Handle(m *lib.SimNode) {

}

func (tv *Tv) Greet(n *lib.SimNode, m *lib.SimNode) {
}

func (tv *Tv) Serialize(w io.Writer) error {
    return nil
}

func (tv *Tv) Deserialize(r io.Reader) error {
    return nil
}

type TvFactory struct {}

var _ lib.EntityFactory = TvFactory{}

func (sf TvFactory) Build(un *lib.UserNode, t lib.Torus) (lib.Entity, error) {
    tv := &Tv{}
    err := decodeEntity(tv, un)
    return tv, err
}