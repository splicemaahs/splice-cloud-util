package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	// "github.com/hashicorp/vault/api"
	"splice-cloud-util/cmd/objects"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// AuthSuccess struct
type AuthSuccess struct {
	/* variables */
	Token string `json:"token"`
}

// AuthError struct
type AuthError struct {
	/* variables */
}

// BasicAuthInfo - user/password
type BasicAuthInfo struct {
	Username string
	Password string
}

var getRepositoriesCmd = &cobra.Command{
	Use:     "dh-repositories",
	Aliases: []string{"repo"},
	Short:   "Get a DockerHub repository listing",
	Long: `EXAMPLE 1
	List repositories for the Organiztion

	#> splice-cloud-util get dh-repositories --org splicemachine

EXAMPLE 2
	This example will add the latest RELEASE TAG and who the last person to update
	the repository was.  This run takes a fair amount of time, be patient.

	#> splice-cloud-util get dh-repositories --org splicemachine --add-release --add-updated
`,
	Run: func(cmd *cobra.Command, args []string) {
		org, _ := cmd.Flags().GetString("org")
		addUpdated, _ := cmd.Flags().GetBool("add-updated")
		addRelease, _ := cmd.Flags().GetBool("add-release")

		repos, err := getDockerHubRepositories(org, addUpdated, addRelease)
		if err != nil {
			logrus.WithError(err).Error("Error getting repositories")
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			repos.ToJSON()
		case "gron":
			repos.ToGRON()
		case "yaml":
			repos.ToYAML()
		case "text", "table":
			repos.ToTEXT(noHeaders, addUpdated, addRelease)
		}

	},
}

func getDockerHubRepositories(org string, updated bool, release bool) (objects.DockerHubRepositoriesList, error) {

	token, err := getToken()
	if err != nil {
		logrus.WithError(err).Error("Error getting token ")
		return objects.DockerHubRepositoriesList{}, err
	}

	repositories := objects.DockerHubRepositoriesList{}

	// curl -s -H "Authorization: JWT ${TOKEN}" https://hub.docker.com/v2/repositories/splicemachine/zeppelin/tags/\?page_size\=10000
	restClient := resty.New()

	apiURI := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s?page_size=100", org)
	for {
		resp, resperr := restClient.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", fmt.Sprintf("JWT %s", token)).
			Get(apiURI)
		if resperr != nil {
			logrus.WithError(resperr).Error(fmt.Sprintf("Error fetching repositories: %s", apiURI))
			return objects.DockerHubRepositoriesList{}, resperr
		}
		repo := &objects.DockerHubRepoResult{}
		marsherr := json.Unmarshal(resp.Body(), &repo)
		if marsherr != nil {
			logrus.WithError(marsherr).Error("Error decoding json")
			return objects.DockerHubRepositoriesList{}, marsherr
		}
		for _, r := range repo.Results {
			rl := objects.DockerHubRepository{}
			rl = r
			if updated {
				lastUser, usererr := getDockerHubRepositoryTags(org, rl.Name, "", 1)
				if usererr == nil {
					if len(lastUser.List) > 0 {
						rl.LastUpdaterUsername = lastUser.List[0].LastUpdaterUsername
					} else {
						rl.LastUpdaterUsername = "none"
					}
				}
			}
			if release {
				lastMaster, mastererr := getDockerHubRepositoryTags(org, rl.Name, "", 50)
				if mastererr == nil {
					for _, t := range lastMaster.List {
						if strings.HasPrefix(t.Name, "master") {
							rl.LatestMaster = t.Name
							rl.LatestMasterUpdated = t.LastUpdated
						}
					}
				}
			}
			repositories.List = append(repositories.List, rl)
		}
		if repo.Next == "" {
			break
		} else {
			apiURI = repo.Next
		}
	}

	return repositories, nil

}

func getToken() (string, error) {
	// TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${UNAME}'", "password": "'${UPASS}'"}' https://hub.docker.com/v2/users/login/ | jq -r .token)
	restClient := resty.New()

	apiURI := "v2/users/login"
	apiHost := "https://hub.docker.com"
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&BasicAuthInfo{Username: dockerUser, Password: dockerPass}).
		SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}).    // or SetError(AuthError{}).
		Post(fmt.Sprintf("%s/%s", apiHost, apiURI))
	if resperr != nil {
		logrus.WithError(resperr).Error(fmt.Sprintf("Error logging in and getting token: %s", apiURI))
		return "BAD", resperr
	}
	token := resp.Result().(*AuthSuccess)

	return token.Token, nil
}

func init() {
	getCmd.AddCommand(getRepositoriesCmd)

	getRepositoriesCmd.Flags().String("org", "", "Specify the dockerhub organization")
	getRepositoriesCmd.Flags().BoolP("add-updated", "u", false, "Add last updater (long run)")
	getRepositoriesCmd.Flags().BoolP("add-release", "r", false, "Add latest release tag info (long run)")
	getRepositoriesCmd.MarkFlagRequired("org")
}
