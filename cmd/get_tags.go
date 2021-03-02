package cmd

import (
	"encoding/json"
	"fmt"
	"splice-cloud-util/cmd/objects"
	"strings"

	// "github.com/hashicorp/vault/api"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var getTagsCmd = &cobra.Command{
	Use:   "dh-tags",
	Short: "Get a list of tags for a DockerHub repository",
	Long: `EXAMPLE 1
	List all the TAGs for an Organization / Repository

	#> splice-cloud-util get dh-tags --org splicemachine --repo sm_k8_kafka-3.0.0

EXAMPLE 2
	List the top N tags for an Organization / Repository

	#> splice-cloud-util get dh-tags --org splicemachine --repo sm_k8_splicectlapi --top 5
`,
	Run: func(cmd *cobra.Command, args []string) {
		org, _ := cmd.Flags().GetString("org")
		repo, _ := cmd.Flags().GetString("repo")
		top, _ := cmd.Flags().GetInt("top")

		tags, err := getDockerHubRepositoryTags(org, repo, top)
		if err != nil {
			logrus.WithError(err).Error("Error getting repositories")
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			tags.ToJSON()
		case "gron":
			tags.ToGRON()
		case "yaml":
			tags.ToYAML()
		case "text", "table":
			tags.ToTEXT(noHeaders)
		}

	},
}

func getDockerHubRepositoryTags(org string, repo string, top int) (objects.DockerHubTagInfoList, error) {

	token, err := getToken()
	if err != nil {
		logrus.WithError(err).Error("Error getting token ")
		return objects.DockerHubTagInfoList{}, err
	}

	tagList, _ := getRepositoryTags(org, repo, token, top)

	// tagJSON, jsonerr := json.MarshalIndent(tagList, "", "    ")
	// if jsonerr != nil {
	// 	return nil, jsonerr
	// }
	// fmt.Println(string(tagJSON[:]))

	return tagList, nil
}

func getRepositoryTags(org string, repo string, token string, top int) (objects.DockerHubTagInfoList, error) {

	tags := objects.DockerHubTagInfoList{}

	// curl -s -H "Authorization: JWT ${TOKEN}" https://hub.docker.com/v2/repositories/splicemachine/zeppelin/tags/\?page_size\=10000
	restClient := resty.New()

	// Top should never exceed 100, if you want more than 100 items off the top
	// then you'll need to re-write to exit the for loop at the proper time
	pageSize := 100
	if top > 0 {
		pageSize = top
	}
	apiURI := fmt.Sprintf("https://hub.docker.com/v2/repositories/splicemachine/%s/tags/?page_size=%d", repo, pageSize)
	for {
		resp, resperr := restClient.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", fmt.Sprintf("JWT %s", token)).
			Get(apiURI)
		if resperr != nil {
			logrus.WithError(resperr).Error(fmt.Sprintf("Error fetching tags: %s", apiURI))
			return tags, resperr
		}
		tag := &objects.DockerHubTagResult{}
		marsherr := json.Unmarshal(resp.Body(), &tag)
		if marsherr != nil {
			logrus.WithError(marsherr).Error("Error decoding json")
			return tags, marsherr
		}
		for _, t := range tag.Results {
			tl := objects.DockerHubTagInfo{}
			tl = t
			tags.List = append(tags.List, tl)
		}
		if tag.Next == "" || top > 0 {
			break
		} else {
			apiURI = tag.Next
		}
	}

	return tags, nil

}

func init() {
	getCmd.AddCommand(getTagsCmd)

	getTagsCmd.Flags().String("org", "", "Specify the dockerhub organization")
	getTagsCmd.Flags().String("repo", "", "Specify the dockerhub repository")
	getTagsCmd.Flags().Int("top", 0, "Specify the top number of TAGs to return")
	getTagsCmd.MarkFlagRequired("org")
	getTagsCmd.MarkFlagRequired("repo")
}
