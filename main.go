package main

// GOOS=linux go build -o iterator && zip -r iterator.zip iterator && aws s3 cp iterator.zip s3://rockset-high-frequency-lambda/iterator.zip
// aws lambda update-function-code --function-name hfl-iterator --s3-bucket rockset-high-frequency-lambda --s3-key iterator.zip

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	svc "github.com/aws/aws-sdk-go/service/lambda"
	"log"
	"os"
)

type request struct {
	Iterator struct {
		Index int `json:"index"`
		Count int `json:"count"`
	} `json:"iterator"`
}

type response struct {
	Index    int  `json:"index"`
	Count    int  `json:"count"`
	Continue bool `json:"continue"`
}

func handler(ctx context.Context, rqst request) (response, error) {
	index := rqst.Iterator.Index + 1
	rsp := response{
		Index:    index,
		Count:    rqst.Iterator.Count,
		Continue: index < rqst.Iterator.Count,
	}

	log.Printf("request: %+v", rqst)
	region := os.Getenv("REGION")
	if region == "" {
		return rsp, errors.New("REGION missing")
	}
	fn := os.Getenv("LAMBDA")
	if region == "" {
		return rsp, errors.New("LAMBDA missing")
	}

	sess, err := session.NewSession()
	if err != nil {
		return rsp, err
	}

	client := svc.New(sess, &aws.Config{Region: aws.String(region)})
	result, err := client.Invoke(&svc.InvokeInput{
		FunctionName:   aws.String(fn),
		InvocationType: aws.String(svc.InvocationTypeEvent),
	})
	if err != nil {
		return rsp, err
	}
	log.Printf("invocation result: %d", *result.StatusCode)

	return rsp, nil
}

func main() {
	lambda.Start(handler)
}
