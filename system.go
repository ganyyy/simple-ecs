package ecs

type IEntitySystem[ID any, I comparable, C CID[ID, I]] interface {
	SetupGroup(*EntityWorld[ID, I, C])
	Tick(*EntityWorld[ID, I, C])
	OnEnable(*EntityWorld[ID, I, C])
	OnDisable(*EntityWorld[ID, I, C])
}

type EntitySystemWrap[ID any, I comparable, C CID[ID, I]] struct {
	enableBase
	inner IEntitySystem[ID, I, C]
}

func NewEntitySystem[ID any, I comparable, C CID[ID, I]](system IEntitySystem[ID, I, C]) *EntitySystemWrap[ID, I, C] {
	return &EntitySystemWrap[ID, I, C]{
		inner:      system,
		enableBase: enableBase(true),
	}
}

func (es *EntitySystemWrap[ID, I, C]) SetInner(inner IEntitySystem[ID, I, C]) {
	es.inner = inner
}

func (es *EntitySystemWrap[ID, I, C]) SetupGroup(world *EntityWorld[ID, I, C]) {
	es.inner.SetupGroup(world)
}

func (es *EntitySystemWrap[ID, I, C]) Tick(world *EntityWorld[ID, I, C]) {
	es.inner.Tick(world)
}

func (es *EntitySystemWrap[ID, I, C]) SetEnable(enable bool, world *EntityWorld[ID, I, C]) {
	if es.Enable() == enable {
		return
	}
	es.enableBase.SetEnable(enable)
	if enable {
		es.inner.OnEnable(world)
	} else {
		es.inner.OnDisable(world)
	}
}
