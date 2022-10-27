package ecs

import "sort"

type IEntityCb[ID any, I comparable, C CID[ID, I]] func(IEntity[ID, I, C]) bool

type EntityGroup[ID any, I comparable, C CID[ID, I]] struct {
	ComponentID[ID, I, C]

	entities  []IEntity[ID, I, C]
	addCb     []IEntityCb[ID, I, C]
	removeCb  []IEntityCb[ID, I, C]
	replaceCb []IEntityCb[ID, I, C]
}

type IdentitySort[T interface{ GetID() IdentityID }] []T

func (a IdentitySort[T]) Len() int           { return len(a) }
func (a IdentitySort[T]) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IdentitySort[T]) Less(i, j int) bool { return a[i].GetID() > a[j].GetID() }

func NewEntityGroup[ID any, I comparable, C CID[ID, I]](entities map[IdentityID]IEntity[ID, I, C], idx ...I) *EntityGroup[ID, I, C] {
	var group EntityGroup[ID, I, C]
	group.ID = GetComponentID[ID, I, C](idx...)
	for _, entity := range entities {
		if !group.Match(entity.GetComponentID()) {
			continue
		}
		group.Add(entity)
	}
	sort.Sort(IdentitySort[IEntity[ID, I, C]](group.entities)) // ?
	return &group
}

func (et *EntityGroup[ID, I, C]) OnAdd(cb IEntityCb[ID, I, C]) {
	et.addCb = append(et.addCb, cb)
}

func (et *EntityGroup[ID, I, C]) OnRemove(cb IEntityCb[ID, I, C]) {
	et.removeCb = append(et.removeCb, cb)
}

func (et *EntityGroup[ID, I, C]) OnReplace(cb IEntityCb[ID, I, C]) {
	et.replaceCb = append(et.replaceCb, cb)
}

func (et *EntityGroup[ID, I, C]) Foreach(cb IEntityCb[ID, I, C]) {
	for i := len(et.entities) - 1; i >= 0; i-- {
		if !cb(et.entities[i]) {
			break
		}
	}
}

func (et *EntityGroup[ID, I, C]) Add(e IEntity[ID, I, C]) {
	et.entities = append(et.entities, e)
	for _, cb := range et.addCb {
		cb(e)
	}
}

func (et *EntityGroup[ID, I, C]) Remove(entity IEntity[ID, I, C]) {
	var found bool
	for i, e := range et.entities {
		if e.GetID() != entity.GetID() {
			continue
		}
		et.entities[i] = et.entities[len(et.entities)-1]
		et.entities = et.entities[:len(et.entities)-1]
		found = true
		break
	}
	if !found {
		return
	}
	for _, cb := range et.removeCb {
		cb(entity)
	}
}

func (et *EntityGroup[ID, I, C]) Replace(e IEntity[ID, I, C]) {
	//TODO
}
