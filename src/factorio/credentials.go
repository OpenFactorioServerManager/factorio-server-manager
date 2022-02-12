package factorio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/OpenFactorioServerManager/factorio-server-manager/bootstrap"
)

type Credentials struct {
	Username string `json:"username"`
	Userkey  string `json:"userkey"`
}

func (credentials *Credentials) Save() error {
	var err error
	config := bootstrap.GetConfig()
	credentialsJson, err := json.Marshal(credentials)
	if err != nil {
		log.Printf("error mashalling the credentials: %s", err)
		return err
	}

	err = ioutil.WriteFile(config.FactorioCredentialsFile, credentialsJson, 0664)
	if err != nil {
		log.Printf("error on saving the credentials. %s", err)
		return err
	}

	return nil
}

func (credentials *Credentials) Load() (bool, error) {
	var err error
	config := bootstrap.GetConfig()
	if _, err := os.Stat(config.FactorioCredentialsFile); os.IsNotExist(err) {
		return false, nil
	}

	fileBytes, err := ioutil.ReadFile(config.FactorioCredentialsFile)
	if err != nil {
		credentials.Del()
		log.Printf("error reading CredentialsFile: %s", err)
		return false, err
	}

	err = json.Unmarshal(fileBytes, credentials)
	if err != nil {
		credentials.Del()
		log.Printf("error on unmarshal credentials_file: %s", err)
		return false, err
	}

	if credentials.Userkey != "" && credentials.Username != "" {
		return true, nil
	} else {
		credentials.Del()
		return false, errors.New("incredients incomplete")
	}
}

func (credentials *Credentials) Del() error {
	var err error
	config := bootstrap.GetConfig()
	err = os.Remove(config.FactorioCredentialsFile)
	if err != nil {
		log.Printf("error delete the credentialfile: %s", err)
		return err
	}

	return nil
}
