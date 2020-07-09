package buckets

import (
	"context"
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	FindById(ctx context.Context, tenantId string, bucketId string) (*Bucket, error)
	FindByName(ctx context.Context, tenantId string, bucketName string) (*Bucket, error)
	Search(ctx context.Context, searchContext SearchContext) ([]Bucket, error)
	CreateOrUpdate(ctx context.Context, bucket Bucket) error
	Delete(ctx context.Context, tenantId string, bucketId string) error
}

type Service interface {
	FindById(ctx context.Context, tenantId string, bucketId string) (*Bucket, error)
	Search(ctx context.Context, searchContext SearchContext) ([]Bucket, error)
	Create(ctx context.Context, tenantId string, bucket Bucket) (string, error)
	Update(ctx context.Context, tenantId string, bucket Bucket) error
	Delete(ctx context.Context, tenantId string, bucketId string) error
}

type Handler interface {
	FindById(ctx *gin.Context)
	Search(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}
