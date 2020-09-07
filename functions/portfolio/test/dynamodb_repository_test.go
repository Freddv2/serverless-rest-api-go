package test

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"
	"portfolio"
	"testing"
)

var testRepo *portfolio.DynamoDBRepository

//A local docker instance of dynamodb must be running on port 8000
func createTable() *dynamodb.DynamoDB {
	cfg := &aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region:   aws.String("ca-central-1"),
	}
	sess := session.Must(session.NewSession(cfg))
	db := dynamodb.New(sess)

	createTableReq := &dynamodb.CreateTableInput{
		TableName: aws.String(portfolio.TableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("tenantId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("tenantId"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("RANGE"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
	}
	_, err := db.CreateTable(createTableReq)
	if err != nil {
		panic(err)
	}
	return db
}

func deleteTable(db *dynamodb.DynamoDB) {
	deleteTableReq := &dynamodb.DeleteTableInput{
		TableName: aws.String(portfolio.TableName),
	}
	_, _ = db.DeleteTable(deleteTableReq)
}

func initTestRepository() {
	db := createTable()
	testRepo = portfolio.NewDynamoDBRepository(db)
}

func createTestPortfolios() {
	if err := testRepo.CreateOrUpdate(context.Background(), testPortfolio1); err != nil {
		panic(err)
	}
	if err := testRepo.CreateOrUpdate(context.Background(), testPortfolio2); err != nil {
		panic(err)
	}
}

func TestDynamoDBRepository_FindById(t *testing.T) {
	initTestRepository()
	defer deleteTable(testRepo.DynamoDB)

	createTestPortfolios()

	actualPortfolio, err := testRepo.FindById(context.Background(), testTenant, testPortfolio1.Id)

	require.NoError(t, err)
	require.Equal(t, testPortfolio1, *actualPortfolio)
}

func TestDynamoDBRepository_FindByName(t *testing.T) {
	initTestRepository()
	defer deleteTable(testRepo.DynamoDB)

	createTestPortfolios()

	actualPortfolio, err := testRepo.FindByName(context.Background(), testTenant, testPortfolio1.Name)

	require.NoError(t, err)
	require.Equal(t, testPortfolio1, *actualPortfolio)
}

func TestDynamoDBRepository_CreateOrUpdate(t *testing.T) {
	initTestRepository()
	defer deleteTable(testRepo.DynamoDB)

	createTestPortfolios()

	actualPortfolio, err := testRepo.FindById(context.Background(), testTenant, testPortfolio1.Id)
	require.NoError(t, err)
	require.Equal(t, testPortfolio1, *actualPortfolio)

	actualPortfolio, err = testRepo.FindById(context.Background(), testTenant, testPortfolio2.Id)
	require.NoError(t, err)
	require.Equal(t, testPortfolio2, *actualPortfolio)
}

func TestDynamoDBRepository_Delete(t *testing.T) {
	initTestRepository()
	defer deleteTable(testRepo.DynamoDB)

	createTestPortfolios()

	err := testRepo.Delete(context.Background(), testTenant, testPortfolio1.Id)

	require.NoError(t, err)

	actualPortfolio, err := testRepo.FindById(context.Background(), testTenant, testPortfolio1.Id)

	require.NoError(t, err)
	require.Nil(t, actualPortfolio)
}

func TestDynamoDBRepository_SearchByTenant(t *testing.T) {
	initTestRepository()
	defer deleteTable(testRepo.DynamoDB)

	createTestPortfolios()

	searchContext := portfolio.SearchContext{
		TenantId:             testPortfolio1.TenantId,
		Name:                 "",
		NbOfReturnedElements: 0,
		NextPageCursor:       "",
		Ids:                  nil,
	}

	actualPortfolios, err := testRepo.Search(context.Background(), searchContext)

	require.NoError(t, err)
	require.Len(t, actualPortfolios, 2)
	require.Equal(t, testPortfolios, actualPortfolios)
}

func TestDynamoDBRepository_SearchWithLimit(t *testing.T) {
	initTestRepository()
	defer deleteTable(testRepo.DynamoDB)

	createTestPortfolios()

	searchContext := portfolio.SearchContext{
		TenantId:             testPortfolio1.TenantId,
		Name:                 "",
		NbOfReturnedElements: 1,
		NextPageCursor:       "",
		Ids:                  nil,
	}

	actualPortfolios, err := testRepo.Search(context.Background(), searchContext)

	require.NoError(t, err)
	require.Len(t, actualPortfolios, 1)
	require.Equal(t, actualPortfolios[0], testPortfolio1)
}

func TestDynamoDBRepository_SearchByName(t *testing.T) {
	initTestRepository()
	defer deleteTable(testRepo.DynamoDB)

	createTestPortfolios()

	searchContext := portfolio.SearchContext{
		TenantId:             testPortfolio1.TenantId,
		Name:                 testPortfolio1.Name,
		NbOfReturnedElements: 0,
		NextPageCursor:       "",
		Ids:                  nil,
	}

	actualPortfolios, err := testRepo.Search(context.Background(), searchContext)

	require.NoError(t, err)
	require.Len(t, actualPortfolios, 1)
	require.Equal(t, actualPortfolios[0], testPortfolio1)
}

func TestDynamoDBRepository_SearchByIds(t *testing.T) {
	initTestRepository()
	defer deleteTable(testRepo.DynamoDB)

	createTestPortfolios()

	searchContext := portfolio.SearchContext{
		TenantId:             testPortfolio1.TenantId,
		Name:                 "",
		NbOfReturnedElements: 0,
		NextPageCursor:       "",
		Ids:                  []string{testPortfolio1.Id},
	}

	actualPortfolios, err := testRepo.Search(context.Background(), searchContext)

	require.NoError(t, err)
	require.Len(t, actualPortfolios, 1)
	require.Equal(t, actualPortfolios[0], testPortfolio1)
}
