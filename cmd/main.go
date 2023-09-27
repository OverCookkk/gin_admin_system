package main

import (
	"context"
	"flag"
	"fmt"
	"gin_admin_system/internal/app"
	"gin_admin_system/internal/app/config"
	"gin_admin_system/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type options struct {
	ConfigFile string
	ModelFile  string
	MenuFile   string
}

type Option func(*options)

func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

func SetModelFile(s string) Option {
	return func(o *options) {
		o.ModelFile = s
	}
}

func SetMenuFile(s string) Option {
	return func(o *options) {
		o.MenuFile = s
	}
}

// @title gin-admin
// @version 8.1.0
// @description RBAC scaffolding based on GIN + GORM + CASBIN + WIRE.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /
// @contact.name LyricTian
// @contact.email tiannianshou@gmail.com
func main() {
	ctx := logger.NewTagContext(context.Background(), "__main__")

	// flag
	configFile := flag.String("config", "", "config file")
	modelFile := flag.String("model", "", "model file")
	menuFile := flag.String("menu", "", "menu file")
	flag.Parse()

	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, SetConfigFile(*configFile), SetModelFile(*modelFile), SetMenuFile(*menuFile))
	if err != nil {
		panic("init failed!")
	}

EXIT:
	for {
		sig := <-sc
		logger.WithContext(ctx).Infof("Receive signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.WithContext(ctx).Infof("Server exit")
	time.Sleep(time.Second)
	os.Exit(state)
	// return nil
}

func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	// 解析配置文件
	config.MustLoad(o.ConfigFile)
	if v := o.ModelFile; v != "" {
		config.C.Casbin.Model = v
	}
	if v := o.MenuFile; v != "" {
		config.C.Menu.Data = v
	}
	config.PrintWithJSON(ctx)

	// 初始化日志
	loggerCleanFunc, err := app.InitLogger()
	if err != nil {
		return nil, err
	}

	// injector 依赖注入容器
	injector, injectorCleanFunc, err := app.BuildWireInject()
	if err != nil {
		return nil, err
	}

	// 直接启动服务
	injector.GinEngine.Run(fmt.Sprintf("%s:%d", config.C.HTTP.Host, config.C.HTTP.Port))

	return func() {
		loggerCleanFunc()
		injectorCleanFunc()
	}, nil
}

/*
func InitHttpServer(ctx context.Context, handler http.Handler) func() {
	cfg := config.C.HTTP
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}
*/
