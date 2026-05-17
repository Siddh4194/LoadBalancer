package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func ReadLogsFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}

	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
		return "", err
	}

	lastN := 20
	start := len(lines) - lastN
	if start < 0 {
		start = 0
	}

	latest := strings.Join(lines[start:], "\n")
	fmt.Println("Latest Logs:")
	fmt.Println(latest)

	return latest, nil
}

func ParseAIResponse(response string) ([]CommandRequest, error) {

	var commands []CommandRequest

	err := json.Unmarshal([]byte(response), &commands)

	if err != nil {
		return nil, err
	}

	return commands, nil
}


func AnalyzeInfrastructure(healthCheckIssue string) {
	logs, err := ReadLogsFromFile("logs/app.log")
	if err != nil {
		fmt.Println("Error reading logs:", err)
		return
	}

	// ai step 1: analyze the logs and issue to build a prompt for Ollama
	prompt := BuildAnalysisPrompt(
		InfrastructureData{
			ServerLogs:        logs,
			LoadBalancingAlgo: "Round Robin",
			HealthIssue:       healthCheckIssue,
		},
	)

	// call Ollama on the specified server and model
	endpoint := "http://192.168.1.49:11434"
	model := "llama3.2:3b"

	resp1, err := CallOllama(endpoint, model, prompt)
	if err != nil {
		fmt.Println("Ollama error:", err)
		return
	}

	fmt.Println("1st step response:")
	fmt.Println(resp1)

	// step 2: print the response from Ollama which contains the analysis and recommendations for the infrastructure issue
	prompt = BuildRemediationPrompt(resp1)
	resp2, err := CallOllama(endpoint, model, prompt)
	if err != nil {
		fmt.Println("Ollama error:", err)
		return
	}

	fmt.Println("2nd step response:")
	fmt.Println(resp2)
	commands,err := ParseAIResponse(resp2)
	if err != nil {
		fmt.Println("Error parsing AI response:", err)
		return
	}

	fmt.Println("Parsed Commands:")
	exeCommands := make([]string, len(commands))
	for i, cmd := range commands {
		exeCommands[i] = cmd.Command
		fmt.Printf("Command: %s\n", cmd.Command)
	}

	results := ExecuteMultipleCommandShells(exeCommands)
	fmt.Println("Command Execution Results:", results)

	reportResponse := BuildTheReportLikeWhatHappened(resp1 + "\n" + resp2)

	msg := SendHTMLReportSMTP("AI Analysis Report", reportResponse)

	fmt.Print(msg)
}

// CallOllama sends the prompt to an Ollama server and returns the response body as text.
func CallOllama(endpoint, model, prompt string) (string, error) {
	url := endpoint + "/api/generate"

	payload := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// handle streaming / chunked responses by reading line-by-line
	var builder strings.Builder
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// some servers prefix streaming lines with "data:"
		if strings.HasPrefix(line, "data:") {
			line = strings.TrimSpace(line[len("data:"):])
			if line == "[DONE]" {
				break
			}
		}

		// try to parse JSON chunk, else append raw
		var chunk interface{}
		if err := json.Unmarshal([]byte(line), &chunk); err == nil {
			// if it's an object, try common fields
			if m, ok := chunk.(map[string]interface{}); ok {
				// common field names: "content", "text", "output", "response"
				if v, ok := m["content"]; ok {
					builder.WriteString(fmt.Sprint(v))
					continue
				}
				if v, ok := m["text"]; ok {
					builder.WriteString(fmt.Sprint(v))
					continue
				}
				if v, ok := m["output"]; ok {
					builder.WriteString(fmt.Sprint(v))
					continue
				}
				if v, ok := m["response"]; ok {
					builder.WriteString(fmt.Sprint(v))
					continue
				}
				// some formats put text in choices -> [{"text": "..."}]
				if choices, ok := m["choices"].([]interface{}); ok {
					for _, c := range choices {
						if cm, ok := c.(map[string]interface{}); ok {
							if t, ok := cm["text"]; ok {
								builder.WriteString(fmt.Sprint(t))
							}
						}
					}
					continue
				}
			}
			// fallback: append the raw unmarshalled value
			builder.WriteString(fmt.Sprint(chunk))
		} else {
			// not JSON, append raw line
			builder.WriteString(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return builder.String(), nil
}
