package copylib

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"
)

var testFiles []fs.DirEntry
var testDestinations []string

func SetupTestFileSystemFunctions(destinations []string) {
	testDestinations = destinations

	Chtimes = ChtimesSuccess
	Close = closeTestFile
	Copy = CopySuccess
	Create = CreateSuccess
	IsNotExist = IsNotExistSuccess
	MkdirAll = MkdirAllSuccess
	Open = OpenSuccess
	ReadAll = ReadAllSuccess
	ReadDir = ReadDirSuccess
	Stat = StatDestFileDoesNotExistSuccess
	Sync = syncTestFile
}

func ReinitializeFileSystemFunctions() {
	Chtimes = os.Chtimes
	Close = closeFile
	Copy = io.Copy
	Create = os.Create
	IsNotExist = os.IsNotExist
	MkdirAll = os.MkdirAll
	Open = os.Open
	ReadAll = io.ReadAll
	ReadDir = os.ReadDir
	Stat = os.Stat
	Sync = syncFile
}

func ChtimesSuccess(name string, atime time.Time, mtime time.Time) error {
	return nil
}

func ChtimesFailure(name string, atime time.Time, mtime time.Time) error {
	return errors.New("failed to change the file times")
}

func closeTestFile(file *os.File) error {
	return nil
}

func CopySuccess(dst io.Writer, src io.Reader) (written int64, err error) {
	return 1000, nil
}

func CopyFailure(dst io.Writer, src io.Reader) (written int64, err error) {
	return 0, errors.New("failed to copy file")
}

func CreateSuccess(name string) (*os.File, error) {
	file := &os.File{}
	return file, nil
}

func CreateFailure(name string) (*os.File, error) {
	return nil, errors.New("failed to create the file")
}

func IsNotExistSuccess(err error) bool {
	return true
}

func IsNotExistFailure(err error) bool {
	return false
}

func MkdirAllSuccess(path string, perm os.FileMode) error {
	return nil
}

func MkDirAllFailure(path string, perm os.FileMode) error {
	return errors.New("failed to create the desired path")
}

func OpenSuccess(name string) (*os.File, error) {
	file := &os.File{}
	return file, nil
}

func ReadAllSuccess(r io.Reader) ([]byte, error) {
	bytes := make([]byte, 128)

	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("error while generating random string of bytes: %s", err)
	}

	return bytes, nil
}

func ReadAllFailure(r io.Reader) ([]byte, error) {
	return nil, errors.New("failed to read the file")
}

func ReadDirSuccess(dirname string) ([]fs.DirEntry, error) {
	var entries []fs.DirEntry

	for _, dirEntry := range testFiles {
		if (strings.Contains(dirname, dirEntry.Name())) && (dirEntry.IsDir()) {
			children := dirEntry.(*testDirEntry).children
			for _, entry := range children {
				entries = append(entries, entry)
			}
			break
		}
	}

	if entries == nil {
		entries = testFiles
	}

	return entries, nil
}

func createDirEntry(name string, size int64, isDir bool) *testDirEntry {
	if isDir {
		size = 0
	}
	dirEntry := &testDirEntry{
		name:  name,
		isDir: isDir,
		fileInfo: testFileInfo{
			name:    name,
			size:    size,
			modTime: time.Now(),
			isDir:   isDir,
		},
	}

	return dirEntry
}

func ReadDirFailure(dirname string) ([]fs.FileInfo, error) {
	return nil, errors.New("failed to find the given path")
}

func OpenFailure(name string) (*os.File, error) {
	return nil, errors.New("failed to open the file")
}

func getTestDirEntry(name string) fs.DirEntry {
	var file fs.DirEntry

	for _, entry := range testFiles {
		if strings.Contains(name, entry.Name()) {
			file = entry
			break
		}
	}

	return file
}

func isDestDir(name string) bool {
	var isDestinationDir = false

	for _, destination := range testDestinations {
		if strings.Contains(name, destination) {
			isDestinationDir = true
			break
		}
	}

	return isDestinationDir
}

func StatDestFileDoesNotExistSuccess(name string) (os.FileInfo, error) {
	if isDestDir(name) {
		return nil, os.ErrExist
	} else {
		entry := getTestDirEntry(name)
		return entry.Info()
	}
}

func StatDestFileExistsSuccess(name string) (os.FileInfo, error) {
	if isDestDir(name) {
		entry := getTestDirEntry(name)
		return entry.Info()
	} else {
		entry := getTestDirEntry(name)
		return entry.Info()
	}
}

func StatFailure(name string) (os.FileInfo, error) {
	return nil, errors.New("failed to get the file stat info")
}

func syncTestFile(file *os.File) error {
	return nil
}

type testDirEntry struct {
	name     string
	isDir    bool
	fileInfo fs.FileInfo
	children []testDirEntry
}

func (testDirEntry testDirEntry) Name() string {
	return testDirEntry.name
}

func (testDirEntry testDirEntry) IsDir() bool {
	return testDirEntry.isDir
}

func (testDirEntry testDirEntry) Type() os.FileMode {
	if testDirEntry.isDir {
		return os.ModeDir
	} else {
		return os.ModePerm
	}
}

func (testDirEntry testDirEntry) Info() (fs.FileInfo, error) {
	return testDirEntry.fileInfo, nil
}

func (testDirEntry *testDirEntry) addChildDirEntry(childDirEntry *testDirEntry) error {
	if testDirEntry.isDir {
		testDirEntry.children = append(testDirEntry.children, *childDirEntry)
		return nil
	} else {
		return fmt.Errorf("this entry, \"%s\", is not a directory", testDirEntry.name)
	}
}

type testFileInfo struct {
	name    string
	size    int64
	modTime time.Time
	isDir   bool
}

func (testFileInfo testFileInfo) Name() string {
	return testFileInfo.name
}

func (testFileInfo testFileInfo) Size() int64 {
	return testFileInfo.size
}
func (testFileInfo testFileInfo) Mode() fs.FileMode {
	return os.ModePerm
}
func (testFileInfo testFileInfo) ModTime() time.Time {
	return testFileInfo.modTime
}
func (testFileInfo testFileInfo) IsDir() bool {
	return testFileInfo.isDir
}
func (testFileInfo testFileInfo) Sys() any {
	return nil
}
