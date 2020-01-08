package appviewx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"terraform-provider-appviewx/appviewx/constants"
)

func printRequest(types, url string, headers map[string]interface{}, requestBody []byte) {
	log.Println("[DEBUG] ***************** NEW HTTP REQUEST **********************")
	log.Println("[DEBUG] TYPE : ", types)
	log.Println("[DEBUG] URL : ", url)
	log.Println("[DEBUG] Headers : ", headers)
	log.Println("[DEBUG] Body : ", string(requestBody))
	log.Println("[DEBUG] *********************************************************")
}

//TODO: cleanup to be done
func GetSession(appviewxUserName, appviewxPassword, appviewxEnvironmentIP, appviewxEnvironmentPort, appviewxGwSource string, appviewxEnvironmentIsHTTPS bool) (output string, err error) {

	log.Println("[INFO] Request received for GetSession")

	payload := make(map[string]interface{})

	headers := make(map[string]interface{})
	headers[constants.CONTENT_TYPE] = constants.APPLICATION_JSON
	headers[constants.ACCEPT] = constants.APPLICATION_JSON
	headers[constants.USERNAME] = appviewxUserName
	headers[constants.PASSWORD] = appviewxPassword

	actionID := constants.APPVIEWX_ACTION_ID_LOGIN

	queryParams := make(map[string]string)
	queryParams[constants.GW_SOURCE] = appviewxGwSource

	url := GetURL(appviewxEnvironmentIP, appviewxEnvironmentPort, actionID, queryParams, appviewxEnvironmentIsHTTPS)

	payloadContents, err := json.Marshal(payload)

	if err != nil {
		log.Println("[ERROR] Error in marshalling the ")
	}
	payloadContentsReader := bytes.NewReader(payloadContents)

	printRequest(constants.POST, url, headers, payloadContents)

	client := &http.Client{Transport: HTTPTransport()}
	req, err := http.NewRequest(constants.POST, url, payloadContentsReader)
	if err != nil {
		log.Println("[ERROR] Error in creating the new reqeust")
		return "", err
	}

	for key, value := range headers {
		value1 := fmt.Sprintf("%v", value)
		key1 := fmt.Sprintf("%v", key)
		req.Header.Add(key1, value1)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] Error in executing the request")
		return "", err
	}
	defer resp.Body.Close()
	responseContents, err := ioutil.ReadAll(resp.Body)

	err = ioutil.WriteFile("/tmp/output_session.json", responseContents, 0666)
	if err != nil {
		log.Println("[ERROR] Error in writing the session output to file")
		return "", err
	}

	map1 := make(map[string]interface{})
	json.Unmarshal(responseContents, &map1)

	if map1[constants.RESPONSE] != nil {
		responseMap := map1[constants.RESPONSE].(map[string]interface{})
		if responseMap != nil && responseMap[constants.SESSION_ID] != nil {
			output = responseMap[constants.SESSION_ID].(string)
		}
	}
	log.Println("[INFO] session retrieval success ")

	return
}

func HTTPTransport() *http.Transport {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return tr
}
