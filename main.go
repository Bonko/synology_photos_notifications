package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	path := os.Getenv("PHOTO_DIR")

	err := filepath.WalkDir(path, genFileInfos)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(filesByOwner["bonko"])
	fmt.Println(len(filesByOwner["bonko"]))

	newFiles(path, "bonko")

}
