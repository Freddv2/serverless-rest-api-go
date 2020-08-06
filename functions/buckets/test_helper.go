package buckets

import "time"

var (
	testTenant  = "dv2"
	testBucket1 = Bucket{
		TenantId:         testTenant,
		BucketId:         "1",
		Name:             "An ETF Stocks bucket",
		Description:      "Desc1",
		Assets:           []Asset{{"SPY"}, {"QQQ"}, {"VFV"}},
		CreationDate:     time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
		LastModifiedDate: time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
	}
	testBucket2 = Bucket{
		TenantId:         testTenant,
		BucketId:         "2",
		Name:             "An ETF Bonds bucket",
		Description:      "Desc2",
		Assets:           []Asset{{"IEF"}, {"SHY"}},
		CreationDate:     time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
		LastModifiedDate: time.Now().Truncate(0), //Truncate(0) means removing the monotonic time which causes problems with assert
	}
	testBuckets = []Bucket{testBucket1, testBucket2}
)
