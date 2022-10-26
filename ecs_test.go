package ecs_test

import (
	"testing"
	"unsafe"

	ecs "github.com/ganyyy/simple-ecs"
	"github.com/stretchr/testify/assert"
)

type ComponentType uint8

const (
	None ComponentType = iota
	HP
	Attack
	A2
	A3
	A4
	A5
	A6
	A7
	Count
)

type (
	ComponentBase = ecs.ComponentBase[ComponentID, ComponentType, *ComponentID]
	EntityBase    = ecs.EntityBase[ComponentID, ComponentType, *ComponentID]
	EntityWorld   = ecs.EntityWorld[ComponentID, ComponentType, *ComponentID]
	EntityType    = ecs.EntityType[ComponentID, ComponentType, *ComponentID]
)

type (
	IEntity       = ecs.IEntity[ComponentID, ComponentType, *ComponentID]
	IEntitySystem = ecs.IEntitySystem[ComponentID, ComponentType, *ComponentID]
	IComponent    = ecs.IComponent[ComponentID, ComponentType, *ComponentID]
)

var (
	GetComponentBinding = ecs.GetComponentBinding[ComponentID, ComponentType, *ComponentID]
	NewComponent        = ecs.NewComponent[ComponentID, ComponentType, *ComponentID]
)

func EntityCreate[ET IEntity](typ *EntityType, id ecs.IdentityID) ET {
	return ecs.EntityCreate[ET](typ, id)
}

func GetComponent[C IComponent](e IEntity) (cc C, valid bool) {
	return ecs.GetEntityComponent[C, ComponentID, ComponentType](e)
}

func GetEntity[E IEntity](w *EntityWorld, id ecs.IdentityID) E {
	return ecs.WorldGetEntity[E](w, id)
}

type Entity struct {
	EntityBase
}

type CompHP struct {
	ComponentBase
	Hp int
}

type CompA struct {
	ComponentBase

	Name string
}

type CompB struct {
	ComponentBase

	Age int
}

type BitBase = uint8

const (
	BitBaseSize = int(unsafe.Sizeof(BitBase(0))) * 8
	BitMapSize  = (int(Count) + BitBaseSize - 1) / BitBaseSize
)

type ComponentID [BitMapSize]BitBase

func (c *ComponentID) Set(ctp ...ComponentType) {
	_ = *c
	for _, typ := range ctp {
		c[typ/ComponentType(BitBaseSize)] |= 1 << BitBase(typ%ComponentType(BitBaseSize))
	}
}

func (c *ComponentID) Match(oc ComponentID) bool {
	_ = *c
	for i, t := range c {
		if oc[i]&t != oc[i] {
			return false
		}
	}
	return true
}

func (c *ComponentID) Less(oc ComponentID) bool {
	_ = *c
	// 奇怪的比较方法
	for i := len(*c) - 1; i >= 0; i-- {
		if c[i] != oc[i] {
			return c[i] < oc[i]
		}
	}
	// 这种属于相等的情况
	return true
}

func (c *ComponentID) Unset(ctp ...ComponentType) {
	_ = *c
	for _, typ := range ctp {
		c[typ/ComponentType(BitBaseSize)] &^= 1 << BitBase(typ%ComponentType(BitBaseSize))
	}
}

func TestLogType(t *testing.T) {

	// var _ ecs.NormalID[ComponentType] = ComponentID{}

	t.Run("base", func(t *testing.T) {
		var id ComponentID
		t.Logf("%+v", id)
		id.Set(HP)
		t.Logf("%+v", id)
		id.Set(A7)
		t.Logf("%+v", id)

		id2 := ecs.GetComponentID[ComponentID](A3, A2)

		t.Logf("%+v", id2)

		t.Log(id.Less(id2))
		t.Log(id.Match(id2))

		id2.Set(HP, A7)
		t.Logf("%+v", id2)
		t.Log(id2.Match(id))
		id2.Unset(HP, A7)
		t.Logf("%+v", id2)
	})

	t.Run("Component", func(t *testing.T) {
		ecs.RegComponentType[*CompHP, ComponentID](HP)
		b, ok := GetComponentBinding(HP)
		assert.True(t, ok)
		comp, ok := NewComponent(b)
		assert.True(t, ok)
		var hp CompHP
		hp.SetID(HP)
		assert.Equal(t, comp.(*CompHP), hp)
	})

	t.Run("Entity", func(t *testing.T) {
		ecs.RegComponentType[*CompA, ComponentID](A2)
		ecs.RegComponentType[*CompB, ComponentID](A3)

		var typ EntityType

		typ.Define(HP, true).Define(A2, true).Define(A3, false)

		e := EntityCreate[*Entity](&typ, 1)
		t.Logf("%+v", e)

		hp, ok := GetComponent[*CompHP](e)
		assert.True(t, ok)
		hp.Hp += 100

	})

	t.Run("World", func(t *testing.T) {

		ecs.RegComponentType[*CompHP, ComponentID](HP)
		ecs.RegComponentType[*CompA, ComponentID](A2)
		ecs.RegComponentType[*CompB, ComponentID](A3)

		w := &EntityWorld{}

		const EntityType = "Player"

		w.RegEntityType(EntityType).
			Define(HP, true).
			Define(A2, true).
			Define(A3, false)

		et := ecs.WorldAddEntity[*Entity](w, EntityType)

		t.Logf("%+v", et)

		ee := GetEntity[*Entity](w, et.GetID())

		assert.Equal(t, et, ee)

		h, _ := GetComponent[*CompHP](et)
		ca, _ := GetComponent[*CompA](et)
		cb, _ := GetComponent[*CompB](et)
		t.Logf("%+v, %+v, %+v", h, ca, cb)
	})
}
