package buckets

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initTestHandler(t *testing.T) (h *handler, s *MockService) {
	ctrl := gomock.NewController(t)
	s = NewMockService(ctrl)
	h = NewHandler(s)

	return h, s
}

func initTestHttpServer(h *handler) *httptest.Server {
	//Start a test http server with the handler
	return httptest.NewServer(NewRouter(h))
}

func TestHandlerCanFindById(t *testing.T) {
	handler, mockService := initTestHandler(t)
	ts := initTestHttpServer(handler)

	mockService.EXPECT().
		FindById(gomock.Any(), testTenant, testBucket1.Id).
		Return(&testBucket1, nil)

	resp, err := http.Get(fmt.Sprintf("%s/buckets/%s/%s", ts.URL, testTenant, testBucket1.Id))

	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var b Bucket
	_ = json.NewDecoder(resp.Body).Decode(&b)
	assert.Equal(t, testBucket1, b)
}
