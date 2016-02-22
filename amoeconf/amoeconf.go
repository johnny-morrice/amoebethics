package main

import (
    "flag"
    "os"
    "strings"

    lib "github.com/johnny-morrice/amoebethics/libamoebethics"
    ext "github.com/johnny-morrice/amoebethics/amoebext"
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

    for i := uint(0); i < params.tvcnt; i++ {
        n := ext.NewTvNode(pkt.Torus)
        pkt.Nodes[i] = lib.SimNode2UserNode(n)
    }

    shmax := params.tvcnt + params.sheeplecnt
    for i := params.tvcnt; i < shmax; i++ {
        n := ext.NewSheepleNode(pkt.Torus)
        pkt.Nodes[i] = lib.SimNode2UserNode(n)
    }

    for i, n := range pkt.Nodes {
        n.Id = i
    }

    lib.WriteSimPkt(pkt, os.Stdout)
}