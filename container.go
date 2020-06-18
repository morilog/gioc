package gioc

import (
	"errors"
	"reflect"
	"sync"
)

var container map[reflect.Type]binding
var mu *sync.Mutex

func init() {
	container = make(map[reflect.Type]binding)
	mu = &sync.Mutex{}
}

type binding struct {
	singleton bool
	resolved  bool
	resolver  interface{}
	instance  interface{}
}

// Singleton binds a resolver to an abstraction once
func Singleton(resolver interface{}) {
	Bind(resolver, true)
}

// Bind binds a resolver to an abstraction
func Bind(resolver interface{}, singleton bool) {
	mu.Lock()
	defer mu.Unlock()
	resolverType := reflect.TypeOf(resolver)
	if resolverType.Kind() != reflect.Func {
		panic("ioc resolver should be a function")
	}

	if resolverType.NumOut() != 2 {
		panic("ioc resolver should have two return types")
	}

	if resolverType.Out(1) != reflect.ValueOf(new(error)).Elem().Type() {
		panic("ioc resolver second return type should be error")
	}

	container[reflect.TypeOf(resolver).Out(0)] = binding{
		singleton: singleton,
		resolved:  false,
		resolver:  resolver,
	}
}

// Make loads given parameter's value from container
func Make(receiver interface{}) error {
	receiverType := reflect.TypeOf(receiver)

	return resolveType(receiverType, receiver)
}

func getBindingFor(t reflect.Type) (binding, bool) {
	b, ok := container[t]
	if !ok {
		if reflect.Ptr == t.Kind() {
			b, ok = container[t.Elem()]
		}
	}

	return b, ok
}

func (b binding) resolve() (interface{}, error) {
	if b.singleton && b.resolved {
		return b.instance, nil
	}
	args, err := getArgumentsOf(b.resolver)
	if err != nil {
		return nil, err
	}

	values := reflect.ValueOf(b.resolver).Call(args)

	if values[1].Interface() != nil {
		return nil, values[1].Interface().(error)
	}

	return values[0].Interface(), nil
}

func getArgumentsOf(function interface{}) ([]reflect.Value, error) {
	fType := reflect.TypeOf(function)
	argsCount := fType.NumIn()
	args := make([]reflect.Value, argsCount)

	for i := 0; i < argsCount; i++ {
		abstract := fType.In(i)
		var instance interface{}
		if err := resolveType(abstract, &instance); err != nil {
			return nil, err
		}
		args[i] = reflect.ValueOf(instance)
	}

	return args, nil
}

func resolveType(t reflect.Type, receiver interface{}) error {
	b, ok := getBindingFor(t)
	if !ok {
		return errors.New("binding not found")
	}

	concrete, err := b.resolve()
	if err != nil {
		return err
	}

	reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(concrete))
	return nil
}
