package libamoenator

var StdShapes = map[string]string {
    "sheeple": "circle",
    "tv": "square",
}

var StdGroupFacts = map[string]EntGroupFact {
    "sheeple": MakeSheepleGroup,
    "tv": MakeTvGroup,
}