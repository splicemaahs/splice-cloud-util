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
	Short: "Run the initialization process, set Jenkins Credentials",
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
	}

	answers := struct {
		User       string
		JenkinsURL string `survey:"jenkins"`
		Key        string
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
	viper.Set("jenkins_user", answers.User)
	viper.Set("jenkins_url", URL)
	viper.Set("jenkins_key", answers.Key)
	verr := viper.WriteConfig()
	if verr != nil {
		logrus.Fatal("Jenkins information was validate, failed to store in config")
	}

	logrus.Info("Jenkins information validated and stored.")
}

func init() {
	rootCmd.AddCommand(initCmd)
}
