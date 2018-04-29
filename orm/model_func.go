package orm

import (
	"fmt"
	"reflect"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

type funcModel struct {
	Model
	fnv  reflect.Value
	fnIn []reflect.Value
}

var _ Model = (*funcModel)(nil)

func newFuncModel(fn interface{}) *funcModel {
	m := &funcModel{
		fnv: reflect.ValueOf(fn),
	}

	fnt := m.fnv.Type()
	if fnt.Kind() != reflect.Func {
		panic(fmt.Sprintf("ForEach expects a %s, got a %s", reflect.Func, fnt.Kind()))
	}

	switch fnt.NumIn() {
	case 1:
	default:
		panic(fmt.Sprintf("ForEach expects 1 arg, got %d", fnt.NumIn()))
	}

	if fnt.NumOut() != 1 {
		panic(fmt.Sprintf("ForEach must return 1 value, got %d", fnt.NumOut()))
	}
	if fnt.Out(0) != errorType {
		panic(fmt.Sprintf("ForEach must return an error, got %T", fnt.Out(0)))
	}

	t0 := fnt.In(0)
	if t0.Kind() == reflect.Ptr {
		t0 = t0.Elem()
	}
	v0 := reflect.New(t0)
	m.Model = newStructTableModelValue(v0.Elem())
	m.fnIn = []reflect.Value{v0}

	return m
}

func (m *funcModel) AddModel(_ ColumnScanner) error {
	out := m.fnv.Call(m.fnIn)
	errv := out[0]
	if !errv.IsNil() {
		return errv.Interface().(error)
	}
	return nil
}
