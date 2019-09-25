package rancher

type Services struct {
	Type           string             `json:"type"`
	ResourceType   string             `json:"resourceType"`
	Links          ServicesLinks      `json:"links"`
	CreateTypes    CreateTypes        `json:"createTypes"`
	Actions        ServicesActions    `json:"actions"`
	Data           []Datum            `json:"data"`
	SortLinks      map[string]*string `json:"sortLinks"`
	Pagination     Pagination         `json:"pagination"`
	Sort           interface{}        `json:"sort"`
	Filters        map[string]*string `json:"filters"`
	CreateDefaults CreateDefaults     `json:"createDefaults"`
}

type ServicesActions struct {
}

type CreateDefaults struct {
	StackID string `json:"stackId"`
}

type CreateTypes struct {
	DNSService           string `json:"dnsService"`
	ExternalService      string `json:"externalService"`
	LoadBalancerService  string `json:"loadBalancerService"`
	NetworkDriverService string `json:"networkDriverService"`
	Service              string `json:"service"`
	StorageDriverService string `json:"storageDriverService"`
}

type Datum struct {
	ID                     string        `json:"id"`
	Type                   string        `json:"type"`
	Links                  DatumLinks    `json:"links"`
	Actions                DatumActions  `json:"actions"`
	BaseType               string        `json:"baseType"`
	Name                   string        `json:"name"`
	State                  string        `json:"state"`
	AccountID              string        `json:"accountId"`
	AssignServiceIPAddress bool          `json:"assignServiceIpAddress"`
	CreateIndex            int64         `json:"createIndex"`
	Created                string        `json:"created"`
	CreatedTS              int64         `json:"createdTS"`
	CurrentScale           int64         `json:"currentScale"`
	Description            interface{}   `json:"description"`
	ExternalID             interface{}   `json:"externalId"`
	FQDN                   interface{}   `json:"fqdn"`
	HealthState            string        `json:"healthState"`
	InstanceIDS            []string      `json:"instanceIds"`
	Kind                   string        `json:"kind"`
	LaunchConfig           LaunchConfig  `json:"launchConfig"`
	LBConfig               interface{}   `json:"lbConfig"`
	LinkedServices         interface{}   `json:"linkedServices"`
	Metadata               Metadata      `json:"metadata"`
	PublicEndpoints        interface{}   `json:"publicEndpoints"`
	Removed                interface{}   `json:"removed"`
	RetainIP               interface{}   `json:"retainIp"`
	Scale                  int64         `json:"scale"`
	ScalePolicy            interface{}   `json:"scalePolicy"`
	SecondaryLaunchConfigs []interface{} `json:"secondaryLaunchConfigs"`
	SelectorContainer      interface{}   `json:"selectorContainer"`
	SelectorLink           interface{}   `json:"selectorLink"`
	StackID                string        `json:"stackId"`
	StartOnCreate          bool          `json:"startOnCreate"`
	System                 bool          `json:"system"`
	Transitioning          string        `json:"transitioning"`
	TransitioningMessage   interface{}   `json:"transitioningMessage"`
	TransitioningProgress  interface{}   `json:"transitioningProgress"`
	Upgrade                Upgrade       `json:"upgrade"`
	UUID                   string        `json:"uuid"`
	Vip                    interface{}   `json:"vip"`
}

type DatumActions struct {
	Upgrade           string `json:"upgrade"`
	Restart           string `json:"restart"`
	Update            string `json:"update"`
	Remove            string `json:"remove"`
	Deactivate        string `json:"deactivate"`
	Removeservicelink string `json:"removeservicelink"`
	Addservicelink    string `json:"addservicelink"`
	Setservicelinks   string `json:"setservicelinks"`
}

type LaunchConfig struct {
	Type                  string      `json:"type"`
	DNS                   []string    `json:"dns"`
	DNSSearch             []string    `json:"dnsSearch"`
	Environment           Environment `json:"environment"`
	ImageUUID             string      `json:"imageUuid"`
	InstanceTriggeredStop string      `json:"instanceTriggeredStop"`
	Kind                  string      `json:"kind"`
	Labels                Labels      `json:"labels"`
	LogConfig             LogConfig   `json:"logConfig"`
	Memory                int64       `json:"memory"`
	NetworkMode           string      `json:"networkMode"`
	Privileged            bool        `json:"privileged"`
	PublishAllPorts       bool        `json:"publishAllPorts"`
	ReadOnly              bool        `json:"readOnly"`
	RunInit               bool        `json:"runInit"`
	StartOnCreate         bool        `json:"startOnCreate"`
	StdinOpen             bool        `json:"stdinOpen"`
	System                bool        `json:"system"`
	TTY                   bool        `json:"tty"`
	Version               string      `json:"version"`
	Vcpu                  int64       `json:"vcpu"`
	DrainTimeoutMS        int64       `json:"drainTimeoutMs"`
}

type Environment struct {
	Goenv                    string  `json:"GOENV"`
	SGeneratealgorithm       *string `json:"S_GENERATEALGORITHM,omitempty"`
	SLogLevel                string  `json:"S_LOG_LEVEL"`
	SSnowflakeconfigEpoch    *string `json:"S_SNOWFLAKECONFIG_EPOCH,omitempty"`
	SSnowflakeconfigNodebits *string `json:"S_SNOWFLAKECONFIG_NODEBITS,omitempty"`
	SSnowflakeconfigStepbits *string `json:"S_SNOWFLAKECONFIG_STEPBITS,omitempty"`
	SClientidHost            *string `json:"S_CLIENTID_HOST,omitempty"`
	SMasterdbPassword        *string `json:"S_MASTERDB_PASSWORD,omitempty"`
	SMasterdbUser            *string `json:"S_MASTERDB_USER,omitempty"`
	SSlavedbPassword         *string `json:"S_SLAVEDB_PASSWORD,omitempty"`
	SSlavedbUser             *string `json:"S_SLAVEDB_USER,omitempty"`
}

type Labels struct {
	BasePath                    string  `json:"base_path"`
	IoRancherContainerPullImage string  `json:"io.rancher.container.pull_image"`
	IoRancherContainerStartOnce string  `json:"io.rancher.container.start_once"`
	IoRancherServiceHash        string  `json:"io.rancher.service.hash"`
	LBG7PayExpose80             string  `json:"lb.g7pay.expose80"`
	ProjectDescription          string  `json:"projects.description"`
	ProjectGroup                string  `json:"projects.group"`
	ProjectName                 string  `json:"projects.name"`
	ProjectVersion              string  `json:"projects.version"`
	Upstreams                   *string `json:"upstreams,omitempty"`
}

type LogConfig struct {
	Type string `json:"type"`
}

type DatumLinks struct {
	Self               string `json:"self"`
	Account            string `json:"account"`
	Consumedbyservices string `json:"consumedbyservices"`
	Consumedservices   string `json:"consumedservices"`
	Instances          string `json:"instances"`
	NetworkDrivers     string `json:"networkDrivers"`
	ServiceExposeMaps  string `json:"serviceExposeMaps"`
	ServiceLogs        string `json:"serviceLogs"`
	Stack              string `json:"stack"`
	StorageDrivers     string `json:"storageDrivers"`
	ContainerStats     string `json:"containerStats"`
}

type Metadata struct {
	IoRancherServiceHash string `json:"io.rancher.service.hash"`
}

type Upgrade struct {
	Type              string            `json:"type"`
	InServiceStrategy InServiceStrategy `json:"inServiceStrategy"`
	ToServiceStrategy interface{}       `json:"toServiceStrategy"`
}

type InServiceStrategy struct {
	Type                           string        `json:"type"`
	BatchSize                      int64         `json:"batchSize"`
	IntervalMillis                 int64         `json:"intervalMillis"`
	LaunchConfig                   LaunchConfig  `json:"launchConfig"`
	PreviousLaunchConfig           LaunchConfig  `json:"previousLaunchConfig"`
	PreviousSecondaryLaunchConfigs []interface{} `json:"previousSecondaryLaunchConfigs"`
	SecondaryLaunchConfigs         []interface{} `json:"secondaryLaunchConfigs"`
	StartFirst                     bool          `json:"startFirst"`
}

type ServicesLinks struct {
	Self string `json:"self"`
}

type Pagination struct {
	First    interface{} `json:"first"`
	Previous interface{} `json:"previous"`
	Next     interface{} `json:"next"`
	Limit    int64       `json:"limit"`
	Total    interface{} `json:"total"`
	Partial  bool        `json:"partial"`
}

