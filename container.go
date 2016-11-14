package goxp

import (
	"fmt"
	"reflect"
)

type Container struct {
	// 根据类型map实际的值
	mappers map[reflect.Type]reflect.Value
}

func New() *Container {
	return &Container{make(map[reflect.Type]reflect.Value)}
}

func (inj *Container) GetType(value interface{}) reflect.Type {
	t := reflect.TypeOf(value)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		panic("不是接口类型")
	}

	return t
}

func (inj *Container) SetMap(value interface{}) {
	inj.mappers[reflect.TypeOf(value)] = reflect.ValueOf(value)
}

func (inj *Container) Get(t reflect.Type) reflect.Value {
	val := inj.mappers[t]
	if val.IsValid() {
		return val
	}


	if t.Kind() == reflect.Interface {
		for k, v := range inj.mappers {
			if k.Implements(t) {
				val = v
				break
			}
		}
	}

	return val
}

func (inj *Container) Invoke(i interface{}) ([]reflect.Value, error) {
	t := reflect.TypeOf(i)

	if t.Kind() != reflect.Func {
		panic("必须是一个函数!")
	}
	inValues := make([]reflect.Value, t.NumIn())
	for k := 0; k < t.NumIn(); k++ {
		val := inj.Get(t.In(k))
		if !val.IsValid() {
			return nil, fmt.Errorf("Value没有这个类型 %v", t.In(k))
		}
		inValues[k] = val
	}
	ret := reflect.ValueOf(i).Call(inValues)
	return ret, nil
}
