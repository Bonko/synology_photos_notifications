package main

import (
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

/*
-- config --
rootpath:
users:
	- name:
      email:
*/

//type DirInfo *map[string][]FileInfo

type FileInfo struct {
	Dir  string
	Name string
}

var filesByOwner = make(map[string][]FileInfo)

func main() {
	config, err := NewConfig("config.yml")
	if err != nil {
		log.Error(err)
	}
	//path := os.Getenv("PHOTO_DIR")

	err = filepath.WalkDir(config.Rootpath, genFileInfos)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(filesByOwner["bonko"])
	fmt.Println(len(filesByOwner["bonko"]))

	for owner := range filesByOwner {
		fmt.Print(owner)
		numNewFiles, err := newFiles(config.Rootpath, owner)
		if err != nil {
			log.Errorf("Error getting number of new files: %q", err)
		}
		notifyUsers(config.Rootpath, owner, numNewFiles)
	}
}

func notifyUsers(path, owner string, numNewFiles int) {
	if numNewFiles < 0 {
		log.Infof("No new files uploaded by user %s", owner)
		return
	}
	msg := fmt.Sprintf("%s uploaded %d new files to folder %s", owner, numNewFiles, path)
	for o := range filesByOwner {
		log.Infof("Notifying %s: %s", o, msg)
		//if o == owner {
		//	// Don't notify user about self-uploaded files
		//	continue
		//}

	}
}
