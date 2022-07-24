package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	awspolicy "github.com/n4ch04/aws-policy"
)

// IAMListUsersAPI defines the interface for the ListUsers function.
// We use this interface to test the function using a mocked service.
type IAMListUsersAPI interface {
	ListUsers(ctx context.Context,
		params *iam.ListUsersInput,
		optFns ...func(*iam.Options)) (*iam.ListUsersOutput, error)
}

// ListUsers retrieves a list of your AWS Identity and Access Management (IAM) users.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a ListUsersOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ListUsers.
func ListUsers(c context.Context, api IAMListUsersAPI, input *iam.ListUsersInput) (*iam.ListUsersOutput, error) {
	return api.ListUsers(c, input)
}

// GetUsers retrieves either all users or just the current user and returns a slice of strings of users.
func GetUsers(users AllUsers, iamclient *iam.Client, stsclient *sts.Client, ctx context.Context, aUsers bool) AllUsers {

	//TODO, this is not desirable
	input := &iam.ListUsersInput{
		MaxItems: aws.Int32(int32((1000))),
	}

	if aUsers {
		// If All Users
		//TODO, have not tested with large environments.  May need to paginate.
		result, err := ListUsers(context.TODO(), iamclient, input)
		if err != nil {
			fmt.Println("Got an error retrieving users:")
			panic(err)
		}
		for _, a := range result.Users {
			u := UserMetaData{
				User: a,
			}
			users.UserMetaData = append(users.UserMetaData, u)
		}
	} else {
		// else single user
		result, err := stsclient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			fmt.Println("Got an error retrieving users:")
			panic(err)
		}
		// Sometimes you don't have perms to list users.
		// This is a best effort for single user creation from STS
		faketime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
		a := types.User{
			Arn: result.Arn,
			// TODO pull from somewhere else if at all
			CreateDate: &faketime,
			Path:       result.Arn,
			UserId:     result.UserId,
			UserName:   result.Arn,
			// TODO pull from somewhere else if at all
			// PasswordLastUsed: (*time.Time)(time.Now().UTC().Location()),
			// PermissionsBoundary: ,
			// Tags: ,
			// noSmithyCocumentSerde
		}
		u := UserMetaData{
			User: a,
		}
		users.UserMetaData = append(users.UserMetaData, u)
	}

	fmt.Printf("%v", users)
	return users
}

// GetGroups returns the groups for the specified user in the AWS type
func GetGroups(ctx context.Context, iamclient *iam.Client, a string) []types.Group {
	iarn, err := arn.Parse(a)
	if err != nil {
		fmt.Println("Got an error retrieving groups:")
		panic(err)
	}
	resource := iarn.Resource

	u := strings.Split(resource, "/")

	// TODO IsTruncated should be dealt with for large groups
	group, err := iamclient.ListGroupsForUser(ctx, &iam.ListGroupsForUserInput{
		UserName: aws.String(u[1]),
	})
	fmt.Println(group.Groups)
	if err != nil {
		fmt.Println("Got an error retrieving groups:")
		panic(err)
	}

	//	fmt.Printf("%v", *i.GroupName)
	return group.Groups

}

// GetGroupPolicies runs ListGroupPolicies and returns the groups associated inline policies
func GetGroupPolicies(ctx context.Context, iamclient *iam.Client, groups []types.Group) []GroupPolicies {

	// var for policy slice
	var gpolicies []GroupPolicies
	for _, group := range groups {
		input := &iam.ListGroupPoliciesInput{GroupName: group.GroupName}
		policiespergroup, err := iamclient.ListGroupPolicies(ctx, input)
		if err != nil {
			fmt.Println("Got an error retrieving inline policies:")
			panic(err)
		}
		fmt.Printf("*** group policy names :%v\n", policiespergroup.PolicyNames)
		gp := GroupPolicies{
			Group:    group,
			Policies: policiespergroup.PolicyNames,
		}
		gpolicies = append(gpolicies, gp)
	}
	return gpolicies
}

// GetAttachedGroupPolicies runs ListAttachedPolicies and returns the attached policies
func GetAttachedGroupPolicies(ctx context.Context, iamclient *iam.Client, groups []types.Group) [][]types.AttachedPolicy {

	// var for policy slice
	var policies [][]types.AttachedPolicy
	for _, group := range groups {
		input := &iam.ListAttachedGroupPoliciesInput{GroupName: group.GroupName}
		policiespergroup, err := iamclient.ListAttachedGroupPolicies(ctx, input)
		if err != nil {
			fmt.Println("Got an error retrieving attached group policies:")
			panic(err)
		}
		fmt.Printf("*** attached group policies :%v\n", policiespergroup.AttachedPolicies)
		policies = append(policies, policiespergroup.AttachedPolicies)
	}
	return policies
}

// GetUserPolicies runs ListUserPolicies and returns the associated inline user policies
func GetUserPolicies(ctx context.Context, iamclient *iam.Client, user types.User) []string {

	input := &iam.ListUserPoliciesInput{UserName: user.UserName}
	policies, err := iamclient.ListUserPolicies(ctx, input)
	if err != nil {
		fmt.Println("Got an error retrieving user inline policies:")
		panic(err)
	}
	fmt.Printf("*** Policy names :%v\n", policies.PolicyNames)

	return policies.PolicyNames
}

// GetAttachedUserPolicies runs ListAttachedPolicies and returns the attached policies
func GetAttachedUserPolicies(ctx context.Context, iamclient *iam.Client, user types.User) [][]types.AttachedPolicy {

	// var for policy slice
	var policies [][]types.AttachedPolicy
	input := &iam.ListAttachedUserPoliciesInput{UserName: user.UserName}
	p, err := iamclient.ListAttachedUserPolicies(ctx, input)
	if err != nil {
		fmt.Println("Got an error retrieving attached group policies:")
		panic(err)
	}
	fmt.Printf("*** attached policies :%v\n", p.AttachedPolicies)
	policies = append(policies, p.AttachedPolicies)

	return policies
}

// GetUserPolicyDocument runs GetUserPolicy and returns the associated inline user policy document
func GetUserPolicyDocument(ctx context.Context, iamclient *iam.Client, user types.User, policies []string) []awspolicy.Policy {
	var outPolicies []awspolicy.Policy
	fmt.Println("fetching policies")
	for _, policy := range policies {
		fmt.Println("fetching policy: " + policy)
		input := &iam.GetUserPolicyInput{
			PolicyName: &policy,
			UserName:   user.UserName}
		rawPolicyDocument, err := iamclient.GetUserPolicy(ctx, input)
		if err != nil {
			fmt.Println("Got an error retrieving user inline policy documents:")
			panic(err)
		}
		fmt.Printf("*** policy : %v\n", &rawPolicyDocument.PolicyDocument)
		p := awspolicy.Policy{
			//[]byte(rawPolicyDocument.PolicyDocument)
			//rawPolicyDocument.PolicyDocument.Sta
		}
		outPolicies = append(outPolicies, p)
		//outPolicies := outPolicies.UnmarshalJSON()
		//outPolicies = append(outPolicies, *outPolicy.PolicyDocument)

	}
	fmt.Printf("fdsa :%v\n", policies)

	return outPolicies
}

// GetUserPolicyDocument runs GetUserPolicy and returns the associated inline user policy document
func GetGroupPolicyDocument(ctx context.Context, iamclient *iam.Client, groups []types.Group, policies []GroupPolicies) []awspolicy.Policy {
	var outPolicies []awspolicy.Policy
	fmt.Println("fetching policies")
	for _, group := range groups {
		for _, groupPolicies := range policies {
			for _, groupPolicy := range groupPolicies.Policies {
				fmt.Println("fetching policy: " + groupPolicy)
				input := &iam.GetGroupPolicyInput{
					PolicyName: &groupPolicy,
					GroupName:  group.GroupName}
				rawPolicyDocument, err := iamclient.GetGroupPolicy(ctx, input)
				if err != nil {
					fmt.Println("Got an error retrieving user inline policy documents:")
					panic(err)
				}
				fmt.Printf("*** policy : %v\n", &rawPolicyDocument.PolicyDocument)
				p := awspolicy.Policy{
					//[]byte(rawPolicyDocument.PolicyDocument)
					//rawPolicyDocument.PolicyDocument.Sta
				}
				outPolicies = append(outPolicies, p)
				//outPolicies := outPolicies.UnmarshalJSON()
				//outPolicies = append(outPolicies, *outPolicy.PolicyDocument)

			}
		}
	}
	fmt.Printf("fdsa :%v\n", policies)

	return outPolicies
}
