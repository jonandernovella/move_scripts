package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
)

type RsyncParameters struct {
	userName   string
	targetHost string
	targetDir  string
	privateKey string
}

type FileInfo struct {
	Name string
	Size int64
}

type Lib struct {
	Name        string
	InputSource io.Reader
	HomeDir     string
}

func (lib Lib) check() {
	dirToMove := lib.getDirectoryToMove()
	fmt.Printf("Checking %s\n\n", dirToMove)
	uncompressedFileExtensions := []string{".sam", ".vcf", ".fq", ".fastq", ".fasta", ".txt", ".fa"}
	lib.findUncompressedFiles(dirToMove, uncompressedFileExtensions)
}

func (lib Lib) gen() {

	fmt.Println("This script will generate a SLURM script to transfer data to Dardel.")

	dirToMove := lib.getDirectoryToMove()

	targetHost := lib.getTargetHost()

	targetDir := lib.getTargetDirectory(targetHost)

	username := lib.getUsername(targetHost)

	privateKey := lib.getPrivateKey()

	rsyncParameters := RsyncParameters{username, targetHost, targetDir, privateKey}

	projectId := lib.getProjectId()

	lib.writeScriptFile(dirToMove, projectId, rsyncParameters)
}

func (lib Lib) findUncompressedFiles(root string, extensions []string) {
	listOfFileInfos := make([]FileInfo, 0)
	sizeSumOfUncompressedFiles := int64(0)
	createUncompressedFileLog := false

	err := filepath.WalkDir(root, func(fileName string, d fs.DirEntry, e error) error {
		if e != nil {
			fmt.Printf("Error accessing a path %q: %v\n", fileName, e)
			return e
		}
		if slices.Contains(extensions, filepath.Ext(d.Name())) {
			createUncompressedFileLog = true
			info, err := d.Info()
			if err != nil {
				fmt.Println("Error getting file size:", err)
				os.Exit(1)
			}
			fileSize := info.Size()
			if fileSize > 1024*1024*1024 {
				fmt.Printf("WARNING: %s is %s. This may take a while to transfer.\n", fileName, formatBytes(fileSize))
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
		fmt.Println("WARNING: Uncompressed files found. A log file containing a list of uncompressed files will be created.")
		sort.Slice(listOfFileInfos, func(i, j int) bool {
			return listOfFileInfos[i].Size > listOfFileInfos[j].Size
		})
		logName := fmt.Sprintf("./%s_%s.uncompressed_files.log", lib.Name, filepath.Base(root))
		createFileLog(listOfFileInfos, logName)
	} else {
		fmt.Println("No uncompressed files found.")
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
	fmt.Print("A list of uncompressed files has been created in ", logName, ".\n")
}

func (lib Lib) getDirectoryToMove() string {
	dirToMove := ""
	for dirToMove == "" {
		dirToMove = lib.collectDirectoryToMove()
	}
	return dirToMove
}

func (lib Lib) collectDirectoryToMove() string {
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %s\n", err.Error())
		os.Exit(1)
	}
	dirToMove := lib.getInput("Which directory should be transferred? [default: this one]", workingDir)
	dirToMove, err = getAbsoluteDirectory(dirToMove)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		dirToMove = ""
	}
	return dirToMove
}

func (lib Lib) getTargetHost() string {
	return lib.getInput("Which system should data be moved to? [default: dardel.pdc.kth.se]", "dardel.pdc.kth.se")
}

func (lib Lib) getProjectId() string {
	return lib.getInput("uppmax project id to run the migration job (ex. nais2023-22-999)", "UPPMAX_PROJECT_ID")
}

func (lib Lib) getPrivateKey() string {
	privateKeyPath := ""
	for privateKeyPath == "" {
		privateKeyPath = lib.collectPrivateKey()
	}
	privateKeyAbsPath, err := filepath.Abs(privateKeyPath)
	if err != nil {
		fmt.Printf("Error converting path to absolute: %s\n", err.Error())
		os.Exit(1)
	}
	return privateKeyAbsPath
}

func (lib Lib) collectPrivateKey() string {
	privateKeyPath := lib.getInput("Which private key would you like to use?", fmt.Sprintf("%s/.ssh/id_rsa", lib.HomeDir))
	_, err := os.Stat(privateKeyPath)
	if os.IsNotExist(err) {
		fmt.Printf("Error: %s\n", errors.New("Private key does not exist: "+privateKeyPath))
		privateKeyPath = ""
	}
	return privateKeyPath
}

func (lib Lib) getTargetDirectory(targetHost string) string {
	targetDir := ""
	for targetDir == "" {
		targetDir = lib.collectTargetDir(targetHost)
	}
	return targetDir
}

func (lib Lib) collectTargetDir(targetHost string) string {
	targetDir := lib.getInput("Where on "+targetHost+" should data be moved to?", "")

	isAbsolute := filepath.IsAbs(targetDir)

	if !isAbsolute {
		fmt.Printf("Error: %s\n", errors.New("Path is not absolute: "+targetDir))
		targetDir = ""
	}
	return targetDir
}

func (lib Lib) getUsername(targetHost string) string {
	username := ""
	for username == "" {
		username = lib.collectUsername(targetHost)
	}
	return username
}

func (lib Lib) collectUsername(targetHost string) string {
	username := lib.getInput("What is your user name on "+targetHost+"?", "")
	if len(username) > 25 {
		fmt.Println("Error: Username must be 25 characters or less.")
		username = ""
	}
	return username
}

func (lib Lib) getInput(prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	scanner := bufio.NewScanner(lib.InputSource)
	scanner.Scan()
	input := scanner.Text()
	if input == "" {
		return defaultValue
	}
	return input
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

func (lib Lib) writeScriptFile(dirToMove, projectId string, transferParameters RsyncParameters) {
	scriptName := lib.Name + "_" + filepath.Base(dirToMove) + ".sh"
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

	scriptFile.WriteString("rsync -cavz -e " + "'ssh -i " + transferParameters.privateKey + "' --progress " + dirToMove + " " + transferParameters.userName + "@" + transferParameters.targetHost + ":" + transferParameters.targetDir + " | tee " + lib.Name + "_" + filepath.Base(dirToMove) + ".rsync_log\n")

	fmt.Println("\nWhen you are ready, edit", scriptName, "to set the correct project ID and run \"sbatch", scriptName, "\".")
}
