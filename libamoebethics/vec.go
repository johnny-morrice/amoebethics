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

type screen struct {
    w uint
    h uint
}

// TorusScreen provides a mapping between a toroidal space and
// the screen.
type TorusScreen struct {
    screen
    t Torus
}

func MakeTorusScreen(t Torus, w, h uint) TorusScreen {
    ts := TorusScreen{}
    ts.t = t
    ts.w = w
    ts.h = h
    return ts
}

func (ts TorusScreen) pixelSize() (float64, float64) {
    fw := float64(ts.w)
    fh := float64(ts.h)

    return fw / ts.t.W, fh / ts.t.H
}

// Project a point on the torus onto the screen.
func (ts TorusScreen) Project(v *mat64.Vector) (uint, uint) {
    xUnit, yUnit := ts.pixelSize()

    reflectComps := []float64{
        1, 0,
        0, -1,
    }
    reflect := mat64.NewDense(2, 2, reflectComps)

    trans := Vec2(float64(ts.t.W) / 2.0, float64(ts.t.H) / 2.0)

    // Scaling matrix
    scaleComps := []float64{
        xUnit, 0,
        0, yUnit,
    }
    scale := mat64.NewDense(2, 2, scaleComps)

    pr := Vec2(0, 0)
    pr.MulVec(reflect, v)
    pr.AddVec(pr, trans)
    pr.MulVec(scale, pr)

    rx := uint(math.Floor(pr.At(0, 0)))
    ry := uint(math.Floor(pr.At(1, 0)))

    return rx, ry
}

// TODO: horrible
func (ts TorusScreen) Scale(r float64) float64 {
    return (float64(ts.w) / ts.t.W) * r
}

// Torus is a rectangular manifold on a toroidal surface.
type Torus struct {
    W float64
    H float64
}

// Map produces a vector that is within the bounds of the
// rectangular manifold of toroidal space, given a vector
// that is on the torus but may be outside these bounds.
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

// Explodes returns true when pos is within radius of an explosion at centre,
// with respect to toroidal topology.
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

// Projections produces a list of alternative positions on the
// rectangular manifold of toroidal space.
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
