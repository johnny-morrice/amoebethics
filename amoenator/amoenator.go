package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
    core "github.com/johnny-morrice/amoebethics/libamoebethics"
    animate "github.com/johnny-morrice/amoebethics/libamoenator"
    ext "github.com/johnny-morrice/amoebethics/amoebext"
)

func main() {
    count := uint(30)
    width := uint(1920)
    height := uint(1080)
    flag.UintVar(&count, "framecnt", count, "Frames generated per simulation packet")
    flag.UintVar(&width, "width", width, "Width of output surface")
    flag.UintVar(&height, "height", height, "Height of output surface")
    flag.Parse()

    yard := ext.StdExtensions()
    shapes := animate.StdShapes
    groups := animate.StdGroupFacts

    fact := animate.RenderFactory{}
    fact.Width = width
    fact.Height = height
    fact.Yard = yard
    fact.Framecnt = count
    fact.EntShapes = shapes
    fact.EntGroups = groups
    frtun := runFrameTunnel(fact)

    for frch := range frtun {
        for fr := range frch {
            werr := animate.WriteFrame(fr, os.Stdout)
            if werr != nil {
                fmt.Fprintf(os.Stderr, "Warning: %v", werr)
            }
        }
    }
}

func runFrameTunnel(fact animate.RenderFactory) <-chan chan animate.Frame {
    outch := make(chan chan animate.Frame)
    go func() {
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            r := strings.NewReader(scanner.Text())
            pkt, perr := core.ReadSimPkt(r)
            if perr != nil {
                fatal(perr)
            }
            frch := make(chan animate.Frame)
            go func() {
                render(frch, fact, pkt)
            }()
            outch<- frch
        }
        close(outch)
    }()
    return outch
}

func render(frch chan<- animate.Frame, fact animate.RenderFactory, pkt core.SimPacket) {
    rend, rerr := fact.Build(pkt, nil)
    if rerr != nil {
        fatal(rerr)
    }

    go func() {
        for _, fr := range rend.Render() {
            frch<- fr
        }
        close(frch)
    }()
}

func fatal(e error) {
    fmt.Fprintf(os.Stderr, "Fatal: %v\n", e)
    os.Exit(1)
}
