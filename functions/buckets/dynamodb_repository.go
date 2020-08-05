package buckets

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
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
	filterExp := ""
	expNames := make(map[string]*string)
	expValues := make(map[string]*dynamodb.AttributeValue)

	filterExp = "#tenantId = :tenantId"
	expNames["#tenantId"] = aws.String("tenantId")
	expValues[":tenantId"] = &dynamodb.AttributeValue{S: aws.String(searchCtx.TenantId)}

	if searchCtx.Name != "" {
		filterExp += "and contains(#name, :name)"
		expNames["#name"] = aws.String("name")
		expValues[":name"] = &dynamodb.AttributeValue{S: aws.String(searchCtx.Name)}
	}
	if len(searchCtx.Ids) > 0 {
		keys := make([]string, len(searchCtx.Ids))
		for i, id := range searchCtx.Ids {
			keys[i] = ":bucketId" + strconv.Itoa(i)
			expValues[keys[i]] = &dynamodb.AttributeValue{S: aws.String(id)}
		}
		filterExp += " and #bucketId in(" + strings.Join(keys, ",") + ")"
		expNames["#bucketId"] = aws.String("bucketId")
	}
	scan := &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		FilterExpression:          aws.String(filterExp),
		ExpressionAttributeNames:  expNames,
		ExpressionAttributeValues: expValues,
	}
	if filterExp != "" {
		scan.FilterExpression = aws.String(filterExp)
	}
	if searchCtx.NbOfReturnedElements > 0 {
		scan.Limit = aws.Int64(int64(searchCtx.NbOfReturnedElements))
	}
	result, err := r.db.ScanWithContext(ctx, scan)
	if err != nil {
		return nil, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &buckets); err != nil {
		return nil, err
	}

	return buckets, nil
}
