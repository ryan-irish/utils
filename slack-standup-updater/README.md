# Slack Standup Updater

A simple command-line tool to post daily standup updates to a Slack thread. Messages appear as sent by you, not a bot.

## Features

- Post standup updates that appear as coming from your Slack account
- Parse Slack message links to easily find the right thread
- Send messages to yourself for testing or drafting
- Post to Slackbot (the most reliable way to message yourself)
- Message any Slack user directly by their ID
- Colorized, user-friendly terminal interface
- Multiple bullet points per question
- Secure token storage between sessions

## Installation

1. Ensure you have Go installed (1.16+)
2. Clone this repository:
```
git clone https://github.com/ryan-irish/utils.git
cd utils/slack-standup-updater
```

3. Build the binary:
```
go build -o standup
```

4. Move the binary to your PATH (optional):
```
# On macOS/Linux
mv standup /usr/local/bin/
# Or add current directory to PATH
```

## Slack App Setup

Before using this tool, you need to create a Slack app:

1. Go to [api.slack.com/apps](https://api.slack.com/apps) and click "Create New App"
2. Choose "From scratch" and give your app a name (e.g., "Standup Updater")
3. Select your workspace

### Configure OAuth Settings

1. In the sidebar, click "OAuth & Permissions"
2. Under "Redirect URLs", add: `http://localhost:1337/callback`
3. Under "User Token Scopes", add:
   - `chat:write` (to post messages as yourself)
   - `channels:read` (optional, helps with channel resolution)
   - `im:write` (required for messaging yourself)
4. Note your "Client ID" and "Client Secret" at the top of the OAuth page

### Distribution & Installation

1. Go to "Manage Distribution" from the sidebar
2. Enable "Remove hard-coded information"
3. Click "Activate Public Distribution"

## Usage

### First-time setup

The first time you run the tool, it will:
1. Ask for your app's Client ID and Client Secret
2. Open your browser to authenticate with Slack
3. Store your token securely for future use

You can also set the Client ID and Secret in environment variables:
```
export SLACK_CLIENT_ID="your_client_id"
export SLACK_CLIENT_SECRET="your_client_secret"
```

### Posting Standups

The tool supports directly pasting a Slack message link to identify which thread to reply to:

When prompted, just paste a link in this format:
```
******.slack.com/archives/C048ECCB75H/p1743724813501239
```

The tool will automatically extract the channel ID and thread timestamp.

### Finding Slack Message Links

To find a thread to reply to:
1. Right-click on any message in the thread
2. Select "Copy link"
3. Paste when prompted by the tool

## How to use

Simply run the command:

```
standup
```

The tool will:
1. Authenticate with Slack as yourself (first time only)
2. Ask where you want to post your standup:
   - Reply to a thread in a channel
   - Message yourself directly (works with user tokens)
   - Post to Slackbot (works with all token types)
   - Message any user directly by their ID
3. If using a channel, ask if you have a Slack message link
4. Ask you each standup question
5. Allow you to provide multiple bullet points per question (enter each on a new line)
6. Press Enter twice to move to the next question
7. Format your answers and post them to the selected destination as yourself

### Messaging Options

#### Reply to a Channel Thread
Best for team standups where others should see your update in context.

#### Message Yourself Directly
Good for drafting your standup before posting to a team thread. This option works best with user tokens, not bot tokens.

#### Post to Slackbot
The most reliable way to message yourself. Works with both user and bot tokens. Slackbot acts as a personal messaging channel that only you can see.

#### Message Any User
Send your standup report directly to a colleague or manager. You'll need their Slack user ID, which starts with "U".

### Finding Slack IDs

#### User IDs
To find a Slack user ID:
1. Open Slack in a browser
2. Click on the user's profile
3. In their profile, click the "..." menu and select "Copy member ID"

#### Slackbot Channel ID
If you need to find your Slackbot channel ID:
1. Open Slack in a browser
2. Click on Slackbot in the sidebar
3. The URL will contain the ID, like: /messages/DXXXXXXXX

#### TODO
- add ability to add token "profiles" to store multiple creds
- prune slackbot option
- test in thread
- fix name and icon
