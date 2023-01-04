package main

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

// TODO: lidar com configuração vazia
type Config struct {
	Editor         string
	Directory_name string
	Filetype       string
}

var cfg Config

func main() {
	errReadConfig := readConfig()
	if errReadConfig != nil {
		panic(errReadConfig)
	}

	if fileExists(fmt.Sprintf("%s/%s", getNoteFileDirectory(), getNoteFileName())) == false {
		errCreateNoteFile := createFileRecursive(getNoteFileDirectory(), getNoteFileName())
		if errCreateNoteFile != nil {
			panic(errCreateNoteFile)
		}
	}

	args := os.Args[1:]

	if len(args) == 0 {
		openNoteFileInEditor()
	} else {
		file, err := os.OpenFile(
			fmt.Sprintf("%s/%s", getNoteFileDirectory(), getNoteFileName()),
			os.O_APPEND|os.O_WRONLY|os.O_CREATE,
			0600,
		)
		if err != nil {
			panic("Unable to open file")
		}
		defer file.Close()

		if _, err = file.WriteString(strings.Join(args, " ") + "\n\n"); err != nil {
			panic(err)
		}
	}
}

func openNoteFileInEditor() {
	editorCmd := exec.Command(
		"nvim",
		fmt.Sprintf("%s/%s", getNoteFileDirectory(), getNoteFileName()),
	)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	errStart := editorCmd.Start()
	if errStart != nil {
		fmt.Println(errStart)
	}

	errWait := editorCmd.Wait()
	if errWait != nil {
		fmt.Println(errWait)
	}
}

func getNoteFileName() string {
	now := time.Now()
	return fmt.Sprintf("%02d%s", now.Day(), cfg.Filetype)
}

func getNoteFileDirectory() string {
	now := time.Now()
	home := os.Getenv("HOME")
	if home == "" {
		panic("Error reading HOME")
	}
	return fmt.Sprintf("%s/%s/%d/%02d", home, cfg.Directory_name, now.Year(), now.Month())
}

func readConfig() error {
	// Read config home
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		return fmt.Errorf("error reading XDG_CONFIG_HOME")
	}

	// Check if config file exists and create one if it doesn't
	if fileExists(configHome+"/notes/config.toml") == false {
		err := createFileRecursive(configHome+"/notes", "config.toml")
		if err != nil {
			panic(err)
		}
	}

	// Read file
	tomlData, errRead := ioutil.ReadFile(configHome + "/notes/config.toml")
	if errRead != nil {
		panic("Error reading config file")
	}

	// Save it as struct to cfg variable
	err := toml.Unmarshal(tomlData, &cfg)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		panic("Error checking if file " + path + " exists")
	}
	return true
}

func createFileRecursive(dirPath string, fileName string) error {
	errDir := os.MkdirAll(dirPath, 0777)
	if errDir != nil {
		return errDir
	}

	file, errFile := os.Create(fmt.Sprintf("%s/%s", dirPath, fileName))
	if errFile != nil {
		return errFile
	}
	defer file.Close()
	return nil
}
