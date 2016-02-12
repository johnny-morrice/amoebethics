package libamoebethics

import (
    "github.com/gonum/matrix/mat64"
    m "math"
)

type Torus struct {
    W float64
    H float64
}

type Vec2 struct {
    X float64
    Y float64
}

var _ mat64.Matrix = Vec2{}

func (v Vec2) Dims() (r, c int) {
    return 2, 1
}

func (v Vec2) At(i, j int) float64 {
    if j != 0 {
        panic("Vec2 column out of bounds")
    }
    switch i {
    case 0:
        return v.X
    case 1:
        return v.Y
    default:
        panic("Vec2 row out of bounds")
        return 0.0 // Should never happen
    }
}

func (v Vec2) T() mat64.Matrix {
    return mat64.Transpose{Matrix: v}
}

func (t Torus) explodes(center *mat64.Dense, radius float64, pos *mat64.Dense) bool {
    diff := mat64.DenseCopyOf(center)
    for _, p := range t.projections(pos) {
        diff.Sub(center, p)
        if mat64.Det(diff) < radius {
            return true
        }
    }
    return false
}

func (t Torus) projections(pos *mat64.Dense) []*mat64.Dense {
    right := Vec2{t.W, 0.0}
    top := Vec2{0.0, t.H}
    left := Vec2{-t.W, 0.0}
    bottom := Vec2{0.0, -t.H}

    const sq = 3
    const pcnt = sq * sq
    pro := make([]*mat64.Dense, pcnt)
    zero := Vec2{0.0, 0.0}
    for w := 0; w < sq; w++ {
        var vw Vec2
        if w == 0 {
            vw = left
        } else if w == 2 {
            vw = right
        } else {
            vw = zero
        }
        for h := 0; h < sq; h++ {
            i := w + (h * sq)
            p := mat64.DenseCopyOf(vw)
            pro[i] = p

            if h == 0 {
                p.Add(p, top)
            } else if h == 2 {
                p.Add(p, bottom)
            }
        }
    }

    return pro
}

func steer(unit *mat64.Dense, theta float64) *mat64.Dense {
    mat := []float64 {
        m.Cos(theta),
        -m.Sin(theta),
        m.Sin(theta),
        m.Cos(theta),
    }
    rot := mat64.NewDense(2, 2, mat)
    turned := mat64.DenseCopyOf(unit)
    turned.Mul(rot, unit)
    return turned
}