package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func authClients(kid string, key string, session string, aUsers bool) (*iam.Client, *sts.Client, context.Context) {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(kid, key, session)))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// move these into an auth.go
	iamclient := iam.NewFromConfig(cfg)
	stsclient := sts.NewFromConfig(cfg)

	return iamclient, stsclient, ctx
}
