package main

import (
	"context"
	"fmt"
	"os"

	"github.com/OPTIC7409/vercel-wrapper/vercel"
)

func main() {
	token := os.Getenv("VERCEL_TOKEN")
	teamID := os.Getenv("VERCEL_TEAM_ID")

	ctx := context.Background()

	client := vercel.New(token, vercel.WithTeamID(teamID))

	// List all projects
	projects, err := client.ListProjects(ctx, 100, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing projects: %v\n", err)
		os.Exit(1)
	}

	// List all domains across all projects
	fmt.Println("=== All Domains ===")
	totalDomains := 0

	for _, project := range projects.Projects {
		domains, err := client.ListDomains(ctx, project.ID)
		if err != nil {
			// Skip projects that fail (might not have domains endpoint access)
			continue
		}

		for _, domain := range domains {
			verified := "No"
			if domain.Verified {
				verified = "Yes"
			}
			fmt.Printf("Project: %s | Domain: %s | Verified: %s\n", project.Name, domain.Name, verified)
			totalDomains++
		}
	}

	if totalDomains == 0 {
		fmt.Println("No domains found")
	} else {
		fmt.Printf("\nTotal domains: %d\n", totalDomains)
	}
}
