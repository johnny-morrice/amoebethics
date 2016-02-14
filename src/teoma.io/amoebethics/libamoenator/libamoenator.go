package libamoenator

import (
    "fmt"
    core "teoma.io/amoebethics/libamoebethics"
    ext "teoma.io/amoebethics/amoebext"
)

type PreFrame struct {
    Beliefs []string
    Palette []Color
}

type Frame struct {
    PreFrame
    Explosions []Explosion
    Shapes [][]ColorBox
}

func MakeFrame(pre PreFrame) Frame {
    fr := Frame{PreFrame: pre}
    fr.Explosions = []Explosion {}
    fr.Shapes = [][]ColorBox {
        []ColorBox {},
        []ColorBox {},
    }

    return fr
}

const CircleI int = 0
const SquareI int = 1

type Explosion struct {
    Radius float64
    Color int
    Intensity uint32
}

type ColorBox struct {
    Radius float64
    Colors []int
    P core.UserVec
}

type Color struct {
    R uint8
    G uint8
    B uint8
}

var red Color = Color{R: ^uint8(0),}
var green Color = Color{G: ^uint8(0),}
var blue Color = Color{B: ^uint8(0),}
var purple Color = Color{R: ^uint8(0), B: ^uint8(0),}
var yellow Color = Color{R: ^uint8(0), G: ^uint8(0),}
var turqouise Color = Color{G: ^uint8(0), B: ^uint8(0),}
var __colors []Color

func Init() {
    __colors = []Color{red, green, blue, purple, yellow, turqouise,}
}

type renderer struct {
    pre PreFrame
    framecnt uint
    fr Frame
    nodes []*core.SimNode
}

func MakeRenderer(pkt core.SimPacket, yard core.EntityYard, framecnt uint) (renderer, error) {
    r := renderer{}
    r.framecnt = framecnt
    r.pre.Palette = __colors
    colcnt := len(r.pre.Palette)
    if len(pkt.Beliefs) > colcnt {
        return r, fmt.Errorf("Only supports", colcnt, "beliefs")
    }
    r.pre.Beliefs = pkt.Beliefs
    r.nodes = make([]*core.SimNode, len(pkt.Nodes))
    for i, un := range pkt.Nodes {
        n, err := core.MakeNode(&un, pkt.SimBase, yard)
        if err != nil {
            return r, err
        }
        r.nodes[i] = n
    }
    return r, nil
}

func (r renderer) Render() []Frame {
    out := make([]Frame, r.framecnt)
    entGroups := r.splitNodes()
    seqmax := float64(r.framecnt)
    for i := uint(0); i < r.framecnt; i++ {
        r.fr = MakeFrame(r.pre)
        tstep := float64(i + 1)
        time := 1.0 / (seqmax * tstep)
        for _, entnodes := range entGroups {
            uns := userNodes(entnodes)
            r.nodeComposite(time, uns, entnodes)
        }
        r.interpolate(time)
        out[i] = r.fr
    }

    return out
}

func (r renderer) splitNodes() [][]*core.SimNode {
    split := make(map[string]*[]*core.SimNode) // At least we return a nicer type
    for _, sn := range r.nodes {
        target, ok := split[sn.Name]
        if ok {
            *target = append(*target, sn)
        } else {
            target = new([]*core.SimNode)
            *target = []*core.SimNode { sn, }
            split[sn.Name] = target
        }
    }
    out := make([][]*core.SimNode, len(split))
    i := 0
    for _, slp := range split {
        out[i] = *slp
    }
    return out
}

func (r renderer) interpolate(time float64) {
    for _, n := range r.nodes {
        n.Interpolate(time)
    }
}

func (r renderer) nodeComposite(time float64, user []core.UserNode, nodes []*core.SimNode) {
    const shpTypeCnt = 2
    name := user[0].Name
    switch name {
    case "sheeple":
        r.renderSheeple(user)
    case "tv":
        r.renderTv(time, user, nodes)
    default:
        panic("Unknown node type: " + name)
    }
}

func (r renderer) renderSheeple(user []core.UserNode) {
    circles := make([]ColorBox, len(user))
    for i, un := range user {
        circles[i] = renderNode(un)
    }
    r.fr.Shapes[CircleI] = append(r.fr.Shapes[CircleI], circles...)
}

func (r renderer) renderTv(time float64, user []core.UserNode, sim []*core.SimNode) {
    squares := make([]ColorBox, len(user))
    for i, un := range user {
        squares[i] = renderNode(un)
        e := sim[i].Entity
        tv := e.(*ext.Tv) // A better deserialization mechanism would allow no typecast
        expls := make([]Explosion, len(un.Expression))
        for i, b := range un.Expression {
            expls[i] = renderExplosion(tv.R, time, b)
        }
        r.fr.Explosions = append(r.fr.Explosions, expls...)
    }
    r.fr.Shapes[SquareI] = append(r.fr.Shapes[SquareI], squares...)
}

func userNodes(nodes []*core.SimNode) []core.UserNode {
    out := make([]core.UserNode, len(nodes))
    for i, n := range nodes {
        out[i] = core.SimNode2UserNode(n)
    }
    return out
}

func renderNode(node core.UserNode) ColorBox {
    const radius float64 = 0.05
    box := ColorBox{}
    box.Colors = make([]int, len(node.Beliefs))
    for i, b := range node.Beliefs {
        box.Colors[i] = b.Id
    }
    box.Radius = radius
    box.P = node.P
    return box
}

func renderExplosion(radius, time float64, b core.Belief) Explosion {
    const max = ^uint32(0)
    const fmax = float64(max)
    ex := Explosion{}
    ex.Radius = radius * (time * 2)
    if time < 0.5 {
        ex.Intensity = max
    } else {
        ex.Intensity = uint32(fmax - ((time - 0.5) * 2))
    }
    ex.Color = b.Id
    return ex
}