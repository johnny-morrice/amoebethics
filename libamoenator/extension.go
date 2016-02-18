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
    for i, un := range base.user {
        boxes[i] = renderNode(un)
    }
    sheepShape := base.render.entshapes[entity]
    fr := base.render.fr
    fr.Shapes[sheepShape] = append(fr.Shapes[sheepShape], boxes...)
}

func renderNode(node core.UserNode) ColorBox {
    const radius float64 = 0.05
    box := ColorBox{}
    box.Colors = make([]int, len(node.Beliefs))
    for i, b := range node.Beliefs {
        if op := b.Op; op == core.IsTrue || op == core.IsFalse {
            box.Colors[i] = belief2color(b.Id, op)
        } else {
            panic("Unsupported opinion")
        }
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