package buckets

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initTestService(t *testing.T) (s *service, r *MockRepository) {
	ctrl := gomock.NewController(t)
	r = NewMockRepository(ctrl)
	s = NewService(r)

	return s, r
}

func TestServiceCanFindById(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindById(context.Background(), testBucket1.TenantId, testBucket1.Id).
		Return(&testBucket1, nil)

	b, err := service.FindById(context.Background(), testBucket1.TenantId, testBucket1.Id)

	assert.NoError(t, err)
	assert.Equal(t, testBucket1, *b)
}

func TestServiceCanSearch(t *testing.T) {
	service, mockRepo := initTestService(t)
	searchContext := SearchContext{
		TenantId:             testTenant,
		Name:                 "",
		NbOfReturnedElements: -1,
		NextPageCursor:       "",
		Ids:                  make([]string, 0),
	}
	mockRepo.EXPECT().Search(context.Background(), searchContext).
		Return(testBuckets, nil)

	b, err := service.Search(context.Background(), searchContext)

	assert.NoError(t, err)
	assert.ElementsMatch(t, b, testBuckets)
}

func TestCreateWhenDoesntExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testBucket1.Name).
		Return(nil, nil)
	mockRepo.EXPECT().CreateOrUpdate(context.Background(), gomock.AssignableToTypeOf(testBucket1))

	id, err := service.Create(context.Background(), testTenant, testBucket1)

	assert.NoError(t, err)
	assert.NotNil(t, id)
}

func TestErrWhenCreateAndDoesntExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testBucket1.Name).
		Return(&testBucket1, nil)

	_, err := service.Create(context.Background(), testTenant, testBucket1)

	assert.Error(t, err)
}

func TestUpdateWhenExists(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testBucket1.Name).
		Return(&testBucket1, nil)
	mockRepo.EXPECT().CreateOrUpdate(context.Background(), gomock.AssignableToTypeOf(testBucket1))

	err := service.Update(context.Background(), testTenant, testBucket1)

	assert.NoError(t, err)
}

func TestErrWhenUpdateAndDoesntExist(t *testing.T) {
	service, mockRepo := initTestService(t)
	mockRepo.EXPECT().FindByName(context.Background(), testTenant, testBucket1.Name).
		Return(nil, nil)

	err := service.Update(context.Background(), testTenant, testBucket1)

	assert.Error(t, err)
}
