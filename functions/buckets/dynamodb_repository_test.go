package buckets

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
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

func TestDynamoDBRepository_FindById(t *testing.T) {
	initTestRepository()
	defer destroyTestDynamoDB(testRepo.db)

	if err := testRepo.CreateOrUpdate(context.Background(), testBucket1); err != nil {
		panic(err)
	}
	if err := testRepo.CreateOrUpdate(context.Background(), testBucket2); err != nil {
		panic(err)
	}

	actualBucket, err := testRepo.FindById(context.Background(), testTenant, testBucket1.BucketId)

	require.NoError(t, err)
	require.Equal(t, testBucket1, *actualBucket)
}

func TestCanSearchByName(t *testing.T) {
	service, mockRepo := initTestService(t)

	searchContext := SearchContext{
		TenantId:             testBucket1.TenantId,
		Name:                 "Stocks",
		NbOfReturnedElements: -1,
		NextPageCursor:       "",
		Ids:                  make([]string, 0),
	}
	mockRepo.EXPECT().Search(context.Background(), searchContext).Return([]Bucket{testBucket1}, nil)
	b, err := service.Search(context.Background(), searchContext)

	assert.NoError(t, err)
	assert.Contains(t, b, testBucket1)
}

func TestCanSearchAndLimitTheNbOfReturnedElements(t *testing.T) {
	service, mockRepo := initTestService(t)

	searchContext := SearchContext{
		TenantId:             testBucket1.TenantId,
		Name:                 "",
		NbOfReturnedElements: 1,
		NextPageCursor:       "",
		Ids:                  make([]string, 0),
	}
	mockRepo.EXPECT().Search(context.Background(), searchContext).Return([]Bucket{testBucket1}, nil)
	b, err := service.Search(context.Background(), searchContext)

	assert.NoError(t, err)
	assert.Contains(t, b, testBucket1)
}
