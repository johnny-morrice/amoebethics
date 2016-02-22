package libamoenator

import (
    "fmt"
    "io"
    "encoding/json"
    core "github.com/johnny-morrice/amoebethics/libamoebethics"
)

type PreFrame struct {
    core.SimBase
    Palette []Color
}

type Frame struct {
    FrameNum uint
    PreFrame
    Explosions []Explosion
    Shapes map[string][]ColorBox
}

func MakeFrame(pre PreFrame) Frame {
    fr := Frame{PreFrame: pre}
    fr.Explosions = []Explosion {}
    fr.Shapes = map[string][]ColorBox {
        "circle": []ColorBox {},
        "square": []ColorBox {},
    }

    return fr
}

func WriteFrame(fr Frame, w io.Writer) error {
    enc := json.NewEncoder(w)
    return enc.Encode(fr)
}

func ReadFrame(r io.Reader) (Frame, error) {
    dec := json.NewDecoder(r)
    fr := Frame{}
    err := dec.Decode(&fr)
    return fr, err
}

type Explosion struct {
    Radius float64
    Color int
    Intensity uint8
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
var __defaultpalette []Color

func init() {
    prime := []Color{red, green, blue, purple, yellow, turqouise,}
    __defaultpalette = make([]Color, 2 * len(prime))
    j := 0
    for _, light := range prime {
        dark := Color{}
        dark.R = light.R / 2
        dark.G = light.G / 2
        dark.B = light.B / 2
        __defaultpalette[j] = light
        __defaultpalette[j + 1] = dark
        j += 2
    }
}

func belief2color(bid int, op core.Opinion) int {
    offset := 0
    if op == core.IsFalse {
        offset = 1
    }
    return (bid * 2) + offset
}

type Renderer struct {
    pre PreFrame
    framecnt uint
    fr *Frame
    nodes []*core.SimNode
    slot float64
    time float64
    entshapes map[string]string
    entgrps map[string]EntGroupFact
    pkt core.SimPacket
}

type RenderFactory struct {
    Yard core.EntityYard
    Framecnt uint
    EntShapes map[string]string
    EntGroups map[string]EntGroupFact
}

func (fact RenderFactory) Build(pkt core.SimPacket, palette []Color) (Renderer, error) {
    r := Renderer{}

    if palette == nil {
        palette = __defaultpalette
    }

    support := len(palette) / 2
    if len(pkt.Beliefs) > support {
        return r, fmt.Errorf("Only supports %v beliefs", support)
    }

    r.pre.SimBase = pkt.SimBase
    r.pre.Palette = palette
    r.pkt = pkt
    seqmax := float64(fact.Framecnt)
    r.slot = 1.0 / seqmax
    r.framecnt = fact.Framecnt
    r.entshapes = fact.EntShapes
    r.entgrps = fact.EntGroups

    r.nodes = make([]*core.SimNode, len(r.pkt.Nodes))
    for i, un := range r.pkt.Nodes {
        n, err := core.MakeNode(&un, i, r.pkt.SimBase, fact.Yard)
        if err != nil {
            return r, err
        }
        r.nodes[i] = n
    }

    return r, nil
}

func (r *Renderer) Render() []Frame {
    out := make([]Frame, r.framecnt)
    for i := uint(0); i < r.framecnt; i++ {
        r.interpolate()
        entGroups := r.nodeGroups()
        fr := MakeFrame(r.pre)
        fr.FrameNum = i
        r.fr = &fr
        for _, group := range entGroups {
            group.Render()
        }
        out[i] = fr
        r.time += r.slot
    }

    return out
}

func (r *Renderer) interpolate() {
    for _, n := range r.nodes {
        n.Interpolate(r.time)
    }
}