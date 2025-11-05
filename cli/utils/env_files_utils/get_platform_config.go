package envfilesutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/atomic-blend/backend/cli/config"
)

type PlatformComponent struct {
	Name    string
	Image   string
	Version string
}

// GetPlatformConfig reads the provided .env file and returns a list of
// PlatformComponent with the Version field populated for every component
// listed in config.PlatformComponents. The docker-compose/dockerfilePath
// parameter is currently unused (kept for future image parsing).
func GetPlatformConfig(envfilePath, dockerfilePath string) ([]PlatformComponent, error) {
	components, err := getComponentsFromEnv(envfilePath)
	if err != nil {
		log.Error().Err(err).Msg("failed to get components from env file")
		return nil, err
	}

	// If a docker-compose path is provided, parse it and extract service images.
	components, err = getDockerImagesFromCompose(dockerfilePath, components)
	if err != nil {
		log.Error().Err(err).Msg("failed to get docker images from compose file")
		return nil, err
	}
	// log components for debug
	for _, comp := range components {
		log.Debug().Str("component", comp.Name).Str("image", comp.Image).Str("version", comp.Version).Msg("Retrieved platform component")
	}

	log.Debug().Msg("Successfully retrieved platform components from env and compose files")

	return components, nil
}

func getComponentsFromEnv(envfilePath string) ([]PlatformComponent, error) {
	// Load env file into a map
	envMap, err := godotenv.Read(envfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read env file %s: %w", envfilePath, err)
	}

	var components []PlatformComponent

	for _, comp := range config.PlatformComponents {
		envVar, ok := config.EnvComponentVersionMapping[comp]
		if !ok {
			return nil, fmt.Errorf("no env var mapping for component %q", comp)
		}

		version := envMap[envVar]

		components = append(components, PlatformComponent{
			Name:    comp,
			Image:   "",
			Version: version,
		})
	}
	return components, nil
}

func getDockerImagesFromCompose(dockerfilePath string, components []PlatformComponent) ([]PlatformComponent, error) {
	var services map[string]interface{}
	if dockerfilePath != "" {
		var err error
		services, err = parseComposeServices(dockerfilePath)
		if err != nil {
			return nil, err
		}

		for i := range components {
			img, err := findServiceImageForComponent(services, components[i].Name)
			if err != nil {
				return nil, err
			}
			components[i].Image = img
		}
	}
	return components, nil
}

// parseComposeServices reads and parses the docker-compose YAML and returns
// the `services` map.
func parseComposeServices(dockerfilePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read docker compose file %s: %w", dockerfilePath, err)
	}

	var compose map[string]interface{}
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, fmt.Errorf("failed to parse docker compose file %s: %w", dockerfilePath, err)
	}

	servicesIface, ok := compose["services"]
	if !ok {
		return nil, fmt.Errorf("no services section in docker compose file %s", dockerfilePath)
	}

	services, ok := servicesIface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected services format in docker compose file %s", dockerfilePath)
	}

	return services, nil
}

// findServiceImageForComponent looks up the service corresponding to compName
// in the provided services map and returns the image (without tag/digest).
func findServiceImageForComponent(services map[string]interface{}, compName string) (string, error) {
	candidateKeys := []string{
		compName,
		strings.ReplaceAll(compName, "-", "_"),
		strings.TrimSuffix(compName, "-app"),
		strings.ReplaceAll(strings.TrimSuffix(compName, "-app"), "-", "_"),
	}

	for _, key := range candidateKeys {
		if key == "" {
			continue
		}
		svcIface, ok := services[key]
		if !ok {
			continue
		}

		svcMap, ok := svcIface.(map[string]interface{})
		if !ok {
			continue
		}

		imgIface, ok := svcMap["image"]
		if !ok {
			continue
		}

		imgStr := fmt.Sprintf("%v", imgIface)
		return stripImageTag(imgStr), nil
	}

	return "", fmt.Errorf("service for component %q not found in docker compose (tried keys: %v)", compName, candidateKeys)
}

func stripImageTag(img string) string {
	// If image contains a digest, strip it first (e.g. "image@sha256:...")
	if idx := strings.Index(img, "@"); idx >= 0 {
		img = img[:idx]
	}

	// docker-compose files may contain env substitutions like
	// "${REGISTRY}:${PORT}/repo:tag". To avoid complex parsing and to
	// ensure we remove tags and substitutions, strip everything from
	// the first colon onwards.
	if idx := strings.Index(img, ":"); idx >= 0 {
		return img[:idx]
	}

	return img
}
