package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PortProtocol string

const (
	TCP PortProtocol = "TCP"
	UDP              = "UDP"
)

// Workload settings describe the configuration for a workload.
type WorkloadSetting struct {
	Name          string        `json:"name"`
	Description   string        `json:"description,omitempty"`
	ParameterType ParameterType `json:"type"`
	Required      bool          `json:"required,omitempty"`
	Default       string        `json:"default,omitempty"`
	FromParam     string        `json:"fromParam,omitempty"`
}

/// CPU describes a CPU resource allocation for a container.
///
/// The minimum number of logical cpus required for running this container.
type CPU struct {
	Required float64 `json:"required"`
}

/// Memory describes the memory allocation for a container.
///
/// The minimum amount of memory in MB required for running this container. The value should be a positive integer, greater than zero.
type Memory struct {
	Required string `json:"required"`
}

// GPU describes a Container's need for a GPU.
//
// The minimum number of gpus required for running this container.
type GPU struct {
	Required float64 `json:"required"`
}

/// Volume describes a path that is attached to a Container.
///
/// It specifies not only the location, but also the requirements.
type Volume struct {
	Name          string        `json:"name"`
	MountPath     string        `json:"mountPath"`
	AccessMode    AccessMode    `json:"accessMode,omitempty"`
	SharingPolicy SharingPolicy `json:"sharingPolicy,omitempty"`
	Disk          *Disk         `json:"disk,omitempty"`
}

/// AccessMode defines the access modes for file systems.
///
/// Currently, only read/write and read-only are supported.
type AccessMode string

const (
	RW AccessMode = "RW"
	RO AccessMode = "RO"
)

/// SharingPolicy defines whether a filesystem can be shared across containers.
///
/// An Exclusive filesystem can only be attached to one container.
type SharingPolicy string

const (
	Shared    SharingPolicy = "Shared"
	Exclusive SharingPolicy = "Exclusive"
)

// Disk describes the disk requirements for backing a Volume.
type Disk struct {
	Required  string `json:"required"`
	Ephemeral bool   `json:"ephemeral,omitempty"`
}

// ExtendedResource give extension ability
type ExtendedResource struct {
	Name     string `json:"name"`
	Required string `json:"required"`
}

/// Resources defines the resources required by a container.
type Resources struct {
	Cpu      CPU                `json:"cpu"`
	Memory   Memory             `json:"memory"`
	Gpu      GPU                `json:"gpu,omitempty"`
	Volumes  []Volume           `json:"volumes,omitempty"`
	Extended []ExtendedResource `json:"extended,omitempty"`
}

// Env describes an environment variable for a container.
type Env struct {
	Name      string `json:"name"`
	Value     string `json:"value,omitempty"`
	FromParam string `json:"fromParam,omitempty"`
}

// Port describes a port on a Container.
type Port struct {
	Name          string       `json:"name"`
	ContainerPort int32        `json:"port"`
	Protocol      PortProtocol `json:"protocol,omitempty"`
}

// Exec describes a shell command, as an array, for execution in a Container.
type Exec struct {
	Command []string `json:"command"`
}

/// HttpHeader describes an HTTP header.
///
/// Headers are not stored as a map of name/value because the same header is allowed
/// multiple times.
type HttpHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

/// HttpGet describes an HTTP GET request used to probe a container.
type HttpGet struct {
	Path        string       `json:"path"`
	Port        int32        `json:"port"`
	HttpHeaders []HttpHeader `json:"httpHeaders"`
}

/// TcpSocket defines a socket used for health probing.
type TcpSocket struct {
	Port int32 `json:"port"`
}

// HealthProbe describes a probe used to check on the health of a Container.
type HealthProbe struct {
	Exec                *Exec      `json:"exec,omitempty"`
	HttpGet             *HttpGet   `json:"httpGet,omitempty"`
	TcpSocket           *TcpSocket `json:"tcpSocket,omitempty"`
	InitialDelaySeconds int32      `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int32      `json:"periodSeconds,omitempty"`
	TimeoutSeconds      int32      `json:"timeoutSeconds,omitempty"`
	SuccessThreshold    int32      `json:"successThreshold,omitempty"`
	FailureThreshold    int32      `json:"failureThreshold,omitempty"`
}

// ConfigFile describes locations to write configuration as files accessible within the container
type ConfigFile struct {
	Path      string `json:"path"`
	Value     string `json:"value,omitempty"`
	FromParam string `json:"fromParam,omitempty"`
}

// Container describes the container configuration for a Component.
type Container struct {
	Name            string       `json:"name"`
	Image           string       `json:"image"`
	Resources       Resources    `json:"resources"`
	Cmd             []string     `json:"cmd,omitempty"`
	Args            []string     `json:"args,omitempty"`
	Env             []Env        `json:"env,omitempty"`
	Config          []ConfigFile `json:"config,omitempty"`
	Ports           []Port       `json:"ports,omitempty"`
	LivenessProbe   *HealthProbe `json:"livenessProbe,omitempty"`
	ReadinessProbe  *HealthProbe `json:"readinessProbe,omitempty"`
	ImagePullSecret string       `json:"imagePullSecret,omitempty"`
}

type ParameterType string

const (
	Boolean ParameterType = "boolean"
	String  ParameterType = "string"
	Number  ParameterType = "number"
	Null    ParameterType = "null"
)

type Parameter struct {
	Name          string        `json:"name"`
	Description   string        `json:"description,omitempty"`
	ParameterType ParameterType `json:"type"`
	Required      bool          `json:"required,omitempty"`
	Default       string        `json:"default,omitempty"`
}

// ComponentSpec defines the desired state of ComponentSchematic
type ComponentSpec struct {
	Parameters       []Parameter       `json:"parameters,omitempty"`
	WorkloadType     string            `json:"workloadType"`
	OsType           string            `json:"osType,omitempty"`
	Arch             string            `json:"arch,omitempty"`
	Containers       []Container       `json:"containers,omitempty"`
	WorkloadSettings []WorkloadSetting `json:"workloadSettings,omitempty"`
}

type ComponentStatus struct {
}

// +genclient

// ComponentSchematic is the Schema for the components API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ComponentSchematic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentSpec   `json:"spec,omitempty"`
	Status ComponentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ComponentSchematicList contains a list of ComponentSchematic
type ComponentSchematicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ComponentSchematic `json:"items"`
}
