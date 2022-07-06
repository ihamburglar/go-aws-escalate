package main

import (
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	awspolicy "github.com/n4ch04/aws-policy"
)

type UserMetaData struct {
	User                     types.User
	Groups                   []types.Group
	InlinePolicies           [][]string
	AttachedPolicies         [][]types.AttachedPolicy
	UserInlinePolicies       []string
	userInlinePolicyDocument []awspolicy.Policy
}

type AllUsers struct {
	UserMetaData []UserMetaData
}
