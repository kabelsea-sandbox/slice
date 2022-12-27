package slicetest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kabelsea-sandbox/slice/pkg/di"

	"github.com/kabelsea-sandbox/slice"
)

// NewTestContainer
func NewTestContainer(t *testing.T) *TestContainer {
	container, _ := di.New(
		di.Provide(NewLogger, di.As(new(slice.Logger))),
	)
	return &TestContainer{
		Container: container,
		t:         t,
	}
}

// TestContainer
type TestContainer struct {
	*di.Container

	t    *testing.T
	errs []error
}

// Provide
func (t *TestContainer) Provide(constructor di.Constructor, options ...di.ProvideOption) {
	if err := t.Container.Provide(constructor, options...); err != nil {
		t.errs = append(t.errs, err)
	}
}

// ShouldNoError
func (t *TestContainer) ShouldNoError() {
	require.Len(t.t, t.errs, 0)
}

// Errors
func (t *TestContainer) Errors() []error {
	return t.errs
}
