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

var drillsCmd = &cobra.Command{
	Use:   "drills",
	Short: "List available drills",
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
		err = driller.ListDrills()
		if err != nil {
			fmt.Printf("ERROR: Failed to list drills: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(drillsCmd)
}
