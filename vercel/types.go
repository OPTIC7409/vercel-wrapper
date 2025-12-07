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
		Type  string `json:"type"`
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
	Name string `json:"name"`
	GitBranch string `json:"gitBranch,omitempty"`
}
