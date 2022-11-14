package machinery

import (
	"context"
	"fmt"
	"reflect"

	machinery "github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/tasks"
	"go.uber.org/zap"
)

type Task func() (string, interface{})

type TaskSignature = tasks.Signature

type Chains map[string][]any

type TaskPool interface {
	// Signature
	Signature(name string, args ...any) *TaskSignature

	// Send
	Send(ctx context.Context, name string, args ...any) error

	// Chain
	Chain(ctx context.Context, chains Chains) error
}

type taskPool struct {
	logger *zap.Logger
	server *machinery.Server
	tasks  map[string]interface{}
}

func NewTaskPool(logger *zap.Logger, server *machinery.Server) TaskPool {
	return &taskPool{
		logger: logger.With(zap.Namespace("automate")),
		server: server,
		tasks:  make(map[string]interface{}),
	}
}

func (t *taskPool) Add(name string, task interface{}) error {
	if _, ok := t.tasks[name]; ok {
		return fmt.Errorf("tasks with signature already exists on the pool, %s", name)
	}

	t.tasks[name] = task
	return nil
}

// Signature implements TaskPool interface
func (t *taskPool) Signature(name string, args ...any) *TaskSignature {
	signature := &tasks.Signature{
		Name: name,
		Args: make([]tasks.Arg, 0, len(args)),
	}

	for _, arg := range args {
		signature.Args = append(signature.Args, tasks.Arg{
			Type:  reflect.TypeOf(arg).String(),
			Value: arg,
		})
	}
	return signature
}

// Send implements TaskPool interface
func (t *taskPool) Send(ctx context.Context, name string, args ...any) error {
	var err error

	defer t.logger.Debug("send",
		zap.String("name", name),
		zap.Any("args", args),
		zap.Error(err),
	)

	_, err = t.server.SendTaskWithContext(ctx, t.Signature(name, args...))
	return err
}

// Chain implements TaskPool interface
func (t *taskPool) Chain(ctx context.Context, chains Chains) error {
	var err error

	defer t.logger.Debug("chain",
		zap.Any("chains", chains),
		zap.Error(err),
	)

	signatures := []*TaskSignature{}

	for name, params := range chains {
		signatures = append(signatures, t.Signature(name, params...))
	}

	chain, err := tasks.NewChain(signatures...)
	if err != nil {
		return err
	}

	_, err = t.server.SendChainWithContext(ctx, chain)
	return err
}
