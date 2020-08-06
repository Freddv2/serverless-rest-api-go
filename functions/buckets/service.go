package buckets

import (
	"context"
	"github.com/segmentio/ksuid"
	"time"
)

type Service struct {
	repo repository
}

func NewService(repository repository) *Service {
	return &Service{repository}
}

func (s *Service) FindById(ctx context.Context, tenantId string, bucketId string) (*Bucket, error) {
	bucket, err := s.repo.FindById(ctx, tenantId, bucketId)
	if err != nil {
		return nil, err
	} else if bucket == nil {
		return nil, ErrNotFound
	} else {
		return bucket, err
	}
}

// Create a bucket and return the generated ID
func (s *Service) Create(ctx context.Context, tenantId string, bucket Bucket) (string, error) {
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
		return id, s.repo.CreateOrUpdate(ctx, bucket)
	}
}

func (s *Service) Update(ctx context.Context, tenantId string, bucket Bucket) error {
	bucketAlreadyExists, err := s.bucketAlreadyExists(ctx, tenantId, bucket.Name)
	if err != nil {
		return err
	} else if !bucketAlreadyExists {
		return ErrNotFound
	} else {
		bucket.LastModifiedDate = time.Now()
		return s.repo.CreateOrUpdate(ctx, bucket)
	}
}

func (s *Service) Delete(ctx context.Context, tenantId string, bucketId string) error {
	return s.repo.Delete(ctx, tenantId, bucketId)
}

func (s *Service) Search(ctx context.Context, searchContext SearchContext) ([]Bucket, error) {
	return s.repo.Search(ctx, searchContext)
}

func (s *Service) bucketAlreadyExists(ctx context.Context, tenantId string, name string) (bool, error) {
	bucket, err := s.repo.FindByName(ctx, tenantId, name)
	return bucket != nil, err
}
