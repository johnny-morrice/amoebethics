package main

import (
    "os"
    lib "teoma.io/amoebethics/libamoebethics"
    ext "teoma.io/amoebethics/amoebext"
)

func main() {
    input, inerr := lib.ReadSimInput(os.Stdin)
    if inerr != nil {
        panic(inerr)
    }

    extensions := ext.StdExtensions()
    outch, simerr := lib.Simulate(input, extensions)

    if simerr != nil {
        panic(simerr)
    }

    for out := range outch {
        werr := lib.WriteSimOutput(out, os.Stdout)
        if werr != nil {
            panic(werr)
        }
    }
}