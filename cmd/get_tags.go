package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"splice-cloud-util/cmd/objects"
	"strings"

	"github.com/blang/semver/v4"
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
	PreRun: func(cmd *cobra.Command, args []string) {
		if dockerUser == "" || dockerPass == "" {
			logrus.Info("Unable to locate Docker login credentials, please run 'splice-cloud-util init' to configure access.")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		org, _ := cmd.Flags().GetString("org")
		repo, _ := cmd.Flags().GetString("repo")
		prefix, _ := cmd.Flags().GetString("prefix")
		returnPrefix, _ := cmd.Flags().GetBool("return-prefix")
		top, _ := cmd.Flags().GetInt("top")
		next, _ := cmd.Flags().GetBool("next")

		if next {
			nextSemVer, err := getNextDockerHubRepositoryTag(org, repo, prefix, returnPrefix)
			if err != nil {
				logrus.WithError(err).Error("Error getting next SemVer")
			}
			fmt.Println(nextSemVer)
		} else {
			tags, err := getDockerHubRepositoryTags(org, repo, prefix, top)
			if err != nil {
				logrus.WithError(err).Error("Error getting tags")
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
		}

	},
}

func getNextDockerHubRepositoryTag(org string, repo string, prefix string, returnPrefix bool) (string, error) {

	token, err := getToken()
	if err != nil {
		logrus.WithError(err).Error("Error getting token ")
		return "", err
	}

	// curl -s -H "Authorization: JWT ${TOKEN}" https://hub.docker.com/v2/repositories/splicemachine/zeppelin/tags/\?page_size\=10000
	restClient := resty.New()

	// Top should never exceed 100, if you want more than 100 items off the top
	// then you'll need to re-write to exit the for loop at the proper time
	pageSize := 100
	apiURI := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s/tags/?page_size=%d", org, repo, pageSize)
	maxSemVer, _ := semver.Parse("0.0.0")
	for {
		resp, resperr := restClient.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", fmt.Sprintf("JWT %s", token)).
			Get(apiURI)
		if resperr != nil {
			logrus.WithError(resperr).Error(fmt.Sprintf("Error fetching tags: %s", apiURI))
			return "", resperr
		}
		tag := &objects.DockerHubTagResult{}
		marsherr := json.Unmarshal(resp.Body(), &tag)
		if marsherr != nil {
			logrus.WithError(marsherr).Error("Error decoding json")
			return "", marsherr
		}
		for _, t := range tag.Results {
			if len(prefix) == 0 || strings.HasPrefix(t.Name, prefix) {
				var sv semver.Version
				var err error
				if len(prefix) > 0 {
					sv, err = semver.Parse(strings.Replace(t.Name, prefix, "", 1))
					if err != nil {
						logrus.Warn(fmt.Sprintf("Error parsing SemVer for %s", t.Name))
						break
					}
				} else {
					sv, err = semver.Parse(t.Name)
					if err != nil {
						logrus.Warn(fmt.Sprintf("Error parsing SemVer for %s", t.Name))
						break
					}
				}
				if sv.GT(maxSemVer) {
					maxSemVer = sv
				}

			}
		}
		if tag.Next == "" {
			break
		} else {
			apiURI = tag.Next
		}
	}

	maxSemVer.Patch++
	if len(prefix) > 0 && returnPrefix {
		return fmt.Sprintf("%s%s", prefix, maxSemVer.String()), nil
	}
	return maxSemVer.String(), nil

}

func getDockerHubRepositoryTags(org string, repo string, prefix string, top int) (objects.DockerHubTagInfoList, error) {

	token, err := getToken()
	if err != nil {
		logrus.WithError(err).Error("Error getting token ")
		return objects.DockerHubTagInfoList{}, err
	}

	tags := objects.DockerHubTagInfoList{}

	// curl -s -H "Authorization: JWT ${TOKEN}" https://hub.docker.com/v2/repositories/splicemachine/zeppelin/tags/\?page_size\=10000
	restClient := resty.New()

	// Top should never exceed 100, if you want more than 100 items off the top
	// then you'll need to re-write to exit the for loop at the proper time
	pageSize := 100
	apiURI := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s/tags/?page_size=%d", org, repo, pageSize)
	tagsCollected := 0
	foundAll := false
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
			if len(prefix) == 0 || strings.HasPrefix(t.Name, prefix) {
				tl := objects.DockerHubTagInfo{}
				tl = t
				tags.List = append(tags.List, tl)
				tagsCollected++
				if top > 0 && tagsCollected == top {
					foundAll = true
					break
				}
			}
		}
		if tag.Next == "" || foundAll {
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
	getTagsCmd.Flags().String("prefix", "", "Specify the prefix for the tag")
	getTagsCmd.Flags().BoolP("return-prefix", "r", false, "Return the prefix as part of the SemVer output")
	getTagsCmd.Flags().Int("top", 0, "Specify the top number of TAGs to return")
	getTagsCmd.Flags().BoolP("next", "n", false, "Get the next valid SEMVER based on the highest one")
	getTagsCmd.MarkFlagRequired("org")
	getTagsCmd.MarkFlagRequired("repo")
}
