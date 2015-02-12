package dockertools

import (
	"github.com/silenteh/gantryos/utils"
)

type ContainerStateWaiting struct {
	// Reason could be pulling image,
	Reason string `json:"reason,omitempty"`
}

type ContainerStateRunning struct {
	StartedAt utils.Time `json:"startedAt,omitempty"`
}

type ContainerStateTerminated struct {
	ExitCode   int        `json:"exitCode"`
	Signal     int        `json:"signal,omitempty"`
	Reason     string     `json:"reason,omitempty"`
	Message    string     `json:"message,omitempty"`
	StartedAt  utils.Time `json:"startedAt,omitempty"`
	FinishedAt utils.Time `json:"finishedAt,omitempty"`
}

// ContainerState holds a possible state of container.
// Only one of its members may be specified.
// If none of them is specified, the default one is ContainerStateWaiting.
type ContainerState struct {
	Waiting     *ContainerStateWaiting    `json:"waiting,omitempty"`
	Running     *ContainerStateRunning    `json:"running,omitempty"`
	Termination *ContainerStateTerminated `json:"termination,omitempty"`
}

type ContainerStatus struct {
	// TODO(dchen1107): Should we rename PodStatus to a more generic name or have a separate states
	// defined for container?
	State ContainerState `json:"state,omitempty"`
	// Note that this is calculated from dead containers.  But those containers are subject to
	// garbage collection.  This value will get capped at 5 by GC.
	RestartCount int `json:"restartCount"`
	// TODO(dchen1107): Deprecated this soon once we pull entire PodStatus from node,
	// not just PodInfo. Now we need this to remove docker.Container from API
	PodIP string `json:"podIP,omitempty"`
	// TODO(dchen1107): Need to decide how to represent this in v1beta3
	Image       string `json:"image"`
	ContainerID string `json:"containerID,omitempty" description:"container's ID in the format 'docker://<container_id>'"`
}

// Container represents a single container that is expected to be run on the host.
type Container struct {
	// Required: This must be a DNS_LABEL.  Each container in a pod must
	// have a unique name.
	Name string `json:"name"`
	// Required.
	Image string `json:"image"`
	// Optional: Defaults to whatever is defined in the image.
	Command []string `json:"command,omitempty"`
	// Optional: Defaults to Docker's default.
	WorkingDir string   `json:"workingDir,omitempty"`
	Ports      []Port   `json:"ports,omitempty"`
	Env        []EnvVar `json:"env,omitempty"`
	// Optional: Defaults to unlimited.
	//Memory resource.Quantity `json:"memory,omitempty"`
	// Optional: Defaults to unlimited.
	//CPU          resource.Quantity `json:"cpu,omitempty"`
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`
	//LivenessProbe *LivenessProbe    `json:"livenessProbe,omitempty"`
	//Lifecycle     *Lifecycle        `json:"lifecycle,omitempty"`
	// Optional: Defaults to /dev/termination-log
	TerminationMessagePath string `json:"terminationMessagePath,omitempty"`
	// Optional: Default to false.
	Privileged bool `json:"privileged,omitempty"`
	// Optional: Policy for pulling images for this container
	ImagePullPolicy PullPolicy `json:"imagePullPolicy"`
}

type Protocol string

// Port represents a network port in a single container
type Port struct {
	// Optional: If specified, this must be a DNS_LABEL.  Each named port
	// in a pod must have a unique name.
	Name string `json:"name,omitempty"`
	// Optional: If specified, this must be a valid port number, 0 < x < 65536.
	HostPort int `json:"hostPort,omitempty"`
	// Required: This must be a valid port number, 0 < x < 65536.
	ContainerPort int `json:"containerPort"`
	// Optional: Supports "TCP" and "UDP".  Defaults to "TCP".
	Protocol Protocol `json:"protocol,omitempty"`
	// Optional: What host IP to bind the external port to.
	HostIP string `json:"hostIP,omitempty"`
}

// EnvVar represents an environment variable present in a Container.
type EnvVar struct {
	// Required: This must be a C_IDENTIFIER.
	Name string `json:"name"`
	// Optional: defaults to "".
	Value string `json:"value,omitempty"`
}

// VolumeMount describes a mounting of a Volume within a container.
type VolumeMount struct {
	// Required: This must match the Name of a Volume [above].
	Name string `json:"name"`
	// Optional: Defaults to false (read-write).
	ReadOnly bool `json:"readOnly,omitempty"`
	// Required.
	MountPath string `json:"mountPath"`
}

// PullPolicy describes a policy for if/when to pull a container image
type PullPolicy string

const (
	// PullAlways means that kubelet always attempts to pull the latest image.  Container will fail If the pull fails.
	PullAlways PullPolicy = "PullAlways"
	// PullNever means that kubelet never pulls an image, but only uses a local image.  Container will fail if the image isn't present
	PullNever PullPolicy = "PullNever"
	// PullIfNotPresent means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
	PullIfNotPresent PullPolicy = "PullIfNotPresent"
)

type PodInfo map[string]ContainerStatus

// DNSPolicy defines how a pod's DNS will be configured.
type DNSPolicy string

const (
	// DNSClusterFirst indicates that the pod should use cluster DNS
	// first, if it is available, then fall back on the default (as
	// determined by kubelet) DNS settings.
	DNSClusterFirst DNSPolicy = "ClusterFirst"

	// DNSDefault indicates that the pod should use the default (as
	// determined by kubelet) DNS settings.
	DNSDefault DNSPolicy = "Default"
)

// PodSpec is a description of a pod
type PodSpec struct {
	Volumes       []Volume      `json:"volumes"`
	Containers    []Container   `json:"containers"`
	RestartPolicy RestartPolicy `json:"restartPolicy,omitempty"`
	// Optional: Set DNS policy.  Defaults to "ClusterFirst"
	DNSPolicy DNSPolicy `json:"dnsPolicy,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Host is a request to schedule this pod onto a specific host.  If it is non-empty,
	// the the scheduler simply schedules this pod onto that host, assuming that it fits
	// resource requirements.
	Host string `json:"host,omitempty"`
}

type RestartPolicyAlways struct{}

// TODO(dchen1107): Define what kinds of failures should restart.
// TODO(dchen1107): Decide whether to support policy knobs, and, if so, which ones.
type RestartPolicyOnFailure struct{}

type RestartPolicyNever struct{}

// RestartPolicy describes how the container should be restarted.
// Only one of the following restart policies may be specified.
// If none of the following policies is specified, the default one
// is RestartPolicyAlways.
type RestartPolicy struct {
	Always    *RestartPolicyAlways    `json:"always,omitempty"`
	OnFailure *RestartPolicyOnFailure `json:"onFailure,omitempty"`
	Never     *RestartPolicyNever     `json:"never,omitempty"`
}

// Volume represents a named volume in a pod that may be accessed by any containers in the pod.
type Volume struct {
	// Required: This must be a DNS_LABEL.  Each volume in a pod must have
	// a unique name.
	Name string `json:"name"`
	// Source represents the location and type of a volume to mount.
	// This is optional for now. If not specified, the Volume is implied to be an EmptyDir.
	// This implied behavior is deprecated and will be removed in a future version.
	Source *VolumeSource `json:"source"`
}

// VolumeSource represents the source location of a valume to mount.
// Only one of its members may be specified.
type VolumeSource struct {
	// HostDir represents a pre-existing directory on the host machine that is directly
	// exposed to the container. This is generally used for system agents or other privileged
	// things that are allowed to see the host machine. Most containers will NOT need this.
	// TODO(jonesdl) We need to restrict who can use host directory mounts and who can/can not
	// mount host directories as read/write.
	HostDir *HostDir `json:"hostDir"`
	// EmptyDir represents a temporary directory that shares a pod's lifetime.
	EmptyDir *EmptyDir `json:"emptyDir"`
	// GCEPersistentDisk represents a GCE Disk resource that is attached to a
	// kubelet's host machine and then exposed to the pod.
	//GCEPersistentDisk *GCEPersistentDisk `json:"persistentDisk"`
	// GitRepo represents a git repository at a particular revision.
	GitRepo *GitRepo `json:"gitRepo"`
}

// HostDir represents bare host directory volume.
type HostDir struct {
	Path string `json:"path"`
}

type EmptyDir struct{}

// GitRepo represents a volume that is pulled from git when the pod is created.
type GitRepo struct {
	// Repository URL
	Repository string `json:"repository"`
	// Commit hash, this is optional
	Revision string `json:"revision"`
	// TODO: Consider credentials here.
}
