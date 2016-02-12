package libamoebethics

type Sheeple struct {

}

var _ Entity = (*Sheeple)(nil)

func (s *Sheeple) Handle(m *SimNode) {

}

func (s *Sheeple) Greet(m *SimNode) bool {
    return false
}

func (s *Sheeple) Serialize() string {
    return ""
}