package libamoebethics

import (
    "encoding/json"
    "fmt"
    "io"
    "sync"
)

func ReadSimInput(r io.Reader) (SimInput, error) {
    dec := json.NewDecoder(r)
    s := SimInput{}
    err := dec.Decode(&s)
    return s, err
}

func WriteSimOutput(o []UserNode, w io.Writer) error {
    enc := json.NewEncoder(w)
    return enc.Encode(o)
}

func Simulate(s SimInput) (<-chan []UserNode, error) {
    err := s.Validate()
    if err != nil {
        return nil, err
    }

    outch := make(chan []UserNode)
    sim := MakeSim(s, outch)
    sim.ForkSim()
    return outch, nil
}

type Sim struct {
    outch chan<- []UserNode
    itercount int
    nodes []*SimNode
    torus Torus
}

func MakeSim(in SimInput, outch chan<- []UserNode) *Sim {
    // We ignore input nodes neighbours
    nodes := make([]*SimNode, len(in.Nodes))
    for i, un := range in.Nodes {
        nodes[i] = makeNode(&un)
    }
    s := &Sim{}
    s.nodes = nodes
    s.itercount = in.Itercount
    s.torus = in.Torus
    return s
}

func (s *Sim) ForkSim() {
    go func() {
        for i := 0; i < s.itercount; i++ {
            s.step()
            s.outch<- s.moment()
        }
    }()
}

func (s *Sim) step() {
    s.attachNodes()
    s.nodeHandlers()
}

func (s *Sim) moment() []UserNode {
    out := make([]UserNode, len(s.nodes))
    s.ceach(func (n Neighbour) {
        un := UserNode{}
        sn := n.node
        un.BaseNode = sn.BaseNode
        uneigh := make([]int, len(sn.neighbours))
        un.Neighbours = uneigh
        for i, m := range sn.neighbours {
            uneigh[i] = m.i
        }
        out[n.i] = un
    })
    return out
}


func (s *Sim) attachNodes() {
    s.ceach(func (n Neighbour) {
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

type SimInput struct {
    Itercount int
    Torus Torus
    Nodes []UserNode
}

func (in SimInput) Validate() error {
    if in.Itercount < 1 {
        return fmt.Errorf("Invalid itercount. Was %v.", in.Itercount)
    }

    if in.Torus.W < 0 || in.Torus.H < 0 {
        return fmt.Errorf("Invalid torus.  Was %v.", in.Torus)
    }

    for i, n := range in.Nodes {
        nerr := n.Validate()
        if nerr != nil {
            return fmt.Errorf("Error at node %v: %v", i, nerr)
        }
    }

    return nil
}