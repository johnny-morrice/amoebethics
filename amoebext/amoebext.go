package amoebext

import (
    "math/rand"
    "strings"
    m "github.com/gonum/matrix/mat64"
    lib "github.com/johnny-morrice/amoebethics/libamoebethics"
)

func decodeEntity(e lib.Entity, un *lib.UserNode) error {
    return e.Deserialize(strings.NewReader(un.Extension))
}

func randPlace(t lib.Torus) *m.Vector {
    x := tmap(rand.Float64(), t.W)
    y := tmap(rand.Float64(), t.H)
    return lib.Vec2(x, y)
}

func tmap(f float64, dim float64) float64 {
    scale := f * dim
    trans := scale - (dim / 2.0)
    return trans
}