package di

import (
	"fmt"
	"reflect"

	"github.com/kabelsea-sandbox/slice/pkg/di/internal/reflection"
)

// ctorType describes types of constructor provider.
type ctorType int

const (
	ctorUnknown      ctorType = iota
	ctorStd                   // (deps) (result)
	ctorError                 // (deps) (result, error)
	ctorCleanup               // (deps) (result, cleanup)
	ctorCleanupError          // (deps) (result, cleanup, error)
)

// providerConstructor is a provider that can handle constructor functions.
// Type of this provider provides type with function call.
type providerConstructor struct {
	name string
	// constructor
	ctorType ctorType
	call     reflection.Func
	// injectable params
	injectable struct {
		// params parsed once
		params []parameter
		// field numbers parsed once
		fields []int
	}
}

// newProviderConstructor creates new constructor provider with name as additional identity key.
func newProviderConstructor(name string, fn reflection.Func) (*providerConstructor, error) {
	ctorType := determineCtorType(fn)
	if ctorType == ctorUnknown {
		return nil, fmt.Errorf("invalid constructor signature, got %s", fn.Type)
	}
	provider := &providerConstructor{
		name:     name,
		call:     fn,
		ctorType: ctorType,
	}
	// result type
	rt := fn.Out(0)
	// constructor result with di.Inject - only addressable pointers
	// anonymous parameters with di.Inject - only struct
	if isInjectable(rt) && rt.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("di.Inject not supported for unaddressable result of constructor, use *%s instead", rt)
	}
	// if struct is injectable, range over injectableFields and parse injectable params
	if isInjectable(rt) {
		provider.injectable.params, provider.injectable.fields = parseInjectableType(rt)
	}
	return provider, nil
}

func (c providerConstructor) Type() reflect.Type {
	return c.call.Out(0)
}

func (c providerConstructor) Name() string {
	return c.name
}

// ParameterList returns constructor parameter list.
func (c *providerConstructor) ParameterList() parameterList {
	// todo: move to constructor
	var pl parameterList
	for i := 0; i < c.call.NumIn(); i++ {
		in := c.call.In(i)
		pl = append(pl, parameter{
			name: "", // constructor parameters could be resolved only with empty ("") name
			typ:  in,
		})
	}
	return append(pl, c.injectable.params...)
}

// Provide provides resultant.
func (c *providerConstructor) Provide(values ...reflect.Value) (reflect.Value, func(), error) {
	// constructor last param index
	clpi := c.call.NumIn()
	if c.call.NumIn() == 0 {
		clpi = 0
	}
	out := reflection.CallResult(c.call.Call(values[:clpi]))
	rv := out.Result()
	if c.ctorType == ctorError && out.Error(1) != nil {
		return rv, nil, out.Error(1)
	}
	if c.ctorType == ctorCleanupError && out.Error(2) != nil {
		return rv, nil, out.Error(2)
	}
	// set injectable fields
	if len(c.injectable.fields) > 0 {
		// result value
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		fields := values[clpi:]
		// field index
		for i, value := range fields {
			// field value
			fv := rv.Field(c.injectable.fields[i])
			fv.Set(value)
		}
	}
	switch c.ctorType {
	case ctorStd:
		return out.Result(), nil, nil
	case ctorError:
		return out.Result(), nil, out.Error(1)
	case ctorCleanup:
		return out.Result(), out.Cleanup(), nil
	case ctorCleanupError:
		return out.Result(), out.Cleanup(), out.Error(2)
	}
	bug()
	return reflect.Value{}, nil, nil
}

// determineCtorType
func determineCtorType(fn reflection.Func) ctorType {
	if fn.NumOut() == 1 {
		return ctorStd
	}
	if fn.NumOut() == 2 {
		if reflection.IsError(fn.Out(1)) {
			return ctorError
		}
		if reflection.IsCleanup(fn.Out(1)) {
			return ctorCleanup
		}
	}
	if fn.NumOut() == 3 && reflection.IsCleanup(fn.Out(1)) && reflection.IsError(fn.Out(2)) {
		return ctorCleanupError
	}
	return ctorUnknown
}
