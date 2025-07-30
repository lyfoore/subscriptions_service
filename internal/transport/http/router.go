package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lyfoore/subscriptions_service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

type Router struct {
	engine  *gin.Engine
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	engine := gin.Default()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	return &Router{
		engine:  engine,
		handler: handler,
	}
}

func (r *Router) SetupRouter() {
	v1 := r.engine.Group("/api/v1")
	{
		v1.GET("/subscriptions/:id", r.handler.GetSubscription)
		v1.GET("/subscriptions", r.handler.GetSubscriptionsList)
		v1.GET("/subscriptions/aggregate", r.handler.GetSubscriptionsAggregate)
		v1.POST("/subscriptions", r.handler.CreateSubscription)
		v1.PATCH("/subscriptions/:id", r.handler.UpdateSubscription)
		v1.DELETE("/subscriptions/:id", r.handler.DeleteSubscription)
	}

	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r *Router) Run(port string) {
	log.Println("HTTP server is starting on port", port)
	r.engine.Run(port)
}
