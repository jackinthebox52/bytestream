package path

import (
	"fmt"
	"os"
	"path"
	"slices"
)

func rootPath() string {
	return "./"
}

func CompileScriptPath(scriptName string) string {
	return path.Join(rootPath(), "internal/scripts", scriptName+".sh")
}

func CompileTemplatePath(templateName string) string {
	return path.Join(rootPath(), "web/templates", templateName+".tmpl")
}

// Returns both the directory and the file path for the HLS stream
func CompileHlsPath(streamId string) (string, string) {
	hlsDir := path.Join(rootPath(), "streams/hls", streamId)
	return hlsDir, path.Join(streamId, "index.m3u8")
}

func compilePidPath(streamId string) string {
	return path.Join(rootPath(), "streams/pids", streamId+".pid")
}

// Returns the video file location for a given video file name in the video directory
// Returns an error if the file extension is not supported, or if the file does not exist
func CompileVideoPath(streamName string, fileExtension string) (string, error) {
	supported := []string{"mkv", "mp4", "webm"}
	if !slices.Contains(supported, fileExtension) {
		return "", fmt.Errorf("File extension %v not supported", fileExtension)
	}
	filePath := path.Join(rootPath(), "streams/videos", streamName+".mkv")
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}
	return "", fmt.Errorf("Video file %v not found", filePath)
}

// Returns the database file location for a given database name, simply calls the preferred hardcoded function (compileSqlitePath)
func CompileDatabasePath(dbName string) (string, error) {
	filePath := compileSqlitePath(dbName)
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}
	return "", fmt.Errorf("Database file %v not found", filePath)
}

func compileSqlitePath(dbName string) string {
	return path.Join(rootPath(), "data/database", dbName+".sqlite3")
}
