This assumes you already have a working CLI and Terraform installation.

 
Install using "akamai install https://github.com/IanCassTwo/akamai-property-to-terraform.git". Answer "Y" when it offers to install a binary
Create a new empty subdirectory to hold your new Terraform configuration

Type "akamai terraform create <name of property manager config>"

Once it's done, do a "terraform init", followed by a "terraform apply

The plugin allows some flags:-

edgerc - define where your edgerc is. Defaults to your home directory
section - which section to use in the edgerc. Defaults to "default"
accountkey - which account key to use. You can leave this blank if you don't want to use account switching
