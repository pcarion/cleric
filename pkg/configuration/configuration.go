package configuration

type Configuration struct {
	claudeConfig *ClaudeDesktopConfig
	clericConfig *ClericConfig
}

func LoadConfiguration() *Configuration {
	claudeConfig := NewClaudeDesktopConfig()
	clericConfig := NewClericConfig()

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
	// we check if the servers from claude are already in cleric
	// if they are, we remove them from cleric
	// if they are not, we add them to cleric
	allServers := []*McpServerDescription{}
	// we take all the servers from claude
	for _, server := range claudeServers {
		if !contains(clericServers, server) {
			// we mark the server as in configuration
			server.InConfiguration = true
			allServers = append(allServers, server)
		}
	}
	// we take all the servers from cleric that are not in claude
	for _, server := range clericServers {
		if !contains(allServers, server) {
			// we mark the server as not in configuration
			server.InConfiguration = false
			allServers = append(allServers, server)
		}
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
