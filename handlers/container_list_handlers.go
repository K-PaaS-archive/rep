package handlers

import (
	"code.cloudfoundry.org/executor"
	"code.cloudfoundry.org/garden"
	GardenClient "code.cloudfoundry.org/garden/client"
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ContainerMetricsMetadata struct {
	Limits            Limits       `json:"limits,omitempty"`
	UsageMetrics      UsageMetrics `json:"usage_metrics,omitempty"`
	Container_Id      string       `json:"container_id,omitempty"`
	Interface_Id      string       `json:"interface_id,omitempty"`
	Application_Id    string       `json:"application_id,omitempty"`
	Application_Index string       `json:"application_index,omitempty"`
	Application_Name  string       `json:"application_name,omitempty"`
	Application_Urls  []string     `json:"application_uris,omitempty"`
}
type Applications struct {
	Limits              Limits   `json:"limits,omitempty"`
	Name                string   `json:"name,omitempty"`
	Application_Id      string   `json:"application_id,omitempty"`
	Application_Version string   `json:"application_version,omitempty"`
	Application_Name    string   `json:"application_name,omitempty"`
	Application_Urls    []string `json:"application_uris,omitempty"`
	Application_Index   int      `json:"application_index,omitempty"`
	Container_Port      uint32   `json:"container_port,omitempty"`
	Space_Name          string   `json:"space_name,omitempty"`
	Space_Id            string   `json:"space_id,omitempty"`
	Uris                []string `json:"uris,omitempty"`

	Container_Id string `json:"container_id,omitempty"`
}
type Limits struct {
	Fds    int32 `json:"fds,omitempty"`
	Memory int32 `json:"mem,omitempty"`
	Disk   int32 `json:"disk,omitempty"`
}
type UsageMetrics struct {
	MemoryUsageInBytes uint64        `json:"memory_usage_in_bytes"`
	DiskUsageInBytes   uint64        `json:"disk_usage_in_bytes"`
	TimeSpentInCPU     time.Duration `json:"time_spent_in_cpu"`
}
type ContainerInfo struct {
	Container_Id    string `json:"container_id,omitempty"`
	Application_Id  string `json:"application_id,omitempty"`
	Organization_Id string `json:"organization_id,omitempty"`
	Space_Id        string `json:"space_id,omitempty"`
	Container_Port  uint32 `json:"container_port,omitempty"`
	Interface_Id    string `json:"interface_id,omitempty"`
}
type ContainerListHandler struct {
	logger         lager.Logger
	executorClient executor.Client
	gardenClient   GardenClient.Client
}

func NewContainerListHandler(logger lager.Logger, executorClient executor.Client, gardenClient GardenClient.Client) *ContainerListHandler {
	return &ContainerListHandler{
		logger:         logger,
		executorClient: executorClient,
		gardenClient:   gardenClient,
	}
}
func (c ContainerListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, logger lager.Logger) {
	var applications []Applications
	var containerInfos []ContainerInfo
	var containermetrics []ContainerMetricsMetadata
	containers, err := c.executorClient.ListContainers(c.logger)
	//=============================== Container Metrics ==================================
	containerBulkMetrics, err := c.executorClient.GetBulkMetrics(c.logger)
	//fmt.Println("======================================================================================================================")
	/*for _, bulkMetrics := range containerBulkMetrics{
	       fmt.Println("##### container_list_handler.go = container Bulk Metrics :  key , value :", bulkMetrics)
	       fmt.Println("##### container_list_handler.go = container Bulk Metrics :  guid :", bulkMetrics.Guid)
	       fmt.Println("##### container_list_handler.go = container Bulk Metrics :  index :", bulkMetrics.Index)
	       fmt.Println("##### container_list_handler.go = container Bulk Metrics :  MemoryUsageInBytes :", bulkMetrics.MemoryUsageInBytes)
	       fmt.Println("##### container_list_handler.go = container Bulk Metrics :  DiskUsageInBytes :", bulkMetrics.DiskUsageInBytes)
	       fmt.Println("##### container_list_handler.go = container Bulk Metrics :  TimeSpentInCPU - seconds :", bulkMetrics.TimeSpentInCPU.Seconds())
	}*/
	//=====================================================================================
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.Error("failed-to-fetch-container", err)
		return
	}
	var app_host_port uint32
	var app_index int
	for i := range containers {
		var application Applications
		container := &containers[i]
		//fmt.Println("##### container_list_handler.go, container :", container.Guid, container)
		//fmt.Println("##### container_list_handler.go, container Env:", container.Env)
		//fmt.Println("##### container_list_handler.go, container Network:", container.Network)
		//fmt.Println("##### container_list_handler.go, container Action:", container.Action)
		//fmt.Println("##### container_list_handler.go, container Monitor:", container.Monitor) // no
		//fmt.Println("##### container_list_handler.go, container Setup:", container.Setup) // no
		//fmt.Println("##### container_list_handler.go, container AllocatedAt:", container.AllocatedAt) // no
		//fmt.Println("##### container_list_handler.go, container CachedDependencies:", container.CachedDependencies) // no
		//fmt.Println("##### container_list_handler.go, container EgressRules:", container.EgressRules) // no
		//fmt.Println("##### container_list_handler.go, container LogConfig:", container.LogConfig) // no
		//fmt.Println("##### container_list_handler.go, container RunResult:", container.RunResult) // no
		//fmt.Println("##### container_list_handler.go, container VolumeMounts:", container.VolumeMounts) // no
		//fmt.Println("##### container_list_handler.go, container Resource:", container.Resource) // no
		//fmt.Println("##### container_list_handler.go, container RunInfo:", container.RunInfo) // no
		//fmt.Println("##### container_list_handler.go, container Tags:", container.Tags) // no
		//fmt.Println("##### container_list_handler.go, container.RootFSPath :", container.RunInfo.RootFSPath) // no
		//fmt.Println("##### container_list_handler.go, container.Resource.DiskMB :",container.Resource.DiskMB) // no
		//fmt.Println("##### container_list_handler.go, container.Resource.MemoryMB :",container.Resource.MemoryMB) // no
		////fmt.Println("##### container_list_handler.go, container.RunInfo.DiskScope :",container.RunInfo.DiskScope) // no
		//fmt.Println("##### container_list_handler.go, container.RunInfo.ExternalIP :",container.ExternalIP) // no
		//fmt.Println("##### container_list_handler.go, container.RunInfo.MemoryMB :",container.MemoryMB) // no
		//fmt.Println("##### container_list_handler.go, container.RunInfo.State :",container.State) // no
		//fmt.Println("##### container_list_handler.go, container.RunInfo.Network :",container.Network.Properties) // no
		//fmt.Println("##### container_list_handler.go, container.RunInfo.ENV :",container.Env) // no
		//fmt.Println("##### container_list_handler.go, container Ports:", container.Ports) // no
		//fmt.Println("##### container_list_handler.go, container.RunInfo.Action:", container.RunInfo.Action) // no
		//fmt.Println("##### container_list_handler.go, container.RunInfo.ActionValue:", container.RunInfo.Action.GetValue()) //  no
		for _, value := range container.Ports {
			if value.ContainerPort != 2222 {
				app_host_port = uint32(value.HostPort)
			}
		}
		for _, value := range container.Env {
			if strings.Contains(value.Name, "INSTANCE_INDEX") {
				app_index, err = strconv.Atoi(value.Value)
			}
		}
		if container.RunInfo.Action.CodependentAction != nil {
			action := container.RunInfo.Action.CodependentAction.GetActions()[0].RunAction
			if action != nil {
				for _, envs := range action.Env {
					if envs.Name == "VCAP_APPLICATION" {
						/*fmt.Println("##### container_list_handler.go, CodependentAction.RunAction.Envs:", envs)*/
						json.Unmarshal([]byte(envs.Value), &application)
						//fmt.Println("##### container_list_handler.go, Application Info - Id :", application.Application_Id)
						//fmt.Println("##### container_list_handler.go, Application Info - Name :", application.Application_Name)
						//fmt.Println("##### container_list_handler.go, Application Info - index :", app_index)
						//fmt.Println("##### container_list_handler.go, Application Info - Limits :", application.Limits.Disk, application.Limits.Memory)
						//fmt.Println("##### container_list_handler.go, Application Info - Uris :", application.Uris)

						application.Container_Id = container.Guid

						application.Container_Port = app_host_port
						application.Application_Index = app_index
						applications = append(applications, application)
					}
				}
			}
		}
	}

	properties := garden.Properties{}
	gardenContainers, err := c.gardenClient.Containers(properties)
	var container_host_port uint32
	for _, gc := range gardenContainers {
		var containerInfo ContainerInfo
		gardenContainerInfo, _ := gc.Info()
		/*fmt.Println("### container_list_handler.go - container info : ", gardenContainerInfo)
		  fmt.Println("### container_list_handler.go - container info - container IP: ", gardenContainerInfo.ContainerIP)
		  fmt.Println("### container_list_handler.go - container info - contaienr Path: ", gardenContainerInfo.ContainerPath)
		  fmt.Println("### container_list_handler.go - container info - Host IP: ", gardenContainerInfo.HostIP)
		  fmt.Println("### container_list_handler.go - container info - Properties: ", gardenContainerInfo.Properties)
		  fmt.Println("### container_list_handler.go - container info - ExternalIP: ", gardenContainerInfo.ExternalIP)
		  fmt.Println("### container_list_handler.go - container info - Events: ", gardenContainerInfo.Events)
		  fmt.Println("### container_list_handler.go - container info - MappedPorts: ", gardenContainerInfo.MappedPorts)
		  fmt.Println("### container_list_handler.go - container info - State: ", gardenContainerInfo.State)
		  fmt.Println("### container_list_handler.go - container info - ProcessIDs: ", gardenContainerInfo.ProcessIDs)*/

		var Container_Ip, External_Ip string
		for key, value := range gardenContainerInfo.Properties {
			// fmt.Println("### container_list_handler.go - container info - Properties: key - value :", key, value)
			// Disable diego CNI
			if strings.HasSuffix(key, "host-interface") {
				fmt.Println("### container_list_handler.go - containerInfo.Interface_Id :", containerInfo.Interface_Id)
				containerInfo.Interface_Id = value
				break
			} else if strings.HasSuffix(key, "container-ip") {
				Container_Ip = value
			} else if strings.HasSuffix(key, "external-ip") {
				External_Ip = value
			}
		}

		// Enable diego CNI
		if containerInfo.Interface_Id == "" {
			if Container_Ip != "" && External_Ip != "" {
				containerInfo.Interface_Id = Container_Ip + "-" + External_Ip
			} else {
				fmt.Println("### container_list_handler.go - Not exist Container_Ip, External_Ip : ")
			}
		}

		//bugfix container info - MappedPorts:  [{61000 8080} {61001 2222} {61002 61001} {61003 61002}]
		//fmt.Println("!!!!!!! before container_list_handler.go - container info - MappedPorts: ", gardenContainerInfo.MappedPorts)
		if len(gardenContainerInfo.MappedPorts) > 0 {
			gardenContainerInfo.MappedPorts = gardenContainerInfo.MappedPorts[:len(gardenContainerInfo.MappedPorts)-2]
		}
		//fmt.Println("!!!!!!! after container_list_handler.go - container info - MappedPorts: ", gardenContainerInfo.MappedPorts)
		for _, value := range gardenContainerInfo.MappedPorts {
			if value.ContainerPort != 2222 {
				container_host_port = value.HostPort
			}
		}
		//extract Container ID from gardenContainerInfo.ContainerPath - separator '/' & last value
		containerIDPaths := strings.Split(gardenContainerInfo.ContainerPath, "/")
		containerInfo.Container_Id = containerIDPaths[len(containerIDPaths)-1]
		for key, props := range gardenContainerInfo.Properties {
			if strings.Contains(key, "app_id") {
				containerInfo.Application_Id = props
			}
		}
		containerInfo.Container_Port = container_host_port
		containerInfos = append(containerInfos, containerInfo)
	}
	fmt.Println("###### applicationInfos :", applications)
	fmt.Println("###### containerInfos :", containerInfos)
	for _, apps := range applications {
		var containermetric ContainerMetricsMetadata
		containermetric.Limits = apps.Limits
		containermetric.Application_Id = apps.Application_Id
		containermetric.Application_Name = apps.Application_Name
		containermetric.Application_Urls = apps.Application_Urls
		for _, bulkMetrics := range containerBulkMetrics {
			if apps.Application_Id == bulkMetrics.Guid && apps.Application_Index == bulkMetrics.Index {
				containermetric.UsageMetrics.MemoryUsageInBytes = bulkMetrics.MemoryUsageInBytes
				containermetric.UsageMetrics.DiskUsageInBytes = bulkMetrics.DiskUsageInBytes
				containermetric.UsageMetrics.TimeSpentInCPU = bulkMetrics.TimeSpentInCPU
				containermetric.Application_Index = strconv.Itoa(apps.Application_Index)
			}
		}
		for _, cons := range containerInfos {
			if apps.Application_Id == cons.Application_Id && apps.Container_Port == cons.Container_Port {
				fmt.Println("!!!!!!! true apps.Application_Id == cons.Application_Id && apps.Container_Port == cons.Container_Port ")
				//containermetric.Container_Id = cons.Container_Id
				containermetric.Container_Id = apps.Container_Id
				containermetric.Interface_Id = cons.Interface_Id

			}
		}
		containermetrics = append(containermetrics, containermetric)
	}
	for _, conmetric := range containermetrics {
		fmt.Println("## container_list_handler.go - App & Container Info :", conmetric)
		fmt.Println("## container_list_handler.go - App & Container Info - Container Id :", conmetric.Container_Id)
		fmt.Println("## container_list_handler.go - App & Container Info - Interface Id :", conmetric.Interface_Id)
		fmt.Println("## container_list_handler.go - App & Container Info - App Id :", conmetric.Application_Id)
		fmt.Println("## container_list_handler.go - App & Container Info - App name :", conmetric.Application_Name)
		fmt.Println("## container_list_handler.go - App & Container Info - App uris :", conmetric.Application_Urls)
	}
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(containermetrics)
	if err != nil {
		c.logger.Error("failed-to-marshalling-containermetrics", err)
	}
	w.Write(b)
}
