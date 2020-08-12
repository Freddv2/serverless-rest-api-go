# Serverless GO Rest API example 

## Compile

### Windows

Install the AWS tool that help us package the executable with the right execution permissions
```
go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip 
```

Then compile & zip
```
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o functions\buckets\build\main functions\buckets\cmd\main.go
build-lambda-zip -o functions\buckets\build\main.zip functions\buckets\build\main
```

## Deployment

This project is using AWS CDK Typescript for the deployment. First, you'll need to install 
Typescript & aws-cdk through npm 

```
npm i -g typescript
npm i -g aws-cdk
```

To install & build the rest api infra stack:

``` 
cd infra
npm i
npm run build
```

To deploy:

```
cd bootstrap
cdk deploy
```