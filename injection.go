package ecs

import "reflect"

// import "fmt"

// var global = make(map[any]Binding)

// type Key struct {
// 	Value any
// }

// type baseKeySource[T any] struct{}

// func (baseKeySource[T]) Key() Key {
// 	return Key{Value: (*T)(nil)}
// }

// func (k Key) Generate() any {
// 	return k.Value
// }

// type KeySource interface {
// 	Key() Key
// }

// type Binding interface {
// 	Instance(init bool) (any, error)
// }

// type BindingSource[T any] interface {
// 	Key() Key
// 	Binding() (Binding, error)
// }

// type bindingSource[T any] struct {
// 	binding   Binding
// 	keySource baseKeySource[T]
// }

// // Binding implements BindingSource
// func (b *bindingSource[T]) Binding() (Binding, error) {
// 	instance, err := b.binding.Instance(false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if _, ok := instance.(T); !ok {
// 		var init T
// 		return nil, fmt.Errorf(`binding is not possible for "%v" and "%v"`, init, instance)
// 	}
// 	return b.binding, nil
// }

// // Key implements BindingSource
// func (b *bindingSource[T]) Key() Key {
// 	return b.keySource.Key()
// }

// func Bind[T any](source BindingSource[T]) error {
// 	key := source.Key()
// 	generate := key.Generate()

// 	binding, err := source.Binding()
// 	if err != nil {
// 		return err
// 	}
// 	global[generate] = binding
// 	return nil
// }

// func NewInstance[T any]() (T, error) {
// 	var empty T
// 	source := &baseKeySource[T]{}
// 	key := source.Key()

// 	generate := key.Generate()

// 	binding, ok := global[generate]
// 	if !ok {
// 		return empty, fmt.Errorf("not found")
// 	}

// 	instance, err := binding.Instance(false)
// 	if err != nil {
// 		return empty, err
// 	}
// 	result, ok := instance.(T)
// 	if !ok {
// 		return empty, fmt.Errorf("error type")
// 	}
// 	return result, nil
// }

var (
	allComponentType   = make(map[any]binding)
	allComponentToType = make(map[any]any)
)

type binding interface{ Instance() any }

type bindingBase[Comp IComponent[ID, I, C], ID any, I comparable, C CID[ID, I]] struct{ id I }

func bindingKey[Comp IComponent[ID, I, C], ID any, I comparable, C CID[ID, I]]() interface{} {
	return *new(Comp)
}

func (b bindingBase[Comp, ID, I, C]) Instance() any {
	var comp Comp // 默认是指针类型的
	c := reflect.New(reflect.TypeOf(comp).Elem()).Interface().(Comp)
	c.SetID(b.id)
	return c
}

func RegComponentType[Comp IComponent[ID, I, C], ID any, I comparable, C CID[ID, I]](id I) {
	allComponentType[id] = bindingBase[Comp, ID, I, C]{
		id: id,
	}
	allComponentToType[bindingKey[Comp, ID, I, C]()] = id
}

func GetComponentBinding[ID any, I comparable, C CID[ID, I]](id I) (binding, bool) {
	binding, ok := allComponentType[id]
	if !ok {
		return nil, false
	}
	return binding, ok
}

func NewComponent[ID any, I comparable, C CID[ID, I]](b binding) (IComponent[ID, I, C], bool) {
	v := b.Instance()
	if c, ok := v.(IComponent[ID, I, C]); ok {
		return c, true
	} else {
		return nil, false
	}
}
