package amoebext

import (
    lib "github.com/johnny-morrice/amoebethics/libamoebethics"
)

func StdExtensions() lib.EntityYard {
    m := map[string]lib.EntityFactory{
        "sheeple": SheepleFactory{},
        "tv": TvFactory{},
    }
    return lib.EntityYard(m)
}