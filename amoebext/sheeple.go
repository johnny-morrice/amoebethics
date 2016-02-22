package amoebext

import (
    "encoding/json"
    "io"
    m "github.com/gonum/matrix/mat64"
    lib "github.com/johnny-morrice/amoebethics/libamoebethics"
)

type UserSheeple struct {
    UserRudder
    LastP lib.UserVec
}

type Sheeple struct {
    Rudder
    LastP *m.Vector
}

var _ lib.Entity = (*Sheeple)(nil)

func NewSheepleNode(t lib.Torus) *lib.SimNode {
    sh := &Sheeple{}
    sh.D = lib.Vec2(1.0, 0)
    sh.S = 1.0
    n := lib.EmptyNode()
    n.Entity = sh
    n.Name = "sheeple"
    n.P = randPlace(t)
    sh.LastP = lib.Vec2(0,0)
    sh.wander(n, t)
    sh.wander(n, t)
    return n
}

func (sheep *Sheeple) Handle(n *lib.SimNode, s *lib.Sim) {
    // Change belief to whatever said
    for _, m := range n.Neighbours {
        for _, b := range m.Expression.Slice() {
            n.Beliefs.HoldBelief(b)
        }
    }
    sheep.wander(n, s.Torus)
}

func (sheep *Sheeple) wander(n *lib.SimNode, t lib.Torus) {
    // Move in random direction
    sheep.LastP.CopyVec(n.P)
    sheep.RandSteer()
    sheep.Move(n.P)
    t.Map(n.P)
}

func (sheep *Sheeple) Greet(n *lib.SimNode, m *lib.SimNode, s *lib.Sim) {
    // Never greet anyone
}

func (s *Sheeple) Serialize(w io.Writer) error {
    enc := json.NewEncoder(w)
    us := UserSheeple{}
    us.UserRudder = Rudder2UserRudder(s.Rudder)
    us.LastP = lib.BlasVec2UserVec(s.LastP)
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
    s.LastP = lib.UserVec2BlasVec(us.LastP)
    return nil
}

func (sh *Sheeple) Interpolate(n *lib.SimNode, time float64) {
    n.P.CopyVec(sh.LastP)
    sh.Jolt(n.P, time)
}

type SheepleFactory struct {}

var _ lib.EntityFactory = SheepleFactory{}

func (sf SheepleFactory) Build(un *lib.UserNode, base lib.SimBase) (lib.Entity, error) {
    sh := &Sheeple{}
    err := decodeEntity(sh, un)
    return sh, err
}