package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ljfranklin/fire-drills/driller"
	"github.com/ljfranklin/fire-drills/state"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start DRILL_INDEX",
	Short: "Start a new fire-drill exercise",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("ERROR: Command `start` expects the task index or name as the only arg")
			os.Exit(1)
		}
		taskName := args[0]

		home, err := homedir.Dir()
		if err != nil {
			fmt.Printf("Failed to locate user's homedir")
			os.Exit(1)
		}
		state := state.State{
			StateFilePath: filepath.Join(home, ".fire-drill-state"),
		}
		err = state.Load()
		if err != nil {
			fmt.Printf("Failed to load statefile: %s\n", err)
			os.Exit(1)
		}

		driller := driller.Driller{
			State: &state,
		}
		err = driller.Start(taskName)
		if err != nil {
			fmt.Printf("ERROR: Failed to start task %s: %s\n", taskName, err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
