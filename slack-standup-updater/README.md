# Slack Standup Updater

A simple command-line tool to post daily standup updates to a Slack thread.

## Installation

1. Ensure you have Go installed (1.16+)
2. Clone this repository:
```
git clone https://github.com/ryanirish/slack-standup-updater.git
cd slack-standup-updater
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

## Usage

### Configuration

The tool requires the following information:

- **SLACK_TOKEN**: Your Slack API token (starts with `xoxb-`)
- **SLACK_CHANNEL_ID**: The ID of the channel containing the thread
- **SLACK_THREAD_TS**: The timestamp of the thread to reply to

You can provide these in three ways:

1. As environment variables
2. In your shell profile
3. Enter them when prompted by the tool

#### Using a Slack Message Link

The tool now supports directly pasting a Slack message link to identify which thread to reply to!

When prompted, just paste a link in this format:
```
******.slack.com/archives/C048ECCB75H/p1743724813501239
```

The tool will automatically extract the channel ID and thread timestamp.

### Adding to your shell profile

Add this to your `.bashrc`, `.zshrc`, or equivalent:

```bash
# Slack Standup Configuration
export SLACK_TOKEN="xoxb-your-token-here"
export SLACK_CHANNEL_ID="C12345678"
export SLACK_THREAD_TS="1234567890.123456"

# Optional: Create an alias
alias standup="/path/to/slack-standup-updater/standup"
```

### Finding Slack IDs

1. **Channel ID**: Right-click on the channel and copy the link - the ID is the last part of the URL
2. **Thread TS**: Open the thread in your browser, the timestamp is in the URL after `thread_ts=`
3. **Easiest Method**: Right-click on any message in the thread, select "Copy link", and use that link directly in the tool when prompted

## How to use

Simply run the command:

```
standup
```

The tool will:
1. Ask if you have a Slack message link, or proceed to manual entry if not
2. Ask you each standup question
3. Allow you to provide multiple bullet points per question (enter each on a new line)
4. Press Enter twice to move to the next question
5. Format your answers and post them to the Slack thread

### Example Input

```
Do you have a Slack message link? (y/n)
y
Enter the Slack message link (e.g., ******.slack.com/archives/C048ECCB75H/p1743724813501239):
******.slack.com/archives/C048ECCB75H/p1743724813501239
Successfully parsed link. Channel ID: C048ECCB75H, Thread TS: 1743724813.501239

Enter multiple bullet points per question. Press Enter twice when done with a question.
1. What did you do yesterday?
(Enter each bullet point on a new line. Press Enter twice when done.)
> TESI-2102 reference volume investigation
> Completed PR review for team members
> Team meeting and planning session
> 

2. What will you do today?
(Enter each bullet point on a new line. Press Enter twice when done.)
> Continue TESI-2257 implementation
> Start documentation for new feature
> 

3. Anything blocking your progress?
(Enter each bullet point on a new line. Press Enter twice when done.)
> Waiting for API credentials from DevOps
> 

Standup posted successfully! 