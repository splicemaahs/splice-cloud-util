// Package provider - Azure Cloud Actions
package provider

import (
	"context"
	"os/user"

	"splice-cloud-util/cmd/objects"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-02-01/containerservice"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// Az - Structure to hold stuff
type Az struct {
	Provider string
	Auth     JenkinsAuth
}

type azCredentialInfo struct {
	TenantID       string `ini:"tenant,omitempty"`
	SubscriptionID string `ini:"subscription_id,omitempty"`
	ClientID       string `ini:"client_id,omitempty"`
	Secret         string `ini:"secret,omitempty"`
}

// GetEnvironments - Azure Environments
func (p *Az) GetEnvironments(env string) (objects.DetailList, error) {

	azCredsFile := getAzCredentialsFileObject("")
	credInfo, _ := getAzCredentialInfo(azCredsFile, "default")

	// Create an AKS client
	aksClient := containerservice.NewManagedClustersClient(credInfo.SubscriptionID)

	a, aerr := auth.NewClientCredentialsConfig(credInfo.ClientID, credInfo.Secret, credInfo.TenantID).Authorizer()
	if aerr != nil {
		logrus.Fatal(aerr)
	}

	aksClient.Authorizer = a

	// Setup context
	ctx := context.Background()
	clusters, cerr := aksClient.List(ctx)
	if cerr != nil {
		logrus.Fatal(cerr)
	}
	buildStatus, _ := getJenkinsJobDetails("azure", p.Auth)
	detail := []objects.Detail{}
	var clusterList []objects.ClusterList
	for _, cluster := range clusters.Values() {
		creator := ""
		dateCreated := ""
		for i, s := range cluster.Tags {
			if i == "Creator" {
				creator = *s
			}
			if i == "Date_Created" {
				dateCreated = *s
			}
		}
		status := ""
		if val, ok := buildStatus[*cluster.Name]; ok {
			status = val.Status
		}
		clusterList = append(clusterList, objects.ClusterList{
			ClusterName: *cluster.Name,
			CreatedBy:   creator,
			DateCreated: dateCreated,
			State:       status,
		})
	}
	nd := objects.Detail{
		EnvironmentName: "Azure",
		ClusterList:     clusterList,
	}
	detail = append(detail, nd)

	// // create a dns client
	// dnsClient := dns.NewRecordSetsClient(subscriptionID)
	// dnsClient.Authorizer = autorest.NewBearerAuthorizer(spt)

	// // create a groups client
	// groupsClient = resources.NewGroupsClient(subscriptionID)
	// groupsClient.Authorizer = autorest.NewBearerAuthorizer(spt)

	// // create a resources client
	// resourcesClient = resources.NewClient(subscriptionID)
	// resourcesClient.Authorizer = autorest.NewBearerAuthorizer(spt)
	// resourcesClient.RequestInspector = withAPIVersion("2018-05-01")

	// // Locate resource group
	// for list, _ := groupsClient.ListComplete(ctx, "", nil); list.NotDone(); err = list.Next() {
	// 	groupName := *list.Value().Name
	// 	//reqLogger.Info(fmt.Sprintf("Group name: %s", groupName))
	// 	resourcesList, _ := resourcesClient.Get(ctx, groupName, "Microsoft.Network", "", "dnszones", certName)
	// 	if err != nil {
	// 		log.Error(err, "Listing resources in resource group failed: ")
	// 	}
	// 	if resourcesList.Name != nil {
	// 		resourceName := *resourcesList.Name
	// 		//reqLogger.Info(fmt.Sprintf("Resource name: %s", resourceName))
	// 		if resourceName != "" && resourceName == certName {
	// 			//reqLogger.Info(fmt.Sprintf("Resource: %s", resourceName))
	// 			resourceGroup = groupName
	// 			//reqLogger.Info(fmt.Sprintf("Group name: %s", resourceGroup))
	// 			break
	// 		}
	// 	}
	// }
	envList := objects.DetailList{
		EnvironmentList: detail,
	}
	return envList, nil
}

func getAzDefaultCredentialsFilePath(pathOverride string) string {
	usr, err := user.Current()
	errCheck(err)
	filePath := usr.HomeDir + "/.azure/credentials"
	if len(pathOverride) > 0 {
		// TODO: stat the file, see if it actually exists.
		// TODO: Make a --credFilePath parameters, make it persistent, store in cobra config
		filePath = pathOverride
	}

	return filePath
}

func getAzCredentialsFileObject(pathOverride string) *ini.File {

	filePath := getAzDefaultCredentialsFilePath(pathOverride)
	iniFile, err := ini.Load(filePath)
	errCheck(err)

	return iniFile
}

func getAzCredentialInfo(iniFile *ini.File, sectionName string) (azCredentialInfo, error) {

	sect, err := iniFile.GetSection(sectionName)
	errCheck(err)
	info := &azCredentialInfo{}
	err = sect.MapTo(info)
	errCheck(err)

	return *info, nil

}
