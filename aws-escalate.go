package main

import (
	"fmt"
	"strings"
)

func main() {

	//TODO get from prompt of whatever
	//access_key_id := "asdf"
	//secret_acces_key := "asdf"
	//session_token := "asdf"
	all_users := true

	// More flexibility with auth

	//TODO what happens when there is no session
	//c := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(access_key_id, secret_acces_key, session_token))

	iamc, stsc, ctx := authClients(all_users)
	var users AllUsers
	users = GetUsers(users, iamc, stsc, ctx, all_users)

	for _, user := range users.UserMetaData {
		groups := GetGroups(ctx, iamc, *user.User.Arn)
		groupInlinePolicies := GetGroupPolicies(ctx, iamc, groups)
		groupInlinePolicyDocuments := GetGroupPolicyDocument(ctx, iamc, groupInlinePolicies)
		groupAttachedPolicies := GetAttachedGroupPolicies(ctx, iamc, groups)
		userInlinePolicies := GetUserInlinePolicies(ctx, iamc, user.User)
		userAttachedPolicies := GetAttachedUserPolicies(ctx, iamc, user.User)
		userInlinePolicyDocument := GetUserPolicyDocument(ctx, iamc, user.User, userInlinePolicies)

		users.UserMetaData = append(users.UserMetaData, UserMetaData{
			User:                       user.User,
			Groups:                     groups,
			GroupInlinePolicies:        groupInlinePolicies,
			GroupAttachedPolicies:      groupAttachedPolicies,
			GroupInlinePolicyDocuments: groupInlinePolicyDocuments,
			UserInlinePolicies:         userInlinePolicies,
			UserAttachedPolicies:       userAttachedPolicies,
			UserInlinePolicyDocument:   userInlinePolicyDocument,
		})

		fmt.Println("=== The user: " + strings.TrimPrefix(*user.User.Arn, "arn:aws:iam::") + " has the following groups: \n")
		for _, group := range groups {
			fmt.Println("* Found group name :" + *group.GroupName)

		}

	}

}
