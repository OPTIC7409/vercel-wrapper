package vercel

// Project represents a Vercel project.
type Project struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Link         interface{} `json:"link,omitempty"`
	Framework    string      `json:"framework,omitempty"`
	CreatedAt    int64       `json:"createdAt,omitempty"`
	UpdatedAt    int64       `json:"updatedAt,omitempty"`
	AccountID    string      `json:"accountId,omitempty"`
	TeamID       string      `json:"teamId,omitempty"`
	PublicSource interface{} `json:"publicSource,omitempty"`
}

// ListProjectsResponse represents the response from listing projects.
type ListProjectsResponse struct {
	Projects   []Project `json:"projects"`
	Pagination struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	} `json:"pagination"`
}

// Deployment represents a Vercel deployment.
type Deployment struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	State      string `json:"state"`
	Target     string `json:"target"`
	CreatedAt  int64  `json:"createdAt"`
	ReadyAt    int64  `json:"readyAt,omitempty"`
	BuildingAt int64  `json:"buildingAt,omitempty"`
	ProjectID  string `json:"projectId,omitempty"`
}

// ListDeploymentsResponse represents the response from listing deployments.
type ListDeploymentsResponse struct {
	Deployments []Deployment `json:"deployments"`
	Pagination  struct {
		Count int `json:"count"`
		Limit int `json:"limit"`
		Next  int `json:"next"`
		Prev  int `json:"prev"`
	} `json:"pagination"`
}

// DeploymentFile represents a file in a deployment.
type DeploymentFile struct {
	File string `json:"file"` // path
	Data string `json:"data"` // base64-encoded file
}

// CreateDeploymentRequest represents a request to create a deployment.
type CreateDeploymentRequest struct {
	Name    string            `json:"name"`
	Project string            `json:"project,omitempty"`
	Target  string            `json:"target,omitempty"` // "production" or "staging"
	Files   []DeploymentFile  `json:"files,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// EnvType represents the type of an environment variable.
type EnvType string

const (
	EnvTypePlain  EnvType = "plain"
	EnvTypeSecret EnvType = "secret"
	EnvTypeSystem EnvType = "system"
)

// EnvTarget represents the target environment for an environment variable.
type EnvTarget string

const (
	EnvTargetProduction  EnvTarget = "production"
	EnvTargetPreview     EnvTarget = "preview"
	EnvTargetDevelopment EnvTarget = "development"
)

// EnvVar represents a Vercel environment variable.
type EnvVar struct {
	ID        string      `json:"id"`
	Key       string      `json:"key"`
	Value     string      `json:"value,omitempty"`
	Type      EnvType     `json:"type"`
	Target    []EnvTarget `json:"target"`
	CreatedAt int64       `json:"createdAt,omitempty"`
	UpdatedAt int64       `json:"updatedAt,omitempty"`
}

// CreateEnvVarRequest represents a request to create an environment variable.
type CreateEnvVarRequest struct {
	Key    string      `json:"key"`
	Value  string      `json:"value"`
	Type   EnvType     `json:"type"`
	Target []EnvTarget `json:"target"`
}

// Domain represents a Vercel domain.
type Domain struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	ServiceType  string   `json:"serviceType,omitempty"`
	Nameservers  []string `json:"nameservers,omitempty"`
	Intent       string   `json:"intent,omitempty"`
	CreatedAt    int64    `json:"createdAt,omitempty"`
	UpdatedAt    int64    `json:"updatedAt,omitempty"`
	Verified     bool     `json:"verified,omitempty"`
	Verification []struct {
		Type   string `json:"type"`
		Domain string `json:"domain"`
		Value  string `json:"value"`
	} `json:"verification,omitempty"`
	ConfigVerifiedAt int64  `json:"configVerifiedAt,omitempty"`
	CDNEnabled       bool   `json:"cdnEnabled,omitempty"`
	GitBranch        string `json:"gitBranch,omitempty"`
	ProjectID        string `json:"projectId,omitempty"`
}

// CreateDomainRequest represents a request to create/add a domain to a project.
type CreateDomainRequest struct {
	Name      string `json:"name"`
	GitBranch string `json:"gitBranch,omitempty"`
}

// UpdateProjectRequest represents a request to update a project.
type UpdateProjectRequest struct {
	Name            string `json:"name,omitempty"`
	Framework       string `json:"framework,omitempty"`
	BuildCommand    string `json:"buildCommand,omitempty"`
	DevCommand      string `json:"devCommand,omitempty"`
	InstallCommand  string `json:"installCommand,omitempty"`
	OutputDirectory string `json:"outputDirectory,omitempty"`
	RootDirectory   string `json:"rootDirectory,omitempty"`
	PublicSource    *bool  `json:"publicSource,omitempty"`
	GitRepository   *struct {
		Type string `json:"type"`
		Repo string `json:"repo"`
	} `json:"gitRepository,omitempty"`
}

// UpdateEnvVarRequest represents a request to update an environment variable.
type UpdateEnvVarRequest struct {
	Value     string      `json:"value,omitempty"`
	Target    []EnvTarget `json:"target,omitempty"`
	GitBranch string      `json:"gitBranch,omitempty"`
}

// DeploymentLog represents a log entry from a deployment.
type DeploymentLog struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	Type      string `json:"type,omitempty"`
}

// DeploymentLogsResponse represents the response from getting deployment logs.
type DeploymentLogsResponse struct {
	Logs []DeploymentLog `json:"logs"`
}

// Team represents a Vercel team.
type Team struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Avatar     string `json:"avatar,omitempty"`
	CreatedAt  int64  `json:"createdAt,omitempty"`
	UpdatedAt  int64  `json:"updatedAt,omitempty"`
	Membership *struct {
		Role string `json:"role"`
	} `json:"membership,omitempty"`
}

// ListTeamsResponse represents the response from listing teams.
type ListTeamsResponse struct {
	Teams []Team `json:"teams"`
}

// TeamMember represents a member of a team.
type TeamMember struct {
	User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name,omitempty"`
		Email    string `json:"email,omitempty"`
		Avatar   string `json:"avatar,omitempty"`
	} `json:"user"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"createdAt,omitempty"`
}

// ListTeamMembersResponse represents the response from listing team members.
type ListTeamMembersResponse struct {
	Members []TeamMember `json:"members"`
}

// Alias represents a Vercel alias.
type Alias struct {
	ID         string `json:"id"`
	Alias      string `json:"alias"`
	Deployment *struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"deployment,omitempty"`
	ProjectID string  `json:"projectId,omitempty"`
	Domain    string  `json:"domain,omitempty"`
	Target    string  `json:"target,omitempty"`
	Redirect  *string `json:"redirect,omitempty"`
	CreatedAt int64   `json:"createdAt,omitempty"`
	UpdatedAt int64   `json:"updatedAt,omitempty"`
}

// ListAliasesResponse represents the response from listing aliases.
type ListAliasesResponse struct {
	Aliases    []Alias `json:"aliases"`
	Pagination struct {
		Count int `json:"count"`
		Next  int `json:"next,omitempty"`
		Prev  int `json:"prev,omitempty"`
	} `json:"pagination"`
}

// CreateAliasRequest represents a request to create an alias.
type CreateAliasRequest struct {
	Alias      string `json:"alias"`
	Deployment string `json:"deploymentId,omitempty"`
	Redirect   string `json:"redirect,omitempty"`
}

// Secret represents a Vercel secret.
type Secret struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Value      string   `json:"value,omitempty"` // Only returned when creating
	TeamID     string   `json:"teamId,omitempty"`
	UserID     string   `json:"userId,omitempty"`
	ProjectIDs []string `json:"projectIds,omitempty"`
	CreatedAt  int64    `json:"createdAt,omitempty"`
	UpdatedAt  int64    `json:"updatedAt,omitempty"`
}

// ListSecretsResponse represents the response from listing secrets.
type ListSecretsResponse struct {
	Secrets []Secret `json:"secrets"`
}

// CreateSecretRequest represents a request to create a secret.
type CreateSecretRequest struct {
	Name       string   `json:"name"`
	Value      string   `json:"value"`
	TeamID     string   `json:"teamId,omitempty"`
	ProjectIDs []string `json:"projectIds,omitempty"`
}
