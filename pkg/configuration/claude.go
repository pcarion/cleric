package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type ClaudeDesktopConfig struct {
	Path string
}

func NewClaudeDesktopConfig() *ClaudeDesktopConfig {
	return &ClaudeDesktopConfig{
		Path: getClaudeDesktopConfigPath(),
	}
}

func (c *ClaudeDesktopConfig) LoadMcpServers() ([]*McpServerDescription, error) {
	mcpServers := []*McpServerDescription{}

	content, err := os.ReadFile(c.Path)
	if err != nil {
		return mcpServers, err
	}

	// structure to hold the raw content
	contentMap := make(map[string]interface{})

	// read content as a json object
	err = json.Unmarshal(content, &contentMap)
	if err != nil {
		return mcpServers, err
	}

	// find the mcp_servers key
	mcpServersValue, ok := contentMap["mcpServers"]
	if !ok {
		return mcpServers, nil
	}

	// ensure that mcpServersMap is a map
	mcpServersMap, ok := mcpServersValue.(map[string]interface{})
	if !ok {
		return mcpServers, nil
	}

	// iterate over the keys in mcpServersMap
	for key, mcpServer := range mcpServersMap {
		entry := McpServerDescription{
			Name: key,
		}
		// read the description of the server
		description, ok := mcpServer.(map[string]interface{})
		if !ok {
			return mcpServers, fmt.Errorf("invalid mcp server description for %s", key)
		}
		// read the command of the server (field is required)
		command, ok := description["command"].(string)
		if !ok {
			return mcpServers, fmt.Errorf("invalid command for mcp server %s", key)
		}
		entry.Configuration.Command = command

		// args	is optional
		args, ok := description["args"].([]string)
		if !ok {
			entry.Configuration.Args = []string{}
		} else {
			entry.Configuration.Args = args
		}

		// env is optional
		env, ok := description["env"].(map[string]string)
		if !ok {
			entry.Configuration.Env = map[string]string{}
		} else {
			entry.Configuration.Env = env
		}

		// add the entry to the list
		mcpServers = append(mcpServers, &entry)
	}

	return mcpServers, nil
}

func (c *ClaudeDesktopConfig) SaveMcpServers(servers []*McpServerDescription) {
	// we read the current content of the file
	content, err := os.ReadFile(c.Path)
	if err != nil {
		panic(fmt.Sprintf("Failed to read claude config file at %s: %v", c.Path, err))
	}
	// we decode the content as a map
	var contentMap map[string]interface{}
	err = json.Unmarshal(content, &contentMap)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal claude config file at %s: %v", c.Path, err))
	}
	// we update the contentMap with the new servers
	mcpServersMap := make(map[string]interface{})
	for _, server := range servers {
		if server.InConfiguration {
			// Convert Configuration struct to map
			serverConfig := map[string]interface{}{
				"command": server.Configuration.Command,
				"args":    server.Configuration.Args,
				"env":     server.Configuration.Env,
			}
			mcpServersMap[server.Name] = serverConfig
		}
	}
	// we update the contentMap with the new mcpServersMap
	contentMap["mcpServers"] = mcpServersMap
	// we encode the contentMap as a json object
	// we use the json marshaller to write the contentMap in a format that is easy to read
	file, err := os.Create(c.Path)
	if err != nil {
		panic(fmt.Sprintf("Failed to create claude config file at %s: %v", c.Path, err))
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(contentMap)

	err = os.WriteFile(c.Path, content, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to write claude config file at %s: %v", c.Path, err))
	}
}

func getClaudeDesktopConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	path := ""

	if runtime.GOOS == "windows" {
		path = filepath.Join(homeDir, "AppData", "Roaming", "Claude", "claude_desktop_config.json")
	} else {
		path = filepath.Join(homeDir, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// we create an empty config
		emptyConfig := make(map[string]interface{})
		// we write the empty config in path
		emptyConfig["mcpServers"] = make(map[string]interface{})
		file, err := os.Create(path)
		if err != nil {
			panic(fmt.Sprintf("Failed to create claude config file at %s: %v", path, err))
		}
		defer file.Close()

		// json marshaller to write the empty config in a format that is easy to read

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		encoder.Encode(emptyConfig)

		fmt.Printf("Created empty claude config file at %s with content %v\n", path, emptyConfig)
	}
	return path
}
