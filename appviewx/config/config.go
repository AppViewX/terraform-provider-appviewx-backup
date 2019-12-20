package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var Config map[string]interface{}

func SyncConfigFile(configFile string) {
	Config = make(map[string]interface{})

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal("Error in opening the file :", configFile)
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Error in reading the contents of the file :", err)
	}
	json.Unmarshal(fileContents, &Config)
}
