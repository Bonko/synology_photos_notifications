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
	return usr.Username, nil
}

func newFiles(path string, owner string) {
	current := len(filesByOwner[owner])
	lastNumFileName := fmt.Sprintf("%s/last_num_files_%s", path, owner)
	last, err := readIntFromFile(lastNumFileName)
	if err != nil && last != -1 {
		fmt.Print(err)
		os.Exit(1)
	}
	if current > 0 && current > last {
		log.Info("%d new files created by %s", current-last, owner)
	}
	if err := updateLastNumFileName(lastNumFileName, current); err != nil {
		log.Error(err)
	}
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
		return -1, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	var result []int
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())

		if err != nil {
			return 0, err
		}
		result = append(result, x)
	}
	return result[0], scanner.Err()
}
