package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
	
	question1 = "1. What did you do yesterday?"
	question2 = "2. What will you do today?"
	question3 = "3. Anything blocking your progress?"
	
	// OAuth configuration
	clientID     = "" // To be filled by user
	clientSecret = "" // To be filled by user
	
	// OAuth scopes needed
	scopes = "chat:write,channels:read,im:write"
	
	// Default token config file location
	configDir  = ".slack-standup-updater"
	configFile = "token.json"
)

// Global configuration
var useColors = true

type TokenConfig struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	TeamID      string `json:"team_id"`
	Expiration  int64  `json:"expiration"`
}

// printInfo prints formatted informational messages
func printInfo(message string) {
	if useColors {
		fmt.Printf("%s%s==> %s%s\n", colorCyan, colorBold, message, colorReset)
	} else {
		fmt.Println("==> " + message)
	}
}

// printQuestion prints a formatted question
func printQuestion(message string) {
	if useColors {
		fmt.Printf("\n%s%s❓ %s%s\n", colorBlue, colorBold, message, colorReset)
	} else {
		fmt.Printf("\n❓ %s\n", message)
	}
}

// printPrompt prints a prompt for user input
func printPrompt(message string) {
	if useColors {
		fmt.Printf("%s%s👉 %s%s ", colorYellow, colorBold, message, colorReset)
	} else {
		fmt.Printf("👉 %s ", message)
	}
}

// printSuccess prints a success message
func printSuccess(message string) {
	if useColors {
		fmt.Printf("%s%s✅ %s%s\n", colorGreen, colorBold, message, colorReset)
	} else {
		fmt.Printf("✅ SUCCESS: %s\n", message)
	}
}

// printError prints an error message
func printError(message string) {
	if useColors {
		fmt.Printf("%s%s❌ Error: %s%s\n", colorRed, colorBold, message, colorReset)
	} else {
		fmt.Printf("❌ ERROR: %s\n", message)
	}
}

// printHeader prints a section header
func printHeader(message string) {
	if useColors {
		fmt.Printf("\n%s%s🔹 === %s ===%s\n", colorMagenta, colorBold, message, colorReset)
	} else {
		fmt.Printf("\n🔹 === %s ===\n", message)
	}
}

// printDivider prints a divider line
func printDivider() {
	if useColors {
		fmt.Printf("%s%s✨ ----------------------------------------- ✨%s\n", colorWhite, colorBold, colorReset)
	} else {
		fmt.Println("✨ ----------------------------------------- ✨")
	}
}

func main() {
	// Check if colors should be disabled
	if _, exists := os.LookupEnv("NO_COLOR"); exists {
		useColors = false
	}
	
	printHeader("Slack Standup Updater 🚀")
	
	// Set up user authentication
	token, err := getUserToken()
	if err != nil {
		printError(fmt.Sprintf("Getting user token: %v", err))
		os.Exit(1)
	}
	
	// Get thread details
	var channelID, threadTS string
	
	printHeader("Thread Selection 🧵")
	printInfo("Where do you want to post your standup?")
	printInfo("1. Reply to a thread in a channel (y)")
	printInfo("2. Message yourself directly (m)")
	printInfo("3. Post to Slackbot (s) - most reliable way to message yourself")
	printInfo("4. Message any user by ID (u) - works with all token types")
	printPrompt(">")
	
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	
	// Initialize Slack API client (moved earlier to use for DM channel lookup)
	api := slack.New(token)
	
	if strings.ToLower(answer) == "u" || strings.ToLower(answer) == "user" {
		printInfo("Sending a direct message to a specific user 👥")
		
		// Get the user ID to message
		userID := getInput("Enter the User ID to message (starts with U)")
		
		// Open a conversation with this user
		channel, _, _, err := api.OpenConversation(&slack.OpenConversationParameters{
			Users: []string{userID},
		})
		
		if err != nil {
			printError(fmt.Sprintf("Opening conversation with user: %v", err))
			printInfo("Falling back to manual channel ID entry...")
			channelID = getInput("Enter the DM channel ID for this user (starts with D)")
		} else {
			channelID = channel.ID
			printSuccess(fmt.Sprintf("Found DM channel with user %s: %s", userID, channelID))
		}
		
		threadTS = "" // No thread, just post to the DM
	} else if strings.ToLower(answer) == "s" || strings.ToLower(answer) == "slackbot" {
		printInfo("Sending to Slackbot 🤖")
		
		// Try to find the Slackbot channel
		slackbotFound := false
		
		// List conversations that include Slackbot
		params := &slack.GetConversationsParameters{
			Types: []string{"im"},
			Limit: 200,
		}
		
		channels, _, err := api.GetConversations(params)
		if err != nil {
			printError(fmt.Sprintf("Getting conversations: %v", err))
		} else {
			// Look for a channel that might be Slackbot
			for _, channel := range channels {
				// Slackbot channel typically has the name "slackbot"
				if strings.ToLower(channel.Name) == "slackbot" {
					channelID = channel.ID
					slackbotFound = true
					printSuccess(fmt.Sprintf("Found Slackbot channel: %s", channelID))
					break
				}
			}
		}
		
		// If we couldn't find Slackbot, ask the user
		if !slackbotFound {
			printInfo("Couldn't automatically find your Slackbot channel.")
			printInfo("To find your Slackbot channel ID:")
			printInfo("1. Open Slack in a browser")
			printInfo("2. Click on Slackbot in the sidebar")
			printInfo("3. The URL will contain your Slackbot channel ID, like: /messages/DXXXXXXXX")
			channelID = getInput("Enter your Slackbot channel ID (starts with D)")
		}
		
		threadTS = "" // No thread, just post to Slackbot
	} else if strings.ToLower(answer) == "m" || strings.ToLower(answer) == "myself" || strings.ToLower(answer) == "me" {
		printInfo("Sending to your Slack DM (yourself) 👤")
		
		// Get the user's own identity
		userInfo, err := api.AuthTest()
		if err != nil {
			printError(fmt.Sprintf("Getting user info: %v", err))
			os.Exit(1)
		}
		
		// Try different approaches to message yourself
		var msgChannel string
		
		// First approach - try to open a DM with yourself
		channel, _, _, err := api.OpenConversation(&slack.OpenConversationParameters{
			Users: []string{userInfo.UserID},
		})
		
		if err != nil {
			// If that fails (likely because it's a bot token), try a different approach
			if strings.Contains(err.Error(), "cannot_dm_bot") || strings.Contains(err.Error(), "missing_scope") {
				printInfo("Bot tokens can't DM themselves. Trying alternative approach...")
				
				// For bot tokens, we'll try to find the user's DM channel ID
				// First, we'll list conversations (DMs) that the bot has access to
				params := &slack.GetConversationsParameters{
					Types: []string{"im"},
					Limit: 200, // Get a reasonable number of DMs
				}
				
				channels, _, err := api.GetConversations(params)
				if err != nil {
					printError(fmt.Sprintf("Getting conversations: %v", err))
					printInfo("Using Slackbot as a fallback option...")
					msgChannel = "D01" // This typically works as a fallback "Slackbot" channel
				} else {
					// Look for a channel that might be the user's DM
					// This is a bit hacky but should work most of the time
					if len(channels) > 0 {
						msgChannel = channels[0].ID
						printInfo(fmt.Sprintf("Found potential DM channel: %s", msgChannel))
					} else {
						printInfo("No DM channels found. Asking for manual input...")
						msgChannel = getInput("Enter your own user/channel ID to message")
					}
				}
			} else {
				// Some other error
				printError(fmt.Sprintf("Opening DM channel: %v", err))
				msgChannel = getInput("Enter your own user/channel ID to message")
			}
		} else {
			msgChannel = channel.ID
		}
		
		channelID = msgChannel
		threadTS = "" // No thread, just post to the DM channel
		
		printSuccess(fmt.Sprintf("Will post to DM. User ID: %s, Channel ID: %s", userInfo.UserID, channelID))
	} else if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" || strings.ToLower(answer) == "channel" {
		printInfo("Do you have a Slack message link? (y/n)")
		printPrompt(">")
		
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		
		if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
			printInfo("Enter the Slack message link (e.g., ******.slack.com/archives/C048ECCB75H/p1743724813501239):")
			printPrompt(">")
			link, _ := reader.ReadString('\n')
			link = strings.TrimSpace(link)
			
			var err error
			channelID, threadTS, err = parseSlackLink(link)
			if err != nil {
				printError(fmt.Sprintf("Parsing Slack link: %v", err))
				printInfo("Falling back to manual entry...")
				channelID = getInput("Enter the channel ID")
				threadTS = getInput("Enter the thread timestamp")
			} else {
				printSuccess(fmt.Sprintf("Successfully parsed link. Channel ID: %s, Thread TS: %s", channelID, threadTS))
			}
		} else {
			channelID = getInput("Enter the channel ID")
			threadTS = getInput("Enter the thread timestamp")
		}
	} else {
		// Default to asking for channel info if the input wasn't recognized
		printInfo("Defaulting to channel thread...")
		channelID = getInput("Enter the channel ID")
		threadTS = getInput("Enter the thread timestamp")
	}

	// Get answers from CLI
	answers := make(map[string]string)
	
	printHeader("Standup Questions 📋")
	printInfo("Enter multiple bullet points per question. Press Enter twice when done with a question.")
	printDivider()
	
	// Ask each question
	answers[question1] = askQuestion(question1)
	answers[question2] = askQuestion(question2)
	answers[question3] = askQuestion(question3)
	
	// Format message
	message := formatStandupMessage(answers)
	
	printHeader("Posting to Slack 💬")
	printInfo("Sending your standup message...")
	
	// Post to Slack
	var postErr error
	
	if threadTS != "" {
		// Posting to a thread
		printInfo("Posting to thread in channel...")
		_, _, postErr = api.PostMessage(
			channelID,
			slack.MsgOptionText(message, false),
			slack.MsgOptionTS(threadTS),
			slack.MsgOptionAsUser(true), // This posts as the authenticated user
		)
	} else {
		// Posting to a DM or channel (not as a thread reply)
		printInfo("Posting direct message...")
		_, _, postErr = api.PostMessage(
			channelID,
			slack.MsgOptionText(message, false),
			slack.MsgOptionAsUser(true), // This posts as the authenticated user
		)
	}
	
	if postErr != nil {
		printError(fmt.Sprintf("Posting message: %v", postErr))
		os.Exit(1)
	}
	
	printDivider()
	printSuccess("Standup posted successfully! 🎉")
}

// getUserToken gets the user token from config or initiates OAuth flow
func getUserToken() (string, error) {
	// Check if we have a valid token file
	config, err := readTokenConfig()
	if err == nil && config.AccessToken != "" && (config.Expiration == 0 || time.Now().Unix() < config.Expiration) {
		printInfo("Using saved authentication token 🔑")
		return config.AccessToken, nil
	}
	
	// Need to get a new token via OAuth
	printHeader("Slack Authentication 🔒")
	printInfo("You need to authenticate with Slack.")
	printInfo("This tool will open your browser to authorize access to your Slack account.")
	printInfo("Please ensure you have provided your clientID and clientSecret in the code or environment variables.")
	
	// Check if using ngrok
	printInfo("For secure HTTPS connection, you have two options:")
	printInfo("1. Use an ngrok URL (enter the https:// URL from ngrok)")
	printInfo("2. Use localhost with HTTPS (enter 'https://localhost:1337')")
	printInfo("3. Use localhost with HTTP (enter 'http://localhost:1337')")
	printPrompt("Enter the callback URL base (or press Enter for https://localhost:1337)")
	reader := bufio.NewReader(os.Stdin)
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)
	
	// Default callback URL
	callbackURL := "https://localhost:1337/callback"
	if baseURL != "" {
		// Make sure it doesn't include the /callback part
		baseURL = strings.TrimSuffix(baseURL, "/")
		baseURL = strings.TrimSuffix(baseURL, "/callback")
		callbackURL = baseURL + "/callback"
	}
	
	printInfo("Using callback URL: " + callbackURL)
	
	// Get clientID and clientSecret from env vars or prompt
	cID := getEnvOrDefault("SLACK_CLIENT_ID", clientID)
	cSecret := getEnvOrDefault("SLACK_CLIENT_SECRET", clientSecret)
	
	if cID == "" {
		cID = getInput("Enter your Slack Client ID")
	}
	
	if cSecret == "" {
		cSecret = getInput("Enter your Slack Client Secret")
	}
	
	// Create a random state for security
	state := fmt.Sprintf("%d", time.Now().UnixNano())
	
	// Start local server to receive the OAuth callback
	tokenChan := make(chan string)
	errorChan := make(chan error)
	go startOAuthServer(tokenChan, errorChan, state, cID, cSecret, callbackURL)
	
	// Construct the authorize URL
	authURL := fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&scope=%s&state=%s&redirect_uri=%s",
		cID, scopes, state, callbackURL,
	)
	
	// Open browser to the authorization URL
	err = openBrowser(authURL)
	if err != nil {
		printError(fmt.Sprintf("Could not open browser: %v", err))
		printInfo(fmt.Sprintf("Please open this URL in your browser:\n%s", authURL))
	}
	
	printInfo("Waiting for authentication... 🔄")
	
	// Wait for the token or error
	select {
	case token := <-tokenChan:
		// Save the token
		config := TokenConfig{
			AccessToken: token,
			Expiration:  0, // No expiration for user tokens
		}
		if err := saveTokenConfig(config); err != nil {
			printInfo(fmt.Sprintf("Warning: Could not save token: %v", err))
		}
		printSuccess("Authentication successful! 🎊")
		return token, nil
	case err := <-errorChan:
		return "", fmt.Errorf("OAuth error: %v", err)
	case <-time.After(5 * time.Minute):
		return "", fmt.Errorf("authentication timed out")
	}
}

// startOAuthServer starts a local HTTP server to handle the OAuth callback
func startOAuthServer(tokenChan chan string, errorChan chan error, state, clientID, clientSecret, callbackURL string) {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Check state to prevent CSRF
		if r.FormValue("state") != state {
			errorChan <- fmt.Errorf("invalid state parameter")
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}
		
		// Check for error
		if r.FormValue("error") != "" {
			errorChan <- fmt.Errorf("authorization error: %s", r.FormValue("error"))
			http.Error(w, fmt.Sprintf("Authorization error: %s", r.FormValue("error")), http.StatusBadRequest)
			return
		}
		
		// Get authorization code
		code := r.FormValue("code")
		if code == "" {
			errorChan <- fmt.Errorf("no code provided")
			http.Error(w, "No code provided", http.StatusBadRequest)
			return
		}
		
		// Exchange code for token
		tokenResp, err := http.PostForm("https://slack.com/api/oauth.v2.access", 
			url.Values{
				"client_id":     {clientID},
				"client_secret": {clientSecret},
				"code":          {code},
				"redirect_uri":  {callbackURL},
			},
		)
		if err != nil {
			errorChan <- fmt.Errorf("token exchange error: %v", err)
			http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
			return
		}
		defer tokenResp.Body.Close()
		
		// Parse token response
		var tokenData struct {
			Ok          bool   `json:"ok"`
			Error       string `json:"error,omitempty"`
			AccessToken string `json:"access_token"`
			Authed_User struct {
				ID string `json:"id"`
			} `json:"authed_user"`
			Team struct {
				ID string `json:"id"`
			} `json:"team"`
		}
		
		body, err := io.ReadAll(tokenResp.Body)
		if err != nil {
			errorChan <- fmt.Errorf("error reading response: %v", err)
			http.Error(w, "Failed to read token response", http.StatusInternalServerError)
			return
		}
		
		if err := json.Unmarshal(body, &tokenData); err != nil {
			errorChan <- fmt.Errorf("error parsing token response: %v", err)
			http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
			return
		}
		
		if !tokenData.Ok {
			errorChan <- fmt.Errorf("Slack API error: %s", tokenData.Error)
			http.Error(w, fmt.Sprintf("Slack API error: %s", tokenData.Error), http.StatusInternalServerError)
			return
		}
		
		// Success! Send the token
		tokenChan <- tokenData.AccessToken
		
		// Return success page
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Authentication Successful</title>
				<style>
					body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
					.success { color: green; }
				</style>
			</head>
			<body>
				<h1 class="success">Authentication Successful!</h1>
				<p>You can now close this window and return to the terminal.</p>
			</body>
			</html>
		`))
	})
	
	// Check if the callback URL uses HTTPS
	isHttps := strings.HasPrefix(callbackURL, "https://")
	
	// Start the server
	go func() {
		var err error
		
		if isHttps && strings.Contains(callbackURL, "localhost") {
			// First check if we have certificates
			_, certErr := os.Stat("certs/cert.pem")
			_, keyErr := os.Stat("certs/key.pem")
			
			// Generate certificates if they don't exist
			if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
				printInfo("Generating self-signed certificates for local HTTPS...")
				
				// Create certs directory if it doesn't exist
				if err := os.MkdirAll("certs", 0755); err != nil {
					errorChan <- fmt.Errorf("failed to create certs directory: %v", err)
					return
				}
				
				// Generate certificates using OpenSSL
				cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", 
					"-keyout", "certs/key.pem", "-out", "certs/cert.pem", 
					"-days", "365", "-nodes", "-subj", "/CN=localhost")
				
				if err := cmd.Run(); err != nil {
					errorChan <- fmt.Errorf("failed to generate certificates: %v\nPlease install OpenSSL or create certificates manually", err)
					return
				}
				
				printSuccess("Certificates generated successfully")
			}
			
			printInfo("Starting HTTPS server on port 1337...")
			err = http.ListenAndServeTLS(":1337", "certs/cert.pem", "certs/key.pem", nil)
		} else {
			printInfo("Starting HTTP server on port 1337...")
			err = http.ListenAndServe(":1337", nil)
		}
		
		if err != nil {
			errorChan <- fmt.Errorf("server error: %v", err)
		}
	}()
}

// openBrowser opens the default browser to the specified URL
func openBrowser(url string) error {
	var err error
	
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	
	return err
}

// readTokenConfig reads the token configuration from disk
func readTokenConfig() (TokenConfig, error) {
	var config TokenConfig
	
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}
	
	// Create full path
	configPath := filepath.Join(homeDir, configDir, configFile)
	
	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	
	// Parse JSON
	err = json.Unmarshal(data, &config)
	return config, err
}

// saveTokenConfig saves the token configuration to disk
func saveTokenConfig(config TokenConfig) error {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	// Create config directory if it doesn't exist
	configPath := filepath.Join(homeDir, configDir)
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return err
	}
	
	// Marshal config to JSON
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	
	// Write to file
	return os.WriteFile(filepath.Join(configPath, configFile), data, 0600)
}

// getEnvOrDefault returns environment variable or default value
func getEnvOrDefault(envName, defaultValue string) string {
	if value, exists := os.LookupEnv(envName); exists && value != "" {
		return value
	}
	return defaultValue
}

// getInput prompts the user for input
func getInput(prompt string) string {
	printInfo(prompt + ":")
	printPrompt(">")
	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil {
		printError(fmt.Sprintf("Reading input: %v", err))
		os.Exit(1)
	}
	
	return strings.TrimSpace(value)
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
	printQuestion(question)
	printInfo("(Enter each bullet point on a new line. Press Enter twice when done.)")
	printPrompt(">")
	
	reader := bufio.NewReader(os.Stdin)
	var lines []string
	
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			printError(fmt.Sprintf("Reading input: %v", err))
			os.Exit(1)
		}
		
		line = strings.TrimSpace(line)
		
		if line == "" {
			break
		}
		
		lines = append(lines, line)
		printPrompt(">")
	}
	
	return strings.Join(lines, "\n")
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