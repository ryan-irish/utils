#!/bin/bash

# install-ryans-theme.sh
# A script to install Ryan's custom theme for Cursor IDE

set -e  # Exit on any error

# Display banner
echo "============================================="
echo "      Installing Ryan's Theme for Cursor     "
echo "============================================="

# Define variables
THEME_NAME="ryans-theme"
PUBLISHER="ryan"
VERSION="0.0.1"
CURSOR_EXTENSIONS_DIR="$HOME/.cursor/extensions"
EXTENSION_DIR="$CURSOR_EXTENSIONS_DIR/$PUBLISHER.$THEME_NAME-$VERSION"

# Create a temporary directory that will be cleaned up on exit
THEME_SOURCE_DIR=$(mktemp -d)
trap 'rm -rf "$THEME_SOURCE_DIR"' EXIT  # Clean up on script exit

mkdir -p "$THEME_SOURCE_DIR/themes"

# Check if the extensions directory exists
if [ ! -d "$CURSOR_EXTENSIONS_DIR" ]; then
    echo "Error: Cursor extensions directory not found at $CURSOR_EXTENSIONS_DIR"
    echo "Make sure Cursor IDE is installed and has been run at least once."
    exit 1
fi

# Create package.json
echo "Creating package.json..."
cat > "$THEME_SOURCE_DIR/package.json" << EOF
{
  "name": "$THEME_NAME",
  "displayName": "Ryan's Theme",
  "description": "Custom dark theme with purple accents",
  "version": "$VERSION",
  "publisher": "$PUBLISHER",
  "engines": {
    "vscode": "^1.60.0"
  },
  "categories": [
    "Themes"
  ],
  "contributes": {
    "themes": [
      {
        "label": "Ryan's Theme",
        "uiTheme": "vs-dark",
        "path": "./themes/ryans-theme.json"
      }
    ]
  }
}
EOF

# Handle theme file creation
echo "Creating theme file..."
cat > "$THEME_SOURCE_DIR/themes/ryans-theme.json" << EOF
{
    "name": "Ryan's Theme",
    "type": "dark",
    "colors": {
        "editor.background": "#1E1E1E",
        "editor.foreground": "#D4D4D4",
        "activityBarBadge.background": "#8B5CF6",
        "sideBarTitle.foreground": "#BBBBBB",
        "statusBar.background": "#8B5CF6",
        "statusBar.foreground": "#FFFFFF",
        "titleBar.activeBackground": "#181818",
        "tab.activeBackground": "#252525",
        "editorGroupHeader.tabsBackground": "#1E1E1E"
    },
    "tokenColors": [
        {
            "name": "Comments",
            "scope": ["comment", "punctuation.definition.comment"],
            "settings": {
                "fontStyle": "italic",
                "foreground": "#6A9955"
            }
        },
        {
            "name": "String",
            "scope": ["string", "string.quoted"],
            "settings": {
                "foreground": "#CE9178"
            }
        },
        {
            "name": "Keywords",
            "scope": ["keyword", "storage.type"],
            "settings": {
                "foreground": "#C586C0"
            }
        },
        {
            "name": "Functions",
            "scope": ["entity.name.function"],
            "settings": {
                "foreground": "#DCDCAA"
            }
        },
        {
            "name": "Variables and Parameters",
            "scope": ["variable", "parameter"],
            "settings": {
                "foreground": "#9CDCFE"
            }
        }
    ]
}
EOF

# Create README.md
echo "Creating README.md..."
cat > "$THEME_SOURCE_DIR/README.md" << EOF
# Ryan's Theme

A custom dark theme with purple accents.

## Features

- Dark theme optimized for readability
- Purple accent colors
- Custom syntax highlighting

## Installation

1. Install this theme through Cursor's extension manager
2. Go to \`Preferences > Color Theme\`
3. Select "Ryan's Theme"

## Feedback

If you have suggestions or issues, please let me know!
EOF

# Create .vsixmanifest
echo "Creating .vsixmanifest..."
cat > "$THEME_SOURCE_DIR/.vsixmanifest" << EOF
<?xml version="1.0" encoding="utf-8"?>
<PackageManifest Version="2.0.0" xmlns="http://schemas.microsoft.com/developer/vsx-schema/2011">
  <Metadata>
    <Identity Language="en-US" Id="$THEME_NAME" Version="$VERSION" Publisher="$PUBLISHER"/>
    <DisplayName>Ryan's Theme</DisplayName>
    <Description xml:space="preserve">Custom dark theme with purple accents</Description>
    <Tags>theme,color-theme</Tags>
    <Categories>Themes</Categories>
    <GalleryFlags>Public</GalleryFlags>
    <Badges></Badges>
    <Properties>
      <Property Id="Microsoft.VisualStudio.Code.Engine" Value="^1.60.0" />
    </Properties>
  </Metadata>
  <Installation>
    <InstallationTarget Id="Microsoft.VisualStudio.Code"/>
  </Installation>
  <Dependencies/>
  <Assets>
    <Asset Type="Microsoft.VisualStudio.Code.Manifest" Path="extension/package.json" Addressable="true" />
    <Asset Type="Microsoft.VisualStudio.Services.Content.Details" Path="extension/README.md" Addressable="true" />
  </Assets>
</PackageManifest>
EOF

# Create the extension directory
echo "Creating extension directory..."
mkdir -p "$EXTENSION_DIR"

# Copy the theme files to the extension directory
echo "Copying theme files..."
cp "$THEME_SOURCE_DIR/package.json" "$EXTENSION_DIR/"
cp "$THEME_SOURCE_DIR/README.md" "$EXTENSION_DIR/"
cp "$THEME_SOURCE_DIR/.vsixmanifest" "$EXTENSION_DIR/"
mkdir -p "$EXTENSION_DIR/themes"
cp "$THEME_SOURCE_DIR/themes/ryans-theme.json" "$EXTENSION_DIR/themes/"

# Update extensions.json
echo "Updating extensions.json..."
EXTENSIONS_JSON="$CURSOR_EXTENSIONS_DIR/extensions.json"

# Backup the original file
cp "$EXTENSIONS_JSON" "$EXTENSIONS_JSON.backup"

# Check if python3 is available
if command -v python3 &>/dev/null; then
    # Update using Python (more reliable for JSON)
    python3 -c "
import json, uuid, time
from pathlib import Path

try:
    extensions_file = Path('$EXTENSIONS_JSON')
    data = json.loads(extensions_file.read_text())

    # Check if the theme entry already exists
    theme_exists = False
    for ext in data:
        if ext.get('identifier', {}).get('id', '') == '$PUBLISHER.$THEME_NAME':
            theme_exists = True
            break

    if not theme_exists:
        # Add the new theme entry
        new_entry = {
            'identifier': {'id': '$PUBLISHER.$THEME_NAME', 'uuid': str(uuid.uuid4())},
            'version': '$VERSION',
            'location': {'$\$mid': 1, 'path': '$EXTENSION_DIR', 'scheme': 'file'},
            'relativeLocation': '$PUBLISHER.$THEME_NAME-$VERSION',
            'metadata': {
                'installedTimestamp': int(time.time() * 1000),
                'pinned': False,
                'source': 'local',
                'publisherDisplayName': '$PUBLISHER',
                'targetPlatform': 'undefined',
                'updated': False,
                'isPreReleaseVersion': False,
                'hasPreReleaseVersion': False
            }
        }
        data.append(new_entry)

        # Write back the updated file
        extensions_file.write_text(json.dumps(data))
        print('Added theme to extensions.json')
    else:
        print('Theme entry already exists in extensions.json')
except Exception as e:
    print(f'Error updating extensions.json: {e}')
    print('You may need to manually add the theme to your extensions list.')
"
else
    echo "Warning: Python 3 not found. Unable to automatically update extensions.json."
    echo "You may need to manually add the theme to your extensions list."
fi

echo "============================================="
echo "Installation complete!"
echo "============================================="
echo ""
echo "To use Ryan's Theme:"
echo "1. Restart Cursor IDE"
echo "2. Go to Preferences > Color Theme (⌘K ⌘T or Ctrl+K Ctrl+T)"
echo "3. Select 'Ryan's Theme' from the dropdown"
echo ""
echo "If you don't see the theme, check the following:"
echo "- Make sure Cursor is completely closed and restarted"
echo "- Verify your theme file exists: cat $EXTENSION_DIR/themes/ryans-theme.json"
echo "- Check if there are any error messages in Cursor's Developer Tools"
echo ""
# No need for explicit cleanup - the trap command will handle it on exit
