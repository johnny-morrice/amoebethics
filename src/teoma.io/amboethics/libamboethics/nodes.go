package libamboethics

import (
    "fmt"
)

type UserNode struct {
    BaseNode
    NodeRef
}

func (un *UserNode) Validate() error {
    if !knownNodeName(un.name) {
        return fmt.Errorf("Unknown node name: %v", un.name)
    }
    return nil
}

type NodeRef struct {
    Neighbours []int
}

type BaseNode struct {
    name string
    beliefs []Belief
    expression []Belief
    pos Vec2
}

type MachineNode struct {
    neighbours []Neighbour
    handle Handler
    greet Greeter
}

type SimNode struct {
    MachineNode
    BaseNode
}

type Handler func(n *SimNode)
type Greeter func(n *SimNode)

type Neighbour struct {
    node *SimNode
    i int
}

func makeNode(un *UserNode) *SimNode {
    sn := &SimNode{}
    h, g := makeEntity(un.name)
    sn.BaseNode = un.BaseNode
    sn.neighbours = []Neighbour{}
    sn.handle = h
    sn.greet = g
    return sn
}

func (n *SimNode) Handshake(m Neighbour) {
    mn := m.node
    if n != mn {
        n.greet(mn)
    }
}

func (n *SimNode) Update() {
    n.handle(n)
}

func knownNodeName(name string) bool {
    return false
}

func makeEntity(name string) (Handler, Greeter) {
    return nil, nil
}