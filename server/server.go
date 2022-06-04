package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"mock_net/middleware"
	"mock_net/server/router"
	"mock_net/setting"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	ctx        context.Context
	cancel     context.CancelFunc
}

func New() *Server {

	return &Server{}
}

var _ ContainerRunnable = &Server{}


func (server *Server) Run(ctx context.Context) error  {
	server.ctx, server.cancel = context.WithCancel(ctx)
	return server.Start()
}


func (server *Server) Start() error {
	setting.LoadLocalConfig()

	r := gin.New()

	r.GET("/reload", func(context *gin.Context) {
		go func() {
			stopChannel<-true
		}()
		result := map[string]string{}
		context.JSON(http.StatusOK, result)
	})

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.NoFundHandle(setting.GetApiInfo()))

	if len(setting.GetStaticFilePath()) != 0 {
		r.StaticFS("/file", gin.Dir(setting.GetStaticFilePath(), true))
	}

	router.InitApi(r, setting.GetApiInfo())

	srv := &http.Server{
		Addr:    setting.GetStartAddress(), // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
		Handler: r,
	}

	server.httpServer = srv

	return server.listenAndServe()
}

func (server *Server) listenAndServe() error {
	var g errgroup.Group
	g.Go(func() error {
		<-server.ctx.Done()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		return server.httpServer.Shutdown(ctx)
	})
	g.Go(func() error {
		log.Println("Server  Start...")
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Error Shutdown Server ...")
			return err
		}
		log.Println("Normal Shutdown Server ...")
		return nil
	})
	return g.Wait()
}


func (server *Server) Close() {
	server.cancel()
}
