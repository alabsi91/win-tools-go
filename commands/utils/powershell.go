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

func (powershell *powershell) GetShellPath() string {

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

func (powershell *powershell) SetEnvVariable(key string, value string, scope string) {
	shellPath := powershell.GetShellPath()

	// Make sure the user has admin privileges when using the scope "Machine"
	if scope == EnvironmentScope.Machine && !powershell.IsAdmin() {
		Log.Error(
			fmt.Sprintf(`You dont have enough privileges to set the system environment variable "%s" with the value "%s".`, key, value),
		)
		return
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
			Log.Fatal(
				fmt.Sprintf(`Failed to set the environment variable "%s" with the value "%s".`, key, value),
			)
		}

		return
	}

	// Add key value
	cmd := exec.Command(
		shellPath,
		"-Command",
		fmt.Sprintf(`[System.Environment]::SetEnvironmentVariable("%s", "%s", "%s")`, key, value, scope),
	)

	_, err := cmd.Output()
	if err != nil {
		Log.Fatal(
			fmt.Sprintf(`Failed to set the environment variable "%s" with the value "%s".`, key, value),
		)
	}
}

func (powershell *powershell) IsAdmin() bool {
	shellPath := powershell.GetShellPath()

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

func (powershell *powershell) RemoveWinPackage(packageName string) {

	shellPath := powershell.GetShellPath()

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
		Log.Fatal(
			fmt.Sprintf(`Failed to remove the package with the name "%s".`, packageName),
		)
	}

	if err := cmd.Wait(); err != nil {
		Log.Fatal(
			fmt.Sprintf(`Failed to remove the package with the name "%s".`, packageName),
		)
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
		Log.Fatal(
			fmt.Sprintf(`Failed to remove the package with the name "%s".`, packageName),
		)
	}

	if err := cmd.Wait(); err != nil {
		Log.Fatal(
			fmt.Sprintf(`Failed to remove the package with the name "%s".`, packageName),
		)
	}
}

func (powershell *powershell) RestartWinExplorer() {
	shellPath := powershell.GetShellPath()

	cmd := exec.Command(shellPath, "-Command", "stop-process", "-name", "explorer", "–force")

	_, err := cmd.Output()

	if err != nil {
		Log.Warning("Failed to restart Windows Explorer.")
	}
}

func (powershell *powershell) RunPathThroughCmd(args ...string) error {
	shell := powershell.GetShellPath()

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
