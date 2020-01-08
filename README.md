>COMMANDS

>> Build
```
> cd ../terraform-provider-appviewx
> go build -o terraform-provider-appviewx

Need to build the plugin with 'terraform-provider-' as prefix

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

>> provider name should be 'appviewx' if the plugin name is 'terraform-provider-appviewx'