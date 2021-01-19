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
		log.Printf("Error opening old leveldb: %s", err)
		panic(err)
	}
	defer oldDB.Close()

	newDB, err := gorm.Open(sqlite.Open(newDBFile), nil)
	if err != nil {
		log.Printf("Error open sqlite and gorm: %s", err)
		panic(err)
	}
	defer func() {
		db, err2 := newDB.DB()
		if err2 != nil {
			log.Printf("Error getting real DB from gorm: %s", err2)
		}
		if db != nil {
			err2 = db.Close()
			if err2 != nil {
				log.Printf("Error closing real DB of gorm: %s", err2)
				panic(err2)
			}
		}
	}()

	err = newDB.AutoMigrate(&User{})
	if err != nil {
		log.Printf("Error autoMigrating sqlite database with user: %s", err)
		panic(err)
	}

	oldUserData, err := oldDB.Get([]byte("httpauth::userdata"), nil)
	if err != nil {
		log.Printf("Error getting `httpauth::userdata` from leveldb: %s", err)
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
		log.Printf("Error unmarshalling old ")
		panic(err)
	}

	for _, datum := range migrationData {
		// check if password is "factorio", which was the default password in the old system
		decodedHash, err := base64.StdEncoding.DecodeString(datum.Hash)
		if err != nil {
			log.Printf("Error decoding base64 hash: %s", err)
			panic(err)
		}

		err = bcrypt.CompareHashAndPassword(decodedHash, []byte("factorio"))
		if err == nil {
			// password is "factorio" .. change it
			newPassword := GenerateRandomPassword()

			bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Error generating has from password: %s", err)
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
	log.Println("Deleting old leveldb database.")
	err = os.RemoveAll(oldDBFile)
	if err != nil {
		log.Printf("Error removing leveldb: %s", err)
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
