package main

type CompID uint32

func (c *CompID) Set(typ ...CompType) {
	for _, t := range typ {
		if t >= CompCount {
			continue
		}
		(*c) |= 1 << CompID(t)
	}
}

func (c *CompID) Unset(typ ...CompType) {
	for _, t := range typ {
		if t >= CompCount {
			continue
		}
		(*c) &^= 1 << CompID(t)
	}
}

func (c *CompID) Match(oc CompID) bool {
	return *c&oc == *c
}

func (c *CompID) Less(oc CompID) bool {
	return *c < oc
}

type CompType uint32

const (
	CompHP CompType = iota
	CompAttack
	CompMagic
	CompCount
)

type CompHPData struct {
	ComponentBase

	Hp int
}

type CompAttackData struct {
	ComponentBase

	Attack int
}

type CompMagicData struct {
	ComponentBase

	Magic int
}

func init() {
	RegComponentType[*CompAttackData](CompAttack)
	RegComponentType[*CompHPData](CompHP)
	RegComponentType[*CompMagicData](CompMagic)
}
