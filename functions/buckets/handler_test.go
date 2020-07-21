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

func TestHandler_FindById(t *testing.T) {
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

func TestHandler_Search(t *testing.T) {
	handler, mockService := initTestHandler(t)
	ts := initTestHttpServer(handler)

	expectedSC := SearchContext{
		TenantId:             testTenant,
		Name:                 testBucket1.Name,
		NextPageCursor:       "",
		NbOfReturnedElements: 0,
		Ids:                  []string{},
	}

	mockService.EXPECT().
		Search(gomock.Any(), expectedSC).
		Return([]Bucket{testBucket1}, nil)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/buckets/%s", ts.URL, testTenant), nil)
	q := req.URL.Query()
	q.Add("name", testBucket1.Name)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
