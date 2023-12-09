// Package internal contains the internal functions and data structures for the bytestream package, and utility functions for the cmd package.
// System.go contains the functions for process management and other system functions.
package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jackinthebox52/bytestream/internal/paths"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}

	return false
}

func CleanOrphanedPIDFiles() error {
	pidsDir := filepath.Join("streams", "pids")
	files, err := os.ReadDir(pidsDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".pid") {
			pid, err := LoadPIDFile(f.Name()[:len(f.Name())-4])
			if err != nil {
				return err
			}
			if !IsProcessRunning(pid) {
				fmt.Printf("Removing orphaned PID file %v\n", f.Name())
				err := os.Remove(filepath.Join(pidsDir, f.Name()))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CleanHlsDir(uuid string) error {
	if hlsDir, _, err := paths.CompileHlsPath(uuid); err == nil {
		fmt.Printf("Cleaning HLS directory for UUID %v\n", hlsDir)
		cmd := exec.Command("rm", "-rf", path.Join(hlsDir))
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("Error cleaning HLS directory for UUID %v", uuid)
		}
		return nil
	}
	return fmt.Errorf("Error cleaning HLS directory for UUID %v", uuid)
}

func CleanOldHlsFiles() error { //TODO rewrite to attempt to find the stream name
	hlsDir := filepath.Join("streams", "hls")
	files, err := os.ReadDir(hlsDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".m3u8") || strings.HasSuffix(f.Name(), ".ts") {
			filePath := filepath.Join(hlsDir, f.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				return err
			}
			if time.Since(fileInfo.ModTime()) > 24*time.Hour {
				fmt.Printf("Removing old HLS file %v\n", f.Name())
				err := os.Remove(filePath)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func LoadPIDFile(uuid string) (int, error) {
	time.Sleep(200 * time.Millisecond) //TODO remove
	pidFilePath := filepath.Join(".", "streams", "pids", uuid+".pid")
	pidBytes, err := os.ReadFile(pidFilePath)
	if err != nil {
		return 0, err
	}
	pidStr := string(pidBytes)
	pidStr = strings.TrimSpace(pidStr) // Strip newline character
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

func CleanPIDFile(uuid string) error {
	if pidFile, err := paths.CompilePidPath(uuid); err == nil {
		err := os.Remove(pidFile)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Error cleaning PID file for UUID %v", uuid)
}
