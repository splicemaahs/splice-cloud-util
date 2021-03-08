package objects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Human readable sizes
const (
	TB = 1000000000000
	GB = 1000000000
	MB = 1000000
	KB = 1000
)

// DockerHubTagInfoList - Array of TagInfo
type DockerHubTagInfoList struct {
	List []DockerHubTagInfo
}

// DockerHubTagResult - structure for repository list
type DockerHubTagResult struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Creator int `json:"creator"`
		ID      int `json:"id"`
		ImageID int `json:"image_id"`
		Images  []struct {
			Architecture string `json:"architecture"`
			Features     string `json:"features"`
			Variant      string `json:"variant"`
			Digest       string `json:"digest"`
			OS           string `json:"os"`
			OSFeatures   string `json:"os_features"`
			OSVersion    string `json:"os_version"`
			Size         int    `json:"size"`
		} `json:"images"`
		LastUpdated         string `json:"last_updated"`
		LastUpdater         int    `json:"last_updater"`
		LastUpdaterUsername string `json:"last_updater_username"`
		Name                string `json:"name"`
		Repository          int    `json:"repository"`
		FullSize            int    `json:"full_size"`
		V2                  bool   `json:"v2"`
	} `json:"results"`
}

// DockerHubTagInfo - tag info
type DockerHubTagInfo struct {
	Creator int `json:"creator"`
	ID      int `json:"id"`
	ImageID int `json:"image_id"`
	Images  []struct {
		Architecture string `json:"architecture"`
		Features     string `json:"features"`
		Variant      string `json:"variant"`
		Digest       string `json:"digest"`
		OS           string `json:"os"`
		OSFeatures   string `json:"os_features"`
		OSVersion    string `json:"os_version"`
		Size         int    `json:"size"`
	} `json:"images"`
	LastUpdated         string `json:"last_updated"`
	LastUpdater         int    `json:"last_updater"`
	LastUpdaterUsername string `json:"last_updater_username"`
	Name                string `json:"name"`
	Repository          int    `json:"repository"`
	FullSize            int    `json:"full_size"`
	V2                  bool   `json:"v2"`
}

// ToJSON - Write the output as JSON
func (tags *DockerHubTagInfoList) ToJSON() error {

	tagsJSON, enverr := json.MarshalIndent(tags, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}
	fmt.Println(string(tagsJSON[:]))

	return nil

}

// ToGRON - Write the output as GRON
func (tags *DockerHubTagInfoList) ToGRON() error {
	listJSON, enverr := json.MarshalIndent(tags, "", "  ")
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
func (tags *DockerHubTagInfoList) ToYAML() error {

	tagsYAML, enverr := yaml.Marshal(tags)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	}
	fmt.Println(string(tagsYAML[:]))

	return nil

}

// ToTEXT - Write the output as TEXT
func (tags *DockerHubTagInfoList) ToTEXT(noHeaders bool) error {

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
	if !noHeaders {
		table.SetHeader([]string{"TAG", "LAST_UPDATED_BY", "ARCH", "SIZE"})
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
	for _, r := range tags.List {
		for _, i := range r.Images {
			row := []string{r.Name, r.LastUpdaterUsername, i.Architecture, lenReadable(i.Size, 2)}
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

func lenReadable(length int, decimals int) (out string) {
	var unit string
	var i int
	var remainder int

	// Get whole number, and the remainder for decimals
	// if length > TB {
	// unit = "TB"
	// i = length / TB
	// remainder = length - (i * TB)
	// } else if length > GB {
	if length > GB {
		unit = "GB"
		i = length / GB
		remainder = length - (i * GB)
	} else if length > MB {
		unit = "MB"
		i = length / MB
		remainder = length - (i * MB)
	} else if length > KB {
		unit = "KB"
		i = length / KB
		remainder = length - (i * KB)
	} else {
		return strconv.Itoa(length) + " B"
	}

	if decimals == 0 {
		return strconv.Itoa(i) + " " + unit
	}

	// This is to calculate missing leading zeroes
	width := 0
	if remainder > GB {
		width = 12
	} else if remainder > MB {
		width = 9
	} else if remainder > KB {
		width = 6
	} else {
		width = 3
	}

	// Insert missing leading zeroes
	remainderString := strconv.Itoa(remainder)
	for iter := len(remainderString); iter < width; iter++ {
		remainderString = "0" + remainderString
	}
	if decimals > len(remainderString) {
		decimals = len(remainderString)
	}

	return fmt.Sprintf("%d.%s %s", i, remainderString[:decimals], unit)
}
