package router

import (
	"github.com/gin-gonic/gin"

	"github.com/aldotp/rate-limiter/config"
	"github.com/aldotp/rate-limiter/internal/handler"
	"github.com/aldotp/rate-limiter/internal/middleware"
)

type Router struct {
	*gin.Engine
	cfg *config.Config
	mdl middleware.MiddlewareInterface
	hdl handler.HandlerInterface
}

func NewRouter(cfg *config.Config, mdl middleware.MiddlewareInterface, hdl handler.HandlerInterface) *Router {
	return &Router{
		cfg: cfg,
		mdl: mdl,
		hdl: hdl,
	}
}

func (r *Router) SetupRouter() {
	if r.cfg.App.Env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Engine = gin.New()
	r.Engine.Use(gin.Recovery())
	r.Engine.Use(r.mdl.Log())

	r.Engine.GET("/ping", r.mdl.RateLimiter(), r.hdl.Ping)
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}

var _ RouterInterface = (*Router)(nil)
