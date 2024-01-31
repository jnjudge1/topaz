package cc

import (
	"fmt"
	"os"
	"runtime"

	"github.com/aserto-dev/topaz/pkg/cli/g"
)

const (
	DefaultValue             string = ""
	defaultContainerRegistry string = "ghcr.io/aserto-dev"
	defaultContainerImage    string = "topaz"
	defaultContainerTag      string = "latest"
	defaultContainerName     string = "topaz"
)

// Container returns the fully qualified container name (registry/image:tag).
func Container(registry, image, tag string) string {
	if container := os.Getenv("CONTAINER"); container != "" {
		return container
	}

	return fmt.Sprintf("%s/%s:%s",

		g.Iff(registry != "", registry, ContainerRegistry()),
		g.Iff(image != "", image, ContainerImage()),
		g.Iff(tag != "", tag, ContainerTag()),
	)
}

// ContainerRegistry returns the container registry (host[:port]/repo).
func ContainerRegistry() string {
	if containerRegistry := os.Getenv("CONTAINER_REGISTRY"); containerRegistry != "" {
		return containerRegistry
	}
	return defaultContainerRegistry
}

// ContainerImage returns the container image name.
func ContainerImage() string {
	if containerImage := os.Getenv("CONTAINER_IMAGE"); containerImage != "" {
		return containerImage
	}
	return defaultContainerImage
}

// ContainerTag returns the container tag (label or semantic version).
func ContainerTag() string {
	if containerTag := os.Getenv("CONTAINER_TAG"); containerTag != "" {
		return containerTag
	}
	return defaultContainerTag
}

// ContainerPlatform, returns the container platform for multi-platform capable servers.
func ContainerPlatform() string {
	if containerPlatform := os.Getenv("CONTAINER_PLATFORM"); containerPlatform != "" {
		return containerPlatform
	}
	return "linux/" + runtime.GOARCH
}

// ContainerName returns the container instance name (docker run --name CONTAINER_NAME).
func ContainerName() string {
	if containerName := os.Getenv("CONTAINER_NAME"); containerName != "" {
		return containerName
	}
	return defaultContainerName
}
