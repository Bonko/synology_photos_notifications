package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

func Test_lastNumFileName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		owner    string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/photos",
			owner:    "alice",
			expected: "/photos/last_num_files_alice",
		},
		{
			name:     "nested path",
			path:     "/photos/folder1",
			owner:    "bob",
			expected: "/photos/folder1/last_num_files_bob",
		},
		{
			name:     "numeric owner",
			path:     "/photos",
			owner:    "999",
			expected: "/photos/last_num_files_999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lastNumFileName(tt.path, tt.owner)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_newFiles(t *testing.T) {
	// Reset global state
	filesByOwner = make(map[string][]FileInfo)

	dir := createTestdataFull(1)
	defer os.RemoveAll(dir)

	err := filepath.WalkDir(dir, genFileInfos)
	assert.NoError(t, err)

	u, err := user.Current()
	assert.NoError(t, err)

	fileNum, err := newFiles(dir, u.Username)
	assert.NoError(t, err)

	assert.Equal(t, 5, fileNum)
}

func Test_newFilesFirstTimeUploader(t *testing.T) {
	// Reset global state
	filesByOwner = make(map[string][]FileInfo)

	dir := createTestFiles()
	defer os.RemoveAll(dir)

	err := filepath.WalkDir(dir, genFileInfos)
	assert.NoError(t, err)

	u, err := user.Current()
	assert.NoError(t, err)

	fileNum, err := newFiles(dir, u.Username)
	assert.NoError(t, err)

	assert.Equal(t, 6, fileNum)
}

func Test_newFiles_withMultipleOwners(t *testing.T) {
	// Reset global state
	filesByOwner = make(map[string][]FileInfo)

	dir := createTestFiles()
	defer os.RemoveAll(dir)

	// Manually populate with multiple owners
	populateFilesByOwnerWithMockOwners(dir)

	// Test newFiles for each owner
	fileNum, err := newFiles(dir, "alice")
	assert.NoError(t, err)
	assert.Equal(t, 2, fileNum, "alice should have 2 new files")

	fileNum, err = newFiles(dir, "bob")
	assert.NoError(t, err)
	assert.Equal(t, 2, fileNum, "bob should have 2 new files")

	fileNum, err = newFiles(dir, "999")
	assert.NoError(t, err)
	assert.Equal(t, 2, fileNum, "user 999 should have 2 new files")
}

func Test_newFiles_noNewFiles(t *testing.T) {
	// Reset global state
	filesByOwner = make(map[string][]FileInfo)

	dir := createTestdataFull(6) // Set last count to 6
	defer os.RemoveAll(dir)

	u, err := user.Current()
	assert.NoError(t, err)

	// Populate with 6 files (same as last count)
	for i := 1; i <= 6; i++ {
		fileInfo := FileInfo{
			Dir:  dir,
			Name: fmt.Sprintf("file%d.jpg", i),
		}
		filesByOwner[u.Username] = append(filesByOwner[u.Username], fileInfo)
	}

	// Should return -2 (no new files, current == last)
	fileNum, err := newFiles(dir, u.Username)
	assert.NoError(t, err)
	assert.Equal(t, -2, fileNum, "should return -2 when current == last")
}

func Test_newFiles_filesDeleted(t *testing.T) {
	// Reset global state
	filesByOwner = make(map[string][]FileInfo)

	dir := createTestdataFull(10) // Set last count to 10
	defer os.RemoveAll(dir)

	u, err := user.Current()
	assert.NoError(t, err)

	// Populate with only 5 files (less than last count of 10)
	for i := 1; i <= 5; i++ {
		fileInfo := FileInfo{
			Dir:  dir,
			Name: fmt.Sprintf("file%d.jpg", i),
		}
		filesByOwner[u.Username] = append(filesByOwner[u.Username], fileInfo)
	}

	// Should return -2 (no new files, current < last)
	fileNum, err := newFiles(dir, u.Username)
	assert.NoError(t, err)
	assert.Equal(t, -2, fileNum, "should return -2 when current < last (files deleted)")
}

func Test_newFiles_emptyOwner(t *testing.T) {
	// Reset global state
	filesByOwner = make(map[string][]FileInfo)

	dir := createTestFiles()
	defer os.RemoveAll(dir)

	// Test with owner that has no files
	fileNum, err := newFiles(dir, "nonexistent")
	assert.NoError(t, err)
	assert.Equal(t, -2, fileNum, "should return -2 when owner has no files")
}

func Test_genFileInfos(t *testing.T) {
	// Reset global state
	filesByOwner = make(map[string][]FileInfo)

	// setup
	dir := createTestdataFull(1)
	defer os.RemoveAll(dir)

	assert.Empty(t, filesByOwner)
	err := filepath.WalkDir(dir, genFileInfos)
	assert.NoError(t, err)

	assert.NotEmpty(t, filesByOwner)

	u, err := user.Current()
	assert.NoError(t, err)

	assert.Equal(t, 6, len(filesByOwner[u.Username]))
}

func Test_genFileInfos_withMultipleOwners(t *testing.T) {
	// Reset global state
	filesByOwner = make(map[string][]FileInfo)

	// Create test directory with files
	dir := createTestFiles()
	defer os.RemoveAll(dir)

	// Manually populate filesByOwner with different owners to simulate
	// files owned by different users (without needing chown)
	populateFilesByOwnerWithMockOwners(dir)

	// Verify files are grouped by owner
	assert.Equal(t, 2, len(filesByOwner["alice"]), "alice should have 2 files")
	assert.Equal(t, 2, len(filesByOwner["bob"]), "bob should have 2 files")
	assert.Equal(t, 2, len(filesByOwner["999"]), "user 999 should have 2 files")
}

func Test_genFileInfos_skipsStateFiles(t *testing.T) {
	// Reset global state
	filesByOwner = make(map[string][]FileInfo)

	dir := createTestFiles()
	defer os.RemoveAll(dir)

	u, err := user.Current()
	assert.NoError(t, err)

	// Create a state file
	stateFilePath := lastNumFileName(dir, u.Username)
	stateFile, err := os.Create(stateFilePath)
	assert.NoError(t, err)
	stateFile.Close()

	// Walk directory
	err = filepath.WalkDir(dir, genFileInfos)
	assert.NoError(t, err)

	// Verify state file is not in filesByOwner
	for _, files := range filesByOwner {
		for _, file := range files {
			assert.NotEqual(t, stateFilePath, fmt.Sprintf("%s/%s", file.Dir, file.Name),
				"state file should not be in filesByOwner")
		}
	}
}

// populateFilesByOwnerWithMockOwners manually populates filesByOwner
// to simulate files owned by different users without requiring chown.
// This allows testing multi-owner scenarios without root privileges.
func populateFilesByOwnerWithMockOwners(baseDir string) {
	// Create files for different owners
	owners := map[string][]string{
		"alice": {"folder1/1.jpg", "folder1/2.jpg"},
		"bob":   {"folder1/3.jpg", "folder2/1.jpg"},
		"999":   {"folder2/2.jpg", "folder2/3.jpg"},
	}

	for owner, filePaths := range owners {
		for _, filePath := range filePaths {
			fullPath := fmt.Sprintf("%s/%s", baseDir, filePath)
			fileInfo := FileInfo{
				Dir:  filepath.Dir(fullPath),
				Name: filepath.Base(fullPath),
			}
			filesByOwner[owner] = append(filesByOwner[owner], fileInfo)
		}
	}
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

// setupNotificationTest creates a test logger with a hook to capture notifications
func setupNotificationTest() *testLogHook {
	hook := &testLogHook{notifications: make(map[string][]string)}
	log.AddHook(hook)
	return hook
}

func Test_notifyUsers_skipsSelfNotification(t *testing.T) {
	// Setup: Create a config with users alice and bob
	config := &Config{
		Rootpath: "/test/photos",
		Users: []User{
			{Name: "alice", Email: "alice@example.com"},
			{Name: "bob", Email: "bob@example.com"},
		},
	}

	// Test: alice uploads files, should notify bob but not alice
	owner := "alice"
	numNewFiles := 5

	// Capture log output to verify notifications
	hook := setupNotificationTest()

	notifyUsers(owner, numNewFiles, config)

	// Verify: bob should be notified, alice should not
	assert.Contains(t, hook.notifications, "bob", "bob should be notified")
	assert.NotContains(t, hook.notifications, "alice", "alice should not be notified about her own uploads")
	assert.Len(t, hook.notifications["bob"], 1, "bob should receive one notification")
}

func Test_notifyUsers_notifiesAllOtherUsers(t *testing.T) {
	// Setup: Create a config with multiple users
	config := &Config{
		Rootpath: "/test/photos",
		Users: []User{
			{Name: "alice", Email: "alice@example.com"},
			{Name: "bob", Email: "bob@example.com"},
			{Name: "charlie", Email: "charlie@example.com"},
		},
	}

	// Test: alice uploads files, should notify bob and charlie but not alice
	owner := "alice"
	numNewFiles := 3

	hook := setupNotificationTest()

	notifyUsers(owner, numNewFiles, config)

	// Verify: bob and charlie should be notified, alice should not
	assert.Contains(t, hook.notifications, "bob", "bob should be notified")
	assert.Contains(t, hook.notifications, "charlie", "charlie should be notified")
	assert.NotContains(t, hook.notifications, "alice", "alice should not be notified about her own uploads")
	assert.Len(t, hook.notifications["bob"], 1, "bob should receive one notification")
	assert.Len(t, hook.notifications["charlie"], 1, "charlie should receive one notification")
}

func Test_notifyUsers_noNewFiles(t *testing.T) {
	// Setup: Create a config with users
	config := &Config{
		Rootpath: "/test/photos",
		Users: []User{
			{Name: "alice", Email: "alice@example.com"},
			{Name: "bob", Email: "bob@example.com"},
		},
	}

	// Test: No new files (negative number)
	owner := "alice"
	numNewFiles := -1

	hook := setupNotificationTest()

	notifyUsers(owner, numNewFiles, config)

	// Verify: No notifications should be sent
	assert.Empty(t, hook.notifications, "No notifications should be sent when numNewFiles < 0")
}

func Test_notifyUsers_ownerNotInConfig(t *testing.T) {
	// Setup: Create a config with users, but owner is not in config
	config := &Config{
		Rootpath: "/test/photos",
		Users: []User{
			{Name: "alice", Email: "alice@example.com"},
			{Name: "bob", Email: "bob@example.com"},
		},
	}

	// Test: charlie (not in config) uploads files, should notify all users
	owner := "charlie"
	numNewFiles := 2

	hook := setupNotificationTest()

	notifyUsers(owner, numNewFiles, config)

	// Verify: Both alice and bob should be notified
	assert.Contains(t, hook.notifications, "alice", "alice should be notified")
	assert.Contains(t, hook.notifications, "bob", "bob should be notified")
	assert.Len(t, hook.notifications["alice"], 1, "alice should receive one notification")
	assert.Len(t, hook.notifications["bob"], 1, "bob should receive one notification")
}

// testLogHook is a logrus hook to capture notification log entries for testing
type testLogHook struct {
	notifications map[string][]string
}

func (h *testLogHook) Levels() []log.Level {
	return []log.Level{log.InfoLevel}
}

func (h *testLogHook) Fire(entry *log.Entry) error {
	// Check if this is a notification log entry
	// Format: "Notifying %s: %q"
	msg := entry.Message
	if len(msg) > 10 && msg[:10] == "Notifying " {
		// Extract username from message like "Notifying alice: \"message\""
		// Find the colon after the username
		colonIdx := -1
		for i := 10; i < len(msg); i++ {
			if msg[i] == ':' {
				colonIdx = i
				break
			}
		}
		if colonIdx > 10 {
			username := msg[10:colonIdx]
			// Extract the notification message (after the colon and space)
			if colonIdx+2 < len(msg) {
				notificationMsg := msg[colonIdx+2:]
				// Remove quotes if present
				if len(notificationMsg) > 0 && notificationMsg[0] == '"' {
					notificationMsg = notificationMsg[1:]
				}
				if len(notificationMsg) > 0 && notificationMsg[len(notificationMsg)-1] == '"' {
					notificationMsg = notificationMsg[:len(notificationMsg)-1]
				}
				h.notifications[username] = append(h.notifications[username], notificationMsg)
			} else {
				h.notifications[username] = append(h.notifications[username], "")
			}
		}
	}
	return nil
}
