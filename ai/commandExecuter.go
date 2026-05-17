package ai

import (
	"bytes"
	"os/exec"
	"sync"
)

// CommandResult holds the output and error info for a single command execution
type CommandResult struct {
	Command string
	Output  string
	Error   string
	Success bool
}

// CommandRequest holds a command to be executed
type CommandRequest struct {
	Command string
	Args    []string
}

// ExecuteCommand executes a single command and returns the result
func ExecuteCommand(cmd string, args []string) CommandResult {
	result := CommandResult{
		Command: cmd,
	}

	command := exec.Command(cmd, args...)
	output, err := command.CombinedOutput()

	result.Output = string(output)
	if err != nil {
		result.Error = err.Error()
		result.Success = false
	} else {
		result.Success = true
	}

	return result
}

// ExecuteCommandShell executes a shell command string and returns the result
func ExecuteCommandShell(cmdStr string) CommandResult {

	result := CommandResult{
		Command: cmdStr,
	}

	password := "1418"

	cmd := exec.Command("sudo", "-S", "bash", "-c", cmdStr)

	// Pass password to sudo
	cmd.Stdin = bytes.NewBufferString(password + "\n")

	output, err := cmd.CombinedOutput()

	result.Output = string(output)

	if err != nil {
		result.Error = err.Error()
		result.Success = false
	} else {
		result.Success = true
	}

	return result
}

// ExecuteMultipleCommands executes multiple commands sequentially and returns all results
func ExecuteMultipleCommands(commands []CommandRequest) []CommandResult {

	results := make([]CommandResult, len(commands))

	for i, cmd := range commands {

		results[i] = ExecuteCommand(cmd.Command, cmd.Args)
	}

	return results
}

// ExecuteMultipleCommandsParallel executes multiple commands in parallel and returns all results
func ExecuteMultipleCommandsParallel(commands []CommandRequest) []CommandResult {
	results := make([]CommandResult, len(commands))
	var wg sync.WaitGroup

	for i, cmd := range commands {
		wg.Add(1)
		go func(index int, command CommandRequest) {
			defer wg.Done()
			results[index] = ExecuteCommand(command.Command, command.Args)
		}(i, cmd)
	}

	wg.Wait()
	return results
}

// ExecuteMultipleCommandShells executes multiple shell command strings sequentially
func ExecuteMultipleCommandShells(cmdStrings []string) []CommandResult {
	results := make([]CommandResult, len(cmdStrings))

	for i, cmdStr := range cmdStrings {
		results[i] = ExecuteCommandShell(cmdStr)
	}

	return results
}

// ExecuteMultipleCommandShellsParallel executes multiple shell commands in parallel
func ExecuteMultipleCommandShellsParallel(cmdStrings []string) []CommandResult {
	results := make([]CommandResult, len(cmdStrings))
	var wg sync.WaitGroup

	for i, cmdStr := range cmdStrings {
		wg.Add(1)
		go func(index int, command string) {
			defer wg.Done()
			results[index] = ExecuteCommandShell(command)
		}(i, cmdStr)
	}

	wg.Wait()
	return results
}
