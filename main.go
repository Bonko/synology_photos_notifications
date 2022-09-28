package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

/*
users:
	- name:
      email:
*/
 */
//type DirInfo *map[string][]FileInfo

type FileInfo struct {
	Dir  string
	Name string
}

var filesByOwner = make(map[string][]FileInfo)

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

	fmt.Println(filesByOwner["bonko"])
	fmt.Println(len(filesByOwner["bonko"]))

	newFiles(path, "bonko")

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
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	if current > last {
		fmt.Printf("%d new files created by %s", current, owner)
	}
}

func readIntFromFile(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}

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
