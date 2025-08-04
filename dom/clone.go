package dom

import (
	"fmt"
	"reflect"
)

func cloneProps(props Props) Props {
	if props == nil {
		return nil
	}
	if cloneable, ok := props.(IPropsCloneable); ok {
		return cloneable.Clone()
	}
	rv := reflect.ValueOf(props)
	if rv.Kind() == reflect.Struct {
		n := reflect.New(rv.Type())
		n.Elem().Set(rv)
		return n.Elem().Interface().(Props)
	}
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		el := rv.Elem()
		if el.Kind() == reflect.Struct {
			n := reflect.New(el.Type())
			n.Elem().Set(el)
			return n.Elem().Interface().(Props)
		}
	}
	panic(fmt.Errorf("unable to clone props: %T", props))
}

type IPropsCloneable interface {
	Clone() Props
}
