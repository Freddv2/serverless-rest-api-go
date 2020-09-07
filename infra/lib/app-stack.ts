import * as cdk from '@aws-cdk/core';
import {Code, Function, Runtime} from '@aws-cdk/aws-lambda';
import {LambdaRestApi} from "@aws-cdk/aws-apigateway";
import {AttributeType, Table} from "@aws-cdk/aws-dynamodb";

export class AppStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const lambda = new Function(this, 'Portfolio-Lambda', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset("../functions/portfolios/build/go_build_main_go_linux.zip"),
      handler: 'main',
      memorySize: 128,
    });
    new LambdaRestApi(this, id + 'Portfolio-API', {
      handler: lambda,
    });
    const table = new Table(this, 'Portfolio-Table', {
      tableName: 'Portfolio',
      partitionKey: {name: 'tenantId', type: AttributeType.STRING},
      sortKey: {name: 'portfolioId', type: AttributeType.STRING},
      readCapacity: 1,
      writeCapacity: 1
    });

    table.grantReadWriteData(lambda)
  }
}
