package machinery

import (
	"go.uber.org/zap"

	"github.com/RichardKnop/machinery/v2/tasks"
)

// Machinery handler interface
type PostTaskHandler func(signature *tasks.Signature)

func postTaskHandler(logger *zap.Logger) PostTaskHandler {
	logger = logger.With(zap.Namespace("machinery"))

	return func(signature *tasks.Signature) {
		logger.Debug("post_task",
			zap.Any("signature", signature),
		)
	}
}

// Machinery handler interface
type ErrorHandler func(error)

func errorHandler(logger *zap.Logger) ErrorHandler {
	logger = logger.With(zap.Namespace("machinery"))

	return func(err error) {
		logger.Debug("error",
			zap.Error(err),
		)
	}
}

// Machinery handler interface
type PreTaskHandler func(signature *tasks.Signature)

func preTaskHandler(logger *zap.Logger) PreTaskHandler {
	logger = logger.With(zap.Namespace("machinery"))

	return func(signature *tasks.Signature) {
		logger.Debug("pre_task",
			zap.Any("signature", signature),
		)
	}
}
