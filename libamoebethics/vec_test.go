package libamoebethics

import (
    "testing"
    "github.com/gonum/matrix/mat64"
)

func TestMap(t *testing.T) {
    tor := Torus{10, 10}
    input := []*mat64.Vector{
        Vec2(-10, 10),
        Vec2(0, 10),
        Vec2(10, 10),
        Vec2(-10, 0),
        Vec2(10, 0),
        Vec2(-10, -10),
        Vec2(0, -10),
        Vec2(10, -10),
    }
    expect := Vec2(0.0, 0.0)
    for i, actual := range input {
        tor.Map(actual)
        for j := 0; j < 2; j++ {
            elem := actual.At(j, 0)
            if elem != expect.At(j, 0) {
                t.Error("Actual", i, "differed from expected at element", j, ":", elem)
            }
        }
    }
}

func TestProjections(t *testing.T) {
    input := Vec2(0.0, 0.0)
    exptab := []*mat64.Vector{
        Vec2(-10, 10),
        Vec2(0, 10),
        Vec2(10, 10),
        Vec2(-10, 0),
        Vec2(0, 0),
        Vec2(10, 0),
        Vec2(-10, -10),
        Vec2(0, -10),
        Vec2(10, -10),
    }
    tor := Torus{10, 10}
    out := tor.Projections(input)
    for i, expect := range exptab {
        actual := out[i]
        for j := 0; j < 2; j++ {
            elem := actual.At(j, 0)
            if elem != expect.At(j, 0) {
                t.Error("Actual", i, "differed from expected at element", j, ":", elem)
            }
        }
    }
}

func TestExplodes(t *testing.T) {
    explodeTest(t,  3, Vec2(0, 0), Vec2(2, 0), Vec2(4, 0))
    explodeTest(t, 3, Vec2(-4, 0), Vec2(5, 0), Vec2(0, 0))
}

func explodeTest(t *testing.T, radius float64, center, inPos, outPos *mat64.Vector) {
    tor := Torus{10, 10}
    if !tor.Explodes(radius, center, inPos) {
        t.Error("Erroneously not exploded")
    }

    if tor.Explodes(radius, center, outPos) {
        t.Error("Erroneously exploded")
    }
}