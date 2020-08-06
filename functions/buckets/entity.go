package buckets

import (
	"time"
)

type Bucket struct {
	TenantId         string    `json:"tenantId"`
	BucketId         string    `json:"bucketId"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Assets           []Asset   `json:"assets"`
	CreationDate     time.Time `json:"creationDate"`
	LastModifiedDate time.Time `json:"lastModifiedDate"`
}

type Asset struct {
	Symbol string `json:"symbol"`
}

type SearchContext struct {
	TenantId             string   `json:"tenantId"`
	Name                 string   `json:"name"`
	NbOfReturnedElements int      `json:"nbOfReturnedElements"`
	NextPageCursor       string   `json:"nextPageCursor"`
	Ids                  []string `json:"ids"`
}
