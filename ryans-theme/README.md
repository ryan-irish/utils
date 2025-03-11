# Ryan's Theme for Cursor IDE

A custom dark theme with purple accents for Cursor IDE.

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/ryans-theme.git
   cd ryans-theme
   ```

2. Run the installation script:
   ```bash
   ./install-ryans-theme.sh
   ```

3. Follow the prompts in the script

4. Restart Cursor IDE after installation

5. Go to Preferences > Color Theme (⌘K ⌘T)

6. Select "Ryan's Theme" from the dropdown

## What's Included

- `install-ryans-theme.sh`: Installation script that handles all the setup
- `themes/ryans-theme.json`: The actual theme color configuration
- Support files that are generated during installation

## How It Works

The installation script:

1. Creates necessary configuration files (`package.json`, `.vsixmanifest`, etc.)
2. Registers the theme with Cursor IDE 
3. Places the theme in the correct location for Cursor to find it

## Troubleshooting

If the theme doesn't appear after installation:

- Make sure Cursor IDE is completely closed and restarted
- Check if there are any error messages in Cursor's Developer Tools (Help > Toggle Developer Tools)

## Customization

To customize the theme:

1. Modify `themes/ryans-theme.json`
2. Run the installation script again
3. Restart Cursor IDE

## Sharing with Others

Feel free to share this theme with your teammates. They just need to follow the installation instructions above. 