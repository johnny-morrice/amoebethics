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
}

type Entity interface {
    Handle(m *SimNode, s *Sim)
    Greet(n *SimNode, m *SimNode, s *Sim)
    Serialize(w io.Writer) error
    Deserialize(r io.Reader) error
}

func MakeNode(un *UserNode, base SimBase, yard EntityYard) (*SimNode, error) {
    sn := &SimNode{}
    ent, err := yard.MakeEntity(un, base)
    if err != nil {
        return nil, err
    }
    sn.Entity = ent
    sn.BaseNode = un.BaseNode
    sn.Neighbours = []*SimNode{}
    sn.Beliefs = MakeBeliefSet(base.BeliefMap, un.Beliefs)
    sn.P = UserVec2BlasVec(un.P)
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

func (n *SimNode) Update(s *Sim) {
    n.Entity.Handle(n, s)
}

func Node2String(n *SimNode) string {
    buff := bytes.Buffer{}
    err := n.Entity.Serialize(&buff)
    if err != nil {
        // Should never happen
        panic(err)
    }
    return buff.String()
}