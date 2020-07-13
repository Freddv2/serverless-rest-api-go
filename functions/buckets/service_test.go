package buckets

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	bucket1 = Bucket{
		TenantId:         "dv2",
		Id:               "1",
		Name:             "An ETF Stocks bucket",
		Description:      "Desc1",
		Asset:            []Asset{{"SPY"}, {"QQQ"}, {"VFV"}},
		CreationDate:     time.Now(),
		LastModifiedDate: time.Now(),
	}
	bucket2 = Bucket{
		TenantId:         "dv2",
		Id:               "2",
		Name:             "An ETF Bonds bucket",
		Description:      "Desc2",
		Asset:            []Asset{{"IEF"}, {"SHY"}},
		CreationDate:     time.Now(),
		LastModifiedDate: time.Now(),
	}
	buckets = []Bucket{bucket1, bucket2}
)

func initMock(t *testing.T) (s *service, r *MockRepository) {

	ctrl := gomock.NewController(t)
	r = NewMockRepository(ctrl)

	s = NewService(r)

	return s, r
}

func TestCanFindBucketById(t *testing.T) {
	service, mockRepo := initMock(t)

	mockRepo.EXPECT().FindById(context.Background(), bucket1.TenantId, bucket1.Id).Return(&bucket1, nil)
	b, err := service.FindById(context.Background(), bucket1.TenantId, bucket1.Id)

	assert.NoError(t, err)
	assert.Equal(t, bucket1, *b)
}

func TestCanFindAllBucketByTenant(t *testing.T) {
	service, mockRepo := initMock(t)

	searchContext := SearchContext{
		TenantId:             bucket1.TenantId,
		Name:                 "",
		NbOfReturnedElements: -1,
		NextPageCursor:       "",
		Ids:                  make([]string, 0),
	}
	mockRepo.EXPECT().Search(context.Background(), searchContext).Return(buckets, nil)
	b, err := service.Search(context.Background(), searchContext)

	assert.NoError(t, err)
	assert.ElementsMatch(t, b, buckets)
}

func TestCanFindBucketsByName(t *testing.T) {
	service, mockRepo := initMock(t)

	searchContext := SearchContext{
		TenantId:             bucket1.TenantId,
		Name:                 "Stocks",
		NbOfReturnedElements: -1,
		NextPageCursor:       "",
		Ids:                  make([]string, 0),
	}
	mockRepo.EXPECT().Search(context.Background(), searchContext).Return([]Bucket{bucket1}, nil)
	b, err := service.Search(context.Background(), searchContext)

	assert.NoError(t, err)
	assert.Contains(t, b, bucket1)
}
