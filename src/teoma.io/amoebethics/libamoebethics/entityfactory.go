package libamoebethics

type EntityYard map[string]EntityFactory

type EntityFactory interface {
     Build(un *UserNode, t Torus) (Entity, error)
}

func (yard EntityYard) MakeEntity(un *UserNode, t Torus) (Entity, error) {
    fact, ok := yard[un.name]
    if !ok {
        panic("Unknown node type: " + un.name)
    }
    return fact.Build(un, t)
}
