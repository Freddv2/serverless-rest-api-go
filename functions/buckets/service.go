package buckets

import "context"

type repository interface {
	Get(ctx context.Context, tenantId string, bucketId string) (*Bucket, error)
	GetByName(ctx context.Context, tenantId string, bucketName string) (*Bucket, error)
	CreateOrUpdate(ctx context.Context, bucket Bucket) error
	Delete(ctx context.Context, tenantId string, bucketId string) error
	Query(ctx context.Context, queryBuckets QueryBuckets) ([]Bucket, error)
}

type Service struct {
	Repository repository
}

func (s *Service) findById(ctx context.Context, tenantId string, bucketId string) (*Bucket, error) {
	return s.Repository.Get(ctx, bucketId, tenantId)
}

func (s *Service) create(ctx context.Context, tenantId string, bucketExists Bucket) error {
	bucketAlreadyExists, err := s.bucketAlreadyExists(ctx, tenantId, bucketExists.Name)
	if err != nil {
		return err
	}
	if bucketAlreadyExists {
		return err()
	}
}

func (s *Service) bucketAlreadyExists(ctx context.Context, tenantId string, name string) (bool, error) {
	bucket, err := s.Repository.GetByName(ctx, tenantId, name)
	return bucket != nil, err
}
