package types

// TODO: better to use docker/docker/image.Image.
// if use docker/docker/image.Image, it's occur dependency errors.

import (
	"time"

	"github.com/docker/go-connections/nat"
)

// Image stores the image configuration
type Image struct {
	V1Image
	History    []History `json:"history,omitempty"`
	OSVersion  string    `json:"os.version,omitempty"`
	OSFeatures []string  `json:"os.features,omitempty"`
}

// V1Image stores the V1 image configuration.
type V1Image struct {
	ID              string    `json:"id,omitempty"`
	Parent          string    `json:"parent,omitempty"`
	Comment         string    `json:"comment,omitempty"`
	Created         time.Time `json:"created"`
	Container       string    `json:"container,omitempty"`
	ContainerConfig Config    `json:"container_config,omitempty"`
	DockerVersion   string    `json:"docker_version,omitempty"`
	Author          string    `json:"author,omitempty"`
	Config          Config    `json:"config,omitempty"`
	Architecture    string    `json:"architecture,omitempty"`
	OS              string    `json:"os,omitempty"`
	Size            int64     `json:",omitempty"`
}

type Config struct {
	Hostname        string              // Hostname
	Domainname      string              // Domainname
	User            string              // User that will run the command(s) inside the container, also support user:group
	AttachStdin     bool                // Attach the standard input, makes possible user interaction
	AttachStdout    bool                // Attach the standard output
	AttachStderr    bool                // Attach the standard error
	ExposedPorts    nat.PortSet         `json:",omitempty"` // List of exposed ports
	Tty             bool                // Attach standard streams to a tty, including stdin if it is not closed.
	OpenStdin       bool                // Open stdin
	StdinOnce       bool                // If true, close stdin after the 1 attached client disconnects.
	Env             []string            // List of environment variable to set in the container
	Cmd             []string            // Command to run when starting the container
	Healthcheck     *HealthConfig       `json:",omitempty"` // Healthcheck describes how to check the container is healthy
	ArgsEscaped     bool                `json:",omitempty"` // True if command is already escaped (Windows specific)
	Image           string              // Name of the image as it was passed by the operator (e.g. could be symbolic)
	Volumes         map[string]struct{} // List of volumes (mounts) used for the container
	WorkingDir      string              // Current directory (PWD) in the command will be launched
	Entrypoint      []string            // Entrypoint to run when starting the container
	NetworkDisabled bool                `json:",omitempty"` // Is network disabled
	MacAddress      string              `json:",omitempty"` // Mac Address of the container
	OnBuild         []string            // ONBUILD metadata that were defined on the image Dockerfile
	Labels          map[string]string   // List of labels set to this container
	StopSignal      string              `json:",omitempty"` // Signal to stop a container
	StopTimeout     *int                `json:",omitempty"` // Timeout (in seconds) to stop a container
	Shell           []string            `json:",omitempty"` // Shell for shell-form of RUN, CMD, ENTRYPOINT
}

// HealthConfig holds configuration settings for the HEALTHCHECK feature.
type HealthConfig struct {
	Test        []string      `json:",omitempty"`
	Interval    time.Duration `json:",omitempty"` // Interval is the time to wait between checks.
	Timeout     time.Duration `json:",omitempty"` // Timeout is the time to wait before considering the check to have hung.
	StartPeriod time.Duration `json:",omitempty"` // The start period for the container to initialize before the retries starts to count down.
	Retries     int           `json:",omitempty"`
}

// History stores build commands that were used to create an image
type History struct {
	Created    time.Time `json:"created"`
	Author     string    `json:"author,omitempty"`
	CreatedBy  string    `json:"created_by,omitempty"`
	Comment    string    `json:"comment,omitempty"`
	EmptyLayer bool      `json:"empty_layer,omitempty"`
}
