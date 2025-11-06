package yamlutils

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// PersistInConfigFile updates (or creates) the config file used by Viper,
// setting `key` to `value` while preserving comments and formatting.
// It determines the config file path using viper.ConfigFileUsed(), or falls
// back to `<directory>/.ab-config.yaml` where `directory` comes from
// viper.GetString("directory") (default ".").
func PersistInConfigFile(key string, value interface{}) error {
	cfgPath := viper.ConfigFileUsed()
	if cfgPath == "" {
		dir := viper.GetString("directory")
		if dir == "" {
			dir = "."
		}
		cfgPath = filepath.Join(dir, ".ab-config.yaml")
		viper.SetConfigFile(cfgPath)
	}

	var doc yaml.Node
	content, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create empty document with mapping
			doc = yaml.Node{Kind: yaml.DocumentNode}
			mapping := yaml.Node{Kind: yaml.MappingNode}
			doc.Content = []*yaml.Node{&mapping}
		} else {
			log.Error().Err(err).Str("path", cfgPath).Msg("failed to read config file")
			// Fallback: try to let viper write if caller previously set values
			if werr := viper.WriteConfigAs(cfgPath); werr != nil {
				log.Error().Err(werr).Msg("failed fallback write to config file")
			}
			return err
		}
	} else {
		if len(content) == 0 {
			doc = yaml.Node{Kind: yaml.DocumentNode}
			mapping := yaml.Node{Kind: yaml.MappingNode}
			doc.Content = []*yaml.Node{&mapping}
		} else {
			if err := yaml.Unmarshal(content, &doc); err != nil {
				log.Error().Err(err).Str("path", cfgPath).Msg("failed to parse existing yaml config file")
				doc = yaml.Node{Kind: yaml.DocumentNode}
				mapping := yaml.Node{Kind: yaml.MappingNode}
				doc.Content = []*yaml.Node{&mapping}
			}
		}
	}

	// Ensure document has a mapping node
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		mapping := yaml.Node{Kind: yaml.MappingNode}
		doc.Content = []*yaml.Node{&mapping}
	}

	mapping := doc.Content[0]

	// Convert value to string representation for scalar nodes
	var valStr string
	switch v := value.(type) {
	case string:
		valStr = v
	default:
		// Let yaml marshal non-string values to preserve types
		tmp, _ := yaml.Marshal(v)
		valStr = string(bytes.TrimSpace(tmp))
	}

	// Find and replace the key if present
	found := false
	for i := 0; i < len(mapping.Content); i += 2 {
		keyNode := mapping.Content[i]
		valNode := mapping.Content[i+1]
		if keyNode.Value == key {
			valNode.Kind = yaml.ScalarNode
			valNode.Tag = "!!str"
			valNode.Value = valStr
			found = true
			break
		}
	}

	if !found {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}
		valNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: valStr}
		mapping.Content = append(mapping.Content, keyNode, valNode)
	}

	// Encode back to YAML preserving comments/formatting
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		enc.Close()
		log.Error().Err(err).Msg("failed to encode yaml document")
		return err
	}
	enc.Close()

	if err := os.WriteFile(cfgPath, buf.Bytes(), 0o644); err != nil {
		log.Error().Err(err).Str("path", cfgPath).Msg("failed to write config file")
		return err
	}

	log.Info().Str("key", key).Str("path", cfgPath).Msg("persisted key to config file")
	return nil
}
