package amoebext

import (
    "bytes"
    lib "teoma.io/amoebethics/libamoebethics"
)

func decodeEntity(e lib.Entity, un *lib.UserNode) error {
    buff := bytes.NewBufferString(un.Extension)
    return e.Deserialize(buff)
}