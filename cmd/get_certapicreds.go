// Package cmd - Get Environments Code
package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"splice-cloud-util/cmd/objects"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CertAPIRecord - ClientID and Key
type CertAPIRecord struct {
	ClientID  string `json:"CERTMGR_CLIENT_ID"`
	SecretKey string `json:"CERTMGR_SECRET_KEY"`
}

// getCmdCertAPICreds - List Active Environments
var getCmdCertAPICreds = &cobra.Command{
	Use:   "certapicreds",
	Short: "Get an output of the Certificate API Server Credentials (Build Environment)",
	Long: `EXAMPLE 1
	#> splice-cloud-util get certapicreds

`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Check our Jenkins Stored Credentials
		if err := validateJenkins(); err != nil {
			logrus.Fatal("Unable to access Jenkins, please run 'splice-cloud-util init' to configure Jenkins access")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		err := getCertAPICreds(outputFormat)
		if err != nil {
			logrus.WithError(err).Error("Error getting the credentials")
		}

	},
}

func getCertAPICreds(format string) error {

	// secret/deployments/k8s/default/infrastructure/
	// List the FOLDERs from above, those are the environments
	// Loop thorugh them, and read these
	// /deployments/k8s/default/infrastructure/{environment}/default
	// CERTMGR_CLIENT_ID=$(vault kv get --format table --field certmgr_client_id secret/deployments/k8s/default/infrastructure/nonprod/default); echo ${CERTMGR_CLIENT_ID}
	// CERTMGR_SECRET_KEY=$(vault kv get --format table --field certmgr_secret_key secret/deployments/k8s/default/infrastructure/nonprod/default); echo ${CERTMGR_SECRET_KEY}

	accountPaths, perr := vaultClient.GetPaths("deployments/k8s/default/infrastructure/")
	if perr != nil {
		logrus.WithError(perr).Error("Could not get paths (accounts) from vault")
	}

	credList := objects.CertAPICredList{}
	for k, v := range accountPaths {
		accountData, derr := vaultClient.GetData(fmt.Sprintf("%sdefault", k))
		if derr != nil {
			logrus.WithError(derr).Error("Could not get data (account) from vault")
		}

		jsonData, jerr := json.MarshalIndent(accountData.Data.Data, "", "  ")
		if jerr != nil {
			logrus.WithError(jerr).Error("Error extracting json")
		}
		certRecord := CertAPIRecord{}
		marshErr := json.Unmarshal([]byte(jsonData), &certRecord)
		if marshErr != nil {
			logrus.Fatal("Could not unmarshall data", marshErr)
		}

		credList.List = append(credList.List, objects.CertAPICred{
			Account:   v.Path,
			ClientID:  certRecord.ClientID,
			SecretKey: certRecord.SecretKey,
		})
	}

	switch strings.ToLower(format) {
	case "json":
		credList.ToJSON()
	case "yaml":
		credList.ToYAML()
	case "text", "table":
		credList.ToTEXT(noHeaders)
	}

	return nil

}

func init() {

	getCmd.AddCommand(getCmdCertAPICreds)

}
