package libamoebethics

type Tv struct {

}

var _ Entity = (*Tv)(nil)

func (tv *Tv) Handle(m *SimNode) {

}

func (tv *Tv) Greet(m *SimNode) bool {
    return false
}

func (tv *Tv) Serialize() string {
    return ""
}