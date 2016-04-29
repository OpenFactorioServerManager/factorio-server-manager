package main

import (
	"log"

	"github.com/apexskier/httpauth"
)

type AuthHTTP struct {
	backend httpauth.LeveldbAuthBackend
	aaa     httpauth.Authorizer
}

func initAuth() *AuthHTTP {
	return &AuthHTTP{}
}

func (auth *AuthHTTP) createRoles() {
	var err error
	roles := make(map[string]httpauth.Role)

	roles["user"] = 30
	roles["admin"] = 80
	auth.aaa, err = httpauth.NewAuthorizer(auth.backend, []byte("topsecretkey"), "user", roles)
	if err != nil {
		log.Printf("Error creating roles: %s", err)
	}
}

func (auth *AuthHTTP) createUser(username, role, password, email string) error {
	user := httpauth.UserData{Username: username, Role: role, Email: email}
	err := auth.backend.SaveUser(user)
	if err != nil {
		log.Printf("Error saving user: %s", err)
		return err
	}

	err = auth.aaa.Update(nil, nil, username, password, "")
	if err != nil {
		log.Printf("Error saving user: %s", err)
		return err
	}

	return nil
}
