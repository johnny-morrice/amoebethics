package libamoenator

import (
    ext "github.com/johnny-morrice/amoebethics/amoebext"
    core "github.com/johnny-morrice/amoebethics/libamoebethics"
)

type TvGroup BaseExtension

var _ EntGroup = TvGroup{}

func MakeTvGroup(base BaseExtension) EntGroup {
    return TvGroup(base)
}

func (gr TvGroup) Render() {
    boxes := make([]ColorBox, len(gr.user))
    fr := gr.render.fr
    for i, un := range gr.user {
        ts := gr.render.torscreen
        mapped := core.MapNode(ts, gr.nodes[i])
        e := mapped.Entity
        tv := e.(*ext.Tv) // A better deserialization mechanism would allow no typecast
        exprad := ts.Scale(tv.R)
        boxes[i] = renderNode(core.SimNode2UserNode(mapped), BaseExtension(gr).noderad())
        expls := make([]Explosion, len(un.Expression))
        for i, b := range un.Expression {
            expls[i] = renderExplosion(core.BlasVec2UserVec(mapped.P), exprad, gr.render.time, b)
        }
        fr.Explosions = append(fr.Explosions, expls...)
    }
    tvShape := gr.render.entshapes["tv"]

    fr.Shapes[tvShape] = append(fr.Shapes[tvShape], boxes...)
}
