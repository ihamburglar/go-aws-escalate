package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func authClients(aUsers bool) (*iam.Client, *sts.Client, context.Context) {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	iamclient := iam.NewFromConfig(cfg)
	stsclient := sts.NewFromConfig(cfg)
	return iamclient, stsclient, ctx
}
