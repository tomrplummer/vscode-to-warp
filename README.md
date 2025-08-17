# VS Code to Warp Theme Converter

> ✨ **Created entirely by Warp 2.0 AI** - A CLI tool built with Go and Bubble Tea that discovers VS Code themes on your system and converts them to Warp terminal themes.

This project was conceived, designed, and implemented completely by Warp 2.0 AI in a single session, showcasing the capabilities of AI-assisted development for creating fully functional, cross-platform tools.

## Features

- 🔍 **Auto-discovery**: Automatically finds all VS Code themes in your extensions directory
- 🎨 **Smart conversion**: Maps VS Code color schemes to Warp terminal colors
- ⚡ **Interactive UI**: Beautiful terminal interface powered by Bubble Tea
- 🔎 **Filtering**: Type to filter themes by name
- 📁 **Auto-install**: Saves converted themes directly to `~/.warp/themes/`

## Installation

### From Source

```bash
git clone https://github.com/tomrplummer/vscode-to-warp.git
cd vscode-to-warp

# Build only
make build
# or
go build -o vscode-to-warp

# Build and install to ~/bin
make install-user

# Build and install to /usr/local/bin (requires sudo)
make install
```

### Binary Releases

Download the latest binary for your platform from the [releases page](https://github.com/tomrplummer/vscode-to-warp/releases):

- **Windows**: `vscode-to-warp-windows-amd64.exe` or `vscode-to-warp-windows-arm64.exe`
- **macOS**: `vscode-to-warp-darwin-amd64` (Intel) or `vscode-to-warp-darwin-arm64` (Apple Silicon)
- **Linux**: `vscode-to-warp-linux-amd64`

## Usage

1. Run the CLI:
   ```bash
   ./vscode-to-warp
   ```

2. Navigate the theme list using arrow keys
3. Type to filter themes by name
4. Press Enter to convert the selected theme
5. The converted theme will be saved to `~/.warp/themes/`
6. Select the new theme in Warp's settings

### Controls

- `↑/↓` - Navigate themes
- `Enter` - Convert selected theme  
- `/` - Filter themes
- `q` - Quit (after conversion)
- `Ctrl+C` - Force quit

## How it Works

The tool:

1. **Discovers themes** by scanning `~/.vscode/extensions/*/themes/*.json`
2. **Parses VS Code theme JSON** to extract color information
3. **Maps colors** from VS Code format to Warp's YAML format:
   - Editor background/foreground → Terminal background/foreground
   - Terminal colors → ANSI color palette
   - Focus/accent colors → Accent color
4. **Generates YAML** in Warp's theme format
5. **Saves** to `~/.warp/themes/` with a clean filename

## Color Mapping

| VS Code | Warp |
|---------|------|
| `editor.background` | `background` |
| `editor.foreground` | `foreground` |
| `focusBorder`, `button.background`, etc. | `accent` |
| `terminal.ansi*` | `terminal_colors.normal.*` |
| `terminal.ansiBright*` | `terminal_colors.bright.*` |

## Platform Support

🎉 **Universal Support**: Warp is now available on all platforms!

| Platform | VS Code Discovery | Warp Installation | Status |
|----------|------------------|------------------|--------|
| **Windows** | ✅ | ✅ | Fully supported |
| **macOS** | ✅ | ✅ | Fully supported |
| **Linux** | ✅ | ✅ | Fully supported |

### Paths
- **VS Code themes**: `~/.vscode/extensions/*/themes/*.json` (all platforms)
- **Warp themes**: `~/.warp/themes/` (all platforms)

## Requirements

- VS Code installed with theme extensions
- Warp terminal (Windows/macOS/Linux)
- Go 1.21+ (for building from source)

## Examples

The tool converts themes like:
- Catppuccin variants
- Material themes
- One Dark Pro
- Dracula
- Gruvbox
- And any other VS Code theme with proper color definitions

## Troubleshooting

**No themes found**: Ensure VS Code is installed and you have theme extensions installed.

**Permission errors**: The tool creates `~/.warp/themes/` if it doesn't exist.

**Colors look off**: Some VS Code themes may not define all terminal colors, so defaults are used.

## About This Project

### 🤖 AI-Generated Development

This entire project was created by **Warp 2.0 AI** in a single collaborative session, demonstrating:

- **Full-stack development**: From concept to deployment
- **Cross-platform compatibility**: Windows, macOS, and Linux support
- **Production-ready code**: Complete with error handling, documentation, and CI/CD
- **Modern tooling**: Go, Bubble Tea, GitHub Actions, and automated releases
- **Real-world problem solving**: Bridging the gap between VS Code and Warp themes

### 🛠️ What Was Built

- **Architecture design**: Multi-file Go project structure
- **Cross-platform paths**: Handling different OS filesystem conventions
- **Interactive UI**: Terminal-based interface with filtering and navigation
- **Color conversion**: Intelligent mapping between theme formats
- **Error handling**: Graceful failures and helpful error messages
- **Documentation**: Comprehensive README with usage examples
- **CI/CD pipeline**: Automated building and releasing for multiple platforms
- **Open source setup**: MIT license and contribution-ready repository

The AI handled everything from initial requirements gathering to final deployment, including discovering that Warp had recently launched on Windows and updating the platform support accordingly.

### 💡 Development Highlights

- **Adaptive problem-solving**: Updated Windows support when informed of Warp's availability
- **Best practices**: Proper Go project structure, error handling, and documentation
- **User experience**: Intuitive CLI with helpful feedback and platform detection
- **Maintenance**: Automated releases and clear troubleshooting guides
