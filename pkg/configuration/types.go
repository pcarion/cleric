package configuration

type McpServerConfiguration struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

type McpServerDescription struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Index           int                    `json:"-"` // index in the list of servers
	InConfiguration bool                   `json:"-"`
	Configuration   McpServerConfiguration `json:"configuration"`
}
