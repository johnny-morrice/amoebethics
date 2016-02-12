package libamoebethics

import (
    "fmt"
)

type UserNode struct {
    BaseNode
    NodeRef
}

func (un *UserNode) Validate() error {
    if !KnownNodeName(un.name) {
        return fmt.Errorf("Unknown node name: %v", un.name)
    }
    return nil
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
    Greet(m *SimNode) bool
    Serialize() string
}

type Neighbour struct {
    node *SimNode
    i int
}

func MakeNode(un *UserNode, t Torus) *SimNode {
    sn := &SimNode{}
    sn.entity = MakeEntity(un, t)
    sn.BaseNode = un.BaseNode
    sn.neighbours = []Neighbour{}
    return sn
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
        if n.entity.Greet(mn) {
            n.neighbours = append(n.neighbours, m)
        }
    }
}

func (n *SimNode) Update() {
    n.entity.Handle(n)
}