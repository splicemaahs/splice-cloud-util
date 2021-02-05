// Package provider - GCP Cloud Actions
package provider

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"splice-cloud-util/cmd/objects"

	// This is the way
	_ "github.com/mattn/go-sqlite3"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	container "google.golang.org/api/container/v1"
	"google.golang.org/api/option"

	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

// Gcp - Structure to hold stuff
type Gcp struct {
	Provider string
	Auth     JenkinsAuth
}

// GetEnvironments - GCP Environments
func (p *Gcp) GetEnvironments(env string) (objects.DetailList, error) {

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	credDb := fmt.Sprintf("%s/%s", home, ".config/gcloud/credentials.db")
	if _, err := os.Stat(credDb); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	database, _ := sql.Open("sqlite3", credDb)
	rows, _ := database.Query("SELECT * FROM credentials")
	var account string
	var authJSON string
	for rows.Next() {
		rows.Scan(&account, &authJSON)
		break
	}
	database.Close()

	detail := []objects.Detail{}

	ctx := context.TODO()
	containerService, cserr := container.NewService(ctx, option.WithCredentialsJSON([]byte(authJSON)))
	if cserr != nil {
		logrus.Fatal(cserr)
	}
	crmService, csmerr := cloudresourcemanager.NewService(ctx, option.WithCredentialsJSON([]byte(authJSON)))
	if csmerr != nil {
		logrus.Fatal(cserr)
	}

	buildStatus, _ := getJenkinsJobDetails("gcp", p.Auth)
	reportList := make(map[string]string)

	projects, perr := crmService.Projects.List().Do()
	if perr != nil {
		logrus.Fatal(perr)
	}
	for _, prj := range projects.Projects {
		if strings.HasSuffix(prj.Name, "-gke") {
			if env == "" || env == prj.ProjectId {
				account := strings.Split(prj.ProjectId, "-")[0]
				// get all clusters list
				clusters, err := container.NewProjectsLocationsClustersService(containerService).List("projects/" + prj.ProjectId + "/locations/-").Do()
				if err != nil {
					logrus.Warn(fmt.Sprintf("Unable to get clusters from project %s", prj.Name))
				} else {
					var clusterList []objects.ClusterList
					for _, cluster := range clusters.Clusters {
						status := EnvironmentInfo{}
						if val, ok := buildStatus[cluster.Name]; ok {
							status = val
						}
						clusterList = append(clusterList, objects.ClusterList{
							ClusterName: cluster.Name,
							CreatedBy:   cluster.ResourceLabels["creator"],
							DateCreated: cluster.ResourceLabels["date_created"],
							State:       status.Status,
						})
						reportList[cluster.Name] = cluster.Name
					}
					for be, bs := range buildStatus {
						if strings.HasPrefix(be, account) {
							if _, ok := reportList[be]; !ok {
								clusterList = append(clusterList, objects.ClusterList{
									ClusterName: be,
									CreatedBy:   bs.CreatedBy,
									State:       bs.Status,
								})
								reportList[be] = be
							}
						}
					}
					nd := objects.Detail{
						EnvironmentName: prj.ProjectId,
						ClusterList:     clusterList,
					}
					detail = append(detail, nd)
				}
			}
		}
	}

	envList := objects.DetailList{
		EnvironmentList: detail,
	}
	return envList, nil

}
