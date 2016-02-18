package libamoebethics

type EntityYard map[string]EntityFactory

type EntityFactory interface {
     Build(un *UserNode, base SimBase) (Entity, error)
}

func (yard EntityYard) MakeEntity(un *UserNode, base SimBase) (Entity, error) {
    fact, ok := yard[un.Name]
    if !ok {
        panic("Unknown node type: " + un.Name)
    }
    return fact.Build(un, base)
}
