package test

import (
	"buckets"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func initTestService(t *testing.T) (s *buckets.Service, r *MockRepository) {
	ctrl := gomock.NewController(t)
	r = NewMockRepository(ctrl)
	s = buckets.NewService(r)

	return s, r
}

func TestService_FindById(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindById(context.Background(), testBucket1.TenantId, testBucket1.BucketId).
		Return(&testBucket1, nil)
	b, err := service.FindById(context.Background(), testBucket1.TenantId, testBucket1.BucketId)

	require.NoError(t, err)
	assert.Equal(t, testBucket1, *b)
}

func TestService_Search(t *testing.T) {
	service, mockRepo := initTestService(t)
	searchContext := buckets.SearchContext{
		TenantId:             testTenant,
		Name:                 "",
		NbOfReturnedElements: -1,
		NextPageCursor:       "",
		Ids:                  make([]string, 0),
	}
	mockRepo.EXPECT().Search(context.Background(), searchContext).
		Return(testBuckets, nil)

	b, err := service.Search(context.Background(), searchContext)

	require.NoError(t, err)
	assert.ElementsMatch(t, b, testBuckets)
}

func TestService_CreateWhenDoesntExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testBucket1.Name).
		Return(nil, nil)
	mockRepo.EXPECT().CreateOrUpdate(context.Background(), gomock.AssignableToTypeOf(testBucket1))

	id, err := service.Create(context.Background(), testTenant, testBucket1)

	require.NoError(t, err)
	assert.NotNil(t, id)
}

func TestService_CreateErrWhenAlreadyExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testBucket1.Name).
		Return(&testBucket1, nil)

	_, err := service.Create(context.Background(), testTenant, testBucket1)

	assert.Error(t, err)
}

func TestService_UpdateWhenExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testBucket1.Name).
		Return(&testBucket1, nil)
	mockRepo.EXPECT().CreateOrUpdate(context.Background(), gomock.AssignableToTypeOf(testBucket1))

	err := service.Update(context.Background(), testTenant, testBucket1)

	assert.NoError(t, err)
}

func TestService_UpdateErrWhenDoesntExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testBucket1.Name).
		Return(nil, nil)

	err := service.Update(context.Background(), testTenant, testBucket1)

	assert.Error(t, err)
}

func TestService_Delete(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().Delete(context.Background(), testTenant, testBucket1.BucketId).Return(nil)

	err := service.Delete(context.Background(), testTenant, testBucket1.BucketId)

	assert.NoError(t, err)
}
