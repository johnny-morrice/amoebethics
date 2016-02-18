package libamoenator

type SheepleGroup BaseExtension

var _ EntGroup = SheepleGroup{}

func MakeSheepleGroup(base BaseExtension) EntGroup {
    return SheepleGroup(base)
}

func (gr SheepleGroup) Render() {
    BaseExtension(gr).DefaultRender("sheeple")
}