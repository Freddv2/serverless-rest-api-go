package portfolio

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

func (s *Service) FindById(ctx context.Context, tenantId string, id string) (*Portfolio, error) {
	portfolio, err := s.repo.FindById(ctx, tenantId, id)
	if err != nil {
		return nil, err
	} else if portfolio == nil {
		return nil, ErrNotFound
	} else {
		return portfolio, err
	}
}

// Create a portfolio and return the generated ID
func (s *Service) Create(ctx context.Context, tenantId string, portfolio Portfolio) (string, error) {
	alreadyExists, err := s.alreadyExists(ctx, tenantId, portfolio.Name)
	if err != nil {
		return "", err
	} else if alreadyExists {
		return "", ErrAlreadyExists
	} else {
		id := ksuid.New().String()
		portfolio.Id = id
		portfolio.CreationDate = time.Now()
		portfolio.LastModifiedDate = time.Now()
		return id, s.repo.CreateOrUpdate(ctx, portfolio)
	}
}

func (s *Service) Update(ctx context.Context, tenantId string, portfolio Portfolio) error {
	alreadyExists, err := s.alreadyExists(ctx, tenantId, portfolio.Name)
	if err != nil {
		return err
	} else if !alreadyExists {
		return ErrNotFound
	} else {
		portfolio.LastModifiedDate = time.Now()
		return s.repo.CreateOrUpdate(ctx, portfolio)
	}
}

func (s *Service) Delete(ctx context.Context, tenantId string, id string) error {
	return s.repo.Delete(ctx, tenantId, id)
}

func (s *Service) Search(ctx context.Context, searchContext SearchContext) ([]Portfolio, error) {
	return s.repo.Search(ctx, searchContext)
}

func (s *Service) alreadyExists(ctx context.Context, tenantId string, name string) (bool, error) {
	portfolio, err := s.repo.FindByName(ctx, tenantId, name)
	return portfolio != nil, err
}
