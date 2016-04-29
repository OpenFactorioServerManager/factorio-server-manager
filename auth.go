package main

import (
	"log"

	"github.com/apexskier/httpauth"
)

type Auth struct {
	backend httpauth.LeveldbAuthBackend
	aaa     httpauth.Authorizer
}

func initAuth() *Auth {
	return &Auth{}
}

func (auth *Auth) createRoles() {
	roles := make(map[string]httpauth.Role)

	roles["user"] = 30
	roles["admin"] = 80
	auth.aaa = httpauth.NewAuthorizer(auth.backend, "topsecretkey", "user", roles)
}

func (auth *Auth) createUser(username, role, password, email string) error {
	user := httpauth.UserData{Username: username, Role: role}
	err = backend.SaveUser(user)
	if err != nil {
		log.Printf("Error saving user: %s", err)
		return err
	}

	err := auth.aaa.Update(nil, nil, username, password, email)
	if err != nil {
		log.Printf("Error saving user: %s", err)
		return err
	}
}
