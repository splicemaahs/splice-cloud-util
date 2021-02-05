package objects

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// DetailList - Collection of Details
type DetailList struct {
	EnvironmentList []Detail `json:"environmentList,omitempty"`
}

// Detail - Return an environment name with cluster list
type Detail struct {
	EnvironmentName string        `json:"environmentName,omitempty"`
	ClusterList     []ClusterList `json:"clusterList,omitempty"`
}

// ClusterList - details about the k8s cluster
type ClusterList struct {
	ClusterName string `json:"clusterName,omitempty"`
	CreatedBy   string `json:"createdBy,omitempty"`
	DateCreated string `json:"dateCreated,omitempty"`
	State       string `json:"state,omitempty"`
}

// ToJSON - Write the output as JSON
func (envDetail *DetailList) ToJSON() error {

	envJSON, enverr := json.MarshalIndent(envDetail, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	} else {
		fmt.Println(string(envJSON[:]))
	}

	return nil

}

// ToYAML - Write the output as YAML
func (envDetail *DetailList) ToYAML() error {

	envYAML, enverr := yaml.Marshal(envDetail)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	} else {
		fmt.Println(string(envYAML[:]))
	}
	return nil

}

// ToTEXT - Write the output as TEXT
func (envDetail *DetailList) ToTEXT(noHeaders bool) error {

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
	if !noHeaders {
		table.SetHeader([]string{"ENVIRONMENT", "CLUSTER_NAME", "CREATED_BY", "DATE_CREATED", "TRANSITION"})
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
	for _, environment := range envDetail.EnvironmentList {
		for _, cluster := range environment.ClusterList {
			row := []string{environment.EnvironmentName, cluster.ClusterName, cluster.CreatedBy, cluster.DateCreated, cluster.State}
			table.Append(row)
		}
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
