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
    Nodes []*SimNode
}

func MakeSim(in SimPacket, yard EntityYard, outch chan<- SimPacket) (*Sim, error) {
    // We ignore input nodes neighbours
    nodes := make([]*SimNode, len(in.Nodes))
    for i, un := range in.Nodes {
        n, err := MakeNode(&un, in.SimBase, yard)
        if err != nil {
            return nil, err
        }
        nodes[i] = n
    }
    s := &Sim{}
    s.Nodes = nodes
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
    usernodes := make([]UserNode, len(s.Nodes))
    s.Ceach(func (sn *SimNode) {
        un := SimNode2UserNode(sn)
        usernodes[sn.Id] = un
    })
    out := SimPacket{}
    out.SimBase = s.SimBase
    out.Nodes = usernodes
    return out
}


func (s *Sim) attachNodes() {
    s.Ceach(func (n *SimNode) {
        n.ClearNeighbours()
        s.Each(func (m *SimNode) {
            n.Handshake(m, s)
        })
    })
}

func (s *Sim) nodeHandlers() {
    s.Ceach(func (n *SimNode) {
        n.Update(s)
    })
    s.Each(func (n *SimNode) {
        n.SwapFrames()
    })
}

func (s *Sim) Each(f func(n *SimNode)) {
    for _, n := range s.Nodes {
        f(n)
    }
}

func (s *Sim) Ceach(f func(n *SimNode)) {
    count := len(s.Nodes)
    hold := sync.WaitGroup{}
    hold.Add(count)
    for _, n := range s.Nodes {
        go func() {
            f(n)
            hold.Done()
        }()
    }
    hold.Wait()
}

type SimBase struct {
    Iteration int
    Itermax int
    Torus Torus
    Beliefs []string
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