package v1

import "github.com/gin-gonic/gin"

type v1Router struct {
	path string
}

func NewRouter() *v1Router {
	return &v1Router{
		path: "/v1",
	}
}

// TODO:
func (router *v1Router) SetupV1Routes(rg *gin.RouterGroup) {
	public := rg.Group(router.path)
	{
		public.GET("/.well-known/jwks.json", nil)
		public.POST("/users/register", nil)
		public.GET("/clients/:id/providers", nil)
		public.GET("/providers/telegram/callback", nil)
		public.GET("/providers/google/callback", nil)
		public.POST("/authorize", nil)
		public.POST("/token", nil)
	}

	protected := rg.Group("/v1")
	protected.Use(nil) // TODO: jwt middleware
	{
		protected.POST("/users/endsession", nil)
	}

	admin := public.Group("/admin")
	admin.Use(nil) // TODO: admin middleware
	{
		admin.GET("/clients", nil)
		admin.POST("/clients", nil)
		admin.PUT("/clients/:id", nil)
		admin.DELETE("/clients/:id", nil)
	}
}
