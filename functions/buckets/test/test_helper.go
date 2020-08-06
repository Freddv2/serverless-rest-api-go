package test

import (
	"buckets"
	"time"
)

var (
	testTenant  = "dv2"
	testBucket1 = buckets.Bucket{
		TenantId:         testTenant,
		BucketId:         "1",
		Name:             "An ETF Stocks bucket",
		Description:      "Desc1",
		Assets:           []buckets.Asset{{"SPY"}, {"QQQ"}, {"VFV"}},
		CreationDate:     time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
		LastModifiedDate: time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
	}
	testBucket2 = buckets.Bucket{
		TenantId:         testTenant,
		BucketId:         "2",
		Name:             "An ETF Bonds bucket",
		Description:      "Desc2",
		Assets:           []buckets.Asset{{"IEF"}, {"SHY"}},
		CreationDate:     time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
		LastModifiedDate: time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
	}
	testBuckets = []buckets.Bucket{testBucket1, testBucket2}
)
