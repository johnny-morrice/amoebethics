package libamoebethics

import (
    "github.com/gonum/matrix/mat64"
)

type Torus struct {
    X int
    Y int
}

func (t Torus) explodes(center mat64.Dense, radius float64, pos mat64.Dense) bool {
    for _, p := range t.projections(pos) {
        diff := mat64.DenseCopyOf(center)
        diff.Sub(center, pos)
        if diff.Det() < radius {
            return true
        }
    }
    return false
}

func (t Torus) projections(pos mat64.Dense) []mat64.Dense {

}