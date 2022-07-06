package main

import "github.com/aws/aws-sdk-go-v2/service/iam/types"

type UserMetaData struct {
	User   types.User
	Groups []types.Group
}

type AllUsers struct {
	UserMetaData []UserMetaData
}
