package libamoebethics

import (
    "github.com/gonum/matrix/mat64"
    m "math"
)



type Torus struct {
    W float64
    H float64
}

func (t Torus) explodes(centre mat64.Dense, radius float64, pos mat64.Dense) bool {
    diff := mat64.DenseCopyOf(centre)

    for _, p := range t.projections() {
        diff.Sub(p, centre)
        if m.Abs(diff.Det()) < radius {
            return true
        }
    }
    return false
}

func (t Torus) projections(pos mat64.Dense) []mat64.Dense {
    // Three-point trick? Transform pos' point along symmetries in 3x3 grid of squares...
}

