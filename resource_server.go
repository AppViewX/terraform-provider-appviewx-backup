package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"terraform-provider-appviewx/appviewx"
	"terraform-provider-appviewx/appviewx/config"
	"terraform-provider-appviewx/appviewx/constants"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"action_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"config_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"payload": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"headers": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
			},
			"master_payload": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	log.Println(" **************** GET OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	// Since the resource is for stateless operation, only nil returned
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println(" **************** UPDATE OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	//Update implementation is empty since this resource is for the stateless generic api invocation
	return errors.New("Update not supported")
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	log.Println(" **************** DELETE OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	// Delete implementation is empty since this resoruce is for the stateless generic api invocation
	return nil
}

//TODO: cleanup to be done
func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

	log.Println("*********************** Request received to create")

	configFile := d.Get(constants.CONFIG_FILE).(string)
	config.SyncConfigFile(configFile)

	appviewxSessionID, err := GetSession()
	if err != nil {
		log.Println("Error in getting the session : ", err)
		return err
	}

	types := strings.ToUpper(d.Get(constants.TYPE).(string))
	if types == constants.POST || types == constants.PUT || types == constants.DELETE || types == constants.GET {

		actionID := d.Get(constants.APPVIEWX_ACTION_ID).(string)
		payloadString := d.Get(constants.PAYLOAD).(string)
		masterPayloadFileName := d.Get(constants.MASTER_PAYLOAD).(string)

		log.Println("Input minimal payload : ", payloadString)

		payloadMinimal := make(map[string]interface{})
		json.Unmarshal([]byte(payloadString), &payloadMinimal)

		masterPayload := appviewx.GetMasterPayloadApplyingMinimalPayload(masterPayloadFileName, payloadMinimal)

		outputFilePath := config.Config[constants.OUTPUT_FILE_PATH].(string)
		appviewxEnvironmentIP := config.Config[constants.APPVIEWX_ENVIRONMENT_IP].(string)
		appviewxEnvironmentPort := config.Config[constants.APPVIEWX_ENVIRONMENT_PORT].(string)
		appviewxEnvironmentIsHTTPS := config.Config[constants.APPVIEWX_ENVIRONMENT_Is_HTTPS].(bool)
		appviewxEnvironmentGwKey := config.Config[constants.APPVIEWX_ENVIRONMENT_GW_KEY].(string)
		appviewxEnvironmentGwSource := config.Config[constants.APPVIEWX_ENVIRONMENT_GW_SOURCE].(string)

		queryParams := make(map[string]string)
		queryParams[constants.GW_KEY] = appviewxEnvironmentGwKey
		queryParams[constants.GW_SOURCE] = appviewxEnvironmentGwSource

		url := appviewx.GetURL(appviewxEnvironmentIP, appviewxEnvironmentPort, actionID, queryParams, appviewxEnvironmentIsHTTPS)

		headers := d.Get(constants.HEADERS).(map[string]interface{})

		client := &http.Client{Transport: HTTPTransport()}
		requestBody, _ := json.Marshal(masterPayload)

		printRequest(types, url, headers, requestBody)

		req, err := http.NewRequest(types, url, bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatalln(err)
		}

		for key, value := range headers {
			value1 := fmt.Sprintf("%v", value)
			key1 := fmt.Sprintf("%v", key)
			req.Header.Add(key1, value1)
		}
		req.Header.Add(constants.SESSION_ID, appviewxSessionID)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		} else {
			log.Println("Request success : url :", url)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		log.Println(string(body))

		err = ioutil.WriteFile(outputFilePath, body, 0666)

		log.Println("API ionvoke success")
		d.SetId(strconv.Itoa(rand.Int()))
		return resourceServerRead(d, m)
	}
	return nil
}

func printRequest(types, url string, headers map[string]interface{}, requestBody []byte) {
	log.Println("***************** NEW HTTP REQUEST **********************")
	log.Println("TYPE : ", types)
	log.Println("URL : ", url)
	log.Println("Headers : ", headers)
	log.Println("Body : ", string(requestBody))
	log.Println("*********************************************************")
}

//TODO: cleanup to be done
func GetSession() (output string, err error) {

	log.Println("Request received for GetSession")

	payload := make(map[string]interface{})

	headers := make(map[string]interface{})
	headers[constants.CONTENT_TYPE] = constants.APPLICATION_JSON
	headers[constants.ACCEPT] = constants.APPLICATION_JSON
	headers[constants.USERNAME] = config.Config[constants.APPVIEWX_USERNAME].(string)
	headers[constants.PASSWORD] = config.Config[constants.APPVIEWX_PASSWORD].(string)

	appviewxEnvironmentIP := config.Config[constants.APPVIEWX_ENVIRONMENT_IP].(string)
	appviewxEnvironmentPort := config.Config[constants.APPVIEWX_ENVIRONMENT_PORT].(string)
	appviewxEnvironmentIsHTTPS := config.Config[constants.APPVIEWX_ENVIRONMENT_Is_HTTPS].(bool)
	appviewxEnvironmentGwKey := config.Config[constants.APPVIEWX_ENVIRONMENT_GW_KEY].(string)
	appviewxEnvironmentGwSource := config.Config[constants.APPVIEWX_ENVIRONMENT_GW_SOURCE].(string)

	actionID := constants.APPVIEWX_ACTION_ID_LOGIN

	queryParams := make(map[string]string)
	queryParams[constants.GW_KEY] = appviewxEnvironmentGwKey
	queryParams[constants.GW_SOURCE] = appviewxEnvironmentGwSource

	url := appviewx.GetURL(appviewxEnvironmentIP, appviewxEnvironmentPort, actionID, queryParams, appviewxEnvironmentIsHTTPS)

	payloadContents, err := json.Marshal(payload)

	if err != nil {
		log.Println("Error in marshalling the ")
	}
	payloadContentsReader := bytes.NewReader(payloadContents)

	printRequest(constants.POST, url, headers, payloadContents)

	client := &http.Client{Transport: HTTPTransport()}
	req, err := http.NewRequest(constants.POST, url, payloadContentsReader)
	if err != nil {
		log.Println("Error in creating the new reqeust")
		return "", err
	}

	for key, value := range headers {
		value1 := fmt.Sprintf("%v", value)
		key1 := fmt.Sprintf("%v", key)
		req.Header.Add(key1, value1)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error in executing the request")
		return "", err
	}
	defer resp.Body.Close()
	responseContents, err := ioutil.ReadAll(resp.Body)

	err = ioutil.WriteFile("/tmp/output_session.json", responseContents, 0666)
	if err != nil {
		fmt.Println("Error in writing the session output to file")
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
	log.Println("session retrieval success ")

	return
}

func HTTPTransport() *http.Transport {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return tr
}
