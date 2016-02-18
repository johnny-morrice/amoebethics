package amoebext

import (
    "strings"
    lib "teoma.io/amoebethics/libamoebethics"
)

func decodeEntity(e lib.Entity, un *lib.UserNode) error {
    return e.Deserialize(strings.NewReader(un.Extension))
}