package commands

import (
	"fmt"
	"os"
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

	err := Powershell.RunPathThroughCmd(
		"Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process -Force;",
		fmt.Sprintf(`&"%s"`, scriptPath),
	)

	if err != nil {
		Log.Fatal("\n"+err.Error(), "\n")
		os.Exit(1)
	}

	Log.Success("\nDone!\n")
}
