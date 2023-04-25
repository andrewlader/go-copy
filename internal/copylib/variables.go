package copylib

import (
	"io"
	"os"
)

var Chtimes = os.Chtimes
var Close = closeFile
var Copy = io.Copy
var Create = os.Create
var IsNotExist = os.IsNotExist
var MkdirAll = os.MkdirAll
var Open = os.Open
var ReadAll = io.ReadAll
var ReadDir = os.ReadDir
var Stat = os.Stat
var Sync = syncFile
