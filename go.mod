module github.com/kabelsea-sandbox/slice

go 1.19

replace (
	cloud.google.com/go/pubsub => cloud.google.com/go/pubsub v1.9.0 // bug on machinery
	github.com/pegasus-kv/thrift v0.14.1 => github.com/pegasus-kv/thrift v0.13.0 // bug on gocache
)

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.2
	firebase.google.com/go/v4 v4.9.0
	github.com/99designs/gqlgen v0.17.20
	github.com/RichardKnop/logging v0.0.0-20190827224416-1a693bdd4fae
	github.com/RichardKnop/machinery/v2 v2.0.11
	github.com/aerospike/aerospike-client-go v4.5.2+incompatible
	github.com/eko/gocache/v3 v3.1.1
	github.com/getsentry/raven-go v0.2.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-pg/pg/v9 v9.2.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/imdario/mergo v0.3.13
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/nats-io/nats.go v1.18.0
	github.com/oklog/run v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.8.2
	github.com/sercand/kuberesolver v2.4.0+incompatible
	github.com/stretchr/testify v1.8.0
	github.com/vmihailenco/msgpack/v5 v5.3.5
	go.mongodb.org/mongo-driver v1.10.3
	go.opencensus.io v0.23.0
	go.uber.org/automaxprocs v1.5.1
	go.uber.org/zap v1.23.0
	golang.org/x/time v0.1.0
	google.golang.org/grpc v1.50.1
)

require (
	cloud.google.com/go v0.104.0 // indirect
	cloud.google.com/go/compute v1.10.0 // indirect
	cloud.google.com/go/firestore v1.8.0 // indirect
	cloud.google.com/go/iam v0.5.0 // indirect
	cloud.google.com/go/kms v1.5.0 // indirect
	cloud.google.com/go/pubsub v1.10.0 // indirect
	cloud.google.com/go/storage v1.27.0 // indirect
	github.com/XiaoMi/pegasus-go-client v0.0.0-20210427083443-f3b6b08bc4c2 // indirect
	github.com/aws/aws-sdk-go v1.37.16 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20220106215444-fb4bf637b56d // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-redsync/redsync/v4 v4.0.4 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/montanaflynn/stats v0.6.6 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pegasus-kv/thrift v0.13.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/exp v0.0.0-20220518171630-0b5c67f07fdf // indirect
	golang.org/x/oauth2 v0.1.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine/v2 v2.0.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
	k8s.io/apimachinery v0.23.5 // indirect
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/codemodus/kace v0.5.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-pg/zerochecker v0.2.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/websocket v1.5.0
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jackc/pgx/v4 v4.17.2
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/nats-io/nats-server/v2 v2.3.4 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.13.0
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/prometheus/statsd_exporter v0.22.8 // indirect
	github.com/segmentio/encoding v0.3.6 // indirect
	github.com/uptrace/bun v1.1.8
	github.com/uptrace/bun/dialect/pgdialect v1.1.8
	github.com/uptrace/bun/extra/bundebug v1.1.8
	github.com/vektah/gqlparser/v2 v2.5.1 // indirect
	github.com/vmihailenco/bufpool v0.1.11 // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/api v0.100.0
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221024183307-1bc688fe9f3e // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mellium.im/sasl v0.3.0 // indirect
)
