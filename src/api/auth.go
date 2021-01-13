package api

import (
	"encoding/base64"
	"github.com/gorilla/sessions"
	"github.com/mroote/factorio-server-manager/bootstrap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type User bootstrap.User

type Auth struct {
	db *gorm.DB
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
		Path:   "/",
		Secure: true,
	}

	auth.db, err = gorm.Open(sqlite.Open(config.SQLiteDatabaseFile), nil)
	if err != nil {
		panic(err)
	}

	auth.db.AutoMigrate(&User{})

	var userCount int64
	auth.db.Model(&User{}).Count(&userCount)

	if userCount == 0 {
		// no user created yet, create a default one
		var password = bootstrap.GenerateRandomPassword()

		var user User
		user.Username = "admin"
		user.Password = password
		user.Role = "admin"

		err := auth.addUser(user)
		if err != nil {
			panic(err)
		}

		log.Println("Created default admin user. Please change it's password as soon as possible.")
		log.Printf("Username: %s", user.Username)
		log.Printf("Password: %s", password)
	}
}

func (a *Auth) checkPassword(username, password string) error {
	var user User
	result := a.db.Where(&User{Username: username}).Take(&user)
	if result.Error != nil {
		// TODO
		return result.Error
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
	result := a.db.Model(&User{}).Where(&User{Username: username}).Delete(&User{})
	if result.Error != nil {
		// TODO
		return result.Error
	}
	return nil
}

func (a *Auth) hasUser(username string) (bool, error) {
	var count int64
	result := a.db.Model(&User{}).Where(&User{Username: username}).Count(&count)
	if result.Error != nil {
		// TODO
		return false, result.Error
	}
	return count == 1, nil
}

func (a *Auth) getUser(username string) (User, error) {
	var user User
	result := a.db.Model(&User{}).Where(&User{Username: username}).Take(&user)
	if result.Error != nil {
		// TODO
		return User{}, result.Error
	}

	return user, nil
}

func (a *Auth) listUsers() ([]User, error) {
	var users []User
	result := a.db.Find(&users)
	if result.Error != nil {
		// TODO
		return nil, result.Error
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

	// add user to db
	result := a.db.Create(&user)
	if result.Error != nil {
		// TODO
		return result.Error
	}

	return nil
}

func (a *Auth) addUserWithHash(user User) error {
	// add user to db
	result := a.db.Create(&user)
	if result.Error != nil {
		// TODO
		return result.Error
	}

	return nil
}

func (a *Auth) removeUser(username string) error {
	result := a.db.Model(&User{}).Where(&User{Username: username}).Delete(&User{})
	if result.Error != nil {
		// TODO
		return result.Error
	}
	return nil
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
