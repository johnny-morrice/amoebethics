# Amoebethics
## the toy epidemiology of political belief
---

### Installation

    $ go get github.com/johnny-morrice/amoebethics
    $ go install github.com/johnny-morrice/amoebethics...

---

### Node types

* Sheeple (fully impressionable, random walk)
* Shy (Expresses only when expressed to)
* Conservative (Fixed belief, limited expression)
* Activist (Very mobile, high expression)
* Contrarian (Expresses opposite of received expressions)
* Politician (Express popular beliefs)
* Rebel (Expresses unpopular belief)
* Chugger (Activist with repulsion field)
* Celebrity (Activist with attraction field)

#### Other:
* Media

---

### TODO Extensions

* Taboos and dog whistles. Cannot express belief x; allow formation of "believe x, express d" pairs
* Impressionability. Some nodes follow, others resist, others negate.
* Simulation of age groups: younger is more impressionable and more extreme
* Endogenous Media: Murdoch-nodes with huge radii.
* Violence. Vicious cycle of polarisation
* Internal inconsistency. Both A & ~A.
* Sim environment variables: poverty, perceived threat to beliefs (from tv)
* Parliament.
* What are the attractors? Do two tribes emerge? Two parties?