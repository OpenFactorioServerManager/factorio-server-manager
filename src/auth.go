package main

import (
	"log"
	"os"

	"github.com/apexskier/httpauth"
)

type AuthHTTP struct {
	backend httpauth.LeveldbAuthBackend
	aaa     httpauth.Authorizer
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

func initAuth() *AuthHTTP {
	return &AuthHTTP{}
}

func (auth *AuthHTTP) CreateAuth(backendFile string, cookieKey string) error {
	var err error
	os.Mkdir(backendFile, 0755)

	auth.backend, err = httpauth.NewLeveldbAuthBackend(backendFile)
	if err != nil {
		log.Printf("Error creating Auth backend: %s", err)
		return err
	}

	roles := make(map[string]httpauth.Role)
	roles["user"] = 30
	roles["admin"] = 80

	auth.aaa, err = httpauth.NewAuthorizer(auth.backend, []byte(cookieKey), "user", roles)
	if err != nil {
		log.Printf("Error creating authorizer: %s", err)
		return err
	}

	return nil
}

func (auth *AuthHTTP) CreateOrUpdateUser(username, password, role, email string) error {
	user := httpauth.UserData{Username: username, Role: role, Email: email}
	err := auth.backend.SaveUser(user)
	if err != nil {
		log.Printf("Error saving user: %s", err)
		return err
	}

	err = auth.aaa.Update(nil, nil, username, password, email)
	if err != nil {
		log.Printf("Error updating user: %s", err)
		return err
	}

	log.Printf("Created user: %s", user.Username)

	return nil
}

func (auth *AuthHTTP) listUsers() ([]User, error) {
	var userResponse []User
	users, err := auth.backend.Users()
	if err != nil {
		log.Printf("Error list users: %s", err)
		return nil, err
	}

	for _, user := range users {
		u := User{Username: user.Username, Role: user.Role, Email: user.Email}
		userResponse = append(userResponse, u)
	}

	log.Printf("listing users: %v found", len(users))
	return userResponse, nil
}

func (auth *AuthHTTP) addUser(username, password, email, role string) error {
	user := httpauth.UserData{Username: username, Hash: []byte(password), Email: email, Role: role}
	err := auth.backend.SaveUser(user)
	if err != nil {
		log.Printf("Error creating user %v: %s", user, err)
	}
	err = auth.aaa.Update(nil, nil, username, password, email)
	if err != nil {
		log.Printf("Error saving user: %s", err)
		return err
	}

	log.Printf("Added user: %v", user)
	return nil
}

func (auth *AuthHTTP) removeUser(username string) error {
	err := auth.backend.DeleteUser(username)
	if err != nil {
		log.Printf("Could not delete user %s, error: %s", username, err)
		return err
	}

	return nil
}
