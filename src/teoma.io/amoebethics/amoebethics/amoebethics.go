package main

import (
    "fmt"
    "os"
    lib "teoma.io/amoebethics/libamoebethics"
    ext "teoma.io/amoebethics/amoebext"
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
            fatal(werr)
        }
    }
}

func fatal(e error) {
    fmt.Fprintf(os.Stderr, "Fatal: %v", e)
    os.Exit(1)
}