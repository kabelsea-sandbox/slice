package machinery

import (
	"context"
	"os"
	"time"

	"slice"
	"slice/pkg/di"
	"slice/pkg/run"

	machinery "github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/config"
	machinerylog "github.com/RichardKnop/machinery/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	redislock "github.com/RichardKnop/machinery/v2/locks/redis"
)

type Bundle struct {
	Broker          string        `envconfig:"MACHINERY_BROKER" default:"redis://redis:6379" required:"true"`
	DefaultQueue    string        `envconfig:"MACHINERY_DEQAULT_QUEUE" default:"machinery_tasks" required:"true"`
	ResultBackend   string        `envconfig:"MACHINERY_RESULT_BACKEND" default:"redis://redis:6379" required:"true"`
	ResultsExpireIn time.Duration `envconfig:"MACHINERY_RESULTS_EXPIRE_IN" default:"60s"`
	Lock            string        `envconfig:"MACHINERY_LOCK" default:"redis://redis:6379" required:"true"`
	Redis           struct {
		MaxIdle        int `envconfig:"MAX_IDLE" default:"3"`
		IdleTimeout    int `envconfig:"IDLE_TIMEOUT" default:"240"`
		ReadTimeout    int `envconfig:"READ_TIMEOUT" default:"15"`
		WriteTimeout   int `envconfig:"WRITE_TIMEOUT" default:"15"`
		ConnectTimeout int `envconfig:"CONNECT_TIMEOUT" default:"15"`
	} `envconfig:"MACHINERY_REDIS"`
	ConsumerEnabled bool `envconfig:"MACHINERY_CONSUMER_ENABLED"`
	Consumer        struct {
		Tag         string `envconfig:"TAG"`
		Concurrency int    `envconfig:"CONCURRENCY" default:"1"`
	} `envconfig:"MACHINERY_CONSUMER"`
}

// Build implements Bundle.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	var _errorHandler ErrorHandler
	if !builder.Has(&_errorHandler) {
		builder.Provide(errorHandler)
	}

	var _preTaskHandler PreTaskHandler
	if !builder.Has(&_preTaskHandler) {
		builder.Provide(preTaskHandler)
	}

	var _postTaskHandler PostTaskHandler
	if !builder.Has(&_postTaskHandler) {
		builder.Provide(postTaskHandler)
	}

	var _tasks []Task
	if !builder.Has(&_tasks) {
		builder.Provide(b.RegisterTasks)
	}

	builder.Provide(b.NewServer)
	builder.Provide(NewTaskPool)

	if b.ConsumerEnabled {
		builder.Provide(b.NewWorker, di.As(new(run.Worker)))
	}
}

// Boot implements Bundle interface.
func (b *Bundle) Boot(_ context.Context, container slice.Container) error {
	if b.Consumer.Tag == "" {
		b.Consumer.Tag = os.Getenv("HOSTNAME")
	}

	var _tasks []Task
	if container.Has(&_tasks) {
		if err := container.Invoke(b.RegisterTasks); err != nil {
			return errors.Wrap(err, "register tasks failed")
		}
	}
	return nil
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) error {
	return nil
}

func (b *Bundle) NewServer(logger slice.Logger, zaplogger *zap.Logger) (*machinery.Server, error) {
	defer logger.Infof("MACHINERY", "register server")

	cfg := &config.Config{
		DefaultQueue:    b.DefaultQueue,
		ResultsExpireIn: int(b.ResultsExpireIn.Seconds()),
		Redis: &config.RedisConfig{
			MaxIdle:        b.Redis.MaxIdle,
			IdleTimeout:    b.Redis.IdleTimeout,
			ReadTimeout:    b.Redis.ReadTimeout,
			WriteTimeout:   b.Redis.WriteTimeout,
			ConnectTimeout: b.Redis.ConnectTimeout,
		},
		NoUnixSignals: true,
	}

	// broker
	brokerOpt, err := redis.ParseURL(b.Broker)
	if err != nil {
		return nil, err
	}

	// result backend
	resultOpt, err := redis.ParseURL(b.ResultBackend)
	if err != nil {
		return nil, err
	}

	// lock
	lockOpt, err := redis.ParseURL(b.Lock)
	if err != nil {
		return nil, err
	}

	machinerylog.Set(NewLogAdapter(zaplogger))

	return machinery.NewServer(
		cfg,
		redisbroker.New(cfg, brokerOpt.Addr, brokerOpt.Password, "", brokerOpt.DB),
		redisbackend.New(cfg, resultOpt.Addr, resultOpt.Password, "", resultOpt.DB),
		redislock.New(cfg, []string{lockOpt.Addr}, lockOpt.DB, 0),
	), nil
}

func (b *Bundle) NewWorker(
	logger slice.Logger,
	server *machinery.Server,
	tasks TaskPool,
	errorHandler ErrorHandler,
	preHandler PreTaskHandler,
	postHandler PostTaskHandler,
) (Worker, error) {
	defer logger.Infof("MACHINERY", "register consumer worker")

	return NewWorker(
		server,
		tasks,
		b.Consumer.Tag,
		errorHandler,
		preHandler,
		postHandler,
	)
}

func (b *Bundle) RegisterTasks(logger slice.Logger, pool TaskPool, tasks []Task) error {
	p := pool.(*taskPool)

	for _, task := range tasks {
		name, fn := task()
		p.Add(name, fn)

		logger.Infof("MACHINERY", "register task - [%s]", name)
	}
	return nil
}
