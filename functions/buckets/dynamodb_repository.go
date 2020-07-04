package buckets

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	tableName = "BUCKET"
)

type dynamoDBRepository struct {
	session *dynamodb.DynamoDB
}

func NewDynamoDBRepository(dynamoDBClient *dynamodb.DynamoDB) *dynamoDBRepository {
	return &dynamoDBRepository{dynamoDBClient}
}

func (r *dynamoDBRepository) FindById(ctx context.Context, tenantId string, bucketId string) (*Bucket, error) {
	bucket := &Bucket{}
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
			":name": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":bucketName": {
				S: aws.String(bucketName),
			},
		},
		FilterExpression: aws.String(":name = :bucketName"),
		Limit:            aws.Int64(1),
	}
	result, err := r.session.QueryWithContext(ctx, input)
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
	_, err = r.session.PutItemWithContext(ctx, input)
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
	_, err := r.session.DeleteItemWithContext(ctx, input)
	return err
}

func (r *dynamoDBRepository) Search(ctx context.Context, searchCtx SearchContext) ([]Bucket, error) {

	buckets := make([]Bucket, 0)
	key := expression.Key("tenantId").Equal(expression.Value(searchCtx.TenantId))

	filter := expression.ConditionBuilder{}
	if len(searchCtx.Ids) > 0 {
		var bucketIdsConditions = make([]expression.OperandBuilder, len(searchCtx.Ids))
		for i, bucketId := range searchCtx.Ids {
			bucketIdsConditions[i] = expression.Value(bucketId)
		}
		bucketIdsFilter := expression.Name("bucketId").In(bucketIdsConditions[0], bucketIdsConditions[1:]...)
		filter.And(bucketIdsFilter)
	}
	if searchCtx.Name != "" {
		bucketNameFilter := expression.Name("name").Contains(searchCtx.Name)
		filter.And(bucketNameFilter)
	}

	expr, err := expression.NewBuilder().WithKeyCondition(key).WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		KeyConditionExpression: expr.KeyCondition(),
		FilterExpression:       expr.Filter(),
		Limit:                  aws.Int64(int64(searchCtx.NbOfReturnedElements)),
		ScanIndexForward:       aws.Bool(false),
	}

	result, err := r.session.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, buckets); err != nil {
		return nil, err
	}

	return buckets, nil
}
