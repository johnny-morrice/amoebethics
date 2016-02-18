package amoebext

import (
    "encoding/json"
    "io"
    lib "github.com/johnny-morrice/amoebethics/libamoebethics"
)

type UserSheeple struct {
    UserRudder
}

type Sheeple struct {
    Rudder
}

var _ lib.Entity = (*Sheeple)(nil)

func (sheep *Sheeple) Handle(n *lib.SimNode, s *lib.Sim) {
    // Change belief to whatever said
    for _, m := range n.Neighbours {
        for _, b := range m.Expression.Slice() {
            n.Beliefs.HoldBelief(b)
        }
    }

    // Move in random direction
    sheep.RandSteer()
    sheep.Move(n.P)
    s.Torus.Map(n.P)
}

func (sheep *Sheeple) Greet(n *lib.SimNode, m *lib.SimNode, s *lib.Sim) {
    // Never greet anyone
}

func (s *Sheeple) Serialize(w io.Writer) error {
    enc := json.NewEncoder(w)
    us := UserSheeple{}
    us.UserRudder = Rudder2UserRudder(s.Rudder)
    return enc.Encode(&us)
}

func (s *Sheeple) Deserialize(r io.Reader) error {
    dec := json.NewDecoder(r)
    us := UserSheeple{}
    err := dec.Decode(&us)
    if err != nil {
        return err
    }
    s.Rudder = UserRudder2Rudder(us.UserRudder)
    return nil
}

func (sh *Sheeple) Interpolate(n *lib.SimNode, time float64) {
    sh.Jolt(n.P, time)
}

type SheepleFactory struct {}

var _ lib.EntityFactory = SheepleFactory{}

func (sf SheepleFactory) Build(un *lib.UserNode, base lib.SimBase) (lib.Entity, error) {
    sh := &Sheeple{}
    err := decodeEntity(sh, un)
    return sh, err
}