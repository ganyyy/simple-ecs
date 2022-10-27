package main

import (
	"math/rand"
	"testing"
)

func TestWorld(t *testing.T) {
	var w EntityWorld

	w.RegEntityType(EntityPlayer).
		Define(CompAttack, true).
		Define(CompMagic, true).
		Define(CompHP, false)

	wp := NewWorldEntity[*Entity](&w, EntityPlayer)

	wp.SetWorld(&w)

	et := wp.Add()
	et.Name = "123"
	wp.Add()
	wp.Add()
	wp.Add()
	wp.Add()
	wp.Add()
	wp.Add()
	wp.Add()

	addLogCB := func(g *EntityGroup, addCb func(IEntityCb), add string) {
		addCb(func(entity IEntity) bool {
			t.Logf("[%v] %s %v", g.ID, add, entity.GetID())
			return true
		})
	}

	addRemove := func(g *EntityGroup) {
		addLogCB(g, g.OnRemove, "Remove")
	}

	addAdd := func(g *EntityGroup) {
		addLogCB(g, g.OnAdd, "Add")
	}

	addLog := func(g *EntityGroup) {
		addRemove(g)
		addAdd(g)
	}

	logEntity := func(g *EntityGroup, e IEntity) {
		t.Logf("[%v] %v", g.ID, e.GetID())
	}

	g := w.GetGroup(CompAttack)
	addLog(g)
	g2 := w.GetGroup(CompAttack, CompMagic)
	addLog(g2)
	g3 := w.GetGroup(CompHP)
	addLog(g3)
	g4 := w.GetGroup(CompMagic, CompHP)
	addLog(g4)

	g.Foreach(func(entity IEntity) bool {
		attack, ok := GetComponent[*CompAttackData](entity)
		if !ok {
			return true
		}
		attack.Attack += 100
		opt := rand.Intn(3)
		if opt == 1 {
			entity.EnableComponent(CompHP)
			hp, ok := GetComponent[*CompHPData](entity)
			if ok {
				hp.Hp += 50
			}
		} else if opt == 2 {
			entity.DisableComponent(CompMagic)
		} else {
			entity.DisableComponent(CompAttack)
		}
		logEntity(g, entity)
		return true
	})

	g.Foreach(func(i IEntity) bool {
		logEntity(g, i)
		return true
	})

	g2.Foreach(func(entity IEntity) bool {
		attack, ok := GetComponent[*CompAttackData](entity)
		if !ok {
			return true
		}
		t.Logf("id %v attack:%v", entity.GetID(), attack.Attack)
		return true
	})

	g3.Foreach(func(entity IEntity) bool {
		hp, ok := GetComponent[*CompHPData](entity)
		if ok {
			t.Logf("ID %v HP %v", entity.GetID(), hp.Hp)
		}
		return true
	})

	g4.Foreach(func(entity IEntity) bool {
		t.Logf("info %+v", entity.GetID())
		return true
	})
}
