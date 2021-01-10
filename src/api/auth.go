package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/apexskier/httpauth"
	"github.com/gorilla/sessions"
	"github.com/mroote/factorio-server-manager/bootstrap"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
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

type Auth struct {
	db *leveldb.DB
}

var (
	sessionStore *sessions.CookieStore
	auth         Auth
)

func SetupAuth() {
	var err error

	config := bootstrap.GetConfig()

	cookieEncryptionKey, err := base64.StdEncoding.DecodeString(config.CookieEncryptionKey)
	if err != nil {
		panic(err)
	}
	sessionStore = sessions.NewCookieStore(cookieEncryptionKey)
	sessionStore.Options = &sessions.Options{
		Path:   "",
		Secure: true,
	}

	auth.db, err = leveldb.OpenFile(config.DatabaseFile, nil)
	if err != nil {
		panic(err)
	}

	// check if db is empty, if so, add default user
	iterator := auth.db.NewIterator(nil, nil)
	if !iterator.Next() {
		var password = generateRandomPassword()

		var user User
		user.Username = "admin"
		user.Password = password
		user.Role = "admin"

		err = auth.addUser(user)
		if err != nil {
			panic(err)
		}

		log.Println("Created default admin user. Please change it's password as soon as possible.")
		log.Printf("Username: %s", user.Username)
		log.Printf("Password: %s", password)
	} else {
		// if first key is userdata, migrate it from old design
		if string(iterator.Key()) == "httpauth::userdata" {
			value := iterator.Value()

			var migrationData map[string]struct {
				Username string
				Email    string
				Hash     string
				Role     string
			}
			err = json.Unmarshal(value, &migrationData)
			if err != nil {
				panic(err)
			}

			for _, user := range migrationData {
				newUser := User{
					Username: user.Username,
					Password: user.Hash,
					Role:     user.Role,
					Email:    user.Email,
				}
				err = auth.addUserWithHash(newUser)
				if err != nil {
					panic(err)
				}
			}

			// remove userdata from db
			err = auth.db.Delete(iterator.Key(), nil)
			if err != nil {
				panic(err)
			}
		}
	}
}

var randLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomPassword() string {
	pass := make([]rune, 24)
	for i := range pass {
		pass[i] = randLetters[rand.Intn(len(randLetters))]
	}
	return string(pass)
}

func (a *Auth) checkPassword(username, password string) error {
	dbUser, err := a.db.Get([]byte(username), nil)
	if err != nil {
		// TODO
		return err
	}

	var user User
	err = json.Unmarshal(dbUser, &user)
	if err != nil {
		// TODO
		return err
	}

	decodedHashPw, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		// TODO
		return err
	}

	err = bcrypt.CompareHashAndPassword(decodedHashPw, []byte(password))
	if err != nil {
		// TODO
		return err
	}

	// Password correct
	return nil
}

func (a *Auth) deleteUser(username string) error {
	err := a.db.Delete([]byte(username), nil)
	if err != nil {
		// TODO
		return err
	}
	return nil
}

func (a *Auth) hasUser(username string) (bool, error) {
	return a.db.Has([]byte(username), nil)
}

func (a *Auth) getUser(username string) (User, error) {
	userJson, err := a.db.Get([]byte(username), nil)
	if err != nil {
		// TODO
		return User{}, err
	}

	var user User
	err = json.Unmarshal(userJson, &user)
	if err != nil {
		// TODO
		return User{}, err
	}

	return user, nil
}

func (a *Auth) listUsers() ([]User, error) {
	var users []User
	iterator := a.db.NewIterator(nil, nil)
	for iterator.Next() {
		userJson := iterator.Value()

		var user User
		err := json.Unmarshal(userJson, &user)
		if err != nil {
			// TODO
			return nil, err
		}
		user.Password = ""

		users = append(users, user)
	}
	return users, nil
}

func (a *Auth) addUser(user User) error {
	// encrypt password
	pwHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		// TODO
		return err
	}

	user.Password = base64.StdEncoding.EncodeToString(pwHash)

	// save user as json
	userJson, err := json.Marshal(user)
	if err != nil {
		// TODO
		return err
	}

	// add user to db
	err = a.db.Put([]byte(user.Username), userJson, nil)
	if err != nil {
		// TODO
		return err
	}

	return nil
}

func (a *Auth) addUserWithHash(user User) error {
	// save user as json
	userJson, err := json.Marshal(user)
	if err != nil {
		// TODO
		return err
	}

	// add user to db
	err = a.db.Put([]byte(user.Username), userJson, nil)
	if err != nil {
		// TODO
		return err
	}

	return nil
}

func (a *Auth) removeUser(username string) error {
	return a.db.Delete([]byte(username), nil)
}

// middleware function, that will be called for every request, that has to be authorized
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessionStore.Get(r, "authentication")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		username, ok := session.Values["username"]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		hasUser, err := auth.hasUser(username.(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if hasUser {
			next.ServeHTTP(w, r)
		} else {
			log.Printf("Unauthenticated request %s %s %s", r.Method, r.Host, r.RequestURI)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	})
}
