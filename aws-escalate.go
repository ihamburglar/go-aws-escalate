package main

import (
	"fmt"
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
		groups := listGroups(ctx, iamc, *user.User.UserId)
		// TODO move iarns into GetUsers
		fmt.Println("=== The user: " + *user.User.UserId + " has the following groups:")
		for _, group := range groups {
			fmt.Println(*group.GroupName)

		}

	}

}
