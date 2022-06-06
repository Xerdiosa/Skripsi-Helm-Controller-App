package kinesis

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

var kinesisClient *kinesis.Client

func GetKinesisClient() (*kinesis.Client, error) {
	if kinesisClient != nil {
		return kinesisClient, nil
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	kinesisClient = kinesis.NewFromConfig(cfg)
	return kinesisClient, nil
}
