package api

import (
	"fmt"
	"net/http"
	"time"

	"ccian.cc/really/arweave-api/pkg"
	"ccian.cc/really/arweave-api/pkg/logging"
	"ccian.cc/really/arweave-api/pkg/setting"
	"ccian.cc/really/arweave-api/router"
	"github.com/DeanThompson/ginpprof"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

type Server struct {
	logger *logrus.Entry
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Start(ctx *cli.Context) error {
	defer func() {
		if err := recover(); err != nil {
			if server.logger != nil {
				server.logger.WithField("stack", err).Error("we got a panic")
				return
			}
			logrus.WithField("stack", err).Error("we got a panic")
		}
	}()

	server.loadGlobalConfig(ctx)
	server.initMainLogger()

	serverConfig := setting.Server()
	r := router.InitRouter(serverConfig, server.logger)
	httpServer := &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%d", serverConfig.Port),
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.Timeout.Read) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.Timeout.Write) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	ginpprof.Wrap(r) // pprof

	server.logger.WithField("port", serverConfig.Port).Info("run api service")
	if err := httpServer.ListenAndServe(); err != nil {
		server.logger.WithError(err).Error("ListenAndServe failed")
		return err
	}

	return nil
}

func (server *Server) loadGlobalConfig(ctx *cli.Context) {
	configFile := ctx.String("config")
	if len(configFile) == 0 {
		logrus.Fatal("config file option must be set")
	}
	if err := pkg.Setup(configFile); err != nil {
		logrus.WithError(err).Fatal("parse config file failed")
	}
}

func (server *Server) initMainLogger() {
	server.logger = logging.NewLogEntry("main")
}

func (server *Server) Stop() error {
	return nil
}
