package buckets

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	FindById(ctx *gin.Context)
	Search(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

//go:generate mockgen -destination=mock_service.go -package=buckets -self_package buckets . Service
type Service interface {
	FindById(ctx context.Context, tenantId string, bucketId string) (*Bucket, error)
	Search(ctx context.Context, searchContext SearchContext) ([]Bucket, error)
	Create(ctx context.Context, tenantId string, bucket Bucket) (string, error)
	Update(ctx context.Context, tenantId string, bucket Bucket) error
	Delete(ctx context.Context, tenantId string, bucketId string) error
}

//go:generate mockgen -destination=mock_repository.go -package=buckets -self_package buckets . Repository
type Repository interface {
	FindById(ctx context.Context, tenantId string, bucketId string) (*Bucket, error)
	FindByName(ctx context.Context, tenantId string, bucketName string) (*Bucket, error)
	Search(ctx context.Context, searchContext SearchContext) ([]Bucket, error)
	CreateOrUpdate(ctx context.Context, bucket Bucket) error
	Delete(ctx context.Context, tenantId string, bucketId string) error
}
