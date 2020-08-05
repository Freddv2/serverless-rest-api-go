package buckets

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"
	"testing"
)

var testRepo *dynamoDBRepository

func initLocalDynamoDB() *dynamodb.DynamoDB {
	cfg := &aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region:   aws.String("ca-central-1"),
	}
	sess := session.Must(session.NewSession(cfg))
	db := dynamodb.New(sess)

	createTableDef := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("tenantId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("bucketId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("tenantId"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("bucketId"),
				KeyType:       aws.String("RANGE"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
	}
	_, err := db.CreateTable(createTableDef)
	if err != nil {
		panic(err)
	}
	return db
}

func destroyTestDynamoDB(db *dynamodb.DynamoDB) {
	deleteTableReq := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}
	_, _ = db.DeleteTable(deleteTableReq)
}

func initTestRepository() {
	db := initLocalDynamoDB()
	testRepo = NewDynamoDBRepository(db)
}

func createTestBuckets() {
	if err := testRepo.CreateOrUpdate(context.Background(), testBucket1); err != nil {
		panic(err)
	}
	if err := testRepo.CreateOrUpdate(context.Background(), testBucket2); err != nil {
		panic(err)
	}
}

func TestDynamoDBRepository_FindById(t *testing.T) {
	initTestRepository()
	defer destroyTestDynamoDB(testRepo.db)

	createTestBuckets()

	actualBucket, err := testRepo.FindById(context.Background(), testTenant, testBucket1.BucketId)

	require.NoError(t, err)
	require.Equal(t, testBucket1, *actualBucket)
}

func TestDynamoDBRepository_FindByName(t *testing.T) {
	initTestRepository()
	defer destroyTestDynamoDB(testRepo.db)

	createTestBuckets()

	actualBucket, err := testRepo.FindByName(context.Background(), testTenant, testBucket1.Name)

	require.NoError(t, err)
	require.Equal(t, testBucket1, *actualBucket)
}

func TestDynamoDBRepository_CreateOrUpdate(t *testing.T) {
	initTestRepository()
	defer destroyTestDynamoDB(testRepo.db)

	createTestBuckets()

	actualBucket, err := testRepo.FindById(context.Background(), testTenant, testBucket1.BucketId)
	require.NoError(t, err)
	require.Equal(t, testBucket1, *actualBucket)

	actualBucket, err = testRepo.FindById(context.Background(), testTenant, testBucket2.BucketId)
	require.NoError(t, err)
	require.Equal(t, testBucket2, *actualBucket)
}

func TestDynamoDBRepository_Delete(t *testing.T) {
	initTestRepository()
	defer destroyTestDynamoDB(testRepo.db)

	createTestBuckets()

	err := testRepo.Delete(context.Background(), testTenant, testBucket1.BucketId)

	require.NoError(t, err)

	actualBucket, err := testRepo.FindById(context.Background(), testTenant, testBucket1.BucketId)

	require.NoError(t, err)
	require.Nil(t, actualBucket)
}

func TestDynamoDBRepository_SearchByTenant(t *testing.T) {
	initTestRepository()
	defer destroyTestDynamoDB(testRepo.db)

	createTestBuckets()

	searchContext := SearchContext{
		TenantId:             testBucket1.TenantId,
		Name:                 "",
		NbOfReturnedElements: 0,
		NextPageCursor:       "",
		Ids:                  nil,
	}

	actualBuckets, err := testRepo.Search(context.Background(), searchContext)

	require.NoError(t, err)
	require.Len(t, actualBuckets, 2)
	require.Equal(t, testBuckets, actualBuckets)
}

func TestDynamoDBRepository_SearchWithLimit(t *testing.T) {
	initTestRepository()
	defer destroyTestDynamoDB(testRepo.db)

	createTestBuckets()

	searchContext := SearchContext{
		TenantId:             testBucket1.TenantId,
		Name:                 "",
		NbOfReturnedElements: 1,
		NextPageCursor:       "",
		Ids:                  nil,
	}

	actualBuckets, err := testRepo.Search(context.Background(), searchContext)

	require.NoError(t, err)
	require.Len(t, actualBuckets, 1)
	require.Equal(t, actualBuckets[0], testBucket1)
}

func TestDynamoDBRepository_SearchByName(t *testing.T) {
	initTestRepository()
	defer destroyTestDynamoDB(testRepo.db)

	createTestBuckets()

	searchContext := SearchContext{
		TenantId:             testBucket1.TenantId,
		Name:                 testBucket1.Name,
		NbOfReturnedElements: 0,
		NextPageCursor:       "",
		Ids:                  nil,
	}

	actualBuckets, err := testRepo.Search(context.Background(), searchContext)

	require.NoError(t, err)
	require.Len(t, actualBuckets, 1)
	require.Equal(t, actualBuckets[0], testBucket1)
}

func TestDynamoDBRepository_SearchByIds(t *testing.T) {
	initTestRepository()
	defer destroyTestDynamoDB(testRepo.db)

	createTestBuckets()

	searchContext := SearchContext{
		TenantId:             testBucket1.TenantId,
		Name:                 "",
		NbOfReturnedElements: 0,
		NextPageCursor:       "",
		Ids:                  []string{testBucket1.BucketId},
	}

	actualBuckets, err := testRepo.Search(context.Background(), searchContext)

	require.NoError(t, err)
	require.Len(t, actualBuckets, 1)
	require.Equal(t, actualBuckets[0], testBucket1)
}
