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
	// TODO put users into group
	for _, user := range users.UserMetaData {
		g := listGroups(ctx, iamc, *user.User.Arn)
		users.UserMetaData = append(users.UserMetaData, UserMetaData{
			User:   user.User,
			Groups: g})
		// TODO move iarns into GetUsers
		fmt.Println("=== The user: " + strings.TrimPrefix(*user.User.Arn, "arn:aws:iam::") + " has the following groups:")
		for _, group := range g {
			fmt.Println(*group.GroupName)

		}

	}

}
