// Package cmd - Cobra Main command
package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"splice-cloud-util/vault"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	semVer    string
	gitCommit string
	buildDate string
	gitRef    string
)

var cfgFile string
var outputFormat string
var verbose bool
var noHeaders bool
var jenkinsUser string
var jenkinsURL string
var jenkinsKey string
var dockerUser string
var dockerPass string
var formatOverridden bool

var vaultClient vault.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "splice-cloud-util",
	Args:  cobra.MinimumNArgs(1),
	Short: `Run various commands against cloud providers`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if dockerUser != "" {
			viper.Set("docker_user", dockerUser)
			verr := viper.WriteConfig()
			if verr != nil {
				logrus.WithError(verr).Info("Failed to write config")
			}
		} else {
			du := viper.Get("docker_user")
			if du != nil {
				dockerUser = du.(string)
			}
		}
		if dockerPass != "" {
			viper.Set("docker_pass", dockerPass)
			verr := viper.WriteConfig()
			if verr != nil {
				logrus.WithError(verr).Info("Failed to write config")
			}
		} else {
			dp := viper.Get("docker_pass")
			if dp != nil {
				dockerPass = dp.(string)
			}
		}
		// Validate global parameters here, BEFORE we start to waste time
		// and run any code.
		if outputFormat != "" {
			outputFormat = strings.ToLower(outputFormat)
			switch outputFormat {
			case "json":
			case "gron":
			case "yaml":
			case "text":
			case "table":
			default:
				fmt.Println("Valid options for -o are [json|gron|[text|table]|yaml]")
				os.Exit(1)
			}
			formatOverridden = true
		} else {
			formatOverridden = false
			outputFormat = "json"
		}

		vaultClient = vault.NewVault()

	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func validateJenkins() error {

	jenkinsUser = fmt.Sprintf("%s", viper.Get("jenkins_user"))
	jenkinsURL = fmt.Sprintf("%s", viper.Get("jenkins_url"))
	jenkinsKey = fmt.Sprintf("%s", viper.Get("jenkins_key"))

	if err := verifyJenkinsAccess(jenkinsUser, jenkinsKey, jenkinsURL); err != nil {
		return err
	}
	return nil
}

func verifyJenkinsAccess(user string, key string, apiServer string) error {
	restClient := resty.New()

	depth := 2
	uri := fmt.Sprintf("computer/api/json?depth=%d", depth)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBasicAuth(user, key).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		return resperr
	}

	if resp.RawResponse.StatusCode != 200 {
		return errors.New("Failed to access Jenkins")
	}
	return nil

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.splice-cloud-util/config.yml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Set Verbose Output")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "output types: json, text (default), yaml")
	rootCmd.PersistentFlags().BoolVar(&noHeaders, "no-headers", false, "Suppress header output in Text output")
	rootCmd.PersistentFlags().StringVar(&dockerUser, "user", "", "Your docker user name, stored in the config file.")
	rootCmd.PersistentFlags().StringVar(&dockerPass, "pass", "", "Your docker password, stored in the config file.")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		if _, err := os.Stat(cfgFile); err != nil {
			if os.IsNotExist(err) {
				createRestrictedConfigFile(cfgFile)
			}
		}
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		directory := fmt.Sprintf("%s/%s", home, ".splice-cloud-util")
		if _, err := os.Stat(directory); err != nil {
			if os.IsNotExist(err) {
				os.Mkdir(directory, os.ModePerm)
			}
		}
		if stat, err := os.Stat(directory); err == nil && stat.IsDir() {
			configFile := fmt.Sprintf("%s/%s", home, ".splice-cloud-util/config.yml")
			createRestrictedConfigFile(configFile)
			viper.SetConfigFile(configFile)
		} else {
			logrus.Info("The ~/.splice-cloud-util path is a file and not a directory, please remove the .splice-cloud-util file.")
			os.Exit(1)
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// couldn't read the config file.
	}
}

func createRestrictedConfigFile(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			file, ferr := os.Create(fileName)
			if ferr != nil {
				logrus.Info("Unable to create the configfile.")
				os.Exit(1)
			}
			if runtime.GOOS != "windows" {
				mode := int(0600)
				if cherr := file.Chmod(os.FileMode(mode)); cherr != nil {
					logrus.Info("Chmod for config file failed, please set the mode to 0600.")
				}
			}
		}
	}
}
