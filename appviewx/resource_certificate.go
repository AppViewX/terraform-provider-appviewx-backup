package appviewx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"terraform-provider-appviewx/appviewx/config"
	"terraform-provider-appviewx/appviewx/constants"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ResourceCertificateServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateServerCreate,
		Read:   resourceCertificateServerRead,
		Update: resourceCertificateServerUpdate,
		Delete: resourceCertificateServerDelete,

		Schema: map[string]*schema.Schema{
			constants.APPVIEWX_ACTION_ID: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			constants.PAYLOAD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.TYPE: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			constants.HEADERS: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			constants.MASTER_PAYLOAD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.QUERY_PARAMS: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}
func resourceCertificateServerRead(d *schema.ResourceData, m interface{}) error {
	log.Println(" **************** GET OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	// Since the resource is for stateless operation, only nil returned
	return nil
}

func resourceCertificateServerUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println(" **************** UPDATE OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	//Update implementation is empty since this resource is for the stateless generic api invocation
	return errors.New("Update not supported")
}

func resourceCertificateServerDelete(d *schema.ResourceData, m interface{}) error {
	log.Println(" **************** DELETE OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	// Delete implementation is empty since this resoruce is for the stateless generic api invocation
	return nil
}

//TODO: cleanup to be done
func resourceCertificateServerCreate(d *schema.ResourceData, m interface{}) error {

	configAppViewXEnvironment := m.(*config.AppViewXEnvironment)

	//
	configAppViewXEnvironmentContent, _ := json.Marshal(configAppViewXEnvironment)
	log.Println("configAppViewXEnvironmentContent : ", string(configAppViewXEnvironmentContent))
	//

	log.Println("*********************** Request received to create")
	appviewxUserName := configAppViewXEnvironment.AppViewXUserName
	appviewxPassword := configAppViewXEnvironment.AppViewXPassword
	appviewxEnvironmentIP := configAppViewXEnvironment.AppViewXEnvironmentIP
	appviewxEnvironmentPort := configAppViewXEnvironment.AppViewXEnvironmentPort
	appviewxEnvironmentIsHTTPS := configAppViewXEnvironment.AppViewXIsHTTPS
	appviewxGwSource := "WEB"

	appviewxSessionID, err := GetSession(appviewxUserName, appviewxPassword, appviewxEnvironmentIP, appviewxEnvironmentPort, appviewxGwSource, appviewxEnvironmentIsHTTPS)
	if err != nil {
		log.Println("Error in getting the session : ", err)
		return err
	}

	types := strings.ToUpper(d.Get(constants.TYPE).(string))
	if types == constants.POST || types == constants.PUT || types == constants.DELETE || types == constants.GET {

		actionID := d.Get(constants.APPVIEWX_ACTION_ID).(string)
		payloadString := d.Get(constants.PAYLOAD).(string)

		var masterPayloadFileName = d.Get(constants.MASTER_PAYLOAD).(string)
		if d.Get(constants.MASTER_PAYLOAD) == "" {
			masterPayloadFileName = "./payload.json"
		}

		log.Println("Input minimal payload : ", payloadString)

		payloadMinimal := make(map[string]interface{})
		json.Unmarshal([]byte(payloadString), &payloadMinimal)

		masterPayload := GetMasterPayloadApplyingMinimalPayload(masterPayloadFileName, payloadMinimal)
		log.Println("masterPayload : ", masterPayload)

		queryParams := make(map[string]string)
		queryParams[constants.GW_SOURCE] = appviewxGwSource

		var queryParamReceived = d.Get(constants.QUERY_PARAMS).(map[string]interface{})
		for k, v := range queryParamReceived {
			queryParams[k] = v.(string)
		}

		url := GetURL(appviewxEnvironmentIP, appviewxEnvironmentPort, actionID, queryParams, appviewxEnvironmentIsHTTPS)

		var headers = d.Get(constants.HEADERS).(map[string]interface{})
		if len(headers) == 0 {
			headers["Content-Type"] = "application/json"
			headers["Accept"] = "application/json"
		}

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

		log.Println("API ionvoke success")
		d.SetId(strconv.Itoa(rand.Int()))
		return resourceCertificateServerRead(d, m)
	}
	return nil
}
