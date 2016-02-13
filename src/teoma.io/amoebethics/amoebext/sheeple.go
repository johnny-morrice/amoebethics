package amoebext

import (
    "io"
    lib "teoma.io/amoebethics/libamoebethics"
)

type Sheeple struct {

}

var _ lib.Entity = (*Sheeple)(nil)

func (s *Sheeple) Handle(m *lib.SimNode) {

}

func (s *Sheeple) Greet(n *lib.SimNode, m *lib.SimNode) {
}

func (s *Sheeple) Serialize(w io.Writer) error {
    return nil
}

func (s *Sheeple) Deserialize(r io.Reader) error {
    return nil
}

type SheepleFactory struct {}

var _ lib.EntityFactory = SheepleFactory{}

func (sf SheepleFactory) Build(un *lib.UserNode, t lib.Torus) (lib.Entity, error) {
    sh := &Sheeple{}
    err := decodeEntity(sh, un)
    return sh, err
}