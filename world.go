package ecs

import "fmt"

type EntityWorld[ID any, I comparable, C CID[ID, I]] struct {
	genID IdentityID

	entities   map[IdentityID]IEntity[ID, I, C]
	entityType map[string]*EntityType[ID, I, C]
	singleton  map[string]any
	systems    []*EntitySystemWrap[ID, I, C]
	groups     map[C]*EntityGroup[ID, I, C]
}

func (ew *EntityWorld[ID, I, C]) RegEntityType(name string) *EntityType[ID, I, C] {
	checkMapInit(&ew.entityType)
	checkMapRepeatAdd(ew.entityType, name)
	var et = NewEntityType[ID, I, C](name)
	ew.entityType[name] = et
	return et
}

func (ew *EntityWorld[ID, I, C]) GetGroup(i ...I) *EntityGroup[ID, I, C] {
	checkMapInit(&ew.groups)
	cid := GetComponentID[ID, I, C](i...)
	g, ok := ew.groups[&cid]
	if !ok {
		g = NewEntityGroup(ew.entities, i...)
		ew.groups[&cid] = g
	}
	return g
}

func (ew *EntityWorld[ID, I, C]) AddSystem(system IEntitySystem[ID, I, C]) {
	var wrap = &EntitySystemWrap[ID, I, C]{}
	wrap.SetInner(system)
	wrap.inner.SetupGroup(ew)
	ew.systems = append(ew.systems, wrap)
}

func (ew *EntityWorld[ID, I, C]) AddSingleton(name string, singleton any) {
	if singleton == nil {
		return
	}
	checkMapInit(&ew.singleton)
	checkMapRepeatAdd(ew.singleton, name)
	ew.singleton[name] = singleton
}

func (ew *EntityWorld[ID, I, C]) Tick() {
	for _, sys := range ew.systems {
		sys.Tick(ew)
	}
}

func (ew *EntityWorld[ID, I, C]) onEnableComponent(e IEntity[ID, I, C], i I) {
	oldID := newComponentID[ID, I, C](e.GetComponentID())
	newID := oldID
	oldID.Unset(i)
	for _, group := range ew.groups {
		if group.ComponentID.Match(newID.ID) && !group.ComponentID.Match(oldID.ID) {
			group.Add(e)
		}
	}
}

func (ew *EntityWorld[ID, I, C]) onDisbaleComponent(e IEntity[ID, I, C], i I) {
	oldID := newComponentID[ID, I, C](e.GetComponentID())
	newID := oldID
	oldID.Set(i)
	for _, group := range ew.groups {
		if group.ComponentID.Match(newID.ID) && !group.ComponentID.Match(oldID.ID) {
			group.Remove(e)
		}
	}
}

func NewWorldEntity[ET IEntity[ID, I, C], ID any, I comparable, C CID[ID, I]](w *EntityWorld[ID, I, C], name string) *WorldEntity[ET, ID, I, C] {
	return &WorldEntity[ET, ID, I, C]{
		world: w,
		name:  name,
	}
}

type WorldEntity[ET IEntity[ID, I, C], ID any, I comparable, C CID[ID, I]] struct {
	world *EntityWorld[ID, I, C]
	name  string
}

func (w *WorldEntity[ET, ID, I, C]) Add() ET {
	return WorldAddEntity[ET](w.world, w.name)
}

func (w *WorldEntity[ET, ID, I, C]) Get(id IdentityID) ET {
	return WorldGetEntity[ET](w.world, id)
}

func (w *WorldEntity[ET, ID, I, C]) SetWorld(world *EntityWorld[ID, I, C]) {
	w.world = world
}

func WorldAddEntity[ET IEntity[ID, I, C], ID any, I comparable, C CID[ID, I]](ew *EntityWorld[ID, I, C], name string) ET {
	var empty ET
	var et, ok = ew.entityType[name]
	if !ok {
		return empty
	}
	ew.genID++
	id := ew.genID
	checkMapInit(&ew.entities)
	e := EntityCreate[ET](et, id)
	ew.entities[id] = e
	e.SetWorld(ew)
	return e
}

func WorldGetEntity[ET IEntity[ID, I, C], ID any, I comparable, C CID[ID, I]](ew *EntityWorld[ID, I, C], id IdentityID) (e ET) {
	e, _ = ew.entities[id].(ET)
	return
}

func WorldGetSingleton[T any, ID any, I comparable, C CID[ID, I]](ew *EntityWorld[ID, I, C], name string) (t T) {
	t, _ = ew.singleton[name].(T)
	return
}

func checkMapRepeatAdd[K comparable, V any](m map[K]V, k K) {
	if _, ok := m[k]; ok {
		panic(fmt.Sprintf("repeate add %v", k))
	}
}

func checkMapInit[K comparable, V any](m *map[K]V) {
	if m == nil {
		return
	}
	if *m != nil {
		return
	}
	*m = make(map[K]V)
}
