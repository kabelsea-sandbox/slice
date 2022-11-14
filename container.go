package slice

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kabelsea-sandbox/slice/pkg/di"
)

//go:generate mockgen -package=slice -destination=./container_mock_test.go -source=./container.go \
// ContainerBuilder Container

// ContainerBuilder builds the container. It is used to providing components in bundles.
type ContainerBuilder interface {
	// Has checks that type exists in container, if not it is return false.
	Has(target interface{}, options ...di.ResolveOption) bool
	// Provide provides a reliable way of component building to the container.
	// The constructor will be invoked lazily on-demand. For more information about
	// constructors see di.Constructor interface. ProvideOption can add additional
	// behavior to the process of type resolving.
	Provide(constructor di.Constructor, options ...di.ProvideOption)
}

// Container is a compiled dependency injection container interface.
type Container interface {
	// Has checks that type exists in container, if not it return false.
	Has(target interface{}, options ...di.ResolveOption) bool
	// Resolve builds instance of target type and fills target pointer.
	Resolve(into interface{}, options ...di.ResolveOption) error
	// Invoke calls provided function.
	Invoke(fn di.Invocation, options ...di.InvokeOption) error
}

// newBundleContainerBuilder creates container builder for bundle.
func newBundleContainerBuilder(container *di.Container) *bundleContainerBuilder {
	return &bundleContainerBuilder{
		container: container,
	}
}

// bundleContainerBuilder
type bundleContainerBuilder struct {
	container *di.Container
	buildErrs []error
}

// Has implements ContainerBuilder.
func (b *bundleContainerBuilder) Has(target interface{}, options ...di.ResolveOption) bool {
	return b.container.Has(target, options...)
}

// Provide implements ContainerBuilder.
func (b *bundleContainerBuilder) Provide(constructor di.Constructor, options ...di.ProvideOption) {
	if err := b.container.Provide(constructor, options...); err != nil {
		b.buildErrs = append(b.buildErrs, err) // append bundle provide error
	}
}

// Error return bundle build error. If bundle build success returns nil.
func (b *bundleContainerBuilder) Error() error {
	if len(b.buildErrs) == 0 {
		return nil
	}
	sb := strings.Builder{}
	for _, err := range b.buildErrs {
		sb.WriteString(fmt.Sprintf("\n\t- %s", err))
	}
	return errors.New(sb.String())
}
