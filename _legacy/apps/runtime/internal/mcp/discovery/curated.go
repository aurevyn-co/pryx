package discovery

// Category represents the type/category of an MCP server
type Category string

const (
	CategoryFilesystem Category = "filesystem"
	CategoryWeb        Category = "web"
	CategoryDatabase   Category = "database"
	CategoryAI         Category = "ai"
	CategoryUtility    Category = "utility"
)

// SecurityLevel represents the security rating of an MCP server
type SecurityLevel string

const (
	SecurityLevelA SecurityLevel = "A" // Verified, official, well-maintained
	SecurityLevelB SecurityLevel = "B" // Community verified, good reputation
	SecurityLevelC SecurityLevel = "C" // Known author, limited verification
	SecurityLevelD SecurityLevel = "D" // Unverified, use with caution
	SecurityLevelF SecurityLevel = "F" // Blocked, known malicious
)

// ToolInfo represents metadata about a tool provided by an MCP server
type ToolInfo struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Required    []string               `json:"required,omitempty"`
}

// CuratedServer represents a curated MCP server entry
type CuratedServer struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	Description          string            `json:"description"`
	Author               string            `json:"author"`
	Version              string            `json:"version"`
	Repository           string            `json:"repository,omitempty"`
	Documentation        string            `json:"documentation,omitempty"`
	Category             Category          `json:"category"`
	Tags                 []string          `json:"tags,omitempty"`
	SecurityLevel        SecurityLevel     `json:"security_level"`
	Verified             bool              `json:"verified"`
	InstallationRequired []string          `json:"installation_required,omitempty"`
	EnvironmentRequired  []string          `json:"environment_required,omitempty"`
	Tools                []ToolInfo        `json:"tools,omitempty"`
	Transport            string            `json:"transport"`
	URL                  string            `json:"url,omitempty"`
	Command              []string          `json:"command,omitempty"`
	RecommendedConfig    map[string]string `json:"recommended_config,omitempty"`
	SecurityWarnings     []string          `json:"security_warnings,omitempty"`
}

// CuratedRegistry holds the list of curated MCP servers
type CuratedRegistry struct {
	Version string          `json:"version"`
	Updated string          `json:"updated"`
	Servers []CuratedServer `json:"servers"`
}

// DefaultCuratedRegistry returns the built-in curated list of popular MCP servers
func DefaultCuratedRegistry() *CuratedRegistry {
	return &CuratedRegistry{
		Version: "1.0.0",
		Updated: "2026-01-31",
		Servers: []CuratedServer{
			{
				ID:            "filesystem",
				Name:          "Filesystem",
				Description:   "Local filesystem access with configurable permissions",
				Author:        "pryx",
				Version:       "1.0.0",
				Category:      CategoryFilesystem,
				Tags:          []string{"files", "io", "local"},
				SecurityLevel: SecurityLevelA,
				Verified:      true,
				Transport:     "bundled",
				Tools: []ToolInfo{
					{Name: "read_file", Description: "Read contents of a file"},
					{Name: "write_file", Description: "Write contents to a file"},
					{Name: "list_directory", Description: "List contents of a directory"},
					{Name: "search_files", Description: "Search for files matching a pattern"},
				},
			},
			{
				ID:            "shell",
				Name:          "Shell",
				Description:   "Execute shell commands with security controls",
				Author:        "pryx",
				Version:       "1.0.0",
				Category:      CategoryUtility,
				Tags:          []string{"shell", "command", "exec"},
				SecurityLevel: SecurityLevelB,
				Verified:      true,
				Transport:     "bundled",
				Tools: []ToolInfo{
					{Name: "execute", Description: "Execute a shell command"},
				},
				SecurityWarnings: []string{"Can execute arbitrary shell commands"},
			},
			{
				ID:            "browser",
				Name:          "Browser Automation",
				Description:   "Web browser automation using Playwright",
				Author:        "pryx",
				Version:       "1.0.0",
				Category:      CategoryWeb,
				Tags:          []string{"browser", "playwright", "web", "automation"},
				SecurityLevel: SecurityLevelB,
				Verified:      true,
				Transport:     "bundled",
				Tools: []ToolInfo{
					{Name: "navigate", Description: "Navigate to a URL"},
					{Name: "screenshot", Description: "Take a screenshot of the current page"},
					{Name: "click", Description: "Click on an element"},
					{Name: "get_text", Description: "Extract text from the page"},
				},
			},
			{
				ID:            "clipboard",
				Name:          "Clipboard",
				Description:   "Read and write clipboard contents",
				Author:        "pryx",
				Version:       "1.0.0",
				Category:      CategoryUtility,
				Tags:          []string{"clipboard", "copy", "paste"},
				SecurityLevel: SecurityLevelA,
				Verified:      true,
				Transport:     "bundled",
				Tools: []ToolInfo{
					{Name: "read", Description: "Read clipboard contents"},
					{Name: "write", Description: "Write to clipboard"},
				},
			},
			{
				ID:                   "github",
				Name:                 "GitHub",
				Description:          "Interact with GitHub API - repositories, issues, pull requests",
				Author:               "github",
				Version:              "1.0.0",
				Repository:           "https://github.com/github/github-mcp-server",
				Documentation:        "https://github.com/github/github-mcp-server/blob/main/README.md",
				Category:             CategoryWeb,
				Tags:                 []string{"github", "git", "api", "collaboration"},
				SecurityLevel:        SecurityLevelA,
				Verified:             true,
				InstallationRequired: []string{"github-mcp-server"},
				EnvironmentRequired:  []string{"GITHUB_TOKEN"},
				Transport:            "stdio",
				Command:              []string{"github-mcp-server"},
				RecommendedConfig: map[string]string{
					"GITHUB_TOKEN": "Your GitHub personal access token",
				},
				Tools: []ToolInfo{
					{Name: "search_repositories", Description: "Search for repositories on GitHub"},
					{Name: "create_issue", Description: "Create an issue in a repository"},
					{Name: "get_issue", Description: "Get details of a specific issue"},
					{Name: "list_issues", Description: "List issues in a repository"},
					{Name: "create_pull_request", Description: "Create a pull request"},
					{Name: "fork_repository", Description: "Fork a repository"},
					{Name: "create_branch", Description: "Create a new branch"},
					{Name: "search_code", Description: "Search code across GitHub"},
				},
			},
			{
				ID:                   "postgresql",
				Name:                 "PostgreSQL",
				Description:          "Query and manage PostgreSQL databases",
				Author:               "modelcontextprotocol",
				Version:              "1.0.0",
				Repository:           "https://github.com/modelcontextprotocol/servers",
				Documentation:        "https://github.com/modelcontextprotocol/servers/tree/main/src/postgres",
				Category:             CategoryDatabase,
				Tags:                 []string{"postgres", "sql", "database", "postgresql"},
				SecurityLevel:        SecurityLevelA,
				Verified:             true,
				InstallationRequired: []string{"@modelcontextprotocol/server-postgres"},
				EnvironmentRequired:  []string{"POSTGRES_URL"},
				Transport:            "stdio",
				Command:              []string{"npx", "-y", "@modelcontextprotocol/server-postgres"},
				SecurityWarnings:     []string{"Full database access - use with caution"},
				Tools: []ToolInfo{
					{Name: "query", Description: "Execute a read-only SQL query"},
					{Name: "execute", Description: "Execute a SQL statement"},
					{Name: "get_schema", Description: "Get database schema information"},
				},
			},
			{
				ID:                   "sqlite",
				Name:                 "SQLite",
				Description:          "Query and manage SQLite databases",
				Author:               "modelcontextprotocol",
				Version:              "1.0.0",
				Repository:           "https://github.com/modelcontextprotocol/servers",
				Documentation:        "https://github.com/modelcontextprotocol/servers/tree/main/src/sqlite",
				Category:             CategoryDatabase,
				Tags:                 []string{"sqlite", "sql", "database"},
				SecurityLevel:        SecurityLevelA,
				Verified:             true,
				InstallationRequired: []string{"@modelcontextprotocol/server-sqlite"},
				Transport:            "stdio",
				Command:              []string{"npx", "-y", "@modelcontextprotocol/server-sqlite", "/path/to/database.db"},
				Tools: []ToolInfo{
					{Name: "query", Description: "Execute a SQL query"},
					{Name: "get_schema", Description: "Get database schema"},
					{Name: "list_tables", Description: "List all tables in the database"},
				},
			},
			{
				ID:                   "fetch",
				Name:                 "Fetch",
				Description:          "Make HTTP requests to any URL",
				Author:               "modelcontextprotocol",
				Version:              "1.0.0",
				Repository:           "https://github.com/modelcontextprotocol/servers",
				Documentation:        "https://github.com/modelcontextprotocol/servers/tree/main/src/fetch",
				Category:             CategoryWeb,
				Tags:                 []string{"http", "fetch", "api", "request"},
				SecurityLevel:        SecurityLevelB,
				Verified:             true,
				InstallationRequired: []string{"@modelcontextprotocol/server-fetch"},
				Transport:            "stdio",
				Command:              []string{"npx", "-y", "@modelcontextprotocol/server-fetch"},
				SecurityWarnings:     []string{"Can make arbitrary HTTP requests"},
				Tools: []ToolInfo{
					{Name: "fetch", Description: "Fetch content from a URL"},
				},
			},
			{
				ID:                   "brave-search",
				Name:                 "Brave Search",
				Description:          "Web search using Brave Search API",
				Author:               "modelcontextprotocol",
				Version:              "1.0.0",
				Repository:           "https://github.com/modelcontextprotocol/servers",
				Documentation:        "https://github.com/modelcontextprotocol/servers/tree/main/src/brave-search",
				Category:             CategoryWeb,
				Tags:                 []string{"search", "web", "brave"},
				SecurityLevel:        SecurityLevelA,
				Verified:             true,
				InstallationRequired: []string{"@modelcontextprotocol/server-brave-search"},
				EnvironmentRequired:  []string{"BRAVE_API_KEY"},
				Transport:            "stdio",
				Command:              []string{"npx", "-y", "@modelcontextprotocol/server-brave-search"},
				Tools: []ToolInfo{
					{Name: "search", Description: "Perform a web search"},
					{Name: "search_images", Description: "Search for images"},
					{Name: "search_news", Description: "Search for news"},
				},
			},
			{
				ID:                   "puppeteer",
				Name:                 "Puppeteer",
				Description:          "Browser automation using Puppeteer",
				Author:               "modelcontextprotocol",
				Version:              "1.0.0",
				Repository:           "https://github.com/modelcontextprotocol/servers",
				Documentation:        "https://github.com/modelcontextprotocol/servers/tree/main/src/puppeteer",
				Category:             CategoryWeb,
				Tags:                 []string{"browser", "puppeteer", "chrome", "automation"},
				SecurityLevel:        SecurityLevelB,
				Verified:             true,
				InstallationRequired: []string{"@modelcontextprotocol/server-puppeteer"},
				Transport:            "stdio",
				Command:              []string{"npx", "-y", "@modelcontextprotocol/server-puppeteer"},
				Tools: []ToolInfo{
					{Name: "navigate", Description: "Navigate to a URL"},
					{Name: "screenshot", Description: "Take a screenshot"},
					{Name: "evaluate", Description: "Evaluate JavaScript in the page"},
				},
			},
		},
	}
}

// GetByID returns a curated server by its ID
func (r *CuratedRegistry) GetByID(id string) (CuratedServer, bool) {
	for _, server := range r.Servers {
		if server.ID == id {
			return server, true
		}
	}
	return CuratedServer{}, false
}

// GetByCategory returns all curated servers in a given category
func (r *CuratedRegistry) GetByCategory(category Category) []CuratedServer {
	var result []CuratedServer
	for _, server := range r.Servers {
		if server.Category == category {
			result = append(result, server)
		}
	}
	return result
}

// GetAll returns all curated servers
func (r *CuratedRegistry) GetAll() []CuratedServer {
	return r.Servers
}
