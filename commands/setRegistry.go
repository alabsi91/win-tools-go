package commands

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/alabsi91/win-tools/commands/utils"
	"github.com/charmbracelet/huh"
)

// askToSelectRegistry prompts the user to select the registry they want to modify
//   - Returns an error if the user cancels the prompt
func askToSelectRegistry() ([]string, error) {
	var selected []string

	var other []string
	var contextMenu []string
	var taskbar []string
	var explorer []string

	err := huh.NewForm(
		// Windows
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select the registry you want to modify").
				Description("Page 1 of 4").
				Options(
					huh.NewOption("Disable Copilot", "DisableCopilot.reg"),
					huh.NewOption("Enable Copilot", "EnableCopilot.reg"),

					huh.NewOption("Disable AI Recall", "DisableAIRecall.reg"),
					huh.NewOption("Enable AI Recall", "EnableAIRecall.reg"),

					huh.NewOption("Disable Telemetry", "DisableTelemetry.reg"),
					huh.NewOption("Enable Telemetry", "EnableTelemetry.reg"),

					huh.NewOption("Disable Windows Suggestions", "DisableWindowsSuggestions.reg"),
					huh.NewOption("Enable Windows Suggestions", "EnableWindowsSuggestions.reg"),

					huh.NewOption("Enable dark mode", "EnableDarkMode.reg"),
					huh.NewOption("Enable light mode", "EnableLightMode.reg"),

					huh.NewOption("Disable Bing Cortana In Search", "DisableBingCortanaInSearch.reg"),
					huh.NewOption("Enable Bing Cortana In Search", "EnableBingCortanaInSearch.reg"),

					huh.NewOption("Disable DVR", "DisableDVR.reg"),
					huh.NewOption("Enable DVR", "EnableDVR.reg"),

					huh.NewOption("Disable Lock screen Tips", "DisableLockscreenTips.reg"),
					huh.NewOption("Enable Lock screen Tips", "EnableLockscreenTips.reg"),

					huh.NewOption("Disable Mouse Enhance Pointer Precision", "DisableEnhancePointerPrecision.reg"),
				).
				Value(&other),
		),

		// Context Menu
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Context Menu").
				Description("Page 2 of 4").
				Options(
					huh.NewOption("Enable old Windows 10 context menu", "EnableWin10Context.reg"),
					huh.NewOption("Disable old Windows 10 context menu", "DisableWin10Context.reg"),

					huh.NewOption("Disable Give access to context menu", "DisableGiveaccesstocontextmenu.reg"),
					huh.NewOption("Enable Give access to context menu", "EnableGiveaccesstocontextmenu.reg"),

					huh.NewOption("Disable Include in library from context menu", "DisableIncludeinlibraryfromcontextmenu.reg"),
					huh.NewOption("Enable Include in library to context menu", "EnableIncludeinlibrarytocontextmenu.reg"),

					huh.NewOption("Disable Share from context menu", "DisableSharefromcontextmenu.reg"),
					huh.NewOption("Enable Share to context menu", "EnableSharetocontextmenu.reg"),

					huh.NewOption("Disable Show More Options Context Menu", "DisableShowMoreOptionsContextMenu.reg"),
					huh.NewOption("Enable Show More Options Context Menu", "EnableShowMoreOptionsContextMenu.reg"),
				).
				Value(&contextMenu),
		),

		// Taskbar
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Taskbar").
				Description("Page 3 of 4").
				Options(
					huh.NewOption("Hide Search Taskbar", "HideSearchTaskbar.reg"),
					huh.NewOption("Show Search Box", "ShowSearchBox.reg"),
					huh.NewOption("Show Search Icon", "ShowSearchIcon.reg"),
					huh.NewOption("Show Search Icon And Label", "ShowSearchIconAndLabel.reg"),

					huh.NewOption("Hide Task view Taskbar", "HideTaskviewTaskbar.reg"),
					huh.NewOption("Show Task view Taskbar", "ShowTaskviewTaskbar.reg"),

					huh.NewOption("Disable Taskbar Widgets", "DisableWidgetsTaskbar.reg"),
					huh.NewOption("Enable Taskbar Widgets", "EnableWidgetsTaskbar.reg"),

					huh.NewOption("Align Taskbar Left", "AlignTaskbarLeft.reg"),
					huh.NewOption("Align Taskbar Center", "AlignTaskbarCenter.reg"),

					huh.NewOption("Disable Chat Taskbar", "DisableChatTaskbar.reg"),
					huh.NewOption("Enable Chat Taskbar", "EnableChatTaskbar.reg"),
				).
				Value(&taskbar),
		),

		// Explorer
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Explorer").
				Description("Page 4 of 4").
				Options(
					huh.NewOption("Hide duplicate removable drives from navigation pane of File Explorer", "HideduplicateremovabledrivesfromnavigationpaneofFileExplorer.reg"),
					huh.NewOption("Show duplicate removable drives from navigation pane of File Explorer", "ShowduplicateremovabledrivesfromnavigationpaneofFileExplorer.reg"),

					huh.NewOption("Hide Extensions For Known File Types", "HideExtensionsForKnownFileTypes.reg"),
					huh.NewOption("Show Extensions for Known File Types", "ShowExtensionsForKnownFileTypes.reg"),

					huh.NewOption("Show Hidden Folders", "ShowHiddenFolders.reg"),
					huh.NewOption("Hide Hidden Folders", "HideHiddenFolders.reg"),

					huh.NewOption("Hide 3DObjects Folder", "Hide3DObjectsFolder.reg"),
					huh.NewOption("Show 3DObjects Folder", "Show3DObjectsFolder.reg"),

					huh.NewOption("Hide Gallery from Explorer", "HideGalleryfromExplorer.reg"),
					huh.NewOption("Show Gallery in Explorer", "ShowGalleryinExplorer.reg"),

					huh.NewOption("Hide Music Folder", "HideMusicFolder.reg"),
					huh.NewOption("Show Music Folder", "ShowMusicFolder.reg"),

					huh.NewOption("Hide One drive Folder", "HideOnedriveFolder.reg"),
					huh.NewOption("Show One drive folder", "ShowOnedrivefolder.reg"),
				).Value(&explorer),
		),
	).
		Run()

	selected = append(selected, other...)
	selected = append(selected, contextMenu...)
	selected = append(selected, taskbar...)
	selected = append(selected, explorer...)

	return selected, err
}

func SetRegistry() {
	selected, err := askToSelectRegistry()

	if err != nil {
		Log.Error("\nfailed to get user selection\n")
		return
	}

	if len(selected) == 0 {
		Log.Warning("\nNo registry selected\n")
		return
	}

	println("")
	for _, registry := range selected {
		Log.Info(fmt.Sprintf(`Setting registry: "%s"`, registry))

		regPath := filepath.Join(AssetsPath, "RegFiles", registry)

		if !utils.IsPathExists(regPath) {
			Log.Error("\ncould not find registry file: ", regPath, "\n")
			continue
		}

		cmd := exec.Command("cmd", "/C", "regedit.exe", "/s", regPath)

		_, err := cmd.Output()
		if err != nil {
			Log.Error("failed to set registry")
			return
		}
	}

	// restart explorer
	Log.Info("\nRestarting Windows Explorer...")
	Powershell.RestartWinExplorer()

	Log.Success("\nDone!\n")
}
