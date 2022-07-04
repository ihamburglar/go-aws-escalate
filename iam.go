package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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
func GetUsers(iamclient *iam.Client, stsclient *sts.Client, ctx context.Context, aUsers bool) []string {

	//TODO, this is not desirable
	input := &iam.ListUsersInput{
		MaxItems: aws.Int32(int32((1000))),
	}

	var users []string
	if aUsers {
		// If All Users
		//TODO, have not tested with large environments.  May need to paginate.
		result, err := ListUsers(context.TODO(), iamclient, input)
		if err != nil {
			fmt.Println("Got an error retrieving users:")
			panic(err)
		}
		for _, i := range result.Users {
			users = append(users, *i.Arn)
		}
	} else {
		// else single user
		result, err := stsclient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			fmt.Println("Got an error retrieving users:")
			panic(err)
		}
		users = append(users, *result.Arn)
	}

	fmt.Printf("%v", users)
	return users
}

func listGroups(ctx context.Context, iamclient *iam.Client, a string) []types.Group {
	iarn, err := arn.Parse(a)
	if err != nil {
		fmt.Println("Got an error retrieving groups:")
		panic(err)
	}
	resource := iarn.Resource

	u := strings.Split(resource, "/")

	group, err := iamclient.ListGroupsForUser(ctx, &iam.ListGroupsForUserInput{
		UserName: aws.String(u[1]),
	})
	if err != nil {
		fmt.Println("Got an error retrieving groups:")
		panic(err)
	}

	//	fmt.Printf("%v", *i.GroupName)
	return group.Groups

}
