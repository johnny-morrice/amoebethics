package libamoebethics

import (
    "bytes"
    "io"
)

type UserNode struct {
    BaseNode
    NodeRef
}

type NodeRef struct {
    Neighbours []int
    Extension string
}

type BaseNode struct {
    name string
    beliefs []Belief
    expression []Belief
    pos Vec2
}

type MachineNode struct {
    neighbours []Neighbour
    entity Entity
}

type SimNode struct {
    MachineNode
    BaseNode
}

type Entity interface {
    Handle(m *SimNode)
    Greet(n *SimNode, m *SimNode)
    Serialize(w io.Writer) error
    Deserialize(r io.Reader) error
}

type Neighbour struct {
    node *SimNode
    i int
}

func MakeNode(un *UserNode, t Torus, yard EntityYard) (*SimNode, error) {
    sn := &SimNode{}
    ent, err := yard.MakeEntity(un, t)
    if err != nil {
        return nil, err
    }
    sn.entity = ent
    sn.BaseNode = un.BaseNode
    sn.neighbours = []Neighbour{}
    return sn, nil
}

func (n *SimNode) clearNeighbours() {
    cnt := len(n.neighbours)
    // Clear to allow garbage collection of pointers, if any
    for i := 0; i < cnt; i++ {
        n.neighbours[i] = Neighbour{}
    }
    // Reslice to avoid allocation
    n.neighbours = n.neighbours[:0]
}

func (n *SimNode) Handshake(m Neighbour) {
    mn := m.node
    if n != mn {
        n.entity.Greet(n, mn)
    }
}

func (n *SimNode) Update() {
    n.entity.Handle(n)
}

func Node2String(n *SimNode) string {
    buff := bytes.Buffer{}
    err := n.entity.Serialize(&buff)
    if err != nil {
        // Should never happen
        panic(err)
    }
    return buff.String()
}