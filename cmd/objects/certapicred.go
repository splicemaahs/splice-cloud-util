package objects

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// CertAPICredList - Array of Creds
type CertAPICredList struct {
	List []CertAPICred
}

// CertAPICred - Collection of Details
type CertAPICred struct {
	Account   string `json:"account"`
	ClientID  string `json:"clientId"`
	SecretKey string `json:"secretKey"`
}

// ToJSON - Write the output as JSON
func (credDetail *CertAPICredList) ToJSON() error {

	credJSON, enverr := json.MarshalIndent(credDetail, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}
	fmt.Println(string(credJSON[:]))

	return nil

}

// ToYAML - Write the output as YAML
func (credDetail *CertAPICredList) ToYAML() error {

	credYAML, enverr := yaml.Marshal(credDetail)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	}
	fmt.Println(string(credYAML[:]))

	return nil

}

// ToTEXT - Write the output as TEXT
func (credDetail *CertAPICredList) ToTEXT(noHeaders bool) error {

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
	if !noHeaders {
		table.SetHeader([]string{"ACCOUNT", "CLIENT_ID", "SECRET_KEY"})
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	}
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	for _, cred := range credDetail.List {
		row := []string{cred.Account, cred.ClientID, cred.SecretKey}
		table.Append(row)
	}
	table.Render()

	// ****************** Go Templates **************************
	// I am going to leave this here, Go Templates are very cool, and clearly
	// Helm Templating is just an extension/use of that, which makes it very
	// familiar.  However, it becomes a little bit cumbersome when you start
	// to add more and more fields and need to determine the length of each
	// field, adjusting the '%-NNs' for the widest text in the field.
	// github.com/olekukonko/tablewriter may be overkill for what we are doing
	// here, though it certainly seems to be popular and very extendable.

	// tmp1 := template.New("Template_1")

	// tmp1, _ = tmp1.Parse(environmentsTemplate(noHeaders))

	// err := tmp1.Execute(os.Stdout, out)
	// if err != nil {
	// 	return err
	// }

	return nil

}
