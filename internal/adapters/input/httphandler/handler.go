package httphandler

import (
	"github.com/gin-gonic/gin"
	"github.com/shanth1/authorization/internal/adapters/input/httphandler/v1"
	"github.com/shanth1/authorization/internal/config"
	"github.com/shanth1/gotools/consts"
)

type httpHandler struct {
	cfg *config.Config
}

func New(cfg *config.Config) *httpHandler {
	return &httpHandler{
		cfg: cfg,
	}
}

func (h *httpHandler) SetupRouter() *gin.Engine {
	if h.cfg.Env == consts.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	e.Static("/static", "./static")

	api := e.Group("/api")

	v1Router := v1.NewRouter()
	v1Router.SetupV1Routes(api)

	return e
}
