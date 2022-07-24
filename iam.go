package main

import (
	"context"
	"fmt"
	"net/url"
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
			fmt.Println("* Found user " + *a.UserName)

		}
	} else {
		// else single user
		result, err := stsclient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			fmt.Println("Got an error retrieving users:")
			panic(err)
		}

		iarn, err := arn.Parse(*result.Arn)
		if err != nil {
			fmt.Println("Got an error parsing arn")
			panic(err)
		}

		//TODO not sure this does the right thing
		uname := iarn.String()
		fmt.Println("* Found user " + uname)

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

	return users
}

// GetGroups returns the groups for the specified user in the AWS type
func GetGroups(ctx context.Context, iamclient *iam.Client, a string) []types.Group {
	iarn, err := arn.Parse(a)
	if err != nil {
		fmt.Println("Got an error parsing arn")
		panic(err)
	}
	resource := iarn.Resource

	u := strings.Split(resource, "/")

	// TODO IsTruncated should be dealt with for large groups
	group, err := iamclient.ListGroupsForUser(ctx, &iam.ListGroupsForUserInput{
		UserName: aws.String(u[1]),
	})
	for _, g := range group.Groups {
		fmt.Println("* Found group: " + *g.GroupName + " for this user")
	}
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
		fmt.Printf("policiespergroup: %v\n", policiespergroup.PolicyNames)
		for _, p := range policiespergroup.PolicyNames {
			fmt.Printf("*** From GetGroupPolicies() group policy names :%v\n", &p)
		}
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
		for _, p := range policiespergroup.AttachedPolicies {
			fmt.Printf("*** attached group policies: %v\n", *p.PolicyName)
		}
		policies = append(policies, policiespergroup.AttachedPolicies)
	}
	return policies
}

// GetUserPolicies runs ListUserPolicies and returns the associated inline user policies
func GetUserInlinePolicies(ctx context.Context, iamclient *iam.Client, user types.User) []string {

	input := &iam.ListUserPoliciesInput{UserName: user.UserName}
	policies, err := iamclient.ListUserPolicies(ctx, input)
	if err != nil {
		fmt.Println("Got an error retrieving user inline policies:")
		panic(err)
	}
	for _, p := range policies.PolicyNames {
		fmt.Printf("*** From GetUserPolicies() Policy names :%v\n", &p)
	}

	return policies.PolicyNames
}

// GetAttachedUserPolicies runs ListAttachedPolicies and returns the attached policies
func GetAttachedUserPolicies(ctx context.Context, iamclient *iam.Client, user types.User) [][]types.AttachedPolicy {

	// var for policy slice
	var policies [][]types.AttachedPolicy
	input := &iam.ListAttachedUserPoliciesInput{UserName: user.UserName}
	ap, err := iamclient.ListAttachedUserPolicies(ctx, input)
	if err != nil {
		fmt.Println("Got an error retrieving attached group policies:")
		panic(err)
	}

	for _, p := range ap.AttachedPolicies {
		fmt.Printf("*** attached policies: %v\n", *p.PolicyName)
	}
	policies = append(policies, ap.AttachedPolicies)

	return policies
}

// GetUserPolicyDocument runs GetUserPolicy and returns the associated inline user policy document
func GetUserPolicyDocument(ctx context.Context, iamclient *iam.Client, user types.User, policies []string) []awspolicy.Policy {
	var outPolicies []awspolicy.Policy
	fmt.Println("fetching user policy documents")
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
		fmt.Printf("*** From GetUserPolicyDocument() policy : %v\n", &rawPolicyDocument.PolicyDocument)
		p := awspolicy.Policy{
			//[]byte(rawPolicyDocument.PolicyDocument)
			//rawPolicyDocument.PolicyDocument.Sta
		}
		outPolicies = append(outPolicies, p)
		//outPolicies := outPolicies.UnmarshalJSON()
		//outPolicies = append(outPolicies, *outPolicy.PolicyDocument)

	}
	return outPolicies
}

// GetGroupPolicyDocument runs GetUserPolicy and returns the associated inline user policy document
func GetGroupPolicyDocument(ctx context.Context, iamclient *iam.Client, policies []GroupPolicies) []awspolicy.Policy {
	var outPolicies []awspolicy.Policy
	//for _, group := range groups {
	for _, groupPolicies := range policies {
		fmt.Printf("Group policies: %v\n", groupPolicies.Policies)
		for _, groupPolicy := range groupPolicies.Policies {
			fmt.Println("*** fetching group policies for group: " + *groupPolicies.Group.GroupName)

			fmt.Println("fetching group policy: " + groupPolicy + " for group: " + *groupPolicies.Group.GroupName)
			input := &iam.GetGroupPolicyInput{
				PolicyName: &groupPolicy,
				GroupName:  groupPolicies.Group.GroupName}
			rawPolicyDocument, err := iamclient.GetGroupPolicy(ctx, input)
			if err != nil {
				fmt.Println("Got an error retrieving user inline policy documents:")
				panic(err)
			}
			pd, err := url.QueryUnescape(*rawPolicyDocument.PolicyDocument)
			if err != nil {
				fmt.Println("Got an error decoding url encoded string	")
				panic(err)
			}
			fmt.Printf("*** From GetGroupPolicyDocument() policy : %v\n", pd)
			p := awspolicy.Policy{
				//[]byte(rawPolicyDocument.PolicyDocument)
				//rawPolicyDocument.PolicyDocument.Sta
			}
			outPolicies = append(outPolicies, p)
			//outPolicies := outPolicies.UnmarshalJSON()
			//outPolicies = append(outPolicies, *outPolicy.PolicyDocument)

		}
	}
	//}
	return outPolicies
}
