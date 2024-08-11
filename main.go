package main

import (
	"github.com/alabsi91/win-tools/commands"
	"github.com/alabsi91/win-tools/commands/utils"
	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/huh"
)

var Log = utils.Log

type RunCMD struct {
	ID      int `arg:"required"`
	Timeout int
}

type ConfigPathArg struct {
	ConfigPath *string `arg:"--config" placeholder:"[PATH]" help:"YAML config file path"`
}

type CreateTemplateArgs struct {
	TemplatePath *string `arg:"--save-path" placeholder:"[PATH]" help:"Output path for the template"`
}

type NoArgs struct{}

type AutoLogonArgs struct {
	Username     *string `arg:"--username" placeholder:"" help:"The username of the user to automatically logon as"`
	Domain       *string `arg:"--domain" placeholder:"" help:"The domain of the user to automatically logon as"`
	AutoLogon    *int    `arg:"--logon-count" placeholder:"" help:"The number of logons that auto logon will be enabled"`
	RemovePrompt *bool   `arg:"--remove-prompt" help:"Removes the system banner to ensure interventionless logon"`
	BackupFile   *string `arg:"--backup-file" placeholder:"" help:"If specified the existing settings such as the system banner text will be backed up to the specified file"`
}

var args struct {
	Backup               *ConfigPathArg      `arg:"subcommand:backup" help:"Create a backup of specified paths as defined in a YAML configuration file."`
	Restore              *ConfigPathArg      `arg:"subcommand:restore" help:"Restore files and directories from a backup using the paths specified in a YAML configuration file."`
	Install              *ConfigPathArg      `arg:"subcommand:choco-install" help:"Install Chocolatey packages according to the list provided in a YAML configuration file."`
	RunScripts           *ConfigPathArg      `arg:"subcommand:run-scripts" help:"Execute a series of scripts defined in a YAML configuration file."`
	SetEnvs              *ConfigPathArg      `arg:"subcommand:set-envs" help:"Set environment variables as defined in a YAML configuration file."`
	CreateConfigTemplate *CreateTemplateArgs `arg:"subcommand:create-template" help:"Generate a new YAML template for configuration, including placeholders for paths, scripts, and environment variables."`
	SetRegistry          *NoArgs             `arg:"subcommand:set-registry" help:"Select a predefined registry to set."`
	CleanStartMenu       *NoArgs             `arg:"subcommand:clean-menu" help:"Clean start menu from all icons."`
	AutoLogon            *AutoLogonArgs      `arg:"subcommand:auto-logon" help:"Enables auto logon when the computer starts."`
	DisableFirewall      *NoArgs             `arg:"subcommand:disable-firewall" help:"Disable Windows firewall, Windows Defender, and Windows Defender Cloud."`
	UninstallBloat       *NoArgs             `arg:"subcommand:uninstall-bloat" help:"Uninstall default Microsoft bloatware."`
}

func main() {

	arg.MustParse(&args)

	if args.Backup != nil {
		commands.BackupData(args.Backup.ConfigPath)
		return
	}
	if args.Restore != nil {
		commands.RestoreData(args.Restore.ConfigPath)
		return
	}
	if args.Install != nil {
		commands.InstallPackages(args.Install.ConfigPath)
		return
	}
	if args.RunScripts != nil {
		commands.RunScripts(args.RunScripts.ConfigPath)
		return
	}
	if args.SetEnvs != nil {
		commands.SetEnvs(args.SetEnvs.ConfigPath)
		return
	}
	if args.CreateConfigTemplate != nil {
		commands.CreateConfigTemplate(args.CreateConfigTemplate.TemplatePath)
		return
	}
	if args.SetRegistry != nil {
		commands.SetRegistry()
		return
	}
	if args.CleanStartMenu != nil {
		commands.CleanStartMenu()
		return
	}
	if args.AutoLogon != nil {
		commands.AutoLogon(args.AutoLogon.Username, args.AutoLogon.Domain, args.AutoLogon.AutoLogon, args.AutoLogon.RemovePrompt, args.AutoLogon.BackupFile)
		return
	}
	if args.DisableFirewall != nil {
		commands.DisableFirewall()
		return
	}
	if args.UninstallBloat != nil {
		commands.UninstallBloat()
		return
	}

	Log.Info("\nRun `win-tools --help` for more information.\n")

	// no command provided, ask for it
	var chosenCommand string

	err := huh.NewSelect[string]().
		Title("\nWhat would you like to do?").
		Description("  Select a command to run:\n").
		Options(
			huh.NewOption("Backup", "backup"),
			huh.NewOption("Restore", "restore"),
			huh.NewOption("Chocolatey install", "install"),
			huh.NewOption("Run scripts", "run-scripts"),
			huh.NewOption("Set environment variables", "set-envs"),
			huh.NewOption("Create config template", "create-template"),
			huh.NewOption("Set registry", "set-registry"),
			huh.NewOption("Clean start menu", "clean-menu"),
			huh.NewOption("Enable auto logon", "auto-logon"),
			huh.NewOption("Disable Windows firewall", "disable-firewall"),
			huh.NewOption("Uninstall bloatware", "uninstall-bloat"),
		).
		Value(&chosenCommand).
		Run()

	if err != nil {
		Log.Error("\nFailed to get user selection\n")
		return
	}

	switch chosenCommand {
	case "backup":
		commands.BackupData(nil)
	case "restore":
		commands.RestoreData(nil)
	case "install":
		commands.InstallPackages(nil)
	case "run-scripts":
		commands.RunScripts(nil)
	case "set-envs":
		commands.SetEnvs(nil)
	case "create-template":
		commands.CreateConfigTemplate(nil)
	case "set-registry":
		commands.SetRegistry()
	case "clean-menu":
		commands.CleanStartMenu()
	case "auto-logon":
		commands.AutoLogon(nil, nil, nil, nil, nil)
	case "disable-firewall":
		commands.DisableFirewall()
	case "uninstall-bloat":
		commands.UninstallBloat()
	}

}
