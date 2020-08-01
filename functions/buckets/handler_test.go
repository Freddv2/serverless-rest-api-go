package buckets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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
		FindById(gomock.Any(), testTenant, testBucket1.BucketId).
		Return(&testBucket1, nil)

	resp, err := http.Get(fmt.Sprintf("%s/buckets/%s/%s", ts.URL, testTenant, testBucket1.BucketId))

	require.NotNil(t, resp)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

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

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/buckets/%s", ts.URL, testTenant), nil)
	q := req.URL.Query()
	q.Add("name", testBucket1.Name)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	require.NotNil(t, resp)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	var actualBuckets []Bucket
	_ = json.NewDecoder(resp.Body).Decode(&actualBuckets)
	assert.Contains(t, actualBuckets, testBucket1)
}

func TestHandler_Create(t *testing.T) {
	handler, mockService := initTestHandler(t)
	ts := initTestHttpServer(handler)

	mockService.EXPECT().
		Create(gomock.Any(), testTenant, testBucket1).
		Return(testBucket1.BucketId, nil)
	bucketJson, _ := json.Marshal(testBucket1)
	resp, err := http.Post(fmt.Sprintf("%s/buckets/%s", ts.URL, testTenant), "JSON", bytes.NewBuffer(bucketJson))

	require.NotNil(t, resp)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	id, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, testBucket1.BucketId, string(id))
}

func TestHandler_Update(t *testing.T) {
	handler, mockService := initTestHandler(t)
	ts := initTestHttpServer(handler)

	mockService.EXPECT().
		Update(gomock.Any(), testTenant, testBucket1).
		Return(nil)

	bucketJson, _ := json.Marshal(testBucket1)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/buckets/%s/%s", ts.URL, testTenant, testBucket1.BucketId), bytes.NewBuffer(bucketJson))
	resp, err := http.DefaultClient.Do(req)

	require.NotNil(t, resp)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestHandler_Delete(t *testing.T) {
	handler, mockService := initTestHandler(t)
	ts := initTestHttpServer(handler)

	mockService.EXPECT().
		Delete(gomock.Any(), testTenant, testBucket1.BucketId).
		Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/buckets/%s/%s", ts.URL, testTenant, testBucket1.BucketId), nil)
	resp, err := http.DefaultClient.Do(req)

	require.NotNil(t, resp)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
