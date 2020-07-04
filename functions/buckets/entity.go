package buckets

import (
	"time"
)

type Bucket struct {
	TenantId         string    `json:"tenantId"`
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Asset            []Asset   `json:"asset"`
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
