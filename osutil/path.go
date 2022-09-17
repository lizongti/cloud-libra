package osutil

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func FindFiles(paths ...string) []string {
	var absPaths []string
	for _, path := range paths {
		if err := StatIsDir(path); err != nil {
			continue
		}
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			absPaths = append(absPaths, path)
			return nil
		}); err != nil {
			log.Panic(err)
		}
	}
	return absPaths
}

func FindFilesBySuffix(suffix string, paths ...string) []string {
	var absPaths []string
	for _, path := range paths {
		if err := StatIsDir(path); err != nil {
			continue
		}
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, fmt.Sprintf(".%s", suffix)) {
				absPaths = append(absPaths, path)
			}
			return nil
		}); err != nil {
			log.Panic(err)
		}
	}
	return absPaths
}

func FindFilesByRegexp(expStr string, paths ...string) []string {
	exp := regexp.MustCompile(expStr)
	var absPaths []string
	for _, path := range paths {
		if err := StatIsDir(path); err != nil {
			continue
		}
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return nil
			}
			if exp.MatchString(path) {
				absPaths = append(absPaths, path)
			}
			return nil
		}); err != nil {
			log.Panic(err)
		}
	}
	return absPaths
}

func FindDirs(path ...string) []string {
	var absPaths []string
	for _, path := range path {
		if err := StatIsDir(path); err != nil {
			continue
		}
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return nil
			}
			absPaths = append(absPaths, path)
			return nil
		}); err != nil {
			log.Panic(err)
		}
	}
	return absPaths
}

func FindDirsBySuffix(suffix string, paths ...string) []string {
	var absPaths []string
	for _, path := range paths {
		if err := StatIsDir(path); err != nil {
			continue
		}
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, fmt.Sprintf(".%s", suffix)) {
				absPaths = append(absPaths, path)
			}
			return nil
		}); err != nil {
			log.Panic(err)
		}
	}
	return absPaths
}

func FindDirsByRegexp(expStr string, paths ...string) []string {
	exp := regexp.MustCompile(expStr)
	var absPaths []string
	for _, path := range paths {
		if err := StatIsDir(path); err != nil {
			continue
		}
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return nil
			}
			if exp.MatchString(path) {
				absPaths = append(absPaths, path)
			}
			return nil
		}); err != nil {
			log.Panic(err)
		}
	}

	return absPaths
}

func ReadFiles(paths ...string) map[string][]byte {
	var (
		files = make(map[string][]byte)
		err   error
	)
	for _, path := range FindFiles(paths...) {
		files[path], err = ReadFile(path)
		if err != nil {
			log.Panic(err)
		}
	}
	return files
}

func ReadFilesBySuffix(suffix string, paths ...string) map[string][]byte {
	var (
		files = make(map[string][]byte)
		err   error
	)
	for _, path := range FindFilesBySuffix(suffix, paths...) {
		files[path], err = ReadFile(path)
		if err != nil {
			log.Panic(err)
		}
	}
	return files
}

func ReadFilesByRegexp(expStr string, paths ...string) map[string][]byte {
	var (
		files = make(map[string][]byte)
		err   error
	)
	for _, path := range FindFilesByRegexp(expStr, paths...) {
		files[path], err = ReadFile(path)
		if err != nil {
			log.Panic(err)
		}
	}
	return files
}

// FindParent returns true if the `path` matches the `parent` directory.
// `path` & `parent“ are in the same style of path, either absolute or relative.
func FindParent(path string, parent string) (string, error) {
	var curPath string
	var base string
	for {
		curPath = filepath.Dir(path)
		if curPath == parent { // parent is a path
			return parent, nil
		}

		base = filepath.Base(path)
		if base == parent { // parent is a dirname or filename
			return path, nil
		}

		// Reach the top of path
		if curPath == path {
			return "", fmt.Errorf(
				"diretory `%s` doesn't match `%s`", path, parent)
		}

		path = curPath
	}
}

// IsParent returns true if the path matches the parent directory.
func IsParent(path string, parent string) bool {
	_, err := FindParent(path, parent)
	return err == nil
}

// GetProjectParth returns the project path.
func GetProjectPath(project string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	if dir, err := FindParent(wd, project); err != nil {
		return wd
	} else {
		return dir
	}
}

func GetWorkPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
		return ""
	}
	return dir
}

func GetExePath() string {
	if _, err := os.Stat(os.Args[0]); err == nil {
		return filepath.Dir(os.Args[0])
	}
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Panic(err)
		return ""
	}
	return filepath.Dir(path)
}

func GetDiskPath(disk string) string {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("%s:\\", strings.ToLower(disk))
	case "linux":
		return fmt.Sprintf("/mnt/%s", strings.ToLower(disk))
	case "darwin":
		return fmt.Sprintf("/Volumes/%s", strings.ToLower(disk))
	default:
		log.Panic(fmt.Errorf("unsupported os: %s", runtime.GOOS))
		return ""
	}
}
