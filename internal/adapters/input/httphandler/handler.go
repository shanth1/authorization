package httphandler

import (
	"github.com/gin-gonic/gin"
	"github.com/shanth1/authorization/internal/config"
)

type httpHandler struct {
	cfg *config.Config
}

func New(cfg *config.Config) *httpHandler {
	return &httpHandler{}
}

func (h *httpHandler) SetupRouter() *gin.Engine {
	e := gin.New()

	e.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	e.Static("/", "./static")

	api := e.Group("/api")
	h.setupV1Routes(api)

	return e
}

// TODO:
func (h *httpHandler) setupV1Routes(r *gin.RouterGroup) {
	common := r.Group("/v1")
	{
		common.GET("/.well-known/jwks.json", nil)
		common.GET("/providers/:id", nil)               // get providers by client id
		common.GET("/providers/telegram/callback", nil) //
		common.GET("/providers/google/callback", nil)   //
		common.POST("/authorize", nil)
		common.POST("/token", nil)
	}

	admin := common.Group("/admin")
	admin.Use(nil) // api key middlewary
	{
		admin.GET("/clients", nil)        // get all clients
		admin.POST("/clients", nil)       // add net client
		admin.PUT("/clients/:id", nil)    // update client by id
		admin.DELETE("/clients/:id", nil) // delete client by id
	}
}
