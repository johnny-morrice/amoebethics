package libamoebethics

import (
    "encoding/json"
    "fmt"
    "io"
    "sync"
)

func ReadSimInput(r io.Reader) (SimPacket, error) {
    dec := json.NewDecoder(r)
    s := SimPacket{}
    err := dec.Decode(&s)
    return s, err
}

func WriteSimOutput(o SimPacket, w io.Writer) error {
    enc := json.NewEncoder(w)
    return enc.Encode(o)
}

func Simulate(s SimPacket, yard EntityYard) (<-chan SimPacket, error) {
    verr := s.Validate()
    if verr != nil {
        return nil, verr
    }

    outch := make(chan SimPacket)
    sim, serr := MakeSim(s, yard, outch)
    if serr != nil {
        return nil, serr
    }
    sim.ForkSim()
    return outch, nil
}

type Sim struct {
    SimBase
    outch chan<- SimPacket
    nodes []*SimNode
}

func MakeSim(in SimPacket, yard EntityYard, outch chan<- SimPacket) (*Sim, error) {
    // We ignore input nodes neighbours
    nodes := make([]*SimNode, len(in.Nodes))
    for i, un := range in.Nodes {
        n, err := MakeNode(&un, in.Torus, yard)
        if err != nil {
            return nil, err
        }
        nodes[i] = n
    }
    s := &Sim{}
    s.nodes = nodes
    s.SimBase = in.SimBase
    return s, nil
}

func (s *Sim) ForkSim() {
    go func() {
        // Cat input out for other programs down the pipeline
        s.outch<- s.moment()
        for i := 0; i < s.Itermax; i++ {
            s.step()
            s.outch<- s.moment()
        }
    }()
}

func (s *Sim) step() {
    s.attachNodes()
    s.nodeHandlers()
    s.Iteration++
}

func (s *Sim) moment() SimPacket {
    usernodes := make([]UserNode, len(s.nodes))
    s.ceach(func (n Neighbour) {
        un := UserNode{}
        sn := n.node
        un.BaseNode = sn.BaseNode
        uneigh := make([]int, len(sn.neighbours))
        un.Neighbours = uneigh
        un.Extension = Node2String(sn)
        for i, m := range sn.neighbours {
            uneigh[i] = m.i
        }
        usernodes[n.i] = un
    })
    out := SimPacket{}
    out.SimBase = s.SimBase
    out.Nodes = usernodes
    return out
}


func (s *Sim) attachNodes() {
    s.ceach(func (n Neighbour) {
        n.node.clearNeighbours()
        s.each(func (m Neighbour) {
            n.node.Handshake(m)
        })
    })
}

func (s *Sim) nodeHandlers() {
    s.ceach(func (n Neighbour) {
        n.node.Update()
    })
}

func (s *Sim) each(f func(n Neighbour)) {
    count := len(s.nodes)
    for i := 0; i < count; i++ {
        n := Neighbour{
            node: s.nodes[i],
            i: i,
        }
        f(n)
    }
}

func (s *Sim) ceach(f func(n Neighbour)) {
    count := len(s.nodes)
    hold := sync.WaitGroup{}
    hold.Add(count)
    for i := 0; i < count; i++ {
        n := Neighbour{
            node: s.nodes[i],
            i: i,
        }
        go func() {
            f(n)
            hold.Done()
        }()
    }
    hold.Wait()
}

type Belief struct {
    opp Opinion
    name string
}

type Opinion uint8

const (
    IsTrue = Opinion(iota)
    IsFalse
    DontKnow
)

type SimBase struct {
    Iteration int
    Itermax int
    Torus Torus
}

type SimPacket struct {
    SimBase
    Nodes []UserNode
}

func (in SimPacket) Validate() error {
    if in.Itermax < 1 {
        return fmt.Errorf("Invalid itermax. Was %v.", in.Itermax)
    }

    if in.Torus.W < 0 || in.Torus.H < 0 {
        return fmt.Errorf("Invalid torus.  Was %v.", in.Torus)
    }

    return nil
}