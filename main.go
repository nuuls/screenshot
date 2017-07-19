package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type config struct {
	CopyCommand   []string
	Dir           string
	UploadCommand []string
}

func loadConfig() config {
	return config{
		CopyCommand:   []string{"pbcopy"},
		Dir:           "/Users/nuuls/Desktop/screenshots",
		UploadCommand: []string{"ni", "{filePath}"},
	}
}

var cfg = loadConfig()

func main() {
	watch(cfg.Dir)
}

func watch(dir string) {
	oldFiles := map[string]time.Time{}
	for range time.Tick(time.Millisecond * 100) {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			modTime := file.ModTime()
			name := file.Name()
			if oldFiles[name].Before(modTime) && modTime.After(time.Now().Add(-time.Second*10)) {
				uploadAndClip(filepath.Join(dir, name))
				oldFiles[name] = modTime
			}
		}
	}
}

func uploadAndClip(path string) {
	log.Println("uploading", path)
	url, err := upload(path)
	if err != nil {
		log.Println(err)
		return
	}
	clip(url)
}

func upload(path string) (string, error) {
	args := cfg.UploadCommand[1:]
	for i, arg := range args {
		if arg == "{filePath}" {
			args[i] = path
		}
	}
	cmd := exec.Command(cfg.UploadCommand[0], args...)
	url, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(url)), err
}

func clip(data string) {
	cmd := exec.Command(cfg.CopyCommand[0], cfg.CopyCommand[1:]...)
	cmd.Stdin = strings.NewReader(data)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
