package libamoebethics

import (
    "github.com/gonum/matrix/mat64"
    "math"
)

type UserVec struct {
    X float64
    Y float64
}

func UserVec2BlasVec(v UserVec) *mat64.Vector {
    return Vec2(v.X, v.Y)
}

func BlasVec2UserVec(v *mat64.Vector) UserVec {
    u := UserVec{}
    u.X = v.At(0, 0)
    u.Y = v.At(1, 0)
    return u
}

type Torus struct {
    W float64
    H float64
}

func (t Torus) Map(v *mat64.Vector) {
    x := v.At(0, 0)
    y := v.At(1, 0)

    remx := x
    right := t.W / 2
    if math.Abs(x) > right  {
        remx = math.Mod(t.W, -x)
    }
    remy := y
    top := t.H / 2
    if math.Abs(y) > top {
        remy = math.Mod(t.H, -y)
    }

    v.SetVec(0, remx)
    v.SetVec(1, remy)
}

func (t Torus) Explodes(radius float64, center, pos *mat64.Vector) bool {
    diff := Vec2(0.0, 0.0)
    diff.CloneVec(center)
    for _, p := range t.Projections(pos) {
        diff.SubVec(center, p)
        if mat64.Norm(diff, 2) < radius {
            return true
        }
    }
    return false
}

func (t Torus) Projections(pos *mat64.Vector) []*mat64.Vector {
    right := Vec2(t.W, 0.0)
    top := Vec2(0.0, t.H)
    left := Vec2(-t.W, 0.0)
    bottom := Vec2(0.0, -t.H)

    const sq = 3
    const pcnt = sq * sq
    pro := make([]*mat64.Vector, pcnt)
    zero := Vec2(0.0, 0.0)
    for w := 0; w < sq; w++ {
        var vw *mat64.Vector
        if w == 0 {
            vw = left
        } else if w == 2 {
            vw = right
        } else {
            vw = zero
        }
        for h := 0; h < sq; h++ {
            i := w + (h * sq)
            p := Vec2(0.0, 0.0)
            p.AddVec(vw, pos)
            pro[i] = p

            if h == 0 {
                p.AddVec(p, top)
            } else if h == 2 {
                p.AddVec(p, bottom)
            }
        }
    }

    return pro
}

func Vec2(x, y float64) *mat64.Vector {
    return mat64.NewVector(2, []float64{x, y})
}