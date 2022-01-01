package slice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kabelsea-sanbox/slice"
)

type FirstBundle struct {
	TestBundleInvocationOrder
}
type SecondBundle struct {
	TestBundleInvocationOrder
}
type ThirdBundle struct {
	TestBundleInvocationOrder
}

// TestBundleInvocationOrder
type TestBundleInvocationOrder struct {
	name  string
	order *[]string
}

func (e *TestBundleInvocationOrder) Build(builder slice.ContainerBuilder) {
	*e.order = append(*e.order, fmt.Sprintf("%s_build", e.name))
}

func (e *TestBundleInvocationOrder) Boot(ctx context.Context, invoker slice.Container) (err error) {
	*e.order = append(*e.order, fmt.Sprintf("%s_boot", e.name))
	return nil
}

func (e *TestBundleInvocationOrder) Shutdown(ctx context.Context, invoker slice.Container) (err error) {
	*e.order = append(*e.order, fmt.Sprintf("%s_shutdown", e.name))
	return nil
}

func TestRun(t *testing.T) {
	t.Run("slice runs invoke function", func(t *testing.T) {
		var invokeCalled bool
		slice.Run(
			slice.SetName("test"),
			slice.SetDispatcher(func() {
				invokeCalled = true
			}),
		)
		require.True(t, invokeCalled)
	})

	t.Run("bundle processes in correct order: build, boot, shutdown", func(t *testing.T) {
		var order []string
		first := &FirstBundle{TestBundleInvocationOrder{name: "first", order: &order}}
		second := &SecondBundle{TestBundleInvocationOrder{name: "second", order: &order}}
		third := &ThirdBundle{TestBundleInvocationOrder{name: "third", order: &order}}
		slice.Run(
			slice.SetName("test"),
			slice.SetDispatcher(func() {}),
			slice.RegisterBundles(
				first,
				second,
				third,
			),
		)
		require.Equal(t, []string{
			"first_build",
			"second_build",
			"third_build",
			"first_boot",
			"second_boot",
			"third_boot",
			"third_shutdown", // shutdown is inverted
			"second_shutdown",
			"first_shutdown",
		}, order)
	})
}
