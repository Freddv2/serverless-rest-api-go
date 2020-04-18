package buckets

import (
	"context"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBRepository struct {
	session   *dynamodb.DynamoDB
	tableName string
}

func (r *DynamoDBRepository) Get(ctx context.Context, tenantId string, bucketId string) (*Bucket, error) {
	bucket := &Bucket{}
	input := &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"tenantId": {
				S: &tenantId,
			},
			"bucketId": {
				S: &bucketId,
			},
		},
	}

	result, err := r.session.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, bucket)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}
