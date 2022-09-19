>COMMANDS

>> Build
```
> cd ../terraform-provider-appviewx
> go build -o terraform-provider-appviewx

Need to build the plugin with 'terraform-provider-' as prefix

```

>> provider name should be 'appviewx' if the plugin name is 'terraform-provider-appviewx'

>> to enable logs  ( TRACE, DEBUG, INFO, WARN or ERROR )
```
	export TF_LOG=TRACE
```

>> Sample .tf file
```
terraform {
  required_providers {
    appviewx = {
      version = "0.2"
      source  = "appviewx.com/provider/appviewx"
    }
  }
}

provider "appviewx"{
  appviewx_username="admin"
	appviewx_password="AppViewX@123"
	appviewx_environment_is_https=true 
	appviewx_environment_ip="192.168.220.129"
	appviewx_environment_port="31443"
}

resource "appviewx_automation" "newcert"{
 payload= <<EOF
 {
  "payload" : {
    "data" : {
      "input" : {
        "requestData" : [ {
          "sequenceNo" : 1,
          "scenario" : "scenario",
          "fieldInfo" : {
            "commonname" : "www.sample.appviewx.com",
            "email" : "vigneshkumar.k@appviewx.com"
          }
        } ]
      },
      "task_action" : 1
    },
    "header" : {
      "workflowName" : "Copy of Generate AppViewX Certificate with Email Approval"
    }
  }
}
EOF
action_id= "visualworkflow-submit-request"

  }

```

>> Keep the .tf file in the current folder

>> keep the "terraform-provider-appviewx" binary file under "~/.terraform.d/plugins/appviewx.com/provider/appviewx/0.2/linux_386"   ( linux_386 is sample, need to change based on the installation system architecture )

>> Run the following commands, to reset and trigger the request
```
	rm -rf ./terraform.tfstate;
	terraform init;
	terraform apply;
```
>>	terraform init  
'''
	To initialize the terraform with the given plugin and validate the .tf files
	( ensure the plugin and .tf or .tf.json files placed in the current folder path )
'''

>> terraform apply
```
	To analyze the local, remote state and carryout the required actions based on the given .tf or .tf.json files
```