package libamoenator

import (
    ext "teoma.io/amoebethics/amoebext"
)

type TvGroup BaseExtension

var _ EntGroup = TvGroup{}

var MakeTvGroup = EntGroupFact(func (base BaseExtension) EntGroup {
    return TvGroup(base)
})

func (gr TvGroup) Render() {
    boxes := make([]ColorBox, len(gr.user))
    fr := gr.render.fr
    for i, un := range gr.user {
        boxes[i] = renderNode(un)
        e := gr.nodes[i].Entity
        tv := e.(*ext.Tv) // A better deserialization mechanism would allow no typecast
        expls := make([]Explosion, len(un.Expression))
        for i, b := range un.Expression {
            expls[i] = renderExplosion(tv.R, gr.render.time, b)
        }
        fr.Explosions = append(fr.Explosions, expls...)
    }
    tvShape := gr.render.entshapes["tv"]

    fr.Shapes[tvShape] = append(fr.Shapes[tvShape], boxes...)
}