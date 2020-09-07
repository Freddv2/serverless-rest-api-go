package test

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
	"portfolio"
	"testing"
)

func initTestHandler(t *testing.T) (h *portfolio.Handler, s *MockService) {
	ctrl := gomock.NewController(t)
	s = NewMockService(ctrl)
	h = portfolio.NewHandler(s)

	return h, s
}

func initTestHttpServer(h *portfolio.Handler) *httptest.Server {
	//Start a test http server with the handler
	return httptest.NewServer(portfolio.NewRouter(h))
}

func TestHandler_FindById(t *testing.T) {
	handler, mockService := initTestHandler(t)
	ts := initTestHttpServer(handler)

	mockService.EXPECT().
		FindById(gomock.Any(), testTenant, testPortfolio1.Id).
		Return(&testPortfolio1, nil)

	resp, err := http.Get(fmt.Sprintf("%s/portfolios/%s/%s", ts.URL, testTenant, testPortfolio1.Id))

	require.NotNil(t, resp)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	var b portfolio.Portfolio
	_ = json.NewDecoder(resp.Body).Decode(&b)
	assert.Equal(t, testPortfolio1, b)
}

func TestHandler_Search(t *testing.T) {
	handler, mockService := initTestHandler(t)
	ts := initTestHttpServer(handler)

	expectedSC := portfolio.SearchContext{
		TenantId:             testTenant,
		Name:                 testPortfolio1.Name,
		NextPageCursor:       "",
		NbOfReturnedElements: 0,
		Ids:                  []string{},
	}

	mockService.EXPECT().
		Search(gomock.Any(), expectedSC).
		Return([]portfolio.Portfolio{testPortfolio1}, nil)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/portfolios/%s", ts.URL, testTenant), nil)
	q := req.URL.Query()
	q.Add("name", testPortfolio1.Name)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	require.NotNil(t, resp)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	var actualBuckets []portfolio.Portfolio
	_ = json.NewDecoder(resp.Body).Decode(&actualBuckets)
	assert.Contains(t, actualBuckets, testPortfolio1)
}

func TestHandler_Create(t *testing.T) {
	handler, mockService := initTestHandler(t)
	ts := initTestHttpServer(handler)

	mockService.EXPECT().
		Create(gomock.Any(), testTenant, testPortfolio1).
		Return(testPortfolio1.Id, nil)
	bucketJson, _ := json.Marshal(testPortfolio1)
	resp, err := http.Post(fmt.Sprintf("%s/portfolios/%s", ts.URL, testTenant), "JSON", bytes.NewBuffer(bucketJson))

	require.NotNil(t, resp)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	id, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, testPortfolio1.Id, string(id))
}

func TestHandler_Update(t *testing.T) {
	handler, mockService := initTestHandler(t)
	ts := initTestHttpServer(handler)

	mockService.EXPECT().
		Update(gomock.Any(), testTenant, testPortfolio1).
		Return(nil)

	bucketJson, _ := json.Marshal(testPortfolio1)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/portfolios/%s/%s", ts.URL, testTenant, testPortfolio1.Id), bytes.NewBuffer(bucketJson))
	resp, err := http.DefaultClient.Do(req)

	require.NotNil(t, resp)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestHandler_Delete(t *testing.T) {
	handler, mockService := initTestHandler(t)
	ts := initTestHttpServer(handler)

	mockService.EXPECT().
		Delete(gomock.Any(), testTenant, testPortfolio1.Id).
		Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/portfolios/%s/%s", ts.URL, testTenant, testPortfolio1.Id), nil)
	resp, err := http.DefaultClient.Do(req)

	require.NotNil(t, resp)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
