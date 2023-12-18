package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type RsyncParameters struct {
	userName   string
	targetHost string
	targetDir  string
	numConns   string
}

type FileInfo struct {
	Name string
	Size int64
}

const MAX_FILES_PER_DIR int = 100000

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: movefiles [check|start]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "check":
		fmt.Println("Running in check mode, data will NOT be transferred.")
		dirToMove := getDirectoryToMove()
		uncompressedFileExtensions := []string{".sam", ".vcf", ".fq", ".fastq", ".fasta", ".txt", ".fa"}
		findUncompressedFiles(dirToMove, uncompressedFileExtensions)
	case "start":
		start()
	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}

func start() {

	fmt.Println("Welcome to this data transfer tool")

	dirToMove := getDirectoryToMove()

	fmt.Printf("Moving %s\n\n", dirToMove)

	fmt.Println("This tool will find all subdirectories with more than", MAX_FILES_PER_DIR, "files in them and package (tar) them before moving.")

	getKeepDirsMessage, keepDirs := getKeepDirs()
	fmt.Println(getKeepDirsMessage)

	autoDelMessage := getAutoDel()
	fmt.Println(autoDelMessage)

	targetHost := getTargetHost()

	targetDir := getTargetDirectory(targetHost)

	username := getUsername(targetHost)

	numConns := getNumConnections()

	rsyncParameters := RsyncParameters{username, targetHost, targetDir, numConns}

	projectId := getProjectId()

	writeScriptFile(dirToMove, keepDirs, projectId, rsyncParameters)
}

func findUncompressedFiles(root string, extensions []string) {
	listOfFileInfos := make([]FileInfo, 0)
	sizeSumOfUncompressedFiles := int64(0)
	createUncompressedFileLog := false

	err := filepath.WalkDir(root, func(fileName string, d fs.DirEntry, e error) error {
		if e != nil {
			fmt.Printf("Error accessing a path %q: %v\n", fileName, e)
			return e
		}
		if slices.Contains(extensions, filepath.Ext(d.Name())) {
			info, err := d.Info()
			if err != nil {
				fmt.Println("Error getting file size:", err)
				os.Exit(1)
			}
			fileSize := info.Size()
			if fileSize > 1024*1024*1024 {
				fmt.Printf("WARNING: %s is %s. This may take a while to transfer.\n", fileName, formatBytes(fileSize))
				createUncompressedFileLog = true
			}
			listOfFileInfos = append(listOfFileInfos, FileInfo{fileName, fileSize})
			sizeSumOfUncompressedFiles += fileSize
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err)
		fmt.Println("Aborting transfer!")
		os.Exit(1)
	}
	if sizeSumOfUncompressedFiles > 1024*1024*1024*100 {
		fmt.Println("WARNING: The total size of the ", len(listOfFileInfos)-1, " uncompressed files to be transferred is ", formatBytes(sizeSumOfUncompressedFiles), ". You might want to compress them before.")
		createUncompressedFileLog = true
	}
	if createUncompressedFileLog {
		sort.Slice(listOfFileInfos, func(i, j int) bool {
			return listOfFileInfos[i].Size > listOfFileInfos[j].Size
		})
		logName := fmt.Sprintf("./transfer_%s.uncompressed_files.log", filepath.Base(root))
		createFileLog(listOfFileInfos, logName)
	}
}

func createFileLog(listOfFileInfos []FileInfo, logName string) {

	umcompressedLog, err := os.Create(logName)
	defer umcompressedLog.Close()
	if err != nil {
		fmt.Println("Error creating file list:", err)
		os.Exit(1)
	}
	for _, fileInfo := range listOfFileInfos {
		umcompressedLog.WriteString(fmt.Sprintf("%s\t%s\n", fileInfo.Name, formatBytes(fileInfo.Size)))
	}
}

func getAutoDel() string {
	autoDel := askForBinaryInput("Do you wish to automatically delete local files after copying them? [y/N]", "N")
	var autoDelMessage string
	if autoDel == "Y" {
		autoDelMessage = "We will delete files that have been copying."
	} else {
		autoDelMessage = "We will keep files here after copying."
	}
	return autoDelMessage
}

func getKeepDirs() (string, string) {
	keepDirs := askForBinaryInput("Should we discard the large subdirectories after packaging? [Y/n]", "Y")

	var keepDirsMessage string
	if keepDirs == "N" {
		keepDirsMessage = "We will discard the big directories after packaging."
	} else {
		keepDirsMessage = "We will keep the big directories."
	}
	return keepDirsMessage, keepDirs
}

func getDirectoryToMove() string {
	dirToMove := ""
	for dirToMove == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting working directory: %s\n", err.Error())
			os.Exit(1)
		}
		dirToMove = getInput("Which directory should be transferred? [default: this one]", workingDir)
		dirToMove, err = getAbsoluteDirectory(dirToMove)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			dirToMove = ""
		}
	}
	return dirToMove
}

func getTargetHost() string {
	return getInput("Which system should data be moved to? [default: dardel.pdc.kth.se]", "dardel.pdc.kth.se")
}

func getProjectId() string {
	return getInput("uppmax project id (ex. nais2023-22-999)", "UPPMAX_PROJECT_ID")
}

func getTargetDirectory(targetHost string) string {
	targetDir := ""
	for targetDir == "" {
		targetDir = getInput("Where on "+targetHost+" should data be moved to?", "")

		isAbsolute := filepath.IsAbs(targetDir)

		if !isAbsolute {
			fmt.Printf("Error: %s\n", errors.New("Path is not absolute: "+targetDir))
			targetDir = ""
		}
	}
	return targetDir
}

func getUsername(targetHost string) string {
	username := ""
	for username == "" {
		username = getInput("What is your user name on "+targetHost+"?", os.Getenv("USER"))
		if len(username) > 25 {
			fmt.Println("Error: Username must be 25 characters or less.")
			username = ""
		}
	}
	return username
}

func getNumConnections() string {
	numConns := ""
	for numConns == "" {
		numConns = getInput("How many parallel rsync connections? [10]", "10")
		if _, err := strconv.Atoi(numConns); err != nil {
			fmt.Println("Invalid input. Please enter a valid number.")
			numConns = ""
		}
	}
	return numConns
}

func getInput(prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	if input == "" {
		return defaultValue
	}
	return input
}

func askForBinaryInput(prompt, defaultValue string) string {
	response := ""
	for {
		response = strings.ToUpper(getInput(prompt, defaultValue))

		if response == "Y" || response == "N" {
			break
		} else {
			fmt.Println("Please enter either 'Y' or 'N'.")
		}
	}
	return response
}

func getAbsoluteDirectory(path string) (string, error) {
	isAbsolute := filepath.IsAbs(path)
	err := error(nil)

	if !isAbsolute {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", errors.New("Error converting path to absolute: " + err.Error())
		}
		path = absPath
	}

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", errors.New("Path does not exist: " + path)
	}

	if !fileInfo.IsDir() {
		return "", errors.New("Path is not a directory: " + path)
	}
	return path, nil
}

func formatBytes(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func writeScriptFile(dirToMove, keepDirs, projectId string, transferParameters RsyncParameters) {
	scriptName := "transfer_" + filepath.Base(dirToMove) + ".sh"
	scriptFile, err := os.Create(scriptName)

	if err != nil {
		fmt.Printf("Error creating script file: %s\n", err)
		os.Exit(1)
	}
	defer scriptFile.Close()

	scriptFile.WriteString("#!/bin/bash\n")
	scriptFile.WriteString("#SBATCH -p core\n")
	scriptFile.WriteString("#SBATCH -n 1\n")
	scriptFile.WriteString("#SBATCH -J " + scriptName + "\n")
	scriptFile.WriteString("#SBATCH -A " + projectId + "\n")
	scriptFile.WriteString("#SBATCH -t 7-00:00:00\n\n")

	scriptFile.WriteString("find " + dirToMove + " -mindepth 1 -maxdepth 2 -not -path '*/.*' -type d -links " + fmt.Sprint(MAX_FILES_PER_DIR) + " > large_directories.txt\n\n")

	scriptFile.WriteString("xargs -a large_directories.txt -I {} tar -czvf {}.tar.gz {}\n\n")

	if keepDirs == "N" {
		scriptFile.WriteString("xargs -a large_directories.txt -I{} rm -rf {}\n\n")
	}

	scriptFile.WriteString("rsync -cavz --progress --parallel=" + fmt.Sprint(transferParameters.numConns) + " --exclude-from=large_directories.txt " + dirToMove + " " + transferParameters.userName + "@" + transferParameters.targetHost + ":" + transferParameters.targetDir + " | tee rsync_log.txt\n")

	fmt.Println("\nWhen you are ready, edit", scriptName, "to set the correct project ID and run \"sbatch", scriptName, "\".")
}
