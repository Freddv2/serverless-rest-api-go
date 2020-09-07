package test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"portfolio"
	"testing"
)

func initTestService(t *testing.T) (s *portfolio.Service, r *MockRepository) {
	ctrl := gomock.NewController(t)
	r = NewMockRepository(ctrl)
	s = portfolio.NewService(r)

	return s, r
}

func TestService_FindById(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindById(context.Background(), testPortfolio1.TenantId, testPortfolio1.Id).
		Return(&testPortfolio1, nil)
	b, err := service.FindById(context.Background(), testPortfolio1.TenantId, testPortfolio1.Id)

	require.NoError(t, err)
	assert.Equal(t, testPortfolio1, *b)
}

func TestService_Search(t *testing.T) {
	service, mockRepo := initTestService(t)
	searchContext := portfolio.SearchContext{
		TenantId:             testTenant,
		Name:                 "",
		NbOfReturnedElements: -1,
		NextPageCursor:       "",
		Ids:                  make([]string, 0),
	}
	mockRepo.EXPECT().Search(context.Background(), searchContext).
		Return(testPortfolios, nil)

	b, err := service.Search(context.Background(), searchContext)

	require.NoError(t, err)
	assert.ElementsMatch(t, b, testPortfolios)
}

func TestService_CreateWhenDoesntExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testPortfolio1.Name).
		Return(nil, nil)
	mockRepo.EXPECT().CreateOrUpdate(context.Background(), gomock.AssignableToTypeOf(testPortfolio1))

	id, err := service.Create(context.Background(), testTenant, testPortfolio1)

	require.NoError(t, err)
	assert.NotNil(t, id)
}

func TestService_CreateErrWhenAlreadyExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testPortfolio1.Name).
		Return(&testPortfolio1, nil)

	_, err := service.Create(context.Background(), testTenant, testPortfolio1)

	assert.Error(t, err)
}

func TestService_UpdateWhenExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testPortfolio1.Name).
		Return(&testPortfolio1, nil)
	mockRepo.EXPECT().CreateOrUpdate(context.Background(), gomock.AssignableToTypeOf(testPortfolio1))

	err := service.Update(context.Background(), testTenant, testPortfolio1)

	assert.NoError(t, err)
}

func TestService_UpdateErrWhenDoesntExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testPortfolio1.Name).
		Return(nil, nil)

	err := service.Update(context.Background(), testTenant, testPortfolio1)

	assert.Error(t, err)
}

func TestService_Delete(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().Delete(context.Background(), testTenant, testPortfolio1.Id).Return(nil)

	err := service.Delete(context.Background(), testTenant, testPortfolio1.Id)

	assert.NoError(t, err)
}
