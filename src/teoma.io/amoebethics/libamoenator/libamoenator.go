package libamoenator

import (
    "github.com/gonum/matrix/mat64"
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

const CircleI int = 0
const SquareI int = 1

type Explosion struct {
    Radius float64
    Color uint
    Intensity uint32
}

type ColorBox struct {
    Radius float64
    Colors []uint
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

func Convert(pkt core.SimPacket, yard core.EntityYard) (PreFrame, []*core.SimNode, error) {
    pre := PreFrame{}
    pre.Palette = __colors
    colcnt := len(pre.Palette)
    belcnt := len(pkt.BeliefMap)
    if belcnt > colcnt {
        return nil, nil, fmt.Errorf("Only supports", colcnt, "beliefs")
    }
    pre.Beliefs = pkt.Beliefs
    nodes := make([]*core.SimNode, len(pkt.Nodes))
    for i, un := range pkt.Nodes {
        n, err = core.MakeNode(un, pkt.SimBase, yard)
        if err != nil {
            return nil, nil, err
        }
        nodes[i] = n
    }
    return pre, nodes, nil
}

func Render(nodes []*core.SimNode, pre PreFrame, framecnt uint) []Frame {
    out := make([]Frame, framecnt)
    entGroups := splitNodes(nodes)
    seqmax := float64(framecnt)
    for i := uint(0); i < framecnt; i++ {
        fr := MakeFrame()
        tstep := float64(i + 1)
        time := 1.0 / (seqmax * tstep)
        for _, entnodes := range entGroups {
            uns := userNodes(entnodes)
            shapes, explosions := nodeComposite(time, uns, entnodes)
            insertFrame(shapes, explosions)
        }
        interpolate(time, nodes)
        out[i] = fr
    }

    return out
}

func insertFrame(shapes [][]ColorBox, explosions []Explosion, fr *Frame) {
    for i, sh := range shapes {
        if sh != nil {
            fr.Shapes[i] = append(fr.Shapes[i], sh...)
        }
    }
    if explosions != nil {
        append(fr.Explosion, explosions...)
    }
}

func interpolate(time float64, nodes []*core.SimNode) {
    for _, n := range nodes {
        n.Interpolate(time)
    }
}

func userNodes(nodes []*core.SimNode) []core.UserNode {
    out := make([]core.UserNode, len(nodes))
    for i, n := range nodes {
        out[i] = core.SimNode2UserNode(n)
    }
    return out
}

func renderNode(node UserNode, pre PreFrame) ColorBox {
    const radius float64 = 0.05
    box := ColorBox{}
    box.Colors = make([]Color, len(node.Beliefs))
    for i, b := range node.Beliefs {
        box.Colors[i] = b.Id
    }
    box.Radius = radius
    box.P = node.P
    return box
}

func nodeComposite(time float64, user []core.UserNode, nodes []*core.SimNode) ([][]ColorBox, []Explosion) {
    const shpTypeCnt = 2
    shapes := make([][]ColorBox, shpTypeCnt)
    explosions := []Explosion {}
    switch user[0].Name {
    case "sheeple":
        circles := make([]ColorBox, len(user))
        for i, un := range user {
            box := renderNode(un, fr.PreFrame)
            circles[i] = box
        }
        shapes[CircleI] = circles
    case "tv":
        for i, un := range user {
            box := renderNode(un, fr.PreFrame)
            shapes[SquareI][i] = box
            e := nodes[i].Entity
            tv := e.(*ext.Tv)
            expls := make([]Explosion, len(un.Expression))
            for i, b := range un.Expression {
                expls[i] := renderExplosion(time, b)
            }
        }
    default:
        panic("Unknown node type: " + name)
    }
    return shapes, explosions
}

func renderExplosion(radius, time float64, b Belief) Explosion {
    const maxIntense float64 = float64(^uint32(0))
    ex := Explosion{}
    ex.Radius = tv.R * (time * 2)
    if time < 0.5 {
        ex.Intensity = maxIntense
    } else {
        ex.Intensity = uint32(maxIntense - ((time - 0.5) * 2))
    }
    ex.Color = b.Id
}