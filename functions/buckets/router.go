package buckets

import (
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"log"
)

func NewLambdaRouter(h Handler) *ginadapter.GinLambda {
	router := NewRouter(h)
	return ginadapter.New(router)
}

func NewRouter(h Handler) *gin.Engine {
	router := gin.Default()
	log.Printf("Defining routes")
	router.GET("/buckets/:tenantId/:bucketId", h.FindById)
	router.GET("/buckets/:tenantId", h.Search)
	router.POST("/buckets/:tenantId", h.Create)
	router.PUT("/buckets/:tenantId/:bucketId", h.Update)
	router.DELETE("/buckets/:tenantId/:bucketId", h.Delete)
	log.Printf("Routes defined")
	return router
}
