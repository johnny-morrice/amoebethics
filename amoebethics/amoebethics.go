package main

import (
    "fmt"
    "os"
    lib "github.com/johnny-morrice/amoebethics/libamoebethics"
    ext "github.com/johnny-morrice/amoebethics/amoebext"
)

func main() {
    input, inerr := lib.ReadSimPkt(os.Stdin)
    if inerr != nil {
        fatal(inerr)
    }

    extensions := ext.StdExtensions()
    outch, simerr := lib.Simulate(input, extensions)

    if simerr != nil {
        fatal(simerr)
    }

    for out := range outch {
        werr := lib.WriteSimPkt(out, os.Stdout)
        if werr != nil {
            break
        }
    }
}

func fatal(e error) {
    fmt.Fprintf(os.Stderr, "Fatal: %v\n", e)
    os.Exit(1)
}