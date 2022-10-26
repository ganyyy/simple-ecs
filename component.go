package ecs

type IComponent[ID any, I comparable, C CID[ID, I]] interface {
	IEable

	SetID(I)

	GetComponentId() I
	GetOwner() IEntity[ID, I, C]
	SetOwner(IEntity[ID, I, C])
}

func GetComponentID[ID any, I comparable, C CID[ID, I]](idx ...I) ID {
	var t ID
	for _, i := range idx {
		C(&t).Set(i)
	}
	return t
}

type ComponentBase[ID any, I comparable, C CID[ID, I]] struct {
	enableBase
	id    I
	owner IEntity[ID, I, C]
}

func (c *ComponentBase[ID, I, C]) SetID(i I)                        { c.id = i }
func (c *ComponentBase[ID, I, C]) GetComponentId() I                { return c.id }
func (c *ComponentBase[ID, I, C]) GetOwner() IEntity[ID, I, C]      { return c.owner }
func (c *ComponentBase[ID, I, C]) SetOwner(owner IEntity[ID, I, C]) { c.owner = owner }

type componentCtor[ID any, I comparable, C CID[ID, I]] struct {
	Binding binding
	Enable  bool
}
