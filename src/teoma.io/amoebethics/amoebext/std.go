package amoebext

import (
    lib "teoma.io/amoebethics/libamoebethics"
)

func StdExtensions() lib.EntityYard {
    m := map[string]lib.EntityFactory{
        "sheeple": SheepleFactory{},
        "tv": TvFactory{},
    }
    return lib.EntityYard(m)
}