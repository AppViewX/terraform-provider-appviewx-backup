package fileops

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func GetFileContentsInMap(fileName string) map[string]interface{} {
	output := make(map[string]interface{})
	log.Println("fileName : ", fileName)
	if fileName == "" {
		log.Println("File name is empty : ", fileName)
		return output
	}
	masterFile, err := os.Open(fileName)
	if err != nil {
		log.Println("Error in opening the file : ", fileName)
		return output
	}
	masterFileContents, err := ioutil.ReadAll(masterFile)
	if err != nil {
		log.Println("Error in reading the file contents")
	}
	json.Unmarshal(masterFileContents, &output)
	return output
}

func WriteContentsToFile(input map[string]interface{}, outputFileName string) error {
	inputContents, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		log.Println("Error in Unmarshalling ", err)
		return err
	}

	err = ioutil.WriteFile(outputFileName, inputContents, 0777)
	if err != nil {
		log.Println("Error in Unmarshalling ", err)
		return err
	}
	return nil
}
