package main

// Issue 01:
// Color is not preserved...

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var VERSION = "1.1.0"

func main() {
	// Parse options (if any)
	flag.Parse()

	// Remaining arguments are the command and its options
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: cmdlog [OPTION...] CMD [CMD OPTION...]")
		fmt.Println("cmdlog prints to stdout and stderr and logs to a logfile. Easy to use,")
		fmt.Println("easy to document.")
		fmt.Printf("\ncmdlog v%s was created by Malte Gejr <malte@gejr.dk>\n", VERSION)
		os.Exit(1)
	}

	// The first argument is the command name
	cmdName := args[0]
	cmdArgs := args[1:]

	// Generate the log filename based on the command and its options
	logFileName := generateLogFileName(cmdName, cmdArgs)

	// Create or open the log file (after ensuring unique file name)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Create a logger that writes to both stdout and the log file
	logger := log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)

	// Print the CMD and CMD OPTIONS to the log file and stdout
	cmdLine := fmt.Sprintf("%s %s", cmdName, strings.Join(cmdArgs, " "))
	fmt.Fprintf(logFile, "%s\n\n", cmdLine)

	// Set up the command to execute
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Fatalf("Error setting up stdout pipe: %v", err)
	}
	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Fatalf("Error setting up stderr pipe: %v", err)
	}
	cmdStdin, err := cmd.StdinPipe() // Create stdin pipe for the command
	if err != nil {
		logger.Fatalf("Error setting up stdin pipe: %v", err)
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		logger.Fatalf("Error starting command: %v", err)
	}

	// Create goroutines to copy the command's stdout and stderr to the logger
	go func() {
		_, _ = io.Copy(logger.Writer(), cmdStdout)
	}()
	go func() {
		_, _ = io.Copy(logger.Writer(), cmdStderr)
	}()

	// Create a goroutine to handle user input and send it to both the command's stdin and the log file
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// Write user input to the log file
			_, err := fmt.Fprintf(logFile, "%s\n", scanner.Text())
			if err != nil {
				logger.Printf("Error writing to log file: %v", err)
				break
			}

			// Write user input to the command's stdin
			_, err = cmdStdin.Write([]byte(scanner.Text() + "\n"))
			if err != nil {
				logger.Printf("Error writing to stdin: %v", err)
				break
			}
		}
		if scanner.Err() != nil {
			logger.Printf("Error reading from stdin: %v", scanner.Err())
		}
	}()

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
}

// generateLogFileName creates a log file name based on the command and its arguments.
// It replaces spaces with underscores and sanitizes the filename to avoid invalid characters.
// It checks if a file already exists and appends a number to the filename if necessary.
func generateLogFileName(cmdName string, cmdArgs []string) string {
	// Combine the command name and its arguments into a single string
	cmdString := cmdName + " " + strings.Join(cmdArgs, " ")

	// Sanitize the filename by replacing invalid characters
	cmdString = sanitizeFilename(cmdString)

	// Append ".log" to create the final log file name
	baseFileName := "cmdlog_" + cmdString + ".log"
	logFileName := baseFileName

	// Check if the log file already exists, and increment the file name if it does
	counter := 1
	for fileExists(logFileName) {
		logFileName = fmt.Sprintf("%s.%d.log", baseFileName[:len(baseFileName)-4], counter)
		counter++
	}

	return logFileName
}

// fileExists checks if a file exists at the given path.
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// sanitizeFilename replaces characters that are invalid in filenames with underscores.
func sanitizeFilename(s string) string {
	// Replace invalid characters with underscores or dashes
	re := regexp.MustCompile(`[\\\/:*?"<>|&^%$#@!~;,+\s]`)
	s = re.ReplaceAllString(s, "_")

	// Ensure that the filename doesn't start with an underscore or dash (if necessary)
	if len(s) > 0 && (s[0] == '_' || s[0] == '-') {
		s = "cmdlog_" + s[1:]
	}

	return s
}
