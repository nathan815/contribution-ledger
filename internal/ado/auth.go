package ado

import (
	"fmt"
	"os/exec"
	"strings"
)

// Auth handles Azure DevOps authentication via `az` CLI
type Auth struct {
	azPath string
}

// NewAuth creates a new Auth handler
func NewAuth() *Auth {
	return &Auth{azPath: "az"}
}

// EnsureLoggedIn checks if user is logged in, prompts if not
func (a *Auth) EnsureLoggedIn() error {
	// Check if already logged in
	cmd := exec.Command(a.azPath, "account", "show")
	if err := cmd.Run(); err == nil {
		return nil // Already logged in
	}

	fmt.Println("⏳ Authenticating with Azure DevOps...")
	fmt.Println("   A browser window will open to sign in.")
	fmt.Println("   You may need to enter your work account credentials.")

	cmd = exec.Command(a.azPath, "login")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	return nil
}

// GetToken retrieves an ADO token from `az` CLI
func (a *Auth) GetToken() (string, error) {
	// Azure DevOps API token resource ID
	// Reference: https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/service-principal-managed-identity
	const adoResourceID = "499b84ac-1321-427f-aa17-267ca6975798"

	cmd := exec.Command(
		a.azPath, "account", "get-access-token",
		"--resource", adoResourceID,
		"--query", "accessToken",
		"--output", "tsv",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w\n%s", err, string(output))
	}

	token := strings.TrimSpace(string(output))
	if token == "" {
		return "", fmt.Errorf("got empty token from az CLI")
	}

	return token, nil
}

// GetCurrentUser gets the logged-in user's email
func (a *Auth) GetCurrentUser() (string, error) {
	cmd := exec.Command(a.azPath, "account", "show", "--query", "user.name", "--output", "tsv")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	user := strings.TrimSpace(string(output))
	if user == "" {
		return "", fmt.Errorf("could not determine current user")
	}

	return user, nil
}

// CheckAzCLI verifies `az` CLI is installed
func (a *Auth) CheckAzCLI() error {
	cmd := exec.Command(a.azPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("az CLI not found or not working: %w\n\nInstall it from: https://learn.microsoft.com/en-us/cli/azure/install-azure-cli", err)
	}
	return nil
}
