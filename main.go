package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

//type DirInfo *map[string][]FileInfo

type FileInfo struct {
	Dir   string
	Name  string
	Owner string
}

var fi []FileInfo

func main() {
	path := os.Getenv("PHOTO_DIR")

	//files, _ := ioutil.ReadDir("./testdata")
	//fmt.Println(len(files))
	//
	//files2, _ := os.ReadDir("./testdata")
	//fmt.Println(len(files2))

	err := filepath.WalkDir(path, files)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(fi)
	fmt.Println(len(fi))
}

func files(path string, di fs.DirEntry, err error) error {
	if di.IsDir() {
		return nil
	}
	owner, err := fileOwner(di)
	if err != nil {
		return err
	}

	current := FileInfo{
		Dir:   filepath.Dir(path),
		Name:  filepath.Base(path),
		Owner: owner,
	}
	fi = append(fi, current)

	return nil
}

func fileOwner(di fs.DirEntry) (string, error) {
	info, err := di.Info()
	if err != nil {
		return "", err
	}
	stat := info.Sys().(*syscall.Stat_t)
	uid := stat.Uid
	u := strconv.FormatUint(uint64(uid), 10)
	usr, err := user.LookupId(u)
	return usr.Username, nil
}
