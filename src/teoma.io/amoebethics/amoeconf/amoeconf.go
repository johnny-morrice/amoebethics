package main

import (
    "flag"
    "math/rand"
    "os"
    "strings"

    lib "teoma.io/amoebethics/libamoebethics"
    ext "teoma.io/amoebethics/amoebext"
)

type args struct {
    sheeplecnt uint
    tvcnt uint
    itermax uint
    width float64
    height float64
    beliefs string
}

func getArgs() args {
    params := args{}
    flag.UintVar(&params.sheeplecnt, "sheeple", 0, "Number of sheeple")
    flag.UintVar(&params.tvcnt, "tv", 0, "Number of TVs")
    flag.UintVar(&params.itermax, "iterations", 100, "Number of iterations")
    flag.Float64Var(&params.width, "width", 10, "Torus width")
    flag.Float64Var(&params.height, "height", 10, "Torus height")
    flag.StringVar(&params.beliefs, "beliefs", "A,B,C", "Comma separated belief list")
    flag.Parse()
    return params
}

func main() {
    params := getArgs()

    nodecnt := params.sheeplecnt + params.tvcnt
    pkt := lib.SimPacket{}
    pkt.Itermax = int(params.itermax) // TODO overflow check
    pkt.Torus.W = params.width
    pkt.Torus.H = params.height
    pkt.Beliefs = strings.Split(params.beliefs, ",")
    pkt.Nodes = make([]lib.UserNode, nodecnt)

    for i := uint(0); i < nodecnt; i++ {
        pkt.Nodes[i] = randNode(pkt.Torus)
    }

    for i := uint(0); i < params.tvcnt; i++ {
        n := &pkt.Nodes[i]
        n.Name = "tv"
        n.Extension = lib.Entity2String(randTv())
    }

    for i := params.tvcnt; i < nodecnt; i++ {
        n := &pkt.Nodes[i]
        n.Name = "sheeple"
        n.Extension = lib.Entity2String(randSheeple())
    }

    lib.WriteSimPkt(pkt, os.Stdout)
}

func randNode(t lib.Torus) lib.UserNode {
    un := lib.UserNode{}
    un.Neighbours = []int {}
    un.P = randPlace(t)
    un.Beliefs = []lib.Belief {}
    un.Expression = []lib.Belief {}
    return un
}

func randTv() *ext.Tv {
    tv := &ext.Tv{}
    tv.R = 3
    tv.InvF = 3
    return tv
}

func randSheeple() *ext.Sheeple {
    sh := &ext.Sheeple{}
    sh.D = lib.Vec2(1.0, 0)
    sh.S = 1.0
    sh.RandSteer()
    return sh
}

func randPlace(t lib.Torus) lib.UserVec {
    x := tmap(rand.Float64(), t.W)
    y := tmap(rand.Float64(), t.H)
    return lib.UserVec{x, y}
}

func tmap(f float64, dim float64) float64 {
    scale := f * dim
    trans := scale - (dim / 2.0)
    return trans
}