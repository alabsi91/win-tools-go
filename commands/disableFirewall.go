package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func DisableFirewall() {
	// has admin privileges
	isAdmin := Powershell.IsAdmin()

	if !isAdmin {
		Log.Error("\nyou need admin privileges to run this command\n")
		Log.Info("Please run this command from an elevated powershell session\n")
		return
	}

	scriptPath := filepath.Join(AssetsPath, "disableFirewall.ps1")

	shell := Powershell.GetShellPath()

	cmd := exec.Command(
		shell,
		"-Command",
		"Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process -Force;",
		fmt.Sprintf(`&"%s"`, scriptPath),
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

	Log.Success("\nDone!\n")
}
