/*
 * Copyright 2018. Akamai Technologies, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
//	"bufio"
	"fmt"
	"strings"
	"text/template"

	"os"
	"io/ioutil"

	"encoding/json"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/janeczku/go-spinner"
) 

type EdgeHostname struct {
	EdgeHostname			string
	ProductName			string
	ID				string
	IPv6				string
	EdgeHostnameResourceName	string
}

type Hostname struct {
	Hostname			string
	EdgeHostnameResourceName	string
}

type Variable struct {
	client.Resource
	Name        string 
	Value       string
	Description string
	Hidden      bool 
	Sensitive   bool
}

type TFData struct {
	GroupName 		string
	PropertyResourceName	string
	PropertyName		string
	CPCodeID		string
	CPCodeName		string
	ProductID		string
	ProductName		string
	RuleFormat		string
	IsSecure		string
	EdgeHostnames		map[string]EdgeHostname
	Hostnames		map[string]Hostname
	Section			string
	Emails			[]string
	Variables		[]Variable
}

func cmdCreate(c *cli.Context) error {

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("property name is required"), 1)
	}

	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	papi.Init(config)

	var tfData TFData
	tfData.EdgeHostnames = make(map[string]EdgeHostname)
	tfData.Hostnames = make(map[string]Hostname)
	tfData.Emails = make([]string, 0)

	tfData.Section = c.GlobalString("section")

	// Get Property
	propertyName := c.Args().First()
	s := spinner.StartNew(fmt.Sprintf("Fetching property...."))

	property := findProperty(propertyName);
	if property == nil {
		s.Stop()
		return cli.NewExitError(color.RedString("Property not found "), 1)
	}

	tfData.PropertyName = property.PropertyName
	tfData.PropertyResourceName = strings.Replace(property.PropertyName, ".", "-", -1);

	s.Stop()
	fmt.Printf("Fetching property...... [%s]\n", color.GreenString("OK"))

	s = spinner.StartNew(fmt.Sprintf("Fetching property rules..."))

	// Get Property Rules
	rules, err := property.GetRules();

	if err != nil {
		s.Stop()
		return cli.NewExitError(color.RedString("Property rules not found: ", err), 1)
	}

	tfData.IsSecure = "false"
	if (rules.Rule.Options.IsSecure) {
		tfData.IsSecure = "true"
	}

	for _, behaviour := range rules.Rule.Behaviors {
		_ = behaviour
		if (behaviour.Name == "cpCode") {
			options := behaviour.Options
			value := options["value"].(map[string]interface{})
			tfData.CPCodeID = fmt.Sprintf("%.0f", value["id"].(float64))
		}
	}

	// Get Variables
	tfData.Variables = make([]Variable, 0)
	for _, variable := range rules.Rule.Variables {
		var v Variable
		v.Name = variable.Name
		v.Value = variable.Value
		v.Description = variable.Description
		v.Hidden = variable.Hidden
		v.Sensitive = variable.Sensitive
		tfData.Variables = append(tfData.Variables, v)
	}

	// Get Rule Format
	tfData.RuleFormat = rules.RuleFormat

	// Save Property Rules
	jsonBody, err := json.MarshalIndent(rules, "", "  ")
	err = ioutil.WriteFile("rules.json", jsonBody, 0644)

	if err != nil {
		s.Stop()
		return cli.NewExitError(color.RedString("Can't write property rules: ", err), 1)
	}

	s.Stop()
	fmt.Printf("Fetching property rules...... [%s]\n", color.GreenString("OK"))

	// Get Group
	s = spinner.StartNew( fmt.Sprintf("Fetching group..."))
	group, err := getGroup(property.GroupID)
	if err != nil {
		s.Stop()
		return cli.NewExitError(color.RedString("Group not found: %s", err), 1)
	}

	tfData.GroupName = group.GroupName;

	s.Stop()
	fmt.Printf("Fetching group...... [%s]\n", color.GreenString("OK"))

	// Get Version
	s = spinner.StartNew(fmt.Sprintf("Fetching property version..."))
	version, err := getVersion(property);
	if err != nil {
		s.Stop()
		return cli.NewExitError(color.RedString("Version not found: %s", err), 1)
	}

	tfData.ProductID = version.ProductID;

	s.Stop()
	fmt.Printf("Fetching property version...... [%s]\n", color.GreenString("OK"))

	// Get Product
	s = spinner.StartNew(fmt.Sprintf("Fetching product name..."))
	product, err := getProduct(tfData.ProductID, property.Contract);
	if err != nil {
		s.Stop()
		return cli.NewExitError(color.RedString("Product not found: %s", err), 1)
	}

	tfData.ProductName = product.ProductName;

	s.Stop()
	fmt.Printf("Fetching product name...... [%s]\n", color.GreenString("OK"))

	// Get Hostnames
	s = spinner.StartNew( fmt.Sprintf("Fetching hostnames..."))
	hostnames, err := getHostnames(property, version);

	if err != nil {
		s.Stop()
		return cli.NewExitError(color.RedString("Hostnames not found: %s", err), 1)
	}

	for _, hostname := range hostnames.Hostnames.Items {
		_ = hostname
		cnameTo := hostname.CnameTo
		cnameFrom := hostname.CnameFrom
		cnameToResource := strings.Replace(cnameTo, ".", "-", -1)


		var edgeHostnameN EdgeHostname
		edgeHostnameN.EdgeHostname = cnameTo
		edgeHostnameN.EdgeHostnameResourceName = cnameToResource
		edgeHostnameN.ProductName = product.ProductName 
		edgeHostnameN.IPv6 = isIPv6(property, hostname.EdgeHostnameID)
		tfData.EdgeHostnames[cnameToResource] = edgeHostnameN

		var hostnamesN Hostname
		hostnamesN.Hostname = cnameFrom
		hostnamesN.EdgeHostnameResourceName = cnameToResource
		tfData.Hostnames[cnameFrom] = hostnamesN

	}

	s.Stop()
	fmt.Printf("Fetching hostnames...... [%s]\n", color.GreenString("OK"))

	// Get CPCode Name
	s = spinner.StartNew( fmt.Sprintf("Fetching CPCode name..."))
	cpcodeName, err := getCPCode(property, tfData.CPCodeID);
	if err != nil {
		s.Stop()
		return cli.NewExitError(color.RedString("Product not found: %s", err), 1)
	}

	tfData.CPCodeName = cpcodeName;

	s.Stop()
	fmt.Printf("Fetching CPCode name...... [%s]\n", color.GreenString("OK"))

	// Get contact details
	s = spinner.StartNew( fmt.Sprintf("Fetching contact details..."))
	activations, err := property.GetActivations()
	if err != nil {
		tfData.Emails = append(tfData.Emails, "")
	} else {
		a, err := activations.GetLatestStagingActivation("")
		if err != nil {
			tfData.Emails = append(tfData.Emails, "")
		} else {
			tfData.Emails = a.NotifyEmails
		}
	}

	s.Stop()
	fmt.Printf("Fetching contact details...... [%s]\n", color.GreenString("OK"))


	// Save file
	s = spinner.StartNew( fmt.Sprintf("Saving TF definition..."))
	err = saveTerraformDefinition(tfData);
	if err != nil {
		s.Stop()
		return cli.NewExitError(color.RedString("Couldn't save tf file: %s", err), 1)
	}
	s.Stop()


	fmt.Printf("Saving TF definition...... [%s]\n", color.GreenString("OK"))

	return nil;


}

func getHostnames(property *papi.Property, version *papi.Version) (*papi.Hostnames, error) {
	hostnames, err := property.GetHostnames(version)
	if (err != nil) {
		return nil, err
	}
	return hostnames, nil
}

func isIPv6(property *papi.Property, ehn string) (string) {
	edgeHostnames, err := papi.GetEdgeHostnames(property.Contract, property.Group, "")
	if (err != nil) {
		return "false"
	}
	for _, edgehostname := range edgeHostnames.EdgeHostnames.Items {
		_ = edgehostname
		if (edgehostname.EdgeHostnameID == ehn) {
			if (edgehostname.IPVersionBehavior == "IPV4") {
				return "false"
			} else {
				return "true"
			}
		}
	}
	return "false"
}

func getCPCode(property *papi.Property, cpCodeID string) (string, error) {
	cpCode := papi.NewCpCodes(property.Contract, property.Group).NewCpCode()
	cpCode.CpcodeID = cpCodeID
	err := cpCode.GetCpCode()
	if err != nil {
		return "", err
	}
	return cpCode.CpcodeName, nil
}

func findProperty(name string) *papi.Property {
	results, err := papi.Search(papi.SearchByPropertyName, name)
	if err != nil {
		return nil
	}

	if err != nil || results == nil {
		return nil
	}

	property := &papi.Property{
		PropertyID: results.Versions.Items[0].PropertyID,
		Group: &papi.Group{
			GroupID: results.Versions.Items[0].GroupID,
		},
		Contract: &papi.Contract{
			ContractID: results.Versions.Items[0].ContractID,
		},
	}

	err = property.GetProperty()
	if err != nil {
		return nil
	}

	return property
}

func getVersion(property *papi.Property) (*papi.Version, error) {


	versions, err := property.GetVersions()
	if err != nil {
		return nil, err
	}

	version, err := versions.GetLatestVersion("")
	if err != nil {
		return nil, err
	}

	return version, nil
}

func getGroup(groupID string) (*papi.Group, error) {
	groups := papi.NewGroups()
	e := groups.GetGroups()
	if e != nil {
		return nil, e
	}

	group, e := groups.FindGroup(groupID)
	if e != nil {
		return nil, e
	}

	return group, nil
}

func getProduct(productID string, contract *papi.Contract) (*papi.Product, error) {
	if contract == nil {
		return nil, nil
	}

	products := papi.NewProducts()
	e := products.GetProducts(contract)
	if e != nil {
		return nil, e
	}

	product, e := products.FindProduct(productID)
	if e != nil {
		return nil, e
	}

	return product, nil
}

func saveTerraformDefinition(data TFData) error {

	template, err := template.New("tf").Parse(
		"provider \"akamai\" {\n" +
  		" edgerc = \"~/.edgerc\"\n" +
  		" papi_section = \"{{.Section}}\"\n" +
		"}\n" +
		"\n" +
		"\n" +
		"data \"akamai_group\" \"group\" {\n" +
  		" name = \"{{.GroupName}}\"\n" +
		"}\n" +
		"\n" +
		"data \"akamai_contract\" \"contract\" {\n" +
		"  group = data.akamai_group.group.name\n" +
		"}\n" +
		"\n" +
		"data \"template_file\" \"rules\" {\n" +
   		" template = file(\"${path.module}/rules.json\")\n" +
		"}\n" +
		"\n" +
		"resource \"akamai_cp_code\" \"{{.PropertyResourceName}}\" {\n" +
    		" product  = \"prd_{{.ProductName}}\"\n" +
    		" contract = data.akamai_contract.contract.id\n" +
    		" group = data.akamai_group.group.id\n" +
    		" name = \"{{.CPCodeName}}\"\n" +
		"}\n" +
		"\n" +


		// Edge hostname loop
		"{{range .EdgeHostnames}}" +
		"resource \"akamai_edge_hostname\" \"{{.EdgeHostnameResourceName}}\" {\n" +
    		" product  = \"prd_{{.ProductName}}\"\n" +
    		" contract = data.akamai_contract.contract.id\n" +
    		" group = data.akamai_group.group.id\n" +
		" ipv6 = {{.IPv6}}\n" +
    		" edge_hostname = \"{{.EdgeHostname}}\"\n" +
		"}\n" +
		"\n" +
		"{{end}}" +

		"{{if .Variables}}" +
		"resource \"akamai_property_variables\" \"variables\" {\n" +
		" variables {\n" +
		"{{range .Variables}}" +
		"  variable {\n" +
    		"    name  = \"{{.Name}}\"\n" +
    		"    value  = \"{{.Value}}\"\n" +
    		"    description  = \"{{.Description}}\"\n" +
    		"    hidden  = \"{{.Hidden}}\"\n" +
    		"    sensitive  = \"{{.Sensitive}}\"\n" +
		"{{end}}" +
		"  }\n" +
		" }\n" +
		"}\n" +
		"\n" +
		"{{end}}" +
	

		"resource \"akamai_property\" \"{{.PropertyResourceName}}\" {\n" +
  		" name = \"{{.PropertyName}}\"\n" +
  		" cp_code = akamai_cp_code.{{.PropertyResourceName}}.id\n" +
  		" contact = [\"\"]\n" +
  		" contract = data.akamai_contract.contract.id\n" +
  		" group = data.akamai_group.group.id\n" +
  		" product = \"prd_{{.ProductName}}\"\n" +
  		" rule_format = \"{{.RuleFormat}}\"\n" +
		" hostnames = {\n" +

		"{{range .Hostnames}}" +
        	"  \"{{.Hostname}}\" = akamai_edge_hostname.{{.EdgeHostnameResourceName}}.edge_hostname\n" +
		"{{end}}" +
		" }\n" +
  		" rules = data.template_file.rules.rendered\n" +
  		" is_secure = {{.IsSecure}}\n" +
		"{{if .Variables}}" +
		" variables = akamai_property_variables.variables.json\n" +
		"{{end}}" +
		"}\n" +
		"\n" +
		"resource \"akamai_property_activation\" \"{{.PropertyResourceName}}\" {\n" +
        	" property = akamai_property.{{.PropertyResourceName}}.id\n" +
        	" contact = [\"{{range $index, $element := .Emails}}{{if $index}},{{end}}{{$element}}{{end}}\"]\n" +
        	" network = upper(var.env)\n" +
        	" activate = true\n" +
		"}\n")

	if err != nil {
		return err;
	}
	f, err := os.Create("property.tf")
	if err != nil {
		return err;
	}

	err = template.Execute(f, data)
	if err != nil {
		return err;
	}
	f.Close()

	// variables
	f, err = os.Create("variables.tf")
	if err != nil {
		return err;
	}

	f.WriteString("variable \"env\" {\n")
	f.WriteString(" default = \"staging\"\n")
	f.WriteString("}\n")
	f.Close()

	return nil;
}

