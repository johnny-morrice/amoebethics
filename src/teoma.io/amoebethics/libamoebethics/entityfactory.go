package libamoebethics

type EntityFactory interface {
     Build(un *UserNode, t Torus) Entity
}

func init() {
    // Add sheeple + tv + activists to __entityFact
    __entityFact = make(map[string]EntityFactory)
    __entityFact["sheeple"] = SheepleFactory{}
    __entityFact["tv"] = TvFactory{}
}

var __entityFact map[string]EntityFactory

func KnownNodeName(name string) bool {
    _, ok := __entityFact[name]
    return ok
}

func MakeEntity(un *UserNode, t Torus) Entity {
    fact, ok := __entityFact[un.name]
    if !ok {
        panic("Unknown node: " + un.name)
    }
    return fact.Build(un, t)
}

type SheepleFactory struct {}

var _ EntityFactory = SheepleFactory{}

func (sf SheepleFactory) Build(un *UserNode, t Torus) Entity {
    sh := Sheeple{}
    return &sh
}

type TvFactory struct {}

var _ EntityFactory = TvFactory{}

func (sf TvFactory) Build(un *UserNode, t Torus) Entity {
    tv := Tv{}
    return &tv
}

