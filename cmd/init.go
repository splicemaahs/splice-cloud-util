package cmd

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Run the initialization process, set Jenkins/DockerHub Credentials",
	Long: `This will prompt for your Jenkins URL, UserId, and KEY along with
	your DockerHub UserName and Password, these details will be stored in the
	~/.splice-cloud-util/config.yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		gatherAndStoreJenkinsCredentials()
	},
}

func gatherAndStoreJenkinsCredentials() {
	var qs = []*survey.Question{
		{
			Name:     "user",
			Prompt:   &survey.Input{Message: "What is your Jenkins user account?"},
			Validate: survey.Required,
		},
		{
			Name:     "jenkins",
			Prompt:   &survey.Input{Message: "What is your Jenkins API URL?"},
			Validate: survey.Required,
		},
		{
			Name:     "key",
			Prompt:   &survey.Password{Message: "Your Jenkins API KEY?"},
			Validate: survey.Required,
		},
		{
			Name:     "dockeruser",
			Prompt:   &survey.Input{Message: "What is your DockerHub User ID?"},
			Validate: survey.Required,
		},
		{
			Name:     "dockerpass",
			Prompt:   &survey.Password{Message: "Your DockerHub Password?"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		User       string `survey:"user"`
		JenkinsURL string `survey:"jenkins"`
		Key        string `survey:"key"`
		DockerUser string `survey:"dockeruser"`
		DockerPass string `survey:"dockerpass"`
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		logrus.Fatal("Please provide all of the configuration information.")
	}

	URL := answers.JenkinsURL
	if !strings.HasSuffix(answers.JenkinsURL, "/") {
		URL = fmt.Sprintf("%s/", answers.JenkinsURL)
	}
	if err := verifyJenkinsAccess(answers.User, answers.Key, URL); err != nil {
		logrus.Fatal("Could not validate Jenkins access")
	}

	dockerUser = answers.DockerUser
	dockerPass = answers.DockerPass
	if _, terr := getToken(); terr != nil {
		logrus.Fatal("Could not validate Docker access")
	}

	viper.Set("jenkins_user", answers.User)
	viper.Set("jenkins_url", URL)
	viper.Set("jenkins_key", answers.Key)
	verr := viper.WriteConfig()
	if verr != nil {
		logrus.Fatal("Jenkins and/or Docker information was valid, failed to store in config")
	}

	logrus.Info("Jenkins and Docker information validated and stored.")
}

func init() {
	rootCmd.AddCommand(initCmd)
}
