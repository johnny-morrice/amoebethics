package libamoenator

type SheepleGroup BaseExtension

var _ EntGroup = SheepleGroup{}

var MakeSheepleGroup = EntGroupFact(func (base BaseExtension) EntGroup {
    return SheepleGroup(base)
})

func (gr SheepleGroup) Render() {
    BaseExtension(gr).DefaultRender("sheeple")
}