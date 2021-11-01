package main

import (
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"log"
	"os"
	"strings"
)

func findPidDir(dirChan chan string, done chan bool) {
	baseDir := `/proc`
	files, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}
	sideCarFile, err := os.ReadFile(`/etc/pid_flag_sc/pid_flag`)
	if err != nil {
		log.Fatal(err)
	}
	sideCarPidContents := strings.TrimSpace(string(sideCarFile))
	for _, f := range files {
		if f.IsDir() && valid.IsInt(f.Name()) {
			fileDir := fmt.Sprintf(`%s/%s/root/etc/pid_flag/pid_flag`, baseDir, f.Name())
			fmt.Println(`Processing ` + fileDir)
			if _, err := os.Stat(fileDir); err == nil {
				file, err := os.ReadFile(fileDir)
				if err != nil {
					return
				}
				fmt.Println(`Found file. Checking contents`)
				fileContents := strings.TrimSpace(string(file))
				if fileContents == sideCarPidContents {
					fmt.Printf("Found file in: %s.", fileDir)
					dirChan <- fmt.Sprintf("%s/%s/root/sys/fs/cgroup", baseDir, f.Name())
					return
				}
			}
		}
	}
	// Done chan is kinda a mess...
	defer close(done)
	fmt.Println(`could not find matching container`)
}
