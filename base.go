package ecs

type IEable interface {
	Enable() bool
	SetEnable(bool)
}

type enableBase bool

func (e *enableBase) Enable() bool     { return bool(*e) }
func (e *enableBase) SetEnable(b bool) { *e = enableBase(b) }

type IName interface{ GetName() string }

type nameBase string

func (n nameBase) GetName() string { return string(n) }

type IIdentity interface {
	IEable
	IName
}

type identityBase struct {
	enableBase
	nameBase
}

type CID[T any, I comparable] interface {
	*T
	Set(...I)
	Unset(...I)
	Less(T) bool  // 排序用
	Match(T) bool // 对比用
}

func newComponentID[ID any, I comparable, C CID[ID, I]](id ID) ComponentID[ID, I, C] {
	return ComponentID[ID, I, C]{ID: id}
}

type ComponentID[ID any, I comparable, C CID[ID, I]] struct{ ID ID }

func (c *ComponentID[ID, I, C]) Set(idx ...I)     { C(&c.ID).Set(idx...) }
func (c *ComponentID[ID, I, C]) Unset(idx ...I)   { C(&c.ID).Unset(idx...) }
func (c *ComponentID[ID, I, C]) Less(id ID) bool  { return C(&c.ID).Less(id) }
func (c *ComponentID[ID, I, C]) Match(id ID) bool { return C(&c.ID).Match(id) }
