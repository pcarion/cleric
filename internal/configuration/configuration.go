package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

type Configuration struct {
	claudeConfig *ClaudeDesktopConfig
	clericConfig *ClericConfig
}

func LoadConfiguration() *Configuration {
	claudeConfig := NewClaudeDesktopConfig(getClaudeDesktopConfigPath())
	clericConfig := NewClericConfig(getClericListMcpServersPath())

	return &Configuration{
		claudeConfig: claudeConfig,
		clericConfig: clericConfig,
	}
}

func (c *Configuration) LoadMcpServers() []*McpServerDescription {
	// we load first the claude config
	claudeServers, err := c.claudeConfig.LoadMcpServers()
	if err != nil {
		panic(err)
	}
	if claudeServers == nil {
		claudeServers = []*McpServerDescription{}
	}

	// then the cleric config
	clericServers, err := c.clericConfig.LoadMcpServers()
	if err != nil {
		panic(err)
	}
	if clericServers == nil {
		clericServers = []*McpServerDescription{}
	}

	// we merge the two lists
	allServers := []*McpServerDescription{}
	// we take all the servers from claude
	for _, server := range claudeServers {
		// we try to find the server in cleric
		var clericServer *McpServerDescription = nil
		for _, s := range clericServers {
			if s.Name == server.Name {
				clericServer = s
				break
			}
		}
		if clericServer != nil {
			// we copy the description from cleric to claude
			server.Description = clericServer.Description
		}
		// we mark the server as in configuration
		server.InConfiguration = true
		allServers = append(allServers, server)
	}
	// we take all the servers from cleric that are not in claude
	for _, server := range clericServers {
		if !contains(allServers, server) {
			// we mark the server as not in configuration
			server.InConfiguration = false
			allServers = append(allServers, server)
		}
	}

	// we set the index for each server
	for _, server := range allServers {
		server.Uuid = uuid.New().String()
	}

	return allServers
}

func (c *Configuration) SaveMcpServers(servers []*McpServerDescription) {
	c.clericConfig.SaveMcpServers(servers)
	c.claudeConfig.SaveMcpServers(servers)
}

func contains(servers []*McpServerDescription, server *McpServerDescription) bool {
	for _, s := range servers {
		if s.Name == server.Name {
			return true
		}
	}
	return false
}

func getClaudeDesktopConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	path := ""

	if runtime.GOOS == "windows" {
		path = filepath.Join(homeDir, "AppData", "Roaming", "Claude", "claude_desktop_config.json")
	} else {
		path = filepath.Join(homeDir, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	}
	return path

}

func getClericListMcpServersPath() string {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, ".cleric.json")

	return path
}

// enum MCP version
type McpVersion string

const (
	McpVersion030 McpVersion = "0.3.0"
	McpVersion031 McpVersion = "0.3.1"
)

func (s *McpServerConfiguration) HasEnvironmentVariables() bool {
	return len(s.Env) > 0
}

func (s *McpServerConfiguration) GetMcpInspectorArgs(mcpVersion McpVersion) []string {
	args := []string{}

	args = append(args, "npx")
	args = append(args, "-y")
	args = append(args, "@modelcontextprotocol/inspector@latest")

	if mcpVersion == McpVersion031 {
		if s.Env != nil {
			for key, value := range s.Env {
				args = append(args, "-e")
				args = append(args, fmt.Sprintf("%s=\"%s\"", key, value))
			}
		}
	}
	// we add the command
	args = append(args, s.Command)

	// we add the args
	// we iterate over the args and add them to the args list
	for _, arg := range s.Args {
		// if an arg contains a space, we add it as a string
		if strings.Contains(arg, " ") {
			args = append(args, fmt.Sprintf("\"%s\"", arg))
		} else {
			args = append(args, arg)
		}
	}
	return args
}
