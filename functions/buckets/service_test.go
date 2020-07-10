package buckets

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	bucket1 = &Bucket{
		TenantId:         "dv2",
		Id:               "1",
		Name:             "An ETF Stocks bucket",
		Description:      "Desc1",
		Asset:            []Asset{{"SPY"}, {"QQQ"}, {"VFV"}},
		CreationDate:     time.Now(),
		LastModifiedDate: time.Now(),
	}
	bucket2 = &Bucket{
		TenantId:         "dv2",
		Id:               "2",
		Name:             "An ETF Bonds bucket",
		Description:      "Desc2",
		Asset:            []Asset{{"IEF"}, {"SHY"}},
		CreationDate:     time.Now(),
		LastModifiedDate: time.Now(),
	}
)

func TestCanFindBucketById(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().FindById(context.Background(), bucket1.TenantId, bucket1.Id).Return(bucket1, nil)

	s := NewService(mockRepo)

	bucket, err := s.FindById(context.Background(), bucket1.TenantId, bucket1.Id)

	assert.NoError(t, err)
	assert.Equal(t, bucket1, bucket)
}
