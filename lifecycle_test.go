package slice

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/kabelsea-sandbox/slice/pkg/di"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestLifecycle_initialization(t *testing.T) {
	t.Run("provide user dependency", func(t *testing.T) {
		c, err := initialization(
			di.Provide(http.NewServeMux),
		)
		require.NoError(t, err)
		var mux *http.ServeMux
		require.True(t, c.Has(&mux))
	})

	t.Run("incorrect option cause error", func(t *testing.T) {
		c, err := initialization(
			di.Provide(func() {}),
		)
		require.Nil(t, c)
		require.Error(t, err)
		require.Contains(t, err.Error(), "lifecycle_test.go:")
		require.Contains(t, err.Error(), ": invalid constructor signature, got func()")
	})
}

func TestLifecycle_configureBundles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("process iterates over all bundles", func(t *testing.T) {
		bundles := []bundle{
			{
				name:   "first-bundle",
				Bundle: NewMockBundle(ctrl),
			},
			{
				name:   "second-bundle",
				Bundle: NewMockBundle(ctrl),
			},
			{
				name:   "third-bundle",
				Bundle: NewMockBundle(ctrl),
			},
		}
		i := 0
		err := configureBundles(func(prefix string, spec interface{}) error {
			require.Equal(t, bundles[i].Bundle, spec)
			i++
			return nil
		}, "", bundles...)
		require.NoError(t, err)
		require.Equal(t, 3, i)
	})

	t.Run("process error causes configure error", func(t *testing.T) {
		bundle := bundle{
			name:   "error-bundle",
			Bundle: NewMockBundle(ctrl),
		}
		err := configureBundles(func(prefix string, spec interface{}) error {
			return errors.New("unexpected error")
		}, "", bundle)
		require.EqualError(t, err, "error-bundle bundle configure failed: unexpected error")
	})
}

func TestLifecycle_buildBundles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("bundle builds in correct order", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		var order []string

		firstBundle := NewMockBundle(ctrl)
		firstBundle.
			EXPECT().
			Build(gomock.All()).
			Do(func(_ ContainerBuilder) {
				order = append(order, "first")
			}).
			Return().
			Times(1)

		secondBundle := NewMockBundle(ctrl)
		secondBundle.
			EXPECT().
			Build(gomock.All()).
			Do(func(_ ContainerBuilder) {
				order = append(order, "second")
			}).
			Return().
			Times(1)

		bundles := []bundle{
			{
				name:   "first-bundle",
				Bundle: firstBundle,
			},
			{
				name:   "second-bundle",
				Bundle: secondBundle,
			},
		}
		err = buildBundles(c, bundles...)
		require.NoError(t, err)

		require.Equal(t, []string{"first", "second"}, order)
	})

	t.Run("bundle build error return as one", func(t *testing.T) {

		mockBundle := NewMockBundle(ctrl)
		mockBundle.
			EXPECT().
			Build(gomock.All()).
			Do(func(builder ContainerBuilder) {
				builder.Provide(func() {})
				builder.Provide(nil)
				builder.Provide(struct{}{})
			}).
			Return().
			Times(1)

		c, err := di.New()
		require.NoError(t, err)
		errorBundle := bundle{
			name:   "error-bundle",
			Bundle: mockBundle,
		}
		err = buildBundles(c, errorBundle)
		require.Error(t, err)
		require.Contains(t, err.Error(), "error-bundle: build failed:")
		require.Contains(t, err.Error(), "invalid constructor signature, got func()")
		require.Contains(t, err.Error(), "invalid constructor signature, got nil")
		require.Contains(t, err.Error(), "invalid constructor signature, got struct {}")
	})
}

func TestLifecycle_boot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("iterates over bundles and run boot function", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		var order []string

		firstBundle := NewMockBootShutdown(ctrl)
		firstBundle.
			EXPECT().
			Boot(gomock.All(), gomock.All()).
			DoAndReturn(func(ctx context.Context, container Container) error {
				order = append(order, "first-bundle")
				return nil
			}).
			Times(1)

		secondBundle := NewMockBootShutdown(ctrl)
		secondBundle.
			EXPECT().
			Boot(gomock.All(), gomock.All()).
			DoAndReturn(func(ctx context.Context, container Container) error {
				order = append(order, "second-bundle")
				return nil
			}).
			Times(1)

		bundles := []bundle{
			{
				name:   "first-bundle",
				Bundle: firstBundle,
			},
			{
				name:   "second-bundle",
				Bundle: secondBundle,
			},
		}
		shutdowns, err := boot(time.Second, c, bundles...)
		require.NoError(t, err)
		require.Len(t, shutdowns, 2)
		require.Equal(t, []string{"first-bundle", "second-bundle"}, order)
	})

	t.Run("bundle boot error causes boot error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)

		mockBundle := NewMockBootShutdown(ctrl)
		mockBundle.
			EXPECT().
			Boot(gomock.All(), gomock.All()).
			DoAndReturn(func(ctx context.Context, container Container) error {
				return errors.New("unexpected errro")
			}).
			Times(1)

		bundle := bundle{
			name:   "error-bundle",
			Bundle: mockBundle,
		}

		shutdowns, err := boot(time.Millisecond, c, bundle)
		require.EqualError(t, err, "error-bundle bundle boot failed: unexpected errro")
		require.Len(t, shutdowns, 0)
	})

	t.Run("shutdowns correct on context deadline", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)

		firstBundle := NewMockBootShutdown(ctrl)
		firstBundle.
			EXPECT().
			Boot(gomock.All(), gomock.All()).
			DoAndReturn(func(ctx context.Context, container Container) error {
				time.Sleep(2 * time.Millisecond)
				return nil
			}).
			Times(1)

		secondBundle := NewMockBootShutdown(ctrl)
		secondBundle.
			EXPECT().
			Boot(gomock.All(), gomock.All()).
			DoAndReturn(func(ctx context.Context, container Container) error {
				return nil
			}).
			Times(0)

		bundles := []bundle{
			{
				name:   "first-bundle",
				Bundle: firstBundle,
			},
			{
				name:   "second-bundle",
				Bundle: secondBundle,
			},
		}
		shutdowns, err := boot(time.Millisecond, c, bundles...)
		require.EqualError(t, err, "boot failed: context deadline exceeded")
		require.Len(t, shutdowns, 1)
	})
}

func TestLifecycle_drun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("resolve dispatcher and run", func(t *testing.T) {
		dispatcher := NewMockDispatcher(ctrl)
		dispatcher.
			EXPECT().
			Run(gomock.All()).
			DoAndReturn(func(ctx context.Context) error {
				return nil
			}).
			Times(1)

		c, _ := di.New(di.Provide(func() *MockDispatcher { return dispatcher }, di.As(new(Dispatcher))))

		require.NotNil(t, c)
		err := drun(c)
		require.NoError(t, err)
	})

	t.Run("undefined dispatcher cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		err = drun(c)
		require.Error(t, err)
		require.Contains(t, err.Error(), "resolve dispatcher failed: ")
		require.Contains(t, err.Error(), "lifecycle.go:")
		require.Contains(t, err.Error(), ": type slice.Dispatcher not exists in container")
	})

	t.Run("run error causes error", func(t *testing.T) {
		dispatcher := NewMockDispatcher(ctrl)
		dispatcher.
			EXPECT().
			Run(gomock.All()).
			DoAndReturn(func(ctx context.Context) error {
				return errors.New("unexpected error")
			}).
			Times(1)

		c, _ := di.New(di.Provide(func() *MockDispatcher { return dispatcher }, di.As(new(Dispatcher))))
		require.NotNil(t, c)
		err := drun(c)
		require.EqualError(t, err, "unexpected error")
	})
}

func TestLifecycle_reverseShutdown(t *testing.T) {
	t.Run("reverse order", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		var order []string
		shutdowns := shutdowns{
			{
				name: "first-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					order = append(order, "first-shutdown")
					return nil
				},
			},
			{
				name: "second-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					order = append(order, "second-shutdown")
					return nil
				},
			},
			{
				name: "third-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					order = append(order, "third-shutdown")
					return nil
				},
			},
		}
		err = reverseShutdown(time.Second, c, shutdowns)
		require.NoError(t, err)
		require.Equal(t, []string{"third-shutdown", "second-shutdown", "first-shutdown"}, order)
	})

	t.Run("shutdown errors returns one", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		shutdowns := shutdowns{
			{
				name: "first-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					return errors.New("first-error")
				},
			},
			{
				name: "second-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					return errors.New("second-error")
				},
			},
			{
				name: "third-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					return errors.New("third-error")
				},
			},
		}
		err = reverseShutdown(time.Second, c, shutdowns)
		require.EqualError(
			t, err,
			"shutdown failed: third-shutdown: third-error; second-shutdown: second-error; first-shutdown: first-error",
		)
	})
}
