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

// finishCmd represents the finish command
var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Reveals the solution for the current drill and cleans up",
	Run: func(cmd *cobra.Command, args []string) {
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
		err = driller.Finish()
		if err != nil {
			fmt.Printf("ERROR: Failed to finish task: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(finishCmd)
}
