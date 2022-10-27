package ecs

import (
	"reflect"
	"sync"
)

type IdentityID uint64

type IEntity[ID any, I comparable, C CID[ID, I]] interface {
	IIdentity
	SetID(IdentityID)
	GetID() IdentityID
	GetComponentID() ID

	SetWorld(w *EntityWorld[ID, I, C])
	GetWorld() *EntityWorld[ID, I, C]

	GetComponent(I) IComponent[ID, I, C]
	AddComponent(comp IComponent[ID, I, C])
	IsEnableComponent(I) bool
	EnableComponent(I)
	DisableComponent(I)
}

type EntityBase[ID any, I comparable, C CID[ID, I]] struct {
	identityBase

	world       *EntityWorld[ID, I, C]
	id          IdentityID
	componentID ComponentID[ID, I, C]
	components  map[I]IComponent[ID, I, C]
}

func (e *EntityBase[ID, I, C]) SetWorld(w *EntityWorld[ID, I, C])     { e.world = w }
func (e *EntityBase[ID, I, C]) GetWorld() *EntityWorld[ID, I, C]      { return e.world }
func (e *EntityBase[ID, I, C]) SetID(id IdentityID)                   { e.id = id }
func (e *EntityBase[ID, I, C]) GetID() IdentityID                     { return e.id }
func (e *EntityBase[ID, I, C]) GetComponentID() ID                    { return e.componentID.ID }
func (e *EntityBase[ID, I, C]) GetComponent(i I) IComponent[ID, I, C] { return e.components[i] }

func (e *EntityBase[ID, I, C]) AddComponent(c IComponent[ID, I, C]) {
	checkMapInit(&e.components)
	if c.Enable() {
		e.componentID.Set(c.GetComponentId())
	}
	e.components[c.GetComponentId()] = c
}

func (e *EntityBase[ID, I, C]) IsEnableComponent(i I) bool {
	c := e.GetComponent(i)
	if c == nil {
		return false
	}
	return c.Enable()
}

func (e *EntityBase[ID, I, C]) EnableComponent(i I) {
	c := e.GetComponent(i)
	if c == nil || c.Enable() {
		return
	}
	c.SetEnable(true)
	e.componentID.Set(c.GetComponentId())
	e.GetWorld().onEnableComponent(e, i)
}

func (e *EntityBase[ID, I, C]) DisableComponent(i I) {
	c := e.GetComponent(i)
	if c == nil || !c.Enable() {
		return
	}
	c.SetEnable(false)
	e.componentID.Unset(i)
	e.GetWorld().onDisbaleComponent(e, i)
}

type EntityType[ID any, I comparable, C CID[ID, I]] struct {
	nameBase
	ComponentID[ID, I, C]
	componentCtors []componentCtor[ID, I, C]
	Pool           sync.Pool //(?)
}

func NewEntityType[ID any, I comparable, C CID[ID, I]](name string) *EntityType[ID, I, C] {
	return &EntityType[ID, I, C]{nameBase: nameBase(name)}
}

func (e *EntityType[ID, I, C]) Define(id I, enable bool) *EntityType[ID, I, C] {
	var setType ComponentID[ID, I, C]
	setType.Set(id)
	if setType.Match(e.ComponentID.ID) {
		return e
	}
	if typ, ok := GetComponentBinding[ID, I, C](id); !ok {
		return e
	} else {
		var ctor = componentCtor[ID, I, C]{
			Binding: typ,
			Enable:  enable,
		}
		e.componentCtors = append(e.componentCtors, ctor)
		e.ComponentID.Set(id)
	}
	return e
}

func EntityCreate[ET IEntity[ID, I, C], ID any, I comparable, C CID[ID, I]](typ *EntityType[ID, I, C], id IdentityID) ET {
	var empty ET
	entity := reflect.New(reflect.TypeOf(empty).Elem()).Interface().(ET)
	entity.SetID(id)
	for _, ctor := range typ.componentCtors {
		var c, ok = NewComponent[ID, I, C](ctor.Binding)
		if !ok {
			continue
		}
		c.SetEnable(ctor.Enable)
		c.SetOwner(entity)
		entity.AddComponent(c)
	}
	return entity
}

func GetEntityComponent[CT IComponent[ID, I, C], ID any, I comparable, C CID[ID, I], ET IEntity[ID, I, C]](e ET) (ret CT, valid bool) {
	id, ok := allComponentToType[bindingKey[CT, ID, I, C]()]
	if !ok {
		return
	}
	i, ok := id.(I)
	if !ok {
		return
	}
	return GetEntityComponentID[CT, ID, I, C](e, i)
}

func GetEntityComponentID[CT IComponent[ID, I, C], ID any, I comparable, C CID[ID, I], ET IEntity[ID, I, C]](e ET, i I) (ret CT, valid bool) {
	c := e.GetComponent(i)
	if c == nil {
		return
	}
	if v, ok := c.(CT); ok {
		return v, true
	}
	return
}
