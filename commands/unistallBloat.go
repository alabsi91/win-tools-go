package commands

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

func askToSelectBloatware() ([]string, error) {
	var selected []string

	err := huh.NewMultiSelect[string]().
		Title("\nSelect the bloatware you want to uninstall").
		Options(
			huh.NewOption("Edge browser", "Microsoft.Edge"),
			huh.NewOption("OneDrive", "Microsoft.OneDrive"),
			huh.NewOption("Clipchamp", "Clipchamp.Clipchamp"),
			huh.NewOption("3D Builder", "Microsoft.3DBuilder"),
			huh.NewOption("Finance", "Microsoft.BingFinance"),
			huh.NewOption("Food and Drink", "Microsoft.BingFoodAndDrink"),
			huh.NewOption("Health and Fitness", "Microsoft.BingHealthAndFitness"),
			huh.NewOption("News", "Microsoft.News"),
			huh.NewOption("Bing News", "Microsoft.BingNews"),
			huh.NewOption("Sports", "Microsoft.BingSports"),
			huh.NewOption("Translator", "Microsoft.BingTranslator"),
			huh.NewOption("Travel", "Microsoft.BingTravel"),
			huh.NewOption("Weather", "Microsoft.BingWeather"),
			huh.NewOption("Messaging", "Microsoft.Messaging"),
			huh.NewOption("3D Viewer", "Microsoft.Microsoft3DViewer"),
			huh.NewOption("Office Hub", "Microsoft.MicrosoftOfficeHub"),
			huh.NewOption("Power BI", "Microsoft.MicrosoftPowerBIForWindows"),
			huh.NewOption("Solitaire Collection", "Microsoft.MicrosoftSolitaireCollection"),
			huh.NewOption("Sticky Notes", "Microsoft.MicrosoftStickyNotes"),
			huh.NewOption("Mixed Reality Portal", "Microsoft.MixedReality.Portal"),
			huh.NewOption("Network Speed Test", "Microsoft.NetworkSpeedTest"),
			huh.NewOption("OneNote", "Microsoft.Office.OneNote"),
			huh.NewOption("Sway", "Microsoft.Office.Sway"),
			huh.NewOption("OneConnect", "Microsoft.OneConnect"),
			huh.NewOption("Print 3D", "Microsoft.Print3D"),
			huh.NewOption("Skype", "Microsoft.SkypeApp"),
			huh.NewOption("To-Do", "Microsoft.Todos"),
			huh.NewOption("Alarms", "Microsoft.WindowsAlarms"),
			huh.NewOption("Feedback Hub", "Microsoft.WindowsFeedbackHub"),
			huh.NewOption("Maps", "Microsoft.WindowsMaps"),
			huh.NewOption("Sound Recorder", "Microsoft.WindowsSoundRecorder"),
			huh.NewOption("Movies & TV", "Microsoft.ZuneVideo"),
			huh.NewOption("Family", "MicrosoftCorporationII.MicrosoftFamily"),
			huh.NewOption("Teams", "MicrosoftTeams"),
			huh.NewOption("Get Help", "Microsoft.GetHelp"),
			huh.NewOption("MS Paint", "Microsoft.MSPaint"),
			huh.NewOption("Paint", "Microsoft.Paint"),
			huh.NewOption("Whiteboard", "Microsoft.Whiteboard"),
			huh.NewOption("Photos", "Microsoft.Windows.Photos"),
			huh.NewOption("Calculator", "Microsoft.WindowsCalculator"),
			huh.NewOption("Camera", "Microsoft.WindowsCamera"),
			huh.NewOption("Your Phone", "Microsoft.YourPhone"),
			huh.NewOption("Music", "Microsoft.ZuneMusic"),
			huh.NewOption("Gaming App", "Microsoft.GamingApp"),
			huh.NewOption("Outlook", "Microsoft.OutlookForWindows"),
			huh.NewOption("People", "Microsoft.People"),
			huh.NewOption("Power Automate Desktop", "Microsoft.PowerAutomateDesktop"),
			huh.NewOption("Mail and Calendar", "Microsoft.windowscommunicationsapps"),
			huh.NewOption("Xbox Game Overlay", "Microsoft.XboxGameOverlay"),
			huh.NewOption("Xbox Gaming Overlay", "Microsoft.XboxGamingOverlay"),
			huh.NewOption("Dev Home", "Windows.DevHome"),
		).
		Value(&selected).
		Run()

	return selected, err
}

func UninstallBloat() {
	// has admin privileges
	isAdmin := Powershell.IsAdmin()

	if !isAdmin {
		Log.Error("\nyou need admin privileges to run this command\n")
		Log.Info("Please run this command from an elevated powershell session\n")
		return
	}

	// ask user to select bloatware
	selected, err := askToSelectBloatware()

	if err != nil {
		Log.Error("\nFailed to get user selection\n")
		return
	}

	if len(selected) == 0 {
		Log.Warning("\nNo bloatware selected\n")
		return
	}

	// loop through selected options
	println("")
	for _, option := range selected {
		Log.Info(fmt.Sprintf(`Uninstalling "%s" ...`, option))

		// special case for Microsoft.Edge
		if option == "Microsoft.Edge" {
			Log.Warning("\nUninstalling Microsoft Edge will uninstall the browser and keep the web view.\nSome Microsoft apps like copilot will not work and web search in taskbar search will not work.\n")

			removeEdgeExePath := filepath.Join(AssetsPath, "RemoveEdgeOnly.exe")

			err := Powershell.RunPathThroughCmd(fmt.Sprintf(`&"%s"`, removeEdgeExePath))

			if err != nil {
				Log.Fatal("\n"+err.Error(), "\n")
			}

			continue
		}

		// special case for Microsoft.OneDrive
		if option == "Microsoft.OneDrive" {
			removeOneDriveScriptPath := filepath.Join(AssetsPath, "uninstallOneDrive.ps1")

			err := Powershell.RunPathThroughCmd(
				"Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process -Force;",
				fmt.Sprintf(`&"%s"`, removeOneDriveScriptPath),
			)

			if err != nil {
				Log.Fatal("\n"+err.Error(), "\n")
			}

			continue
		}

		// remove app
		Powershell.RemoveWinPackage(option)
	}

	Log.Success("\nUninstall complete\n")
}
