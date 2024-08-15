package main

import (
	"reflect"
	"strings"

	"github.com/alabsi91/win-tools/commands"
	"github.com/alabsi91/win-tools/commands/utils"
	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/huh"
)

var Log = utils.Log

type ConfigPathArg struct {
	ConfigPath *string `arg:"--config" placeholder:"[PATH]" help:"YAML config file path"`
}

type CreateTemplateArgs struct {
	TemplatePath *string `arg:"--save-path" placeholder:"[PATH]" help:"Output path for the template"`
}

type AutoLogonArgs struct {
	Username     *string `arg:"--username" placeholder:"" help:"The username of the user to automatically logon as"`
	Domain       *string `arg:"--domain" placeholder:"" help:"The domain of the user to automatically logon as"`
	AutoLogon    *int    `arg:"--logon-count" placeholder:"" help:"The number of logons that auto logon will be enabled"`
	RemovePrompt *bool   `arg:"--remove-prompt" help:"Removes the system banner to ensure interventionless logon"`
	BackupFile   *string `arg:"--backup-file" placeholder:"" help:"If specified the existing settings such as the system banner text will be backed up to the specified file"`
}

type NoArgs struct{}

type ArgsType struct {
	Backup               *ConfigPathArg      `arg:"subcommand:backup" help:"Create a backup of specified paths as defined in a YAML configuration file."`
	Restore              *ConfigPathArg      `arg:"subcommand:restore" help:"Restore files and directories from a backup using the paths specified in a YAML configuration file."`
	Install              *ConfigPathArg      `arg:"subcommand:choco-install" help:"Install Chocolatey packages according to the list provided in a YAML configuration file."`
	RunScripts           *ConfigPathArg      `arg:"subcommand:run-scripts" help:"Execute a series of scripts defined in a YAML configuration file."`
	SetEnvs              *ConfigPathArg      `arg:"subcommand:set-envs" help:"Set environment variables as defined in a YAML configuration file."`
	CreateConfigTemplate *CreateTemplateArgs `arg:"subcommand:create-template" help:"Generate a new YAML template for configuration, including placeholders for paths, scripts, and environment variables."`
	SetRegistry          *NoArgs             `arg:"subcommand:set-registry" help:"Select multiple predefined registry keys to set."`
	CleanStartMenu       *NoArgs             `arg:"subcommand:clean-menu" help:"Clean start menu from all icons."`
	AutoLogon            *AutoLogonArgs      `arg:"subcommand:auto-logon" help:"Enables auto logon when the computer starts."`
	DisableFirewall      *NoArgs             `arg:"subcommand:disable-firewall" help:"Disable Windows firewall, Windows Defender, and Windows Defender Cloud."`
	UninstallBloat       *NoArgs             `arg:"subcommand:uninstall-bloat" help:"Select multiple predefined Windows apps to uninstall."`
}

var args ArgsType

func main() {

	parsedArg := arg.MustParse(&args)

	enteredSubcommands := parsedArg.SubcommandNames()
	if len(enteredSubcommands) > 0 {
		runCommand(enteredSubcommands[0], &args)
		return
	}

	// * No command provided, ask to select one
	Log.Info("\nRun `win-tools --help` for more information.\n")
	chosenCommand, err := askToSelectCommand()
	if err != nil {
		Log.Error("failed to get user selection\n")
		return
	}

	runCommand(chosenCommand, &args)
}

func runCommand(command string, args *ArgsType) {

	switch command {

	case "backup":
		if args.Backup == nil {
			commands.BackupData(nil)
			break
		}
		commands.BackupData(args.Backup.ConfigPath)

	case "restore":
		if args.Restore == nil {
			commands.RestoreData(nil)
			break
		}
		commands.RestoreData(args.Restore.ConfigPath)

	case "install":
		if args.Install == nil {
			commands.InstallPackages(nil)
			break
		}
		commands.InstallPackages(args.Install.ConfigPath)

	case "run-scripts":
		if args.RunScripts == nil {
			commands.RunScripts(nil)
			break
		}
		commands.RunScripts(args.RunScripts.ConfigPath)

	case "set-envs":
		if args.SetEnvs == nil {
			commands.SetEnvs(nil)
			break
		}
		commands.SetEnvs(args.SetEnvs.ConfigPath)

	case "create-template":
		if args.CreateConfigTemplate == nil {
			commands.CreateConfigTemplate(nil)
			break
		}
		commands.CreateConfigTemplate(args.CreateConfigTemplate.TemplatePath)

	case "set-registry":
		commands.SetRegistry()

	case "clean-menu":
		commands.CleanStartMenu()

	case "auto-logon":
		if args.AutoLogon == nil {
			commands.AutoLogon(nil, nil, nil, nil, nil)
			break
		}
		commands.AutoLogon(args.AutoLogon.Username, args.AutoLogon.Domain, args.AutoLogon.AutoLogon, args.AutoLogon.RemovePrompt, args.AutoLogon.BackupFile)

	case "disable-firewall":
		commands.DisableFirewall()

	case "uninstall-bloat":
		commands.UninstallBloat()
	}
}

// getHelpMessages
// parses the help message for a given subcommand
// from the 'args' struct tag.
//
// Returns: help message string
func getHelpMessages(command string) string {
	v := reflect.TypeOf(args)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		subcommand, _ := strings.CutPrefix(field.Tag.Get("arg"), "subcommand:")
		help := field.Tag.Get("help")

		if subcommand == command {
			return help
		}
	}

	return ""
}

// askToSelectCommand prompts the user to select a command
//   - Returns an error if the user cancels the prompt
func askToSelectCommand() (string, error) {
	var chosenCommand string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Description("Select a command to run:").
				Options(
					huh.NewOption("Backup", "backup"),
					huh.NewOption("Restore", "restore"),
					huh.NewOption("Chocolatey install", "choco-install"),
					huh.NewOption("Run scripts", "run-scripts"),
					huh.NewOption("Set environment variables", "set-envs"),
					huh.NewOption("Create config template", "create-template"),
					huh.NewOption("Set registry", "set-registry"),
					huh.NewOption("Clean start menu", "clean-menu"),
					huh.NewOption("Enable auto logon", "auto-logon"),
					huh.NewOption("Disable Windows firewall", "disable-firewall"),
					huh.NewOption("Uninstall bloatware", "uninstall-bloat"),
				).
				Value(&chosenCommand),

			// Print help message for selected subcommand
			huh.NewNote().
				DescriptionFunc(
					func() string {
						return getHelpMessages(chosenCommand)
					},
					&chosenCommand,
				),
		),
	).Run()

	return chosenCommand, err
}
