package main

import (
	"flag"
	"os"

	"ragx/app/pkg/encoder"
	logging "ragx/app/pkg/logger"

	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/sirupsen/logrus"

	"ragx/app/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func main() {
	flag.Parse()
	logger := log.With(logging.NewLogger(logrus.StandardLogger()),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		//"service.id", id,
		//"service.name", Name,
		//"service.version", Version,
		"trace.id", tracing.TraceID(),
		//"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
			env.NewSource(),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	app, cleanup, err := wireApp(bc.Server, bc.Data, &bc, logger)
	if err != nil {
		panic(err)
	}

	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func init() {
	// 全局覆盖 protojson 序列化器，protojson默认会将int64序列化为字符串，这里统一序列化为数字
	encoding.RegisterCodec(encoder.JsonCodec{})
	// 本地debug时可以使用-conf ./configs/local
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf configs")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			//gs,
			hs,
		),
	)
}
