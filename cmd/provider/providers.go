// Package provider - Cloud Specific Actions
package provider

import (
	"encoding/json"
	"fmt"

	"splice-cloud-util/cmd/objects"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/sirupsen/logrus"
)

// EnvironmentInfo - Store and pass our environment info around
type EnvironmentInfo struct {
	Account     string `json:"account"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
	DateCreated string `json:"datecreated"`
	CreatedBy   string `json:"createdby"`
}

// ExecutorStatus - Return from Jenkins
type ExecutorStatus struct {
	Class         string `json:"_class"`
	BusyExecutors int    `json:"busyExecutors"`
	Computer      []struct {
		Class   string `json:"_class"`
		Actions []struct {
		} `json:"actions"`
		AssignedLabels []struct {
			Actions        []interface{} `json:"actions"`
			BusyExecutors  int           `json:"busyExecutors"`
			Clouds         []interface{} `json:"clouds"`
			Description    interface{}   `json:"description"`
			IdleExecutors  int           `json:"idleExecutors"`
			LoadStatistics struct {
				Class              string `json:"_class"`
				AvailableExecutors struct {
				} `json:"availableExecutors"`
				BusyExecutors struct {
				} `json:"busyExecutors"`
				ConnectingExecutors struct {
				} `json:"connectingExecutors"`
				DefinedExecutors struct {
				} `json:"definedExecutors"`
				IdleExecutors struct {
				} `json:"idleExecutors"`
				OnlineExecutors struct {
				} `json:"onlineExecutors"`
				QueueLength struct {
				} `json:"queueLength"`
				TotalExecutors struct {
				} `json:"totalExecutors"`
			} `json:"loadStatistics"`
			Name  string `json:"name"`
			Nodes []struct {
				Class          string `json:"_class"`
				AssignedLabels []struct {
					Name string `json:"name"`
				} `json:"assignedLabels"`
				Mode            string      `json:"mode"`
				NodeDescription string      `json:"nodeDescription"`
				NodeName        string      `json:"nodeName"`
				NumExecutors    int         `json:"numExecutors"`
				Description     interface{} `json:"description"`
				Jobs            []struct {
					Class string `json:"_class"`
					Name  string `json:"name"`
					URL   string `json:"url"`
					Color string `json:"color,omitempty"`
				} `json:"jobs"`
				OverallLoad struct {
				} `json:"overallLoad"`
				PrimaryView struct {
					Class string `json:"_class"`
					Name  string `json:"name"`
					URL   string `json:"url"`
				} `json:"primaryView"`
				QuietingDown   bool `json:"quietingDown"`
				SlaveAgentPort int  `json:"slaveAgentPort"`
				UnlabeledLoad  struct {
					Class string `json:"_class"`
				} `json:"unlabeledLoad"`
				URL         string `json:"url"`
				UseCrumbs   bool   `json:"useCrumbs"`
				UseSecurity bool   `json:"useSecurity"`
				Views       []struct {
					Class string `json:"_class"`
					Name  string `json:"name"`
					URL   string `json:"url"`
				} `json:"views"`
			} `json:"nodes"`
			Offline        bool          `json:"offline"`
			TiedJobs       []interface{} `json:"tiedJobs"`
			TotalExecutors int           `json:"totalExecutors"`
			PropertiesList []interface{} `json:"propertiesList"`
		} `json:"assignedLabels"`
		Description string `json:"description"`
		DisplayName string `json:"displayName"`
		Executors   []struct {
			CurrentExecutable interface{} `json:"currentExecutable"`
			Idle              bool        `json:"idle"`
			LikelyStuck       bool        `json:"likelyStuck"`
			Number            int         `json:"number"`
			Progress          int         `json:"progress"`
		} `json:"executors"`
		Icon            string `json:"icon"`
		IconClassName   string `json:"iconClassName"`
		Idle            bool   `json:"idle"`
		JnlpAgent       bool   `json:"jnlpAgent"`
		LaunchSupported bool   `json:"launchSupported"`
		LoadStatistics  struct {
			Class              string `json:"_class"`
			AvailableExecutors struct {
				Hour struct {
				} `json:"hour"`
				Min struct {
				} `json:"min"`
				Sec10 struct {
				} `json:"sec10"`
			} `json:"availableExecutors"`
			BusyExecutors struct {
				Hour struct {
				} `json:"hour"`
				Min struct {
				} `json:"min"`
				Sec10 struct {
				} `json:"sec10"`
			} `json:"busyExecutors"`
			ConnectingExecutors struct {
				Hour struct {
				} `json:"hour"`
				Min struct {
				} `json:"min"`
				Sec10 struct {
				} `json:"sec10"`
			} `json:"connectingExecutors"`
			DefinedExecutors struct {
				Hour struct {
				} `json:"hour"`
				Min struct {
				} `json:"min"`
				Sec10 struct {
				} `json:"sec10"`
			} `json:"definedExecutors"`
			IdleExecutors struct {
				Hour struct {
				} `json:"hour"`
				Min struct {
				} `json:"min"`
				Sec10 struct {
				} `json:"sec10"`
			} `json:"idleExecutors"`
			OnlineExecutors struct {
				Hour struct {
				} `json:"hour"`
				Min struct {
				} `json:"min"`
				Sec10 struct {
				} `json:"sec10"`
			} `json:"onlineExecutors"`
			QueueLength struct {
				Hour struct {
				} `json:"hour"`
				Min struct {
				} `json:"min"`
				Sec10 struct {
				} `json:"sec10"`
			} `json:"queueLength"`
			TotalExecutors struct {
				Hour struct {
				} `json:"hour"`
				Min struct {
				} `json:"min"`
				Sec10 struct {
				} `json:"sec10"`
			} `json:"totalExecutors"`
		} `json:"loadStatistics"`
		ManualLaunchAllowed bool `json:"manualLaunchAllowed"`
		MonitorData         struct {
			HudsonNodeMonitorsSwapSpaceMonitor struct {
				Class                   string `json:"_class"`
				AvailablePhysicalMemory int64  `json:"availablePhysicalMemory"`
				AvailableSwapSpace      int    `json:"availableSwapSpace"`
				TotalPhysicalMemory     int64  `json:"totalPhysicalMemory"`
				TotalSwapSpace          int    `json:"totalSwapSpace"`
			} `json:"hudson.node_monitors.SwapSpaceMonitor"`
			HudsonNodeMonitorsTemporarySpaceMonitor struct {
				Class     string `json:"_class"`
				Timestamp int64  `json:"timestamp"`
				Path      string `json:"path"`
				Size      int64  `json:"size"`
			} `json:"hudson.node_monitors.TemporarySpaceMonitor"`
			HudsonNodeMonitorsDiskSpaceMonitor struct {
				Class     string `json:"_class"`
				Timestamp int64  `json:"timestamp"`
				Path      string `json:"path"`
				Size      int64  `json:"size"`
			} `json:"hudson.node_monitors.DiskSpaceMonitor"`
			ComSynopsysArcJenkinsPluginsOwnershipNodesOwnershipNodeMonitor struct {
				Class     string `json:"_class"`
				Timestamp int64  `json:"timestamp"`
			} `json:"com.synopsys.arc.jenkins.plugins.ownership.nodes.OwnershipNodeMonitor"`
			HudsonNodeMonitorsArchitectureMonitor string `json:"hudson.node_monitors.ArchitectureMonitor"`
			HudsonNodeMonitorsResponseTimeMonitor struct {
				Class     string `json:"_class"`
				Timestamp int64  `json:"timestamp"`
				Average   int    `json:"average"`
			} `json:"hudson.node_monitors.ResponseTimeMonitor"`
			HudsonNodeMonitorsClockMonitor struct {
				Class string `json:"_class"`
				Diff  int    `json:"diff"`
			} `json:"hudson.node_monitors.ClockMonitor"`
		} `json:"monitorData"`
		NumExecutors       int         `json:"numExecutors"`
		Offline            bool        `json:"offline"`
		OfflineCause       interface{} `json:"offlineCause"`
		OfflineCauseReason string      `json:"offlineCauseReason"`
		OneOffExecutors    []struct {
			CurrentExecutable struct {
				Class   string `json:"_class"`
				Actions []struct {
					Class      string `json:"_class,omitempty"`
					Parameters []struct {
						Class string `json:"_class"`
						Name  string `json:"name"`
						Value string `json:"value,omitempty"`
					} `json:"parameters,omitempty"`
					Causes []struct {
						Class            string `json:"_class"`
						ShortDescription string `json:"shortDescription"`
						UserID           string `json:"userId"`
						UserName         string `json:"userName"`
					} `json:"causes,omitempty"`
					BuildsByBranchName struct {
						Master struct {
							Class       string      `json:"_class"`
							BuildNumber int         `json:"buildNumber"`
							BuildResult interface{} `json:"buildResult"`
							Marked      struct {
								SHA1   string `json:"SHA1"`
								Branch []struct {
									SHA1 string `json:"SHA1"`
									Name string `json:"name"`
								} `json:"branch"`
							} `json:"marked"`
							Revision struct {
								SHA1   string `json:"SHA1"`
								Branch []struct {
									SHA1 string `json:"SHA1"`
									Name string `json:"name"`
								} `json:"branch"`
							} `json:"revision"`
						} `json:"master"`
					} `json:"buildsByBranchName,omitempty"`
					LastBuiltRevision struct {
						SHA1   string `json:"SHA1"`
						Branch []struct {
							SHA1 string `json:"SHA1"`
							Name string `json:"name"`
						} `json:"branch"`
					} `json:"lastBuiltRevision,omitempty"`
					RemoteUrls []string `json:"remoteUrls,omitempty"`
					ScmName    string   `json:"scmName,omitempty"`
				} `json:"actions"`
				Artifacts         []interface{} `json:"artifacts"`
				Building          bool          `json:"building"`
				Description       interface{}   `json:"description"`
				DisplayName       string        `json:"displayName"`
				Duration          int           `json:"duration"`
				EstimatedDuration int           `json:"estimatedDuration"`
				Executor          struct {
					Class string `json:"_class"`
				} `json:"executor"`
				FullDisplayName string        `json:"fullDisplayName"`
				ID              string        `json:"id"`
				KeepLog         bool          `json:"keepLog"`
				Number          int           `json:"number"`
				QueueID         int           `json:"queueId"`
				Result          interface{}   `json:"result"`
				Timestamp       int64         `json:"timestamp"`
				URL             string        `json:"url"`
				ChangeSets      []interface{} `json:"changeSets"`
				Culprits        []interface{} `json:"culprits"`
				NextBuild       interface{}   `json:"nextBuild"`
				PreviousBuild   struct {
					Number int    `json:"number"`
					URL    string `json:"url"`
				} `json:"previousBuild"`
			} `json:"currentExecutable"`
			Idle        bool `json:"idle"`
			LikelyStuck bool `json:"likelyStuck"`
			Number      int  `json:"number"`
			Progress    int  `json:"progress"`
		} `json:"oneOffExecutors"`
		TemporarilyOffline bool `json:"temporarilyOffline"`
	} `json:"computer"`
	DisplayName    string `json:"displayName"`
	TotalExecutors int    `json:"totalExecutors"`
}

// JenkinsAuth - Structure to hold Jenkins Credential Info
type JenkinsAuth struct {
	User       string
	JenkinsURL string
	Key        string
}

// Provider = The main interface used to describe appliances
type Provider interface {
	GetEnvironments(env string) (objects.DetailList, error)
}

//Our Cloud Provider Types
const (
	AWS   = "aws"
	AZURE = "az"
	GCP   = "gcp"
)

// GetProvider - Function to create the appliances
func GetProvider(t string, user string, key string, url string) (Provider, error) {
	//Use a switch case to switch between types, if a type exist then error is nil (null)
	switch t {
	case AWS:
		return &Aws{
			Provider: AWS,
			Auth: JenkinsAuth{
				User:       user,
				Key:        key,
				JenkinsURL: url,
			},
		}, nil
	case AZURE:
		return &Az{
			Provider: AWS,
			Auth: JenkinsAuth{
				User:       user,
				Key:        key,
				JenkinsURL: url,
			},
		}, nil
	case GCP:
		return &Gcp{
			Provider: AWS,
			Auth: JenkinsAuth{
				User:       user,
				Key:        key,
				JenkinsURL: url,
			},
		}, nil
	default:
		logrus.Info("Defaulting CSP to GCP")
		return &Gcp{
			Provider: AWS,
			Auth: JenkinsAuth{
				User:       user,
				Key:        key,
				JenkinsURL: url,
			},
		}, nil
	}
}

func errCheck(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func getJenkinsJobDetails(csp string, auth JenkinsAuth) (map[string]EnvironmentInfo, error) {

	status := ""
	restClient := resty.New()
	buildList := make(map[string]EnvironmentInfo)

	apiServer := auth.JenkinsURL
	depth := 2

	uri := fmt.Sprintf("computer/api/json?depth=%d", depth)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBasicAuth(auth.User, auth.Key).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		return map[string]EnvironmentInfo{}, nil
	}

	var response ExecutorStatus

	marsherr := json.Unmarshal(resp.Body(), &response)
	if marsherr != nil {
		return map[string]EnvironmentInfo{}, nil
	}

	for _, c := range response.Computer {
		for _, o := range c.OneOffExecutors {

			if strings.Contains(o.CurrentExecutable.URL, "KubernetesDeploy") {
				if strings.Contains(o.CurrentExecutable.URL, csp) {
					action := ""
					account := ""
					environment := ""
					createdby := ""
					for _, a := range o.CurrentExecutable.Actions {
						for _, p := range a.Parameters {
							if p.Name == "account" {
								account = p.Value
							}
							if p.Name == "action" {
								action = p.Value
								if action == "create" {
									status = "Jenkins Creating"
								} else if action == "destroy" {
									status = "Jenkins Destroying"
								}
							}
							if p.Name == "environment" {
								environment = p.Value
							}
							if p.Name == "createdby" {
								createdby = p.Value
							}
						}
						for _, c := range a.Causes {
							if createdby == "" {
								createdby = strings.Split(c.UserID, "@")[0]
							}
						}
					}
					if action != "" {
						buildList[fmt.Sprintf("%s-%s", account, environment)] = EnvironmentInfo{
							Account:     account,
							Environment: environment,
							Status:      status,
							CreatedBy:   createdby,
						}
					}
				}
			}
		}
	}
	return buildList, nil
}
