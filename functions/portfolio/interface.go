package portfolio

import (
	"context"
	"github.com/gin-gonic/gin"
)

type handler interface {
	FindById(ctx *gin.Context)
	Search(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

//go:generate mockgen -destination=mock_service.go -package=portfolio -self_package buckets . Service
type service interface {
	FindById(ctx context.Context, tenantId string, bucketId string) (*Portfolio, error)
	Search(ctx context.Context, searchContext SearchContext) ([]Portfolio, error)
	Create(ctx context.Context, tenantId string, bucket Portfolio) (string, error)
	Update(ctx context.Context, tenantId string, bucket Portfolio) error
	Delete(ctx context.Context, tenantId string, bucketId string) error
}

//go:generate mockgen -destination=mock_repository.go -package=C:\Users\Fredd\IdeaProjects\okiok\serverless-rest-api-go\functions\portfolio\interface.go -self_package buckets . Repository
type repository interface {
	FindById(ctx context.Context, tenantId string, bucketId string) (*Portfolio, error)
	FindByName(ctx context.Context, tenantId string, bucketName string) (*Portfolio, error)
	Search(ctx context.Context, searchContext SearchContext) ([]Portfolio, error)
	CreateOrUpdate(ctx context.Context, bucket Portfolio) error
	Delete(ctx context.Context, tenantId string, bucketId string) error
}
