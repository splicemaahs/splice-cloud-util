// Package provider - AWS Cloud Actions
package provider

import (
	"fmt"
	"os/user"
	"time"

	"splice-cloud-util/cmd/objects"

	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type credentialInfo struct {
	AssumedRole        string `ini:"assumed_role,omitempty"`
	AssumedRoleARN     string `ini:"assumed_role_arn,omitempty"`
	AwsAccessKeyID     string `ini:"aws_access_key_id,omitempty"`
	AwsMFADevice       string `ini:"aws_mfa_device,omitempty"`
	AwsSecretAccessKey string `ini:"aws_secret_access_key,omitempty"`
	AwsSecurityToken   string `ini:"aws_security_token,omitempty"`
	AwsSessionToken    string `ini:"aws_session_token,omitempty"`
	Expiration         string `ini:"expiration,omitempty"`
	Region             string `ini:"region,omitempty"`
}

// Aws - Structure to hold stuff
type Aws struct {
	Provider string
	Auth     JenkinsAuth
}

// GetEnvironments - AWS Environments
func (p *Aws) GetEnvironments(env string) (objects.DetailList, error) {

	detail := []objects.Detail{}

	fetchList := []string{
		"build",
		"cs",
		"nonprod",
		"pd",
		"prod",
		"qa",
		"sales",
	}

	var diff time.Duration
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
	now := time.Now().In(loc)
	awsCredsFile := getAwsCredentialsFileObject("")
	for _, profileSuffix := range fetchList {
		awsProfile := fmt.Sprintf("splice-%s", profileSuffix)
		credInfo, _ := getCredentialInfo(awsCredsFile, awsProfile)
		// go based aws-mfa - 2020-08-27T12:13:39Z
		// python based aws-mfa - 2020-05-12 03:50:53
		expire, err := dateparse.ParseLocal(credInfo.Expiration)
		if err != nil {
			logrus.Warn("Unable to parse the time for profile", awsProfile)
			diff = 0
		} else {
			expires := (expire).In(loc)
			diff = expires.Sub(now)
		}
		if diff > 0 {
			// Perform Lookup of K8s environments
			// Setup aws session
			mySession := session.Must(session.NewSession(&aws.Config{
				Region:      aws.String(credInfo.Region),
				Credentials: credentials.NewStaticCredentials(credInfo.AwsAccessKeyID, credInfo.AwsSecretAccessKey, credInfo.AwsSessionToken),
			}))
			svc := eks.New(mySession)
			clusters, err := svc.ListClusters(&eks.ListClustersInput{})
			if err != nil {
				logrus.Fatal(err, "Error listing clusters")
			}

			buildStatus, _ := getJenkinsJobDetails("aws", p.Auth)

			var clusterList []objects.ClusterList
			for _, cluster := range clusters.Clusters {
				describeInput := &eks.DescribeClusterInput{
					Name: cluster,
				}

				c, err := svc.DescribeCluster(describeInput)
				if err != nil {
					logrus.Fatal("DescribeClusters error:", err)
				}

				creator := ""
				dateCreated := ""
				for i, s := range c.Cluster.Tags {
					if i == "Creator" {
						creator = *s
					}
					if i == "Date_Created" {
						dateCreated = *s
					}
				}
				status := ""
				if val, ok := buildStatus[*cluster]; ok {
					status = val.Status
				}
				clusterList = append(clusterList, objects.ClusterList{
					ClusterName: *cluster,
					CreatedBy:   creator,
					DateCreated: dateCreated,
					State:       status,
				})
			}
			nd := objects.Detail{
				EnvironmentName: awsProfile,
				ClusterList:     clusterList,
			}
			detail = append(detail, nd)
		} else {
			var clusterList []objects.ClusterList
			clusterList = append(clusterList, objects.ClusterList{
				ClusterName: "MFA EXPIRED",
				CreatedBy:   "",
			})
			nd := objects.Detail{
				EnvironmentName: awsProfile,
				ClusterList:     clusterList,
			}
			detail = append(detail, nd)
		}
	}

	envList := objects.DetailList{
		EnvironmentList: detail,
	}
	return envList, nil

}

func getDefaultCredentialsFilePath(pathOverride string) string {
	usr, err := user.Current()
	errCheck(err)
	filePath := usr.HomeDir + "/.aws/credentials"
	if len(pathOverride) > 0 {
		// TODO: stat the file, see if it actually exists.
		// TODO: Make a --credFilePath parameters, make it persistent, store in cobra config
		filePath = pathOverride
	}

	return filePath
}

func getAwsCredentialsFileObject(pathOverride string) *ini.File {

	filePath := getDefaultCredentialsFilePath(pathOverride)
	iniFile, err := ini.Load(filePath)
	errCheck(err)

	return iniFile
}

func getCredentialInfo(iniFile *ini.File, sectionName string) (credentialInfo, error) {

	sect, err := iniFile.GetSection(sectionName)
	errCheck(err)
	info := &credentialInfo{}
	err = sect.MapTo(info)
	errCheck(err)

	return *info, nil

}
