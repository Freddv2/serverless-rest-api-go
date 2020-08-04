package buckets

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strings"
)

const (
	tableName = "BUCKET"
)

type dynamoDBRepository struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBRepository(dynamoDBClient *dynamodb.DynamoDB) *dynamoDBRepository {
	return &dynamoDBRepository{dynamoDBClient}
}

func (r *dynamoDBRepository) FindById(ctx context.Context, tenantId string, bucketId string) (*Bucket, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"tenantId": {
				S: &tenantId,
			},
			"bucketId": {
				S: &bucketId,
			},
		},
	}

	result, err := r.db.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	bucket := &Bucket{}
	err = dynamodbattribute.UnmarshalMap(result.Item, bucket)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func (r *dynamoDBRepository) FindByName(ctx context.Context, tenantId string, bucketName string) (*Bucket, error) {
	bucket := &Bucket{}
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"tenantId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(tenantId),
					},
				},
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#name": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":bucketName": {
				S: aws.String(bucketName),
			},
		},
		FilterExpression: aws.String("#name = :bucketName"),
		Limit:            aws.Int64(1),
	}
	result, err := r.db.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) > 0 {
		err = dynamodbattribute.UnmarshalMap(result.Items[0], bucket)
		if err != nil {
			return nil, err
		}
	}

	return bucket, nil
}

func (r *dynamoDBRepository) CreateOrUpdate(ctx context.Context, bucket Bucket) error {
	item, err := dynamodbattribute.MarshalMap(bucket)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}
	_, err = r.db.PutItemWithContext(ctx, input)
	return err
}

func (r *dynamoDBRepository) Delete(ctx context.Context, tenantId string, bucketId string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"tenantId": {
				S: aws.String(tenantId),
			},
			"bucketId": {
				S: aws.String(bucketId),
			},
		},
	}
	_, err := r.db.DeleteItemWithContext(ctx, input)
	return err
}

func (r *dynamoDBRepository) Search(ctx context.Context, searchCtx SearchContext) ([]Bucket, error) {

	buckets := make([]Bucket, 0)

	var filterExp string
	expNames := make(map[string]*string)
	expValues := make(map[string]*dynamodb.AttributeValue)
	if searchCtx.Name != "" {
		filterExp += "#name contains :name and"
		expNames["#name"] = aws.String("name")
		expValues[":name"] = &dynamodb.AttributeValue{
			S: aws.String(searchCtx.Name),
		}
	}
	if len(searchCtx.Ids) > 0 {
		filterExp += " #bucketId in :bucketIds"
		expNames["#bucketId"] = aws.String("name")
		expValues[":bucketId"] = &dynamodb.AttributeValue{
			NS: aws.StringSlice(searchCtx.Ids),
		}
	}
	filterExp = strings.TrimSuffix(filterExp, "and")

	query := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"tenantId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(searchCtx.TenantId),
					},
				},
			},
		},
	}
	if filterExp != "" {
		query.ExpressionAttributeNames = expNames
		query.ExpressionAttributeValues = expValues
		query.FilterExpression = aws.String(filterExp)
	}
	if searchCtx.NbOfReturnedElements > 0 {
		query.Limit = aws.Int64(int64(searchCtx.NbOfReturnedElements))
	}
	result, err := r.db.QueryWithContext(ctx, query)
	if err != nil {
		return nil, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &buckets); err != nil {
		return nil, err
	}

	return buckets, nil
}
