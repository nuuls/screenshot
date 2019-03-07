package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
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

var fileNameRe = regexp.MustCompile(`.+\.png$`)

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
			if !fileNameRe.MatchString(name) || !oldFiles[name].IsZero() {
				continue
			}
			if modTime.After(time.Now().Add(-time.Second * 3)) {
				uploadAndClip(filepath.Join(dir, name))
				oldFiles[name] = modTime
				break
			}
		}
	}
}

func uploadAndClip(path string) {
	log.Println("uploading", path)
	url, err := upload(path)
	if err != nil {
		log.Println(err)
		notify("upload failed\n" + err.Error())
		return
	}
	clip(url)
	notify("upload complete\n" + url)
}

func upload(path string) (string, error) {
	args := []string{}
	for _, arg := range cfg.UploadCommand {
		if arg == "{filePath}" {
			args = append(args, path)
		} else {
			args = append(args, arg)
		}
	}
	cmd := exec.Command(args[0], args[1:]...)
	url, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error uploading file: %v %s", err, url)
	}
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

func notify(text string) {
	// TODO: escape text
	err := exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "Screenshot"`, text)).Run()
	if err != nil {
		log.Println("error showing notification", err)
	}
}
