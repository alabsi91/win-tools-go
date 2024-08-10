package commands

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

func SetRegistry() {

	// Create a new multi-select prompt
	var selected []string
	huh.NewMultiSelect[string]().
		Title("\nSelect the registry you want to modify").
		Options(
			huh.NewOption("Enable old Windows 10 context menu", "EnableWin10Context.reg"),
			huh.NewOption("Disable old Windows 10 context menu", "DisableWin10Context.reg"),

			huh.NewOption("Disable Copilot", "DisableCopilot.reg"),
			huh.NewOption("Enable Copilot", "EnableCopilot.reg"),

			huh.NewOption("Disable Telemetry", "DisableTelemetry.reg"),
			huh.NewOption("Enable Telemetry", "EnableTelemetry.reg"),

			huh.NewOption("Disable Taskbar Widgets", "DisableWidgetsTaskbar.reg"),
			huh.NewOption("Enable Taskbar Widgets", "EnableWidgetsTaskbar.reg"),

			huh.NewOption("Disable Windows Suggestions", "DisableWindowsSuggestions.reg"),
			huh.NewOption("Enable Windows Suggestions", "EnableWindowsSuggestions.reg"),

			huh.NewOption("Show Extensions for Known File Types", "ShowExtensionsForKnownFileTypes.reg"),

			huh.NewOption("Show Hidden Folders", "ShowHiddenFolders.reg"),

			huh.NewOption("Disable Mouse Enhance Pointer Precision", "DisableEnhancePointerPrecision.reg"),

			huh.NewOption("Enable dark mode", "EnableDarkMode.reg"),
			huh.NewOption("Enable light mode", "EnableLightMode.reg"),
		).
		Value(&selected).Run()

	if len(selected) == 0 {
		Log.Warning("\nNo registry selected\n")
		return
	}

	println("")
	for _, registry := range selected {
		Log.Info(fmt.Sprintf(`Setting registry: "%s"`, registry))

		regPath := filepath.Join(AssetsPath, registry)

		cmd := exec.Command("cmd", "/C", "regedit.exe", "/s", regPath)

		_, err := cmd.Output()
		if err != nil {
			Log.Error("Failed to set registry")
			return
		}
	}

	// restart explorer
	Log.Info("\nRestarting Windows Explorer...")
	Powershell.RestartWinExplorer()

	Log.Success("\nDone!\n")
}
