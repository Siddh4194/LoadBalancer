package ai

import "fmt"

type PromptType struct {
	HealthIssue  string
	Optimization string
}

type InfrastructureData struct {
	ServerLogs        string
	LoadBalancingAlgo string
	HealthIssue       string
}

func BuildAnalysisPrompt(data InfrastructureData) string {

	prompt := fmt.Sprintf(`
You're a load balancer log analyzer and panick recovery assistant. Your task is to analyze the provided infrastructure data, identify potential issues, and suggest improvements for load balancing and system performance.
Health Issue:
%s

Server Logs:
%s

Current Load Balancing Algorithm:
%s

Tasks:
1. Detect infrastructure issues
2. Identify unhealthy servers
3. Suggest load balancing improvements / remidication or why that server is unhealthy
`,
		data.HealthIssue,
		data.ServerLogs,
		data.LoadBalancingAlgo,
	)

	return prompt
}

func BuildRemediationPrompt(data string) string {
	prompt := fmt.Sprintf(`
		You're a load balancer log analyzer and panick recovery assistant. Your task is to analyze the provided infrastructure data, identify potential issues, and suggest improvements for load balancing and system performance.erver1
		commands to restart the server 1 "sudo systemctl restart server1_lb"
		commands to restart the server 2 "sudo systemctl restart server2_lb"
		commands to restart the server 3 "sudo systemctl restart server3_lb"

		Identified Issue:
		%s

		Tasks:
		1. By understanding the above issue please provide the respected command to start the server and the output should be like this 
		example output :{
			command:"command"
		}

		Rules : 
		1. Return JSON output only
		2. Do not provide any explanation or text other than the JSON output
		3. The commands should be in the arrays like below example 
		[
			{
				"command": "sudo systemctl restart server1_lb",
			},
			{
				"command": "sudo systemctl restart server2_lb",
			},
			{
				"command": "sudo systemctl restart server3_lb",
			}
		]
			4. Only add the commands which are necessary to fix the issue and do not add unnecessary commands
			4. Return the commands which are required not needed any extra commands so please add the required commands only
`,
		data,
	)

	return prompt
}



func BuildTheReportLikeWhatHappened(data string) string {
	prompt := fmt.Sprintf(`
		You're a load balancer auto action reporter

		Actions taken:
		%s

		Tasks:
		1. Give the proper report of what actions were taken and give the report response in the html format so that Directly sent over the mail

		Rules : 
		1. Return HTML only
		2. Do not provide any explanation or text other than the HTML output
		
`,
		data,
	)

	return prompt
}
