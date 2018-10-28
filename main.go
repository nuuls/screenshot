package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
		CopyCommand:   []string{"xclip", "-sel", "clip"},
		Dir:           "/home/nuuls/Pictures",
		UploadCommand: []string{"ni", "{filePath}"},
	}
}

var cfg = loadConfig()

func main() {
	watch(cfg.Dir)
}

var fileNameRe = regexp.MustCompile(`^Screenshot .+\.png$`)

func watch(dir string) {
	oldFiles := map[string]time.Time{}
	for range time.Tick(time.Millisecond * 300) {
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

func fileSize(path string) int64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return stat.Size()
}

func uploadAndClip(path string) {
	size := fileSize(path)
	for {
		time.Sleep(time.Millisecond * 200)
		newSize := fileSize(path)
		if size == newSize {
			break
		}
		size = newSize
	}
	log.Println("uploading", path)
	url, err := upload(path)
	if err != nil {
		log.Println(err)
		notify("upload failed\n" + err.Error())
		return
	}
	clip(url)
	notify(url)
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
	err := exec.Command("notify-send", "Screenshot uploaded", text).Run()
	if err != nil {
		log.Println("error showing notification", err)
	}
}
