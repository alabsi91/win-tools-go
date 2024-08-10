package commands

import (
	"errors"
	"fmt"
	"sync"

	"github.com/charmbracelet/huh"
)

func askToInstallChocolatey() bool {
	var answer bool = false

	huh.NewConfirm().
		Title("Chocolatey is not installed. Do you want to install it?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&answer).Run()

	return answer
}

func askForNumberOfWorkers() int {
	var answer int = 0

	validate := func(s string) error {
		_, err := fmt.Sscan(s, &answer)
		if err != nil {
			return errors.New("please enter a valid number")
		}
		return nil
	}

	huh.NewInput().
		Title("How many workers do you want to use?").
		Placeholder("4").
		Validate(validate).
		Run()

	return answer
}

func worker(packages []string, tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		index := task - 1
		packageName := packages[index]
		Log.Info("\n", fmt.Sprintf(`%d. installing package: "%s"`, task, packageName))
		Chocolatey.InstallChocolateyPackage(packageName)
	}
}

func InstallPackages(configFilePath *string, numberOfWorkers *int) {
	// has admin privileges
	isAdmin := Powershell.IsAdmin()

	if !isAdmin {
		Log.Error("\nyou need admin privileges to install packages\n")
		Log.Info("Please run this command from an elevated powershell session\n")
		return
	}

	// no config file path provided, ask for it
	if configFilePath == nil {
		answer := AskForConfigFilePath()
		configFilePath = &answer

		// when the user exit the prompt using CTRL + C
		if !Utils.IsPathExists(answer) {
			Log.Error("\nfile not found. Please enter a valid path\n")
			return
		}
	}

	// config file path provided does not exist, ask for a new one
	if !Utils.IsPathExists(*configFilePath) {
		Log.Error("\nfile not found. Please enter a valid path\n")
		answer := AskForConfigFilePath()

		// when the user exit the prompt using CTRL + C
		if !Utils.IsPathExists(answer) {
			Log.Error("\nfile not found. Please enter a valid path\n")
			return
		}

		configFilePath = &answer
	}

	yamlData := ReadConfigFile(*configFilePath)

	// packages is empty, exit
	if len(yamlData.Packages) == 0 {
		Log.Error("\nthe YAML file does not contain any packages\n")
		return
	}

	// check if chocolatey is installed
	isChocolateyInstalled := Chocolatey.IsChocolateyInstalled()
	if !isChocolateyInstalled {
		answer := askToInstallChocolatey()
		if !answer {
			return
		}

		// install chocolatey
		Chocolatey.InstallChocolatey()
	}

	numWorkers := 4
	if numberOfWorkers == nil {
		numWorkers = askForNumberOfWorkers()
	}

	numTasks := len(yamlData.Packages)

	Log.Info("\n" + fmt.Sprintf(`Found "%d" packages`, len(yamlData.Packages)))

	var wg sync.WaitGroup
	tasks := make(chan int, numTasks) // Buffered channel for tasks

	// Start the workers
	for i := 0; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(yamlData.Packages, tasks, &wg)
	}

	// Send tasks to the workers
	for i := 1; i <= numTasks; i++ {
		tasks <- i
	}

	close(tasks) // Close the tasks channel to signal no more tasks

	// Wait for all workers to finish
	wg.Wait()

	Log.Success("\nDone\n")
}
