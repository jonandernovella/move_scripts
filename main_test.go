package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func getLib(reader io.Reader) Lib {
	return Lib{
		Name:        "darsync",
		HomeDir:     os.Getenv("HOME"),
		InputSource: reader,
	}
}

var lib = getLib(strings.NewReader(""))

func TestValidateDirectoryToMove(t *testing.T) {
	badDir := "/test/dir/bad"
	dirToMove := lib.validateDirectoryToMove(badDir)
	if dirToMove != "" {
		t.Errorf("getDirectoryToMove() = %s; want %s", dirToMove, badDir)
	}
	goodDir := "/dev"
	dirToMove = lib.validateDirectoryToMove(goodDir)
	if dirToMove != goodDir {
		t.Errorf("getDirectoryToMove() = %s; want %s", dirToMove, goodDir)
	}
	file := "/dev/null"
	dirToMove = lib.validateDirectoryToMove(file)
	if dirToMove != "" {
		t.Errorf("getDirectoryToMove() = %s; want %s", dirToMove, "")
	}
}

func TestValidateTargetDir(t *testing.T) {
	absolutePath := "/good/target_dir"
	targetDir := lib.validateTargetDir(absolutePath)
	if targetDir != absolutePath {
		t.Errorf("getTargetDirectory() = %s; want %s", targetDir, absolutePath)
	}
	relativePath := "bad_target_dir"
	targetDir = lib.validateTargetDir(relativePath)

	if targetDir != "" {
		t.Errorf("getTargetDirectory() = %s; want %s", targetDir, "")
	}
}

func TestValidateUsername(t *testing.T) {
	goodInput := "good_username"
	username := lib.validateUsername(goodInput)
	if username != goodInput {
		t.Errorf("getUsername() = %s; want %s", username, goodInput)
	}
	longUsername := "this_is_a_really_long_username_that_is_not_allowed"
	username = lib.validateUsername(longUsername)
	if username != "" {
		t.Errorf("getUsername() = %s; want %s", username, "")
	}
}

func TestValidatePrivateKey(t *testing.T) {
	goodInput := "/dev/stdout"
	privateKey := lib.validatePrivateKey(goodInput)
	if privateKey != goodInput {
		t.Errorf("getPrivateKey() = %s; want %s", privateKey, goodInput)
	}
	badInput := "~/bad/path"
	privateKey = lib.validatePrivateKey(badInput)
	if privateKey != "" {
		t.Errorf("getPrivateKey() = %s; want %s", privateKey, "")
	}
}

func TestGetAbsoluteDirectory(t *testing.T) {
	absPath, _ := filepath.Abs(".")
	path, err := getAbsoluteDirectory(absPath)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if path != absPath {
		t.Errorf("Expected %v, got %v", absPath, path)
	}

	path, err = getAbsoluteDirectory(".")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if path != absPath {
		t.Errorf("Expected %v, got %v", absPath, path)
	}

	_, err = getAbsoluteDirectory("/path/does/not/exist")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	_, err = getAbsoluteDirectory("/dev/null")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGetInput(t *testing.T) {
	lib.InputSource = strings.NewReader("test")
	input := lib.getInput("test prompt", "")
	if input != "test" {
		t.Errorf("getInput() = %s; want %s", input, "test")
	}
	lib.InputSource = strings.NewReader("")
	input = lib.getInput("test prompt", "default")
	if input != "default" {
		t.Errorf("getInput() = %s; want %s", input, "default")
	}
}

func TestFormatBytes(t *testing.T) {
	bytes := int64(1024)
	formattedBytes := formatBytes(bytes)
	if formattedBytes != "1.0 KB" {
		t.Errorf("formatBytes() = %s; want %s", formattedBytes, "1.0 KB")
	}
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}
