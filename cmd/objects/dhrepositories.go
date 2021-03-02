package objects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// DockerHubRepositoriesList - Array of Creds
type DockerHubRepositoriesList struct {
	List []DockerHubRepository
}

// DockerHubRepository - Collection of Details
type DockerHubRepository struct {
	User                string `json:"user"`
	Name                string `json:"name"`
	Namespace           string `json:"namespace"`
	RepositoryType      string `json:"repository_type"`
	Status              int    `json:"status"`
	Description         string `json:"description"`
	IsPrivate           bool   `json:"is_private"`
	IsAutomated         bool   `json:"is_automated"`
	CanEdit             bool   `json:"can_edit"`
	StarCount           int    `json:"star_count"`
	PullCount           int    `json:"pull_count"`
	LastUpdated         string `json:"last_updated"`
	IsMigrated          bool   `json:"is_migrated"`
	LastUpdaterUsername string `json:"last_updater_username"`
	LatestMaster        string `json:"latest_master"`
	LatestMasterUpdated string `json:"latest_master_updated"`
}

// DockerHubRepoResult - structure for repository list
type DockerHubRepoResult struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		User                string `json:"user"`
		Name                string `json:"name"`
		Namespace           string `json:"namespace"`
		RepositoryType      string `json:"repository_type"`
		Status              int    `json:"status"`
		Description         string `json:"description"`
		IsPrivate           bool   `json:"is_private"`
		IsAutomated         bool   `json:"is_automated"`
		CanEdit             bool   `json:"can_edit"`
		StarCount           int    `json:"star_count"`
		PullCount           int    `json:"pull_count"`
		LastUpdated         string `json:"last_updated"`
		IsMigrated          bool   `json:"is_migrated"`
		LastUpdaterUsername string `json:"last_updater_username"`
		LatestMaster        string `json:"latest_master"`
		LatestMasterUpdated string `json:"latest_master_updated"`
	} `json:"results"`
}

// ToJSON - Write the output as JSON
func (repo *DockerHubRepositoriesList) ToJSON() error {

	repoJSON, enverr := json.MarshalIndent(repo, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}
	fmt.Println(string(repoJSON[:]))

	return nil

}

// ToGRON - Write the output as GRON
func (repo *DockerHubRepositoriesList) ToGRON() error {
	listJSON, enverr := json.MarshalIndent(repo, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}

	subReader := strings.NewReader(string(listJSON[:]))
	subValues := &bytes.Buffer{}
	ges := gron.NewGron(subReader, subValues)
	ges.SetMonochrome(false)
	serr := ges.ToGron()
	if serr != nil {
		logrus.Error("Problem generating gron syntax", serr)
		return serr
	}
	fmt.Println(string(subValues.Bytes()))

	return nil

}

// ToYAML - Write the output as YAML
func (repo *DockerHubRepositoriesList) ToYAML() error {

	repoYAML, enverr := yaml.Marshal(repo)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	}
	fmt.Println(string(repoYAML[:]))

	return nil

}

// ToTEXT - Write the output as TEXT
func (repo *DockerHubRepositoriesList) ToTEXT(noHeaders bool, updated bool, release bool) error {

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
	if !noHeaders {
		if updated && release {
			table.SetHeader([]string{"REPOSITORY", "PRIVATE", "LAST_UPDATED_BY", "LATEST_MASTER", "MASTER_UPDATED"})
		} else {
			if updated {
				table.SetHeader([]string{"REPOSITORY", "PRIVATE", "LAST_UPDATED_BY"})
			} else {
				if release {
					table.SetHeader([]string{"REPOSITORY", "PRIVATE", "LATEST_MASTER", "MASTER_UPDATED"})
				} else {
					table.SetHeader([]string{"REPOSITORY", "PRIVATE", "LAST_UPDATED"})
				}
			}
		}
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
	for _, r := range repo.List {
		row := []string{}
		if updated && release {
			row = []string{r.Name, fmt.Sprintf("%t", r.IsPrivate), r.LastUpdaterUsername, r.LatestMaster, r.LatestMasterUpdated}
		} else {
			if updated {
				row = []string{r.Name, fmt.Sprintf("%t", r.IsPrivate), r.LastUpdaterUsername}
			} else {
				if release {
					row = []string{r.Name, fmt.Sprintf("%t", r.IsPrivate), r.LatestMaster, r.LatestMasterUpdated}
				} else {
					row = []string{r.Name, fmt.Sprintf("%t", r.IsPrivate), r.LastUpdated}
				}
			}
		}
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
