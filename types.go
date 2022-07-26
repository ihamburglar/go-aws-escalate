package main

import (
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	awspolicy "github.com/n4ch04/aws-policy"
)

type UserMetaData struct {
	User                       types.User
	Groups                     []types.Group
	GroupInlinePolicies        []GroupPolicies
	GroupAttachedPolicies      [][]types.AttachedPolicy
	GroupInlinePolicyDocuments []awspolicy.Policy
	UserInlinePolicies         []string
	UserAttachedPolicies       [][]types.AttachedPolicy
	UserInlinePolicyDocument   []awspolicy.Policy
}

// TODO change to lower?
type AllUsers struct {
	UserMetaData []UserMetaData
}

//TODO change to lower?
type GroupPolicies struct {
	Group    types.Group
	Policies []string
}

type method struct {
	MethodName string
	PolicySets []policySet
}

type policySet struct {
	PolicyName string
	PolicyBool bool
}
