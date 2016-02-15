package libamoebethics

import (
    "bytes"
    "io"
    "github.com/gonum/matrix/mat64"
)

type UserNode struct {
    BaseNode
    NodeRef
}

type NodeRef struct {
    Neighbours []int
    Extension string
    P UserVec
    Beliefs []Belief
    Expression []Belief
}

type BaseNode struct {
    Name string
    Id int
}

type MachineNode struct {
    Neighbours []*SimNode
    Entity Entity
    P *mat64.Vector
    Beliefs BeliefSet
    Expression BeliefSet
}

type SimNode struct {
    MachineNode
    BaseNode
    changeBuffer *SimNode
}

type Entity interface {
    Handle(m *SimNode, s *Sim)
    Greet(n *SimNode, m *SimNode, s *Sim)
    Interpolate(n *SimNode, time float64)
    Serialize(w io.Writer) error
    Deserialize(r io.Reader) error
}

func MakeNode(un *UserNode, id int, base SimBase, yard EntityYard) (*SimNode, error) {
    sn, err := makeNodePart(un, id, base, yard)

    if err != nil {
        return nil, err
    }

    sn.changeBuffer, _ = makeNodePart(un, id, base, yard)

    return sn, nil
}

func makeNodePart(un *UserNode, id int, base SimBase, yard EntityYard) (*SimNode, error) {
    sn := &SimNode{}
    ent, err := yard.MakeEntity(un, base)
    if err != nil {
        return nil, err
    }
    sn.Entity = ent
    sn.BaseNode = un.BaseNode
    sn.Neighbours = []*SimNode{}
    sn.Beliefs = MakeBeliefSet(base.Beliefs, un.Beliefs)
    sn.Expression = MakeBeliefSet(base.Beliefs, un.Expression)
    sn.P = UserVec2BlasVec(un.P)
    sn.Id = id

    return sn, nil
}

func (n *SimNode) ClearNeighbours() {
    cnt := len(n.Neighbours)
    // Clear to allow garbage collection of pointers, if any
    for i := 0; i < cnt; i++ {
        n.Neighbours[i] = nil
    }
    // Reslice to avoid allocation
    n.Neighbours = n.Neighbours[:0]
}

func (n *SimNode) AddNeigbour(m *SimNode) {
    n.Neighbours = append(n.Neighbours, m)
}

func (n *SimNode) Handshake(m *SimNode, s *Sim) {
    if n != m {
        n.Entity.Greet(n, m, s)
    }
}

func (n *SimNode) Interpolate(time float64) {
    n.Entity.Interpolate(n, time)
}

func (n *SimNode) Update(s *Sim) {
    n.Entity.Handle(n.changeBuffer, s)
}

func (n *SimNode) WriteChange() {
    n.BaseNode = n.changeBuffer.BaseNode
    n.ClearNeighbours()
    for _, m := range n.changeBuffer.Neighbours {
        n.Neighbours = append(n.Neighbours, m)
    }
    n.P.CloneVec(n.changeBuffer.P)
    n.Beliefs.Copy(n.changeBuffer.Beliefs)
    n.Expression.Copy(n.changeBuffer.Expression)
}

func SimNode2UserNode(sn *SimNode) UserNode {
    un := UserNode{}
    un.BaseNode = sn.BaseNode
    uneigh := make([]int, len(sn.Neighbours))
    un.Neighbours = uneigh
    un.Extension = Entity2String(sn.Entity)
    for i, m := range sn.Neighbours {
        uneigh[i] = m.Id
    }
    un.Beliefs = sn.Beliefs.Slice()
    un.Expression = sn.Beliefs.Slice()
    un.P = BlasVec2UserVec(sn.P)
    return un
}

func Entity2String(e Entity) string {
    buff := bytes.Buffer{}
    err := e.Serialize(&buff)
    if err != nil {
        // Should never happen
        panic(err)
    }
    return buff.String()
}