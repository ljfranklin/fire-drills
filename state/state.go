package state

import (
	"errors"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type State struct {
	StateFilePath string

	backing internalState
}

type internalState struct {
	CurrentTask string `yaml:"current_task"`
	CurrentHint int    `yaml:"current_hint"`
}

func (s *State) Load() error {
	if _, err := os.Stat(s.StateFilePath); os.IsNotExist(err) {
		return nil
	}
	statefileContents, err := ioutil.ReadFile(s.StateFilePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal([]byte(statefileContents), &s.backing)
}

func (s *State) SaveCurrentDrill(drillName string) error {
	s.backing.CurrentTask = drillName
	return s.persistBacking()
}

func (s State) CurrentDrill() (string, error) {
	if s.backing.CurrentTask == "" {
		return "", errors.New("No drill in progress")
	}
	return s.backing.CurrentTask, nil
}

func (s *State) SaveCurrentHint(hintIndex int) error {
	s.backing.CurrentHint = hintIndex
	return s.persistBacking()
}

func (s State) CurrentHint() (int, error) {
	return s.backing.CurrentHint, nil
}

func (s State) persistBacking() error {
	newContents, err := yaml.Marshal(s.backing)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.StateFilePath, newContents, 0755)
}
