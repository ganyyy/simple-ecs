package main

import (
	"testing"
)

func TestWorld(t *testing.T) {
	var w EntityWorld

	w.RegEntityType(EntityPlayer).
		Define(CompAttack, true).
		Define(CompMagic, true).
		Define(CompHP, false)

	wp := NewWorldEntity[*Entity](&w, EntityPlayer)

	wp.Add()
	wp.Add()
	wp.Add()

	g := w.GetGroup(CompAttack)
	g2 := w.GetGroup(CompAttack, CompMagic)

	g.Foreach(func(i IEntity) bool {
		attack, ok := GetComponent[*CompAttackData](i)
		if !ok {
			return true
		}
		attack.Attack += 100
		return true
	})

	g2.Foreach(func(i IEntity) bool {
		attack, ok := GetComponent[*CompAttackData](i)
		if !ok {
			return true
		}
		t.Logf("id %v attack:%v", i.GetID(), attack.Attack)
		return true
	})
}
