// Package internal contains the internal functions and data structures for the bytestream package, and utility functions for the cmd package.
// System.go contains the functions for process management and other system functions.
package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

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

func LoadPIDFile(uuid string) (int, error) {
	time.Sleep(1 * time.Second) //TODO remove
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
