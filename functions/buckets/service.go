package buckets

import (
	"context"
	"github.com/segmentio/ksuid"
	"time"
)

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindById(ctx context.Context, tenantId string, bucketId string) (*Bucket, error) {
	bucket, err := s.repository.FindById(ctx, tenantId, bucketId)
	if err != nil {
		return nil, err
	} else if bucket == nil {
		return nil, ErrNotFound
	} else {
		return bucket, err
	}
}

// Create a bucket and return the generated ID
func (s *service) Create(ctx context.Context, tenantId string, bucket Bucket) (string, error) {
	bucketAlreadyExists, err := s.bucketAlreadyExists(ctx, tenantId, bucket.Name)
	if err != nil {
		return "", err
	} else if bucketAlreadyExists {
		return "", ErrAlreadyExists
	} else {
		id := ksuid.New().String()
		bucket.BucketId = id
		bucket.CreationDate = time.Now()
		bucket.LastModifiedDate = time.Now()
		return id, s.repository.CreateOrUpdate(ctx, bucket)
	}
}

func (s *service) Update(ctx context.Context, tenantId string, bucket Bucket) error {
	bucketAlreadyExists, err := s.bucketAlreadyExists(ctx, tenantId, bucket.Name)
	if err != nil {
		return err
	} else if !bucketAlreadyExists {
		return ErrNotFound
	} else {
		bucket.LastModifiedDate = time.Now()
		return s.repository.CreateOrUpdate(ctx, bucket)
	}
}

func (s *service) Delete(ctx context.Context, tenantId string, bucketId string) error {
	return s.repository.Delete(ctx, tenantId, bucketId)
}

func (s *service) Search(ctx context.Context, searchContext SearchContext) ([]Bucket, error) {
	return s.repository.Search(ctx, searchContext)
}

func (s *service) bucketAlreadyExists(ctx context.Context, tenantId string, name string) (bool, error) {
	bucket, err := s.repository.FindByName(ctx, tenantId, name)
	return bucket != nil, err
}
