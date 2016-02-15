package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
    "sync"
    core "teoma.io/amoebethics/libamoebethics"
    animate "teoma.io/amoebethics/libamoenator"
    ext "teoma.io/amoebethics/amoebext"
)

func main() {
    yard := ext.StdExtensions()

    count := uint(30)
    flag.UintVar(&count, "framecnt", 30, "Frames generated per simulation packet")
    flag.Parse()

    frtun := runFrameTunnel(yard, count)

    for frch := range frtun {
        for fr := range frch {
            werr := animate.WriteFrame(fr, os.Stdout)
            if werr != nil {
                fatal(werr)
            }
        }
    }
}

func runFrameTunnel(yard core.EntityYard, count uint) <-chan chan animate.Frame {
    outch := make(chan chan animate.Frame)
    go func() {
        hold := sync.WaitGroup{}
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            r := strings.NewReader(scanner.Text())
            pkt, perr := core.ReadSimPkt(r)
            if perr != nil {
                fatal(perr)
            }
            hold.Add(1)
            go func() {
                frch := render(pkt, yard, count)
                outch<- frch
                hold.Done()
            }()
        }
        hold.Wait()
        close(outch)
    }()
    return outch
}

func render(pkt core.SimPacket, yard core.EntityYard, count uint) chan animate.Frame {
    frch := make(chan animate.Frame)
    rend, rerr := animate.MakeRenderer(pkt, yard, count)
    if rerr != nil {
        fatal(rerr)
    }
    go func() {
        for _, fr := range rend.Render() {
            frch<- fr
        }
        close(frch)
    }()

    return frch
}

func fatal(e error) {
    fmt.Fprintf(os.Stderr, "Fatal: %v\n", e)
    os.Exit(1)
}