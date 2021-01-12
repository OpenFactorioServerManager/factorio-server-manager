package bootstrap

import (
	"encoding/base64"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"os"
)

type User struct {
	gorm.Model
	Username string `json:"username",gorm:"uniqueIndex,not null"`
	Password string `json:"password",gorm:"not null"`
	Role     string `json:"role",gorm:"not null"`
	Email    string `json:"email"`
}

func MigrateLevelDBToSqlite(oldDBFile, newDBFile string) {
	oldDB, err := leveldb.OpenFile(oldDBFile, nil)
	if err != nil {
		panic(err)
	}

	newDB, err := gorm.Open(sqlite.Open(newDBFile), nil)
	if err != nil {
		panic(err)
	}

	newDB.AutoMigrate(&User{})

	oldUserData, err := oldDB.Get([]byte("httpauth::userdata"), nil)
	if err != nil {
		panic(err)
	}

	var migrationData map[string]struct {
		Username string
		Email    string
		Hash     string
		Role     string
	}
	err = json.Unmarshal(oldUserData, &migrationData)
	if err != nil {
		panic(err)
	}

	for _, datum := range migrationData {
		// check if password is "factorio", which was the default password in the old system
		decodedHash, err := base64.StdEncoding.DecodeString(datum.Hash)
		if err != nil {
			panic(err)
		}

		err = bcrypt.CompareHashAndPassword(decodedHash, []byte("factorio"))
		if err == nil {
			// password is "factorio" .. change it
			newPassword := GenerateRandomPassword()

			bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}

			datum.Hash = base64.StdEncoding.EncodeToString(bcryptPassword)

			log.Println(`Migrated user in database. It still had default password "factorio" set. New credentials:`)
			log.Printf("Username: %s", datum.Username)
			log.Printf("Password: %s", newPassword)
		}

		user := &User{
			Username: datum.Username,
			Password: datum.Hash,
			Role:     datum.Role,
			Email:    datum.Email,
		}

		newDB.Create(user)
	}

	oldDB.Close()

	// delete oldDB
	err = os.RemoveAll(oldDBFile)
	if err != nil {
		panic(err)
	}
}

var randLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateRandomPassword() string {
	pass := make([]rune, 24)
	for i := range pass {
		pass[i] = randLetters[rand.Intn(len(randLetters))]
	}
	return string(pass)
}
