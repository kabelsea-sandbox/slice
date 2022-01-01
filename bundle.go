package slice

import (
	"context"
	"reflect"
)

//go:generate mockgen -package=slice -destination=./bundle_mock_test.go -source=./bundle.go Bundle BootShutdown DependOn

// A Bundle is a main configuration component. The main purpose of a bundle is a
// registering application dependency on a build stage. The main reason for
// creating a bundle is an existing part of an application that you may reuse in other services.
// A bundle is a configuration first. It will be processed via envconfig (https://github.com/kelseyhightower/envconfig).
// Therefore, the bundle implementation must be a struct pointer.
//
// type LoggerBundle  struct {
//   Level string `envconfig:"level"` // will be loaded from env variable LEVEL
// }
//
// func(b *LoggerBundle) Build(builder di.ContainerBuilder) {
//   builder.Provide(b.provideLogger)
// }
//
// func(b *LoggerBundle) provideLogger() *PreferredLogger {
//   return &PreferredLogger{
//     Level: b.Level,
//   }
// }
//
// Also, bundles may provide some boot and shutdown functions before slice invoke a
// dispatch function and after this, respectively. See BootShutdown and Shutdowner interfaces.
type Bundle interface {
	// Build adds bundle types into container via di.ContainerBuilder.
	Build(builder ContainerBuilder)
}

// A BootShutdown is a bundle that provide some initialization functions.
type BootShutdown interface {
	Bundle
	// Boot provides way to interact with dependency injection container on the
	// boot stage. On boot stage main dependencies already provided to container.
	// And on this stage bundle can interact with them.
	// Boot can return error if process failed. It will be handled by Slice.
	Boot(ctx context.Context, container Container) (err error)
	// Shutdown provides way to interact with dependency injection container
	// on shutdown stage. It can compensate things that was be made on boot stage.
	// Shutdown can return error if process failed. It will be handled by Slice.
	Shutdown(ctx context.Context, container Container) (err error)
}

// A DependOn describe that bundle depends on another bundle.
type DependOn interface {
	Bundle
	// DependOn returns dependent bundle.
	DependOn() []Bundle
}

const (
	temporary = 1
	permanent = 2
)

func sortBundles(bundles []Bundle) ([]Bundle, bool) {
	var sorted []Bundle
	marks := map[reflect.Type]int{}
	for _, b := range bundles {
		if !visit(b, marks, &sorted) {
			return sorted, false
		}
	}
	return sorted, true
}

func visit(b Bundle, marks map[reflect.Type]int, sorted *[]Bundle) bool {
	typ := reflect.TypeOf(b)
	if marks[typ] == permanent {
		return true
	}
	if marks[typ] == temporary {
		// acyclic
		return false
	}
	dependOn, ok := b.(DependOn)
	if !ok {
		marks[typ] = permanent
		*sorted = append(*sorted, b)
		return true
	}
	marks[typ] = temporary
	deps := dependOn.DependOn()
	for _, dep := range deps {
		if !visit(dep, marks, sorted) {
			return false
		}
	}
	marks[typ] = permanent
	*sorted = append(*sorted, b)
	return true
}

// prepareBundles prepares bundles.
func prepareBundles(bundles ...Bundle) (result []bundle) {
	for _, b := range bundles {
		result = append(result, prepareBundle(b))
	}
	return result
}

// prepareBundle prepares bundle.
func prepareBundle(b Bundle) bundle {
	return bundle{
		name:   getBundleName(b),
		Bundle: b,
	}
}

type bundle struct {
	Bundle
	name string
}

// getBundleName gets bundle string representation.
func getBundleName(bundle Bundle) string {
	return reflect.TypeOf(bundle).String()
}
