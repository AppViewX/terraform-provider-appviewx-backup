resource "appviewx_stateless_api" "random-number-014" {				
	payload= <<EOF
	{
  "caConnectorInfo": {
    "certificateAuthority": "AppViewX",
    "caSettingName": "AppViewX CA",
    "name": "AppViewX CA connector",
    "csrParameters": {
      "commonName": "terraform014.appviewx.com",
      "hashFunction": "SHA256",
      "keyType": "RSA",
      "bitLength": "2048",
      "certificateCategories": [
        "Server",
        "Client"
      ],
      "enhancedSANTypes": {
        "dNSNames": [
          "terraform014.appviewx.com"
        ]
      }
    },
    "validityInDays": 365
  }
 }

	EOF
	headers= {
		"Content-Type":"application/json",
		"Accept":"application/json"
	}
	action_id= "certificate/create"
	type=  "post"
	config_file = "./config.json"
	master_payload ="./payload.json"
}

