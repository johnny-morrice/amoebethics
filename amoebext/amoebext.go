package amoebext

import (
    "strings"
    lib "github.com/johnny-morrice/amoebethics/libamoebethics"
)

func decodeEntity(e lib.Entity, un *lib.UserNode) error {
    return e.Deserialize(strings.NewReader(un.Extension))
}