package main

import (
	"bytes"
	"os"
	"testing"
)

func TestGetDirectoryToMove(t *testing.T) {
	badDir := "/test/dir/bad"
	lib.InputSource = bytes.NewBufferString(badDir)
	dirToMove := collectDirectoryToMove()
	if dirToMove != "" {
		t.Errorf("getDirectoryToMove() = %s; want %s", dirToMove, badDir)
	}
	goodDir := "/dev"
	lib.InputSource = bytes.NewBufferString(goodDir)
	dirToMove = collectDirectoryToMove()
	if dirToMove != goodDir {
		t.Errorf("getDirectoryToMove() = %s; want %s", dirToMove, goodDir)
	}
	file := "/dev/null"
	lib.InputSource = bytes.NewBufferString(file)
	dirToMove = collectDirectoryToMove()
	if dirToMove != "" {
		t.Errorf("getDirectoryToMove() = %s; want %s", dirToMove, "")
	}
}
func TestGetTargetHost(t *testing.T) {
	testInput := "test_host"
	lib.InputSource = bytes.NewBufferString(testInput)
	host := getTargetHost()
	if host != testInput {
		t.Errorf("getTargetHost() = %s; want %s", host, testInput)
	}
	emptyInput := ""
	lib.InputSource = bytes.NewBufferString(emptyInput)
	host = getTargetHost()
	if host != "dardel.pdc.kth.se" {
		t.Errorf("getTargetHost() = %s; want %s", host, "dardel.pdc.kth.se")
	}
}

func TestGetProjectId(t *testing.T) {
	goodInput := "good_project_id"
	lib.InputSource = bytes.NewBufferString(goodInput)
	projectId := getProjectId()

	if projectId != goodInput {
		t.Errorf("getProjectId() = %s; want %s", projectId, goodInput)
	}

}

func TestGetTargetDirectory(t *testing.T) {
	absolutePath := "/good/target_dir"
	lib.InputSource = bytes.NewBufferString(absolutePath)
	targetDir := collectTargetDir("test_host")
	if targetDir != absolutePath {
		t.Errorf("getTargetDirectory() = %s; want %s", targetDir, absolutePath)
	}
	relativePath := "bad_target_dir"
	lib.InputSource = bytes.NewBufferString(relativePath)
	targetDir = collectTargetDir("test_host")

	if targetDir != "" {
		t.Errorf("getTargetDirectory() = %s; want %s", targetDir, "")
	}
}

func TestGetUsername(t *testing.T) {
	goodInput := "good_username"
	lib.InputSource = bytes.NewBufferString(goodInput)
	username := collectUsername("test_host")
	if username != goodInput {
		t.Errorf("getUsername() = %s; want %s", username, goodInput)
	}
	longUsername := "this_is_a_really_long_username_that_is_not_allowed"
	lib.InputSource = bytes.NewBufferString(longUsername)
	username = collectUsername("test_host")
	if username != "" {
		t.Errorf("getUsername() = %s; want %s", username, "")
	}
}

func TestGetPrivateKey(t *testing.T) {
	goodInput := "/dev/stdout"
	lib.InputSource = bytes.NewBufferString(goodInput)
	privateKey := collectPrivateKey()
	if privateKey != goodInput {
		t.Errorf("getPrivateKey() = %s; want %s", privateKey, goodInput)
	}
	badInput := "~/bad/path"
	lib.InputSource = bytes.NewBufferString(badInput)
	privateKey = collectPrivateKey()
	if privateKey != "" {
		t.Errorf("getPrivateKey() = %s; want %s", privateKey, "")
	}
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}
