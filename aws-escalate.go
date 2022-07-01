package main

func main() {

	//TODO get from prompt of whatever
	access_key_id := "asdf"
	secret_acces_key := "asdf"
	session_token := "asdf"
	all_users := true

	// More flexibility with auth
	iamc, stsc, ctx := authClients(access_key_id, secret_acces_key, session_token, all_users)
	GetUsers(iamc, stsc, ctx, all_users)
}
