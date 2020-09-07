package portfolio

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
	"strings"
)

const (
	TableName = "PORTFOLIO"
)

type DynamoDBRepository struct {
	DynamoDB *dynamodb.DynamoDB
}

func NewDynamoDBRepository(dynamoDBClient *dynamodb.DynamoDB) *DynamoDBRepository {
	return &DynamoDBRepository{dynamoDBClient}
}

func (r *DynamoDBRepository) FindById(ctx context.Context, tenantId string, id string) (*Portfolio, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"tenantId": {
				S: &tenantId,
			},
			"id": {
				S: &id,
			},
		},
	}

	result, err := r.DynamoDB.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	portfolio := &Portfolio{}
	err = dynamodbattribute.UnmarshalMap(result.Item, portfolio)
	if err != nil {
		return nil, err
	}

	return portfolio, nil
}

func (r *DynamoDBRepository) FindByName(ctx context.Context, tenantId string, name string) (*Portfolio, error) {
	var portfolio *Portfolio
	input := &dynamodb.QueryInput{
		TableName: aws.String(TableName),
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
			":name": {
				S: aws.String(name),
			},
		},
		FilterExpression: aws.String("#name = :name"),
		Limit:            aws.Int64(1),
	}
	result, err := r.DynamoDB.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) > 0 {
		portfolio = &Portfolio{}
		err = dynamodbattribute.UnmarshalMap(result.Items[0], portfolio)
		if err != nil {
			return nil, err
		}
	}

	return portfolio, nil
}

func (r *DynamoDBRepository) CreateOrUpdate(ctx context.Context, portfolio Portfolio) error {
	item, err := dynamodbattribute.MarshalMap(portfolio)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	}
	_, err = r.DynamoDB.PutItemWithContext(ctx, input)
	return err
}

func (r *DynamoDBRepository) Delete(ctx context.Context, tenantId string, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"tenantId": {
				S: aws.String(tenantId),
			},
			"id": {
				S: aws.String(id),
			},
		},
	}
	_, err := r.DynamoDB.DeleteItemWithContext(ctx, input)
	return err
}

func (r *DynamoDBRepository) Search(ctx context.Context, searchCtx SearchContext) ([]Portfolio, error) {

	portfolios := make([]Portfolio, 0)
	filterExp := ""
	expNames := make(map[string]*string)
	expValues := make(map[string]*dynamodb.AttributeValue)

	filterExp = "#tenantId = :tenantId"
	expNames["#tenantId"] = aws.String("tenantId")
	expValues[":tenantId"] = &dynamodb.AttributeValue{S: aws.String(searchCtx.TenantId)}

	if searchCtx.Name != "" {
		filterExp += " and contains(#name, :name)"
		expNames["#name"] = aws.String("name")
		expValues[":name"] = &dynamodb.AttributeValue{S: aws.String(searchCtx.Name)}
	}
	if len(searchCtx.Ids) > 0 {
		keys := make([]string, len(searchCtx.Ids))
		for i, id := range searchCtx.Ids {
			keys[i] = ":portfolioId" + strconv.Itoa(i)
			expValues[keys[i]] = &dynamodb.AttributeValue{S: aws.String(id)}
		}
		filterExp += " and #portfolioId in(" + strings.Join(keys, ",") + ")"
		expNames["#portfolioId"] = aws.String("portfolioId")
	}
	scan := &dynamodb.ScanInput{
		TableName:                 aws.String(TableName),
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
	result, err := r.DynamoDB.ScanWithContext(ctx, scan)
	if err != nil {
		return nil, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &portfolios); err != nil {
		return nil, err
	}

	return portfolios, nil
}
