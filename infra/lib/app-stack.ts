import * as cdk from '@aws-cdk/core';
import {Code, Function, Runtime} from '@aws-cdk/aws-lambda';
import {LambdaRestApi} from "@aws-cdk/aws-apigateway";
import {AttributeType, Table} from "@aws-cdk/aws-dynamodb";

export class AppStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const lambda = new Function(this, 'Bucket-Lambda', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset("../functions/buckets/build/go_build_main_go_linux.zip"),
      handler: 'main',
      memorySize: 128,
    });
    new LambdaRestApi(this, id + 'Bucket-API', {
      handler: lambda,
    });
    const table = new Table(this, 'Bucket-Table', {
      tableName: 'BUCKET',
      partitionKey: {name: 'tenantId', type: AttributeType.STRING},
      sortKey: {name: 'bucketId', type: AttributeType.STRING},
      readCapacity: 1,
      writeCapacity: 1
    });

    table.grantReadWriteData(lambda)
  }
}
