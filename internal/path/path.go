package path

import (
	"fmt"
	"os"
	"path"
	"slices"
)

func rootPath() string {
	return "./" //TODO make this configuarble, ENV variable?
}

func CompileScriptPath(scriptName string) (string, error) {
	filePath := path.Join(rootPath(), "internal/script", scriptName+".sh")
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("Script file %v does not exist", filePath)
	}
	return filePath, nil
}

func CompileTemplatePath(templateName string) (string, error) {
	filePath := path.Join(rootPath(), "web/templates", templateName+".tmpl")
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("Template file %v not found", filePath)
	}
	return filePath, nil
}

// Returns both the directory and the file path for the HLS stream
// Returns an error if the directory does not exist. NO ERROR is returned if the chunklist.m3u8 file does not exist
func CompileHlsPath(streamId string) (string, string, error) {
	hlsDir := path.Join(rootPath(), "streams/hls", streamId)
	hlsFile := path.Join(hlsDir, "chunklist.m3u8")
	if _, err := os.Stat(hlsDir); err != nil {
		return "", "", fmt.Errorf("HLS directory %v not found", hlsDir)
	}
	return hlsDir, hlsFile, nil
}

func CompileHlsBase() (string, error) {
	hlsDir := path.Join(rootPath(), "streams/hls")
	if _, err := os.Stat(hlsDir); err != nil {
		return "", fmt.Errorf("streams/hls directory %v not found", hlsDir)
	}
	return hlsDir, nil
}

func compilePidPath(streamId string) (string, error) {
	filePath := path.Join(rootPath(), "streams/pids", streamId+".pid")
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("PID file %v not found", filePath)
	}
	return filePath, nil
}

// Returns the video file location for a given video file name in the video directory
// Returns an error if the file extension is not supported, or if the file does not exist
func CompileVideoPath(streamName string, fileExtension string) (string, error) {
	supported := []string{"mkv", "mp4", "webm"}
	if !slices.Contains(supported, fileExtension) {
		return "", fmt.Errorf("File extension %v not supported", fileExtension)
	}
	filePath := path.Join(rootPath(), "streams/videos", streamName+".mkv")
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("Video file %v not found", filePath)
	}
	return filePath, nil
}

// Returns the database file location for a given database name, simply calls the preferred hardcoded function (compileSqlitePath)
func CompileDatabasePath(dbName string) (string, error) {
	filePath := compileSqlitePath(dbName)
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("Database file %v not found", filePath)
	}
	return filePath, nil
}

func compileSqlitePath(dbName string) string {
	return path.Join(rootPath(), "data/database", dbName+".sqlite3")
}
