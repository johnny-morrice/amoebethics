package libamoebethics

import (
    "github.com/gonum/matrix/mat64"
    m "math"
)

func steer(unit *mat64.Dense, theta float64) *mat64.Dense {
    mat := []float64 {
        m.Cos(theta), -m.Sin(theta), m.Sin(theta), m.Cos(theta)
    }
    rot := mat64.NewDense(2, 2, mat)
    turned := mat64.DenseCopyOf(unit)
    turned.Mul(rot, unit)
    return turned
}