package main

import (
    "os"
    lib "teoma.io/amboethics/libamboethics"
)

func main() {
    input, inerr := lib.ReadSimInput(os.Stdin)
    if inerr != nil {
        panic(inerr)
    }

    outch, simerr := lib.Simulate(input)

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