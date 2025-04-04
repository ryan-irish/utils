package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
)

const (
	question1 = "1. What did you do yesterday?"
	question2 = "2. What will you do today?"
	question3 = "3. Anything blocking your progress?"
)

func main() {
	// Get Slack token and thread details from environment or user input
	token := getEnvOrPrompt("SLACK_TOKEN", "Enter your Slack API token: ")
	
	// Check if user wants to use a Slack link
	var channelID, threadTS string
	fmt.Println("Do you have a Slack message link? (y/n)")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	
	if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
		fmt.Println("Enter the Slack message link (e.g., ******.slack.com/archives/C048ECCB75H/p1743724813501239):")
		link, _ := reader.ReadString('\n')
		link = strings.TrimSpace(link)
		
		var err error
		channelID, threadTS, err = parseSlackLink(link)
		if err != nil {
			fmt.Printf("Error parsing Slack link: %v\n", err)
			fmt.Println("Falling back to manual entry...")
			channelID = getEnvOrPrompt("SLACK_CHANNEL_ID", "Enter the channel ID: ")
			threadTS = getEnvOrPrompt("SLACK_THREAD_TS", "Enter the thread timestamp: ")
		} else {
			fmt.Printf("Successfully parsed link. Channel ID: %s, Thread TS: %s\n", channelID, threadTS)
		}
	} else {
		channelID = getEnvOrPrompt("SLACK_CHANNEL_ID", "Enter the channel ID: ")
		threadTS = getEnvOrPrompt("SLACK_THREAD_TS", "Enter the thread timestamp: ")
	}

	// Initialize Slack API client
	api := slack.New(token)
	
	// Get answers from CLI
	answers := make(map[string]string)
	
	fmt.Println("Enter multiple bullet points per question. Press Enter twice when done with a question.")
	
	// Ask each question
	answers[question1] = askQuestion(question1)
	answers[question2] = askQuestion(question2)
	answers[question3] = askQuestion(question3)
	
	// Format message
	message := formatStandupMessage(answers)
	
	// Post to Slack thread
	_, _, err := api.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
		slack.MsgOptionTS(threadTS),
	)
	
	if err != nil {
		fmt.Printf("Error posting message: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Standup posted successfully!")
}

// parseSlackLink extracts channel ID and thread timestamp from a Slack link
func parseSlackLink(link string) (string, string, error) {
	// Match pattern like: ******.slack.com/archives/C048ECCB75H/p1743724813501239
	re := regexp.MustCompile(`/archives/([A-Z0-9]+)/p(\d+)`)
	matches := re.FindStringSubmatch(link)
	
	if len(matches) != 3 {
		return "", "", fmt.Errorf("invalid Slack link format")
	}
	
	channelID := matches[1]
	timestampStr := matches[2]
	
	// Format timestamp from p1743724813501239 to 1743724813.501239
	if len(timestampStr) < 7 {
		return "", "", fmt.Errorf("invalid timestamp format")
	}
	
	// Insert a period 6 characters from the end
	mainPart := timestampStr[:len(timestampStr)-6]
	fractionPart := timestampStr[len(timestampStr)-6:]
	threadTS := mainPart + "." + fractionPart
	
	return channelID, threadTS, nil
}

// askQuestion prompts the user with a question and returns the answer
func askQuestion(question string) string {
	fmt.Println(question)
	fmt.Println("(Enter each bullet point on a new line. Press Enter twice when done.)")
	fmt.Print("> ")
	
	reader := bufio.NewReader(os.Stdin)
	var lines []string
	
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		
		line = strings.TrimSpace(line)
		
		if line == "" {
			break
		}
		
		lines = append(lines, line)
		fmt.Print("> ")
	}
	
	return strings.Join(lines, "\n")
}

// getEnvOrPrompt returns environment variable or prompts user for input
func getEnvOrPrompt(envName, prompt string) string {
	if value, exists := os.LookupEnv(envName); exists && value != "" {
		return value
	}
	
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}
	
	return strings.TrimSpace(value)
}

// formatStandupMessage formats the answers into a Slack message
func formatStandupMessage(answers map[string]string) string {
	var builder strings.Builder
	
	builder.WriteString(question1 + "\n")
	for _, line := range strings.Split(answers[question1], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Check if line already starts with a bullet point
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			builder.WriteString(line + "\n")
		} else {
			builder.WriteString("- " + line + "\n")
		}
	}
	
	builder.WriteString("\n" + question2 + "\n")
	for _, line := range strings.Split(answers[question2], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Check if line already starts with a bullet point
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			builder.WriteString(line + "\n")
		} else {
			builder.WriteString("- " + line + "\n")
		}
	}
	
	builder.WriteString("\n" + question3 + "\n")
	for _, line := range strings.Split(answers[question3], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Check if line already starts with a bullet point
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			builder.WriteString(line + "\n")
		} else {
			builder.WriteString("- " + line + "\n")
		}
	}
	
	return builder.String()
} 