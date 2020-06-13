package buckets

import "time"

type Bucket struct {
	TenantId         string    `json:"tenantId"`
	BucketId         string    `json:"bucketId"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Asset            []Asset   `json:"asset"`
	CreationDate     time.Time `json:"creationDate"`
	LastModifiedDate time.Time `json:"lastModifiedDate"`
}

type Asset struct {
	Symbol string `json:"symbol"`
}

type QueryBuckets struct {
	TenantId             string   `json:"tenantId"`
	BucketName           string   `json:"bucketName"`
	NbOfReturnedElements int64    `json:"nbOfReturnedElements"`
	NextPageCursor       string   `json:"nextPageCursor"`
	BucketIds            []string `json:"bucketIds"`
}

type NotFoundError struct{}

func (error *NotFoundError) Error() string {
	return "Not Found"
}

type AlreadyExistsError struct{}

func (error *AlreadyExistsError) Error() string {
	return "Already Exists"
}
