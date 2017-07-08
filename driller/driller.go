package driller

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ljfranklin/fire-drills/state"

	"gopkg.in/yaml.v2"
)

const WorkingDir = ".fire-drills-workspace"

type DrillConfig struct {
	Prompt          string   `yaml:"prompt"`
	Solution        string   `yaml:"solution"`
	SetupCmd        string   `yaml:"setup_cmd"`
	TeardownCmd     string   `yaml:"teardown_cmd"`
	Hints           []string `yaml:"hints"`
	RequiredEnvVars []string `yaml:"required_env_vars"`
	Summary         string   `yaml:"summary"`
}

type Driller struct {
	State *state.State
}

func (d Driller) ListDrills() error {
	drillFiles, err := d.listDrillDirs()
	if err != nil {
		return err
	}

	d.printFred()

	fmt.Println("Here's the list of fires just waitin' to be put out:")
	for i, name := range drillFiles {
		fmt.Printf("[%d] %s\n", i+1, name)
	}
	fmt.Printf("\nRun `%s start DRILL_INDEX` to get started!\n", os.Args[0])

	return nil
}

func (d Driller) Start(taskName string) error {
	taskName, err := d.translateTaskName(taskName)
	if err != nil {
		return err
	}

	config, err := d.loadConfig(taskName)
	if err != nil {
		return err
	}

	err = d.State.SaveCurrentHint(0)
	if err != nil {
		return err
	}

	d.printFred()

	fmt.Printf("Hold on there partner, before we begin lets make sure you're all set up:\n\n%s\n", config.Summary)
	fmt.Printf("You'll also need to have the following environment variables set: %s\n", strings.Join(config.RequiredEnvVars, ", "))
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nHit Enter to start the exercise ")
	_, _ = reader.ReadString('\n')

	fmt.Println("\n\n#################### STARTING SETUP #########################")
	fmt.Println("Running setup script...")
	err = d.validateEnvVars(config)
	if err != nil {
		return err
	}

	err = d.runScript(config, taskName, config.SetupCmd)
	if err != nil {
		fmt.Println("Setup failed, running teardown...")
		_ = d.teardown(config, taskName)
		fmt.Println("Finished teardown")

		d.printFred()
		fmt.Printf("Oh no, we ran into trouble setting up the drill! Try to fix the issue and re-run `%s start %s` to try again.\n", os.Args[0], taskName)
		return err
	}
	fmt.Println("Finished setup")
	fmt.Print("#################### FINISHED SETUP #########################\n\n")

	d.printFred()

	fmt.Printf("\nHere's the issue we'll be tackling today:\n%s", config.Prompt)

	fmt.Printf("\nIf you're having trouble, run `%s hint` to get some help. Good luck!\n", os.Args[0])

	err = d.State.SaveCurrentDrill(taskName)
	if err != nil {
		return err
	}

	return nil
}

func (d Driller) Finish() error {
	taskName, err := d.State.CurrentDrill()
	if err != nil {
		return err
	}

	config, err := d.loadConfig(taskName)
	if err != nil {
		return err
	}

	d.printFred()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nAre you ready to see the solution? Hit Enter to see the answer ")
	_, _ = reader.ReadString('\n')

	fmt.Printf("\nCongratulations! Let's go ahead and check your solution:\n\n%s", config.Solution)

	fmt.Print("\nHit Enter to finish and run cleanup ")
	_, _ = reader.ReadString('\n')

	fmt.Println("\n\n#################### STARTING TEARDOWN #########################")
	fmt.Println("Running teardown script...")
	err = d.validateEnvVars(config)
	if err != nil {
		return err
	}

	err = d.teardown(config, taskName)
	if err != nil {
		return err
	}
	fmt.Println("Finished teardown")
	fmt.Println("\n\n#################### FINISHED TEARDOWN #########################")

	d.printFred()
	fmt.Println("See ya real soon!")

	err = d.State.SaveCurrentDrill("")
	if err != nil {
		return err
	}

	return nil
}

func (d Driller) ProvideHint() error {
	taskName, err := d.State.CurrentDrill()
	if err != nil {
		return err
	}

	config, err := d.loadConfig(taskName)
	if err != nil {
		return err
	}

	currentHint, err := d.State.CurrentHint()
	if err != nil {
		return err
	}

	d.printFred()
	fmt.Println("Need some help?")
	for i := 0; i <= currentHint && i < len(config.Hints); i++ {
		fmt.Printf("Hint %d: %s\n", i+1, config.Hints[i])
	}
	if currentHint >= len(config.Hints) {
		fmt.Println("I'm all out of hints.  You can do it!!")
	}
	if currentHint < len(config.Hints) {
		err = d.State.SaveCurrentHint(currentHint + 1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d Driller) HelpText() string {
	help := ""
	help += fred

	help += "\nHowdy y'all! I'm Foundry Fred and I'm here to help teach you how to put out production fires before they happen.\n"
	help += "I can run you through some fire drills to simulate a production issue in order to better understand the system.\n"
	help += "Try some of the commands below to get started!\n"

	return help
}

func (d Driller) loadConfig(taskName string) (DrillConfig, error) {
	drillsDir, err := d.checkForDrillsDir()
	if err != nil {
		return DrillConfig{}, err
	}
	configContents, err := ioutil.ReadFile(filepath.Join(drillsDir, taskName, "drill.yml"))
	if err != nil {
		return DrillConfig{}, err
	}

	config := DrillConfig{}
	err = yaml.Unmarshal([]byte(configContents), &config)
	if err != nil {
		return DrillConfig{}, err
	}

	return config, nil
}

func (d Driller) teardown(config DrillConfig, taskName string) error {
	err := d.runScript(config, taskName, config.TeardownCmd)
	if err != nil {
		return err
	}

	return os.RemoveAll(WorkingDir)
}

func (d Driller) runScript(config DrillConfig, taskName string, scriptName string) error {
	workingDir, err := d.checkForWorkingDir()
	if err != nil {
		return err
	}

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd := exec.Command(filepath.Join(currDir, "drills", taskName, scriptName))
	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (d Driller) checkForDrillsDir() (string, error) {
	if _, err := os.Stat("drills"); err != nil {
		return "", errors.New("Expected to see a `drills` folder in the current directory. Please `cd` to the `fire-drills` project.")
	}

	return "drills", nil
}

func (d Driller) checkForWorkingDir() (string, error) {
	err := os.MkdirAll(WorkingDir, 0755)
	if err != nil {
		return "", err
	}
	return WorkingDir, nil
}

func (d Driller) validateEnvVars(config DrillConfig) error {
	missingEnvVars := []string{}
	for _, envVar := range config.RequiredEnvVars {
		if os.Getenv(envVar) == "" {
			missingEnvVars = append(missingEnvVars, envVar)
		}
	}

	if len(missingEnvVars) > 0 {
		return fmt.Errorf("Missing required env vars: %s", strings.Join(missingEnvVars, ", "))
	}
	return nil
}

func (d Driller) translateTaskName(userTaskName string) (string, error) {
	taskIndex, err := strconv.Atoi(userTaskName)
	if err != nil {
		return userTaskName, nil // user passed a taskName rather than an index
	}

	drillNames, err := d.listDrillDirs()
	if err != nil {
		return "", err
	}
	if taskIndex-1 >= len(drillNames) || taskIndex <= 0 {
		return "", fmt.Errorf("Invalid drill index %d", taskIndex)
	}

	return drillNames[taskIndex-1], nil
}

func (d Driller) listDrillDirs() ([]string, error) {
	drillsDir, err := d.checkForDrillsDir()
	if err != nil {
		return nil, err
	}

	drillFiles, err := ioutil.ReadDir(drillsDir)
	if err != nil {
		return nil, err
	}

	filenames := []string{}
	for _, drillFile := range drillFiles {
		filenames = append(filenames, drillFile.Name())
	}
	return filenames, nil
}

const fred = `
         .'.---.'.
        //   ,   \\
       ||   '|    ||
       ||    |    ||
       ||   -'-   ||
  .-"''-.,_     _,.-''"-.
 / .'--,___'"""'___,--'. \
 |  /:////_'---'_\\\\:\  |
  \|:|// "_     _" \\|:|/
   '-/| (◕)     (◕) |\-'
     \\     | |     //
      '|   (._.)   |'
       |           |
       \     ‿     /
        '--.___.--'
       --------------
      | FOUNDRY FRED |
       --------------
`

func (d Driller) printFred() {
	fmt.Println(fred)
}
