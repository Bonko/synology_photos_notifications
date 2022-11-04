package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"testing"
)

func Test_updateLastNumFileName(t *testing.T) {
	expectedNumber := 5
	file, err := os.CreateTemp("testdata", "updateLastNumFileName")
	defer os.Remove(file.Name())
	assert.NoError(t, err)

	err = updateLastNumFileName(file.Name(), expectedNumber)
	assert.NoError(t, err)

	i, err := readIntFromFile(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, expectedNumber, i)
}

func Test_readIntFromFile(t *testing.T) {
	i, err := readIntFromFile("testdata/last_num_files_test")

	assert.NoError(t, err)
	assert.Equal(t, 1, i)
}

func Test_newFiles(t *testing.T) {
	dir := createTestdataFull(1)
	defer os.RemoveAll(dir)

	err := filepath.WalkDir(dir, genFileInfos)

	u, err := user.Current()
	assert.NoError(t, err)

	fileNum, err := newFiles(dir, u.Username)
	assert.NoError(t, err)

	assert.Equal(t, 5, fileNum)
}

func Test_newFilesFirstTimeUploader(t *testing.T) {
	dir := createTestFiles()
	defer os.RemoveAll(dir)

	err := filepath.WalkDir(dir, genFileInfos)

	u, err := user.Current()
	assert.NoError(t, err)

	fileNum, err := newFiles(dir, u.Username)
	assert.NoError(t, err)

	assert.Equal(t, 6, fileNum)
}

func Test_genFileInfos(t *testing.T) {
	// setup
	dir := createTestdataFull(1)
	defer os.RemoveAll(dir)

	assert.Empty(t, filesByOwner)
	err := filepath.WalkDir(dir, genFileInfos)
	assert.NoError(t, err)

	assert.NotEmpty(t, filesByOwner)

	u, err := user.Current()
	assert.NoError(t, err)

	//assert.Equal(t, 2, len(filesByOwner["999"]))
	assert.Equal(t, 6, len(filesByOwner[u.Username]))
}

func createTestdataFull(lastNumFiles int) string {
	dir := createTestFiles()
	createTestdataLastNumFiles(dir, lastNumFiles)
	return dir
}
func createTestFiles() string {
	/*
		testdata
		├── folder1
		│	├── 1.jpg
		│	├── 2.jpg
		│	└── 3.jpg
		├── folder2
		│	 ├── 1.jpg
		│	 ├── 2.jpg
		│	 └── 3.jpg
		├── last_num_files_999
		├── last_num_files_<current_user>
	*/
	tmp, err := os.MkdirTemp("", "testdata")
	if err != nil {
		log.Fatal(err)
	}
	testDir1 := fmt.Sprintf("%s/%s", tmp, "folder1")
	testDir2 := fmt.Sprintf("%s/%s", tmp, "folder2")
	files := []string{"1.jpg", "2.jpg", "3.jpg"}

	for _, dir := range []string{testDir1, testDir2} {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {

			filename := fmt.Sprintf("%s/%s", dir, file)
			_, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			// chown requires root privileges :(
			//owner := os.Getuid()
			//if file == "2.jpg" {
			//	owner = 999
			//}
			//err = os.Chown(filename, owner, os.Getgid())
			//if err != nil {
			//	log.Fatal(err)
			//}
		}
	}
	return tmp
}
func createTestdataLastNumFiles(dir string, num int) {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(lastNumFileName(dir, u.Username))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(num))
	if err != nil {
		log.Fatal(err)
	}
}
