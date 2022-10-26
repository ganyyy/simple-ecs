package main

import ecs "github.com/ganyyy/simple-ecs"

// 基类

type (
	ComponentBase = ecs.ComponentBase[CompID, CompType, *CompID]
	EntityBase    = ecs.EntityBase[CompID, CompType, *CompID]
	EntityType    = ecs.EntityType[CompID, CompType, *CompID]
	EntityGroup   = ecs.EntityGroup[CompID, CompType, *CompID]
	EntityWorld   = ecs.EntityWorld[CompID, CompType, *CompID]
)

// 接口

type (
	IComponent    = ecs.IComponent[CompID, CompType, *CompID]
	IEntity       = ecs.IEntity[CompID, CompType, *CompID]
	IEntitySystem = ecs.IEntitySystem[CompID, CompType, *CompID]
)

// 函数

// 辅助函数/结构体

func NewEntityType(name string) *EntityType {
	return ecs.NewEntityType[CompID, CompType](name)
}

func NewWorldEntity[ET IEntity](w *EntityWorld, name string) *ecs.WorldEntity[ET, CompID, CompType, *CompID] {
	return ecs.NewWorldEntity[ET](w, name)
}

func RegComponentType[C IComponent](typ CompType) {
	ecs.RegComponentType[C, CompID](typ)
}

func EntityCreate[ET IEntity](typ *EntityType, id ecs.IdentityID) ET {
	return ecs.EntityCreate[ET](typ, id)
}

func GetComponent[C IComponent](e IEntity) (cc C, valid bool) {
	return ecs.GetEntityComponent[C, CompID, CompType](e)
}
