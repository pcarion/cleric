package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

type ClericJsonConfig struct {
	Version    string                  `json:"version"`
	McpServers []*McpServerDescription `json:"mcpServers"`
}

type ClericConfig struct {
	Path   string
	Config *ClericJsonConfig
}

func NewClericConfig(path string) *ClericConfig {
	checkClericPath(path)
	return &ClericConfig{
		Path: path,
		Config: &ClericJsonConfig{
			Version:    "1.0.0",
			McpServers: []*McpServerDescription{},
		},
	}
}

func (c *ClericConfig) LoadMcpServers() ([]*McpServerDescription, error) {
	// read the file
	file, err := os.Open(c.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open cleric config file %s: %v", c.Path, err)
	}
	defer file.Close()

	// decode the file into the config
	err = json.NewDecoder(file).Decode(&c.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cleric config file %s: %v", c.Path, err)
	}

	return c.Config.McpServers, nil
}

func (c *ClericConfig) SaveMcpServers(servers []*McpServerDescription) {
	c.Config.McpServers = servers
	// we encode the config in a format that is easy to read
	file, err := os.Create(c.Path)
	if err != nil {
		panic(fmt.Sprintf("Failed to create cleric config file at %s: %v", c.Path, err))
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(c.Config)
}

func checkClericPath(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err = os.Create(path)
		if err != nil {
			panic(fmt.Sprintf("Failed to create cleric config file at %s: %v", path, err))
		}
		// we create an empty config
		emptyConfig := &ClericJsonConfig{
			Version:    "1.0.0",
			McpServers: []*McpServerDescription{},
		}
		// we write the empty config in path
		file, err := os.Create(path)
		if err != nil {
			panic(fmt.Sprintf("Failed to create cleric config file at %s: %v", path, err))
		}
		defer file.Close()
		// we encode the empty config in a format that is easy to read
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		encoder.Encode(emptyConfig)
	}
}
