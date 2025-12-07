# Vercel Go SDK

A production-ready Go client library for the [Vercel REST API](https://vercel.com/docs/rest-api). This SDK provides a clean, idiomatic Go interface for interacting with Vercel's API, including projects, deployments, and environment variables.

## Features

- ✅ **Projects**: List and retrieve project information
- ✅ **Deployments**: List, create, get, and cancel deployments
- ✅ **Environment Variables**: List, create, and delete environment variables
- ✅ **Type-safe**: Full type definitions for all API responses
- ✅ **Error handling**: Custom error types with detailed API error information
- ✅ **Context support**: All methods support Go contexts for cancellation and timeouts
- ✅ **Team support**: Optional team ID configuration for team-scoped operations

## Requirements

- Go 1.21 or later
- A Vercel account with an API token

## Installation

```bash
go get github.com/<my-username>/vercel-wrapper
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/<my-username>/vercel-wrapper/vercel"
)

func main() {
    // Create a client with your API token
    client := vercel.New(os.Getenv("VERCEL_TOKEN"))
    
    // Optionally set a team ID
    // client := vercel.New(os.Getenv("VERCEL_TOKEN"), vercel.WithTeamID("team-123"))
    
    ctx := context.Background()
    
    // List projects
    projects, err := client.ListProjects(ctx, 10, 0)
    if err != nil {
        panic(err)
    }
    
    for _, project := range projects.Projects {
        fmt.Printf("Project: %s (ID: %s)\n", project.Name, project.ID)
    }
}
```

## Authentication

The SDK requires a Vercel API token for authentication. You can obtain a token from your [Vercel account settings](https://vercel.com/account/tokens).

Set the token when creating a client:

```go
client := vercel.New("your-api-token")
```

### Team ID

If you're working with a team, you can set the team ID using the `WithTeamID` option:

```go
client := vercel.New("your-api-token", vercel.WithTeamID("team-123"))
```

Alternatively, the team ID can be passed as a query parameter in individual requests, but setting it on the client is more convenient for team-scoped operations.

## Usage Examples

### Projects

```go
// List all projects
projects, err := client.ListProjects(ctx, 10, 0)
if err != nil {
    log.Fatal(err)
}

// Get a specific project
project, err := client.GetProject(ctx, "project-id-or-name")
if err != nil {
    log.Fatal(err)
}
```

### Deployments

```go
// List deployments for a project
deployments, err := client.ListDeployments(ctx, "project-id", 10, 0)
if err != nil {
    log.Fatal(err)
}

// Get a specific deployment
deployment, err := client.GetDeployment(ctx, "deployment-id")
if err != nil {
    log.Fatal(err)
}

// Create a new deployment
req := vercel.CreateDeploymentRequest{
    Name:    "my-deployment",
    Project: "project-id",
    Target:  "production",
    Files: []vercel.DeploymentFile{
        {
            File: "index.html",
            Data: base64.StdEncoding.EncodeToString([]byte("<html>...</html>")),
        },
    },
}
deployment, err := client.CreateDeployment(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Cancel a deployment
err := client.CancelDeployment(ctx, "deployment-id")
if err != nil {
    log.Fatal(err)
}
```

### Environment Variables

```go
// List environment variables for a project
envVars, err := client.ListEnvVars(ctx, "project-id")
if err != nil {
    log.Fatal(err)
}

// Create a new environment variable
req := vercel.CreateEnvVarRequest{
    Key:    "API_KEY",
    Value:  "secret-value",
    Type:   vercel.EnvTypeSecret,
    Target: []vercel.EnvTarget{vercel.EnvTargetProduction},
}
envVar, err := client.CreateEnvVar(ctx, "project-id", req)
if err != nil {
    log.Fatal(err)
}

// Delete an environment variable
err := client.DeleteEnvVar(ctx, "project-id", "env-var-id")
if err != nil {
    log.Fatal(err)
}
```

### Domains

```go
// List all domains for a project
domains, err := client.ListDomains(ctx, "project-id")
if err != nil {
    log.Fatal(err)
}

// Get a specific domain
domain, err := client.GetDomain(ctx, "project-id", "example.com")
if err != nil {
    log.Fatal(err)
}

// Add a domain to a project
req := vercel.CreateDomainRequest{
    Name:      "example.com",
    GitBranch: "main", // optional
}
domain, err := client.CreateDomain(ctx, "project-id", req)
if err != nil {
    log.Fatal(err)
}

// Remove a domain from a project
err := client.DeleteDomain(ctx, "project-id", "example.com")
if err != nil {
    log.Fatal(err)
}
```

## Error Handling

The SDK returns typed errors for API failures. You can check for API errors and inspect their details:

```go
deployment, err := client.GetDeployment(ctx, "deployment-id")
if err != nil {
    if apiErr, ok := vercel.IsAPIError(err); ok {
        fmt.Printf("API Error: %s (code: %s, status: %d)\n",
            apiErr.Message, apiErr.Code, apiErr.StatusCode)
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
    return
}
```

## Client Options

The client supports several configuration options:

```go
// Custom base URL (useful for testing)
client := vercel.New("token", vercel.WithBaseURL("https://custom-api.com"))

// Custom HTTP client
httpClient := &http.Client{Timeout: 60 * time.Second}
client := vercel.New("token", vercel.WithHTTPClient(httpClient))

// Multiple options
client := vercel.New("token",
    vercel.WithTeamID("team-123"),
    vercel.WithBaseURL("https://api.vercel.com"),
    vercel.WithHTTPClient(customClient),
)
```

## Example CLI

The repository includes an example CLI application in `cmd/example/main.go`:

```bash
export VERCEL_TOKEN=your-token
export VERCEL_TEAM_ID=your-team-id  # optional

go run ./cmd/example
go run ./cmd/example -project project-id
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with verbose output to see logged responses:

```bash
go test -v ./vercel/...
```

The tests use `httptest` to mock the Vercel API, so no actual API calls are made during testing. The `TestAllEndpoints_WithLogging` test logs all cleaned responses in formatted JSON for easy debugging and verification.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This is not an official Vercel SDK. It is a community-maintained wrapper around the Vercel REST API. For official Vercel tools and SDKs, please visit [vercel.com](https://vercel.com).

## API Coverage

This SDK currently supports:

- ✅ Projects (list, get)
- ✅ Deployments (list, get, create, cancel)
- ✅ Environment Variables (list, create, delete)
- ✅ Domains (list, get, create, delete)

Additional endpoints can be added as needed. The client architecture makes it easy to extend with new API methods.

