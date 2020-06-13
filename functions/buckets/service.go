package buckets

import (
	"context"
)

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

func (s *Service) FindById(ctx context.Context, tenantId string, bucketId string) (*Bucket, error) {
	bucket, err := s.Repository.Get(ctx, bucketId, tenantId)
	if err != nil {
		return nil, err
	} else if bucket == nil {
		return nil, &NotFoundError{}
	} else {
		return bucket, err
	}
}

func (s *Service) Create(ctx context.Context, tenantId string, bucket Bucket) error {
	bucketAlreadyExists, err := s.bucketAlreadyExists(ctx, tenantId, bucket.Name)
	if err != nil {
		return err
	}
	if bucketAlreadyExists {
		return &AlreadyExistsError{}
	}
	return nil
}

func (s *Service) Update(ctx context.Context, tenantId string, bucket Bucket) error {
	bucketAlreadyExists, err := s.bucketAlreadyExists(ctx, tenantId, bucket.Name)
	if err != nil {
		return err
	} else if !bucketAlreadyExists {
		return &NotFoundError{}
	} else {
		return nil
	}
}

func (s *Service) Delete(ctx context.Context, tenantId string, bucketId string) error {
	return s.Repository.Delete(ctx, tenantId, bucketId)
}

func (s *Service) bucketAlreadyExists(ctx context.Context, tenantId string, name string) (bool, error) {
	bucket, err := s.Repository.GetByName(ctx, tenantId, name)
	return bucket != nil, err
}
