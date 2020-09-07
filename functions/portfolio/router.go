package portfolio

import (
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"log"
)

func NewLambdaRouter(h handler) *ginadapter.GinLambda {
	router := NewRouter(h)
	return ginadapter.New(router)
}

func NewRouter(h handler) *gin.Engine {
	router := gin.Default()
	log.Printf("Defining routes")
	router.GET("/portfolios/:tenantId/:id", h.FindById)
	router.GET("/portfolios/:tenantId", h.Search)
	router.POST("/portfolios/:tenantId", h.Create)
	router.PUT("/portfolios/:tenantId/:id", h.Update)
	router.DELETE("/portfolios/:tenantId/:id", h.Delete)
	log.Printf("Routes defined")
	return router
}
