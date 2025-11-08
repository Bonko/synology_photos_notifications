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

/*
TODO:
- send mail notification
- new files per subdir
*/
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

	//fmt.Println(filesByOwner["bonko"])
	//fmt.Println(len(filesByOwner["bonko"]))

	for owner := range filesByOwner {
		fmt.Print(owner)
		numNewFiles, err := newFiles(config.Rootpath, owner)
		if err != nil {
			log.Errorf("Error getting number of new files: %q", err)
		}
		notifyUsers(owner, numNewFiles, config)
	}
}

func notifyUsers(owner string, numNewFiles int, config *Config) {
	if numNewFiles < 0 {
		log.Infof("No new files uploaded by user %s", owner)
		return
	}
	msg := fmt.Sprintf("%s uploaded %d new files to folder %s", owner, numNewFiles, config.Rootpath)
	for _, user := range config.Users {
		//if user.name == owner {
		//	// Don't notify user about self-uploaded files
		//	continue
		//}
		log.Infof("Notifying %s: %q", user.name, msg)
	}
}
