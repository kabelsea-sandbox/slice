package zaplog

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const StampMilli = "2006/01/02 15:04:05"

const (
	MessageKey = "message"
	LevelKey   = "level"
	TimeKey    = "time"
)

// RFC3339TimeEncoder serializes a time.Time to an RFC3339-formatted string.
func StampMilliEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(StampMilli))
}

// SecondsDurationEncoder serializes a time.Duration to a floating-point number of seconds elapsed.
func MilliDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

// NewProductionEncoderConfig returns new production zap encoder.
func NewProductionEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.RFC3339TimeEncoder
	config.MessageKey = MessageKey
	config.LevelKey = LevelKey
	config.TimeKey = TimeKey
	config.EncodeDuration = MilliDurationEncoder
	return zapcore.NewJSONEncoder(config)
}

// NewDevelopmentEncoder returns new development zap encoder.
func NewDevelopmentEncoder() zapcore.Encoder {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = StampMilliEncoder
	config.MessageKey = MessageKey
	config.LevelKey = LevelKey
	config.TimeKey = TimeKey
	config.EncodeDuration = MilliDurationEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(config)
}

// NewProductionConfig
func NewProductionConfig(level zapcore.Level) zap.Config {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.MessageKey = MessageKey
	config.EncoderConfig.LevelKey = LevelKey
	config.EncoderConfig.TimeKey = TimeKey
	config.EncoderConfig.EncodeDuration = MilliDurationEncoder
	return config
}

// NewDevelopmentConfig
func NewDevelopmentConfig(level zapcore.Level) zap.Config {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = StampMilliEncoder
	config.EncoderConfig.MessageKey = MessageKey
	config.EncoderConfig.LevelKey = LevelKey
	config.EncoderConfig.TimeKey = TimeKey
	config.EncoderConfig.EncodeDuration = MilliDurationEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config
}
