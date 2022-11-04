package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func genFileInfos(path string, di fs.DirEntry, err error) error {
	if di.IsDir() {
		return nil
	}
	owner, err := fileOwner(di)
	if err != nil {
		return err
	}
	current := FileInfo{
		Dir:  filepath.Dir(path),
		Name: filepath.Base(path),
	}
	if path == lastNumFileName(current.Dir, owner) {
		return nil
	}

	filesByOwner[owner] = append(filesByOwner[owner], current)

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
	if err != nil {
		log.Infof("User lookup failed, returning UID instead: %s", u)
		return u, nil
	}
	return usr.Username, nil
}

func lastNumFileName(path, owner string) string {
	return fmt.Sprintf("%s/last_num_files_%s", path, owner)
}

func newFiles(path string, owner string) (int, error) {
	current := len(filesByOwner[owner])
	last, err := readIntFromFile(lastNumFileName(path, owner))
	if err != nil && last != 0 {
		return -1, err
	}
	newFiles := -2
	if current > 0 && current > last {
		newFiles = current - last
		log.Infof("%d new files created by %s", newFiles, owner)
	}
	if err := updateLastNumFileName(lastNumFileName(path, owner), current); err != nil {
		return newFiles, err
	}
	return newFiles, nil
}

func updateLastNumFileName(path string, fileNum int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(fileNum))
	if err != nil {
		return err
	}

	return nil
}

func readIntFromFile(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	var result []int
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())

		if err != nil {
			return -1, err
		}
		result = append(result, x)
	}
	return result[0], scanner.Err()
}
