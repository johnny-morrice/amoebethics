package amoebext

import (
    m "math"
    "math/rand"
    "github.com/gonum/matrix/mat64"
    lib "teoma.io/amoebethics/libamoebethics"

)

type UserRudder struct {
    D lib.UserVec
    S float64
}

func UserRudder2Rudder(ur UserRudder) Rudder {
    r := Rudder{}
    r.D = lib.UserVec2BlasVec(ur.D)
    r.S = ur.S
    return r
}

func Rudder2UserRudder(r Rudder) UserRudder {
    u := UserRudder{}
    u.D = lib.BlasVec2UserVec(r.D)
    u.S = r.S
    return u
}

type Rudder struct {
    D *mat64.Vector
    S float64
}

func (r Rudder) Move(pos *mat64.Vector) {
    pos.AddScaledVec(pos, r.S, r.D)
}

func (r Rudder) Steer(theta float64) {
    mat := []float64 {
        m.Cos(theta),
        -m.Sin(theta),
        m.Sin(theta),
        m.Cos(theta),
    }
    rot := mat64.NewDense(2, 2, mat)
    d := r.D
    d.MulVec(rot, d)
}

func (r Rudder) RandSteer() {
    r.Steer(rand.Float64())
}