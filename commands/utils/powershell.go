package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type powershell struct {
	shellPath *string
}

type environmentScope struct {
	User    string
	Machine string
}

var Powershell = &powershell{}

var EnvironmentScope = environmentScope{"User", "Machine"}

// GetShellName returns the powershell executable name
//   - First will try if "pwsh" exists (Powershell version 7.1+)
//   - If not, will try "powershell"
//   - If neither exists, will exit with a fatal error
//   - The result will be cached, so it will only be executed once per session
//
// Returns: "pwsh" or "powershell"
func (powershell *powershell) GetShellName() string {

	// try get from cache
	if powershell.shellPath != nil {
		return *powershell.shellPath
	}

	// first try "pwsh"
	cmd := exec.Command("cmd", "/C", "where", "pwsh")
	_, err := cmd.Output()
	if err == nil {
		shellPath := "pwsh"
		powershell.shellPath = &shellPath
		return shellPath
	}

	// try "powershell"
	cmd = exec.Command("cmd", "/C", "where", "powershell")
	_, err = cmd.Output()
	if err == nil {
		shellPath := "powershell"
		powershell.shellPath = &shellPath
		return shellPath
	}

	// Failure
	Log.Fatal("\nFailed to get Powershell executable path\n")

	return "powershell"
}

// SetEnvVariable sets an environment variable in the given scope
//   - key: name of the environment variable
//   - value: value of the environment variable
//   - scope: "User" or "Machine"
//
// Returns: error if any
func (powershell *powershell) SetEnvVariable(key string, value string, scope string) error {
	shellPath := powershell.GetShellName()

	// Make sure the user has admin privileges when using the scope "Machine"
	if scope == EnvironmentScope.Machine && !powershell.IsAdmin() {
		return fmt.Errorf(`you dont have enough privileges to set the system environment variable "%s" with the value "%s"`, key, value)
	}

	// Add to a new path
	if key == "PATH" {
		cmd := exec.Command(
			shellPath,
			"-Command",
			fmt.Sprintf(`$tempPathVar = [System.Environment]::GetEnvironmentVariable("PATH", "%s");`, scope),
			fmt.Sprintf(`[System.Environment]::SetEnvironmentVariable("PATH",  $tempPathVar + ";%s", "%s")`, value, scope),
		)

		_, err := cmd.Output()
		if err != nil {
			return fmt.Errorf(`failed to set the environment variable "%s" with the value "%s"`, key, value)
		}

		return nil
	}

	// Add key value
	cmd := exec.Command(
		shellPath,
		"-Command",
		fmt.Sprintf(`[System.Environment]::SetEnvironmentVariable("%s", "%s", "%s")`, key, value, scope),
	)

	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf(`failed to set the environment variable "%s" with the value "%s"`, key, value)
	}

	return nil
}

// IsAdmin checks if the current user has admin privileges.
//
// returns false when encountering an error.
//
// Returns: true if the user has admin privileges
func (powershell *powershell) IsAdmin() bool {
	shellPath := powershell.GetShellName()

	cmd := exec.Command(
		shellPath,
		"-Command",
		`(New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)`,
	)

	output, err := cmd.Output()

	if err != nil {
		return false
	}

	outputStr := strings.Trim(string(output), "\r\n ")
	return outputStr == "True"
}

// RemoveWinPackage removes a Windows package (bloatware)
func (powershell *powershell) RemoveWinPackage(packageName string) error {
	shellPath := powershell.GetShellName()

	cmd := exec.Command(
		shellPath,
		"-Command",
		"Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process -Force;",
		fmt.Sprintf(`Get-AppxPackage -Name "%s" -AllUsers | Remove-AppxPackage`, packageName),
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf(`failed to remove the package with the name "%s"`, packageName)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf(`failed to remove the package with the name "%s"`, packageName)
	}

	cmd = exec.Command(
		shellPath,
		"-Command",
		"Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process -Force;",
		fmt.Sprintf(`Get-AppxProvisionedPackage -Online | Where-Object { $_.PackageName -like "%s" } | ForEach-Object { Remove-ProvisionedAppxPackage -Online -AllUsers -PackageName $_.PackageName }`, packageName),
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf(`failed to remove the package with the name "%s"`, packageName)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf(`failed to remove the package with the name "%s"`, packageName)
	}

	return nil
}

// RestartWinExplorer restarts Windows Explorer
//   - Logs a warning if it fails
func (powershell *powershell) RestartWinExplorer() {
	shellPath := powershell.GetShellName()

	cmd := exec.Command(shellPath, "-Command", "stop-process", "-name", "explorer", "–force")

	_, err := cmd.Output()

	if err != nil {
		Log.Warning("Failed to restart Windows Explorer.")
	}
}

// RunPathThroughCmd runs a powershell command and streams the output to the console
func (powershell *powershell) RunPathThroughCmd(args ...string) error {
	shell := powershell.GetShellName()

	cmd := exec.Command(shell, append([]string{"-Command"}, args...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
