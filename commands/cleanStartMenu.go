package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CleanStartMenu() {
	menuTemplatePath := filepath.Join(AssetsPath, "start2.bin")
	targetPath := fmt.Sprintf(`C:\Users\%s\AppData\Local\Packages\Microsoft.Windows.StartMenuExperienceHost_cw5n1h2txyewy\LocalState`, os.Getenv("USERNAME"))
	targetPath = filepath.Clean(targetPath)

	Log.Info("\nCleaning the start menu ...\n")

	cmd := exec.Command(
		Powershell.GetShellPath(),
		"-Command",
		"Copy-Item",
		"-Path", fmt.Sprintf(`"%s"`, menuTemplatePath),
		"-Destination", fmt.Sprintf(`"%s"`, targetPath),
		"-Force",
	)

	_, err := cmd.Output()
	if err != nil {
		Log.Error("Failed to clean the start menu")
		return
	}

	// restart explorer
	Powershell.RestartWinExplorer()

	Log.Success("\nDone!\n")
}
