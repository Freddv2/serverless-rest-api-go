package buckets

import "time"

type Bucket struct {
	TenantId         string  `json:"tenantId"`
	BucketId         string  `json:"bucketId"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Asset            []Asset `json:"assets"`
	CreationDate     time.Time
	LastModifiedDate time.Time
}

type Asset struct {
	Symbol string `json:"symbol"`
}
