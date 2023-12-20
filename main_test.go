package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func getLib(inputString string) Lib {
	return Lib{
		Name:        "darsync",
		HomeDir:     os.Getenv("HOME"),
		InputSource: bytes.NewBufferString(inputString),
	}
}

func TestGetDirectoryToMove(t *testing.T) {
	badDir := "/test/dir/bad"
	badLib := getLib(badDir)
	dirToMove := badLib.collectDirectoryToMove()
	if dirToMove != "" {
		t.Errorf("getDirectoryToMove() = %s; want %s", dirToMove, badDir)
	}
	goodDir := "/dev"
	goodLib := getLib(goodDir)
	dirToMove = goodLib.collectDirectoryToMove()
	if dirToMove != goodDir {
		t.Errorf("getDirectoryToMove() = %s; want %s", dirToMove, goodDir)
	}
	file := "/dev/null"
	fileLib := getLib(file)
	dirToMove = fileLib.collectDirectoryToMove()
	if dirToMove != "" {
		t.Errorf("getDirectoryToMove() = %s; want %s", dirToMove, "")
	}
}
func TestGetTargetHost(t *testing.T) {
	testInput := "test_host"
	testLib := getLib(testInput)
	host := testLib.getTargetHost()
	if host != testInput {
		t.Errorf("getTargetHost() = %s; want %s", host, testInput)
	}
	emptyInput := ""
	emptyLib := getLib(emptyInput)
	host = emptyLib.getTargetHost()
	if host != "dardel.pdc.kth.se" {
		t.Errorf("getTargetHost() = %s; want %s", host, "dardel.pdc.kth.se")
	}
}

func TestGetProjectId(t *testing.T) {
	goodInput := "good_project_id"
	lib := getLib(goodInput)
	projectId := lib.getProjectId()

	if projectId != goodInput {
		t.Errorf("getProjectId() = %s; want %s", projectId, goodInput)
	}

}

func TestGetTargetDirectory(t *testing.T) {
	absolutePath := "/good/target_dir"
	absLib := getLib(absolutePath)
	targetDir := absLib.collectTargetDir("test_host")
	if targetDir != absolutePath {
		t.Errorf("getTargetDirectory() = %s; want %s", targetDir, absolutePath)
	}
	relativePath := "bad_target_dir"
	relLib := getLib(relativePath)
	targetDir = relLib.collectTargetDir("test_host")

	if targetDir != "" {
		t.Errorf("getTargetDirectory() = %s; want %s", targetDir, "")
	}
}

func TestGetUsername(t *testing.T) {
	goodInput := "good_username"
	goodLib := getLib(goodInput)
	username := goodLib.collectUsername("test_host")
	if username != goodInput {
		t.Errorf("getUsername() = %s; want %s", username, goodInput)
	}
	longUsername := "this_is_a_really_long_username_that_is_not_allowed"
	longLib := getLib(longUsername)
	username = longLib.collectUsername("test_host")
	if username != "" {
		t.Errorf("getUsername() = %s; want %s", username, "")
	}
}

func TestGetPrivateKey(t *testing.T) {
	goodInput := "/dev/stdout"
	goodLib := getLib(goodInput)
	privateKey := goodLib.collectPrivateKey()
	if privateKey != goodInput {
		t.Errorf("getPrivateKey() = %s; want %s", privateKey, goodInput)
	}
	badInput := "~/bad/path"
	badLib := getLib(badInput)
	badLib.InputSource = bytes.NewBufferString(badInput)
	privateKey = badLib.collectPrivateKey()
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

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}
