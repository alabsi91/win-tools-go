package app

import (
	"fmt"
	"os"
	"os/exec"
)

const ChocolateyInstallPath = "C:\\ProgramData\\chocolatey\\bin\\choco.exe"

type chocolatey struct{}

var Chocolatey = &chocolatey{}

func (chocolatey) IsChocolateyInstalled() bool {
	cmd := exec.Command("cmd", "/C", "choco", "--version")

	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return true
	}

	return Utils.IsPathExists(ChocolateyInstallPath)
}

func (chocolatey) GetChocolateyPath() string {
	// first try to the default path
	if Utils.IsPathExists(ChocolateyInstallPath) {
		return ChocolateyInstallPath
	}

	// try to get the path from where chocolatey is installed
	cmd := exec.Command("cmd", "/C", "where", "choco")

	output, err := cmd.Output()
	if err != nil {
		Log.Fatal("Failed to get chocolatey executable path")
		os.Exit(1)
	}

	// verify if the path exists
	path := string(output)

	if !Utils.IsPathExists(path) {
		Log.Fatal("Failed to get chocolatey executable path")
		os.Exit(1)
	}

	return path
}

func (chocolatey) InstallChocolatey() {
	powershell := Powershell.GetShellPath()

	cmd := exec.Command(
		powershell,
		"-Command",
		"Set-ExecutionPolicy Bypass -Scope Process -Force;",
		"[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072;",
		"iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))",
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		Log.Fatal(err.Error())
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		Log.Fatal(err.Error())
		os.Exit(1)
	}
}

func (chocolatey *chocolatey) InstallChocolateyPackage(packageName string) {
	chocolateyPath := chocolatey.GetChocolateyPath()
	powershell := Powershell.GetShellPath()

	cmd := exec.Command(
		powershell,
		"-Command",
		fmt.Sprintf(`&"%s"`, chocolateyPath),
		"install",
		packageName,
		"-yf",
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		Log.Fatal("Failed to install chocolatey package:", packageName)
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		Log.Fatal("Failed to install chocolatey package:", packageName)
		os.Exit(1)
	}
}
