package libamoenator

import (
    core "github.com/johnny-morrice/amoebethics/libamoebethics"
)

type EntGroup interface {
    Render()
}

type EntGroupFact func(b BaseExtension) EntGroup

type BaseExtension struct {
    render *Renderer
    nodes []*core.SimNode
    user []core.UserNode
}

func (be BaseExtension) noderad() float64 {
    const baserad float64 = 0.5
    return be.render.torscreen.Scale(baserad)
}

func (r *Renderer) nodeGroups() []EntGroup {
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
    out := make([]EntGroup, len(split))
    i := 0
    for name, slp := range split {
        base := BaseExtension{}
        base.render = r
        base.nodes = *slp
        base.user = r.userNodes(base.nodes)
        fact := r.entgrps[name]
        out[i] = fact(base)
        i++
    }
    return out
}

func (r *Renderer) userNodes(nodes []*core.SimNode) []core.UserNode {
    out := make([]core.UserNode, len(nodes))
    for i, n := range nodes {
        out[i] = core.SimNode2UserNode(n)
    }
    return out
}

func (base BaseExtension) DefaultRender(entity string) {
    boxes := make([]ColorBox, len(base.user))
    for i, sn := range base.nodes {
        mapped := core.MapNode(base.render.torscreen, sn)
        un := core.SimNode2UserNode(mapped)
        boxes[i] = renderNode(un, base.noderad())
    }
    sheepShape := base.render.entshapes[entity]
    fr := base.render.fr
    fr.Shapes[sheepShape] = append(fr.Shapes[sheepShape], boxes...)
}

func renderNode(node core.UserNode, radius float64) ColorBox {
    box := ColorBox{}
    box.Colors = make([]Coldex, len(node.Beliefs))
    for i, b := range node.Beliefs {
        if op := b.Op; op == core.IsTrue || op == core.IsFalse {
            box.Colors[i] = belief2color(b)
        } else {
            panic("Unsupported opinion")
        }
    }
    box.Radius = radius
    box.P = node.P
    return box
}

func renderExplosion(p core.UserVec, radius, time float64, b core.Belief) Explosion {
    const max = ^uint8(0)
    const fmax = float64(max)
    ex := Explosion{}
    ex.Radius = radius * (time * 2)
    if time < 0.5 {
        ex.Intensity = max
    } else {
        ex.Intensity = uint8(fmax - ((time - 0.5) * 2))
    }
    ex.Color = belief2color(b)
    ex.P = p
    return ex
}

