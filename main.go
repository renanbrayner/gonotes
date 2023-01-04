package main

import (
  "os"
  "os/exec"
  "fmt"
  "time"
  "io/ioutil"
  "github.com/pelletier/go-toml/v2"
)

type Config struct {
  Editor string
  Directory_name string
  Filetype string
}

var cfg Config

func main() {
  errReadConfig := readConfig()
  if errReadConfig != nil {
    panic(errReadConfig)
  }

  errCreateNote := createNoteFile()
  if errCreateNote != nil {
    panic(errReadConfig)
  }
}

// TODO: criar função de conseguir nome do arquivo e caminho do arquivo para poder criar e ler
// TODO: ele de forma mais fácil
// func getNoteFileName() string {
//   now := time.Now()
// }

func createNoteFile() error {
  now := time.Now()

  home := os.Getenv("HOME")
  if home == "" {
    return fmt.Errorf("error reading XDG_CONFIG_HOME")
  }

  year := now.Year()
  month := fmt.Sprintf("%02d", now.Month())
  day := fmt.Sprintf("%02d", now.Day())

  createFileRecursive(
    fmt.Sprintf("%s/%s/%d/%s", home, cfg.Directory_name, year, month),
    fmt.Sprintf("%s%s", day, cfg.Filetype),
  )

  return nil
}

func readConfig() error {
  // Read config home
  configHome := os.Getenv("XDG_CONFIG_HOME")
  if configHome == "" {
    return fmt.Errorf("error reading XDG_CONFIG_HOME")
  }

  // Check if config file exists and create one if it doesn't
  if fileExists(configHome + "/notes/config.toml") == false {
    err := createFileRecursive(configHome + "/notes", "config.toml")
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
  fmt.Printf("Creating file %s/%s\n", dirPath, fileName)
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
