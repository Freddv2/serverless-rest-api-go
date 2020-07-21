package buckets

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanSearchByName(t *testing.T) {
	service, mockRepo := initTestService(t)

	searchContext := SearchContext{
		TenantId:             testBucket1.TenantId,
		Name:                 "Stocks",
		NbOfReturnedElements: -1,
		NextPageCursor:       "",
		Ids:                  make([]string, 0),
	}
	mockRepo.EXPECT().Search(context.Background(), searchContext).Return([]Bucket{testBucket1}, nil)
	b, err := service.Search(context.Background(), searchContext)

	assert.NoError(t, err)
	assert.Contains(t, b, testBucket1)
}

func TestCanSearchAndLimitTheNbOfReturnedElements(t *testing.T) {
	service, mockRepo := initTestService(t)

	searchContext := SearchContext{
		TenantId:             testBucket1.TenantId,
		Name:                 "",
		NbOfReturnedElements: 1,
		NextPageCursor:       "",
		Ids:                  make([]string, 0),
	}
	mockRepo.EXPECT().Search(context.Background(), searchContext).Return([]Bucket{testBucket1}, nil)
	b, err := service.Search(context.Background(), searchContext)

	assert.NoError(t, err)
	assert.Contains(t, b, testBucket1)
}
