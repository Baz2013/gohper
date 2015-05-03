package path2

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/cosiner/gohper/os2"
	"github.com/cosiner/gohper/unibyte"
)

const UNKNOWN rune = ' '

// Home return current user's Home dir
func Home() string {
	u, _ := user.Current()
	return u.HomeDir
}

// ExpandHome expand ~ to user's home dir
func ExpandHome(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}
	u, _ := user.Current()
	return u.HomeDir + path[1:]
}

// ExpandAbs expand path to absolute path
func ExpandAbs(path string) string {
	path, _ = filepath.Abs(ExpandHome(path))
	return path
}

// ProgramDir return dir of program use os.Args[0]
func ProgramDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// LastDir return last dir of path,
// if path is dir, return itself
// else return path's contain dir name
func LastDir(path string) (string, error) {
	var dir string
	absPath, err := filepath.Abs(path)
	if err == nil {
		info, err := os.Stat(absPath)
		if err == nil {
			if info.IsDir() {
				_, dir = filepath.Split(absPath)
			} else {
				dir = filepath.Dir(absPath)
				_, dir = filepath.Split(dir)
			}
		}
	}
	return dir, err
}

// IsRelative check whether a path is relative
// In these condition: path is empty, start with '[.~][/\]', '/', "[a-z]:\"
func IsRelative(path string) bool {
	return !(strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, ".\\") ||
		strings.HasPrefix(path, "~/") ||
		strings.HasPrefix(path, "~\\") ||
		strings.HasPrefix(path, "/") ||
		IsWinRoot(path))
}

// IsWinRoot check whether a path is windows absolute path with disk letter
func IsWinRoot(path string) bool {
	if path == "" {
		return false
	}
	return unibyte.IsLetter(path[0]) && strings.HasPrefix(path[1:], ":\\")
}

// IsRoot check wether or not path is root of filesystem
func IsRoot(path string) bool {
	if l := len(path); l > 0 {
		switch os2.OS() {
		case os2.WINDOWS:
			return IsWinRoot(path)
		case os2.LINUX, os2.DARWIN, os2.FREEBSD, os2.SOLARIS, os2.ANDROID:
			return l == 1 && path[0] == '/'
		}
	}
	return false
}

// EnvSeperator return seperator of env variable "PATH"
func EnvSeperator() rune {
	switch os2.OS() {
	case os2.LINUX, os2.DARWIN, os2.SOLARIS, os2.FREEBSD, os2.ANDROID:
		return ':'
	case os2.WINDOWS:
		return ';'
	default:
		return UNKNOWN
	}
}
