package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// Model represents the application state
type Model struct {
	list         list.Model
	textInput    textinput.Model
	themes       []ThemeInfo
	filteredThemes []ThemeInfo
	choice       string
	quitting     bool
	converting   bool
	converted    bool
	errorMsg     string
	filterMode   bool
	filterText   string
}

// item represents a theme item in the list
type item struct {
	theme ThemeInfo
}

func (i item) FilterValue() string {
	return i.theme.DisplayName
}

func (i item) Title() string {
	return i.theme.DisplayName
}

func (i item) Description() string {
	typeStr := "Unknown"
	if i.theme.Type == "dark" {
		typeStr = "Dark theme"
	} else if i.theme.Type == "light" {
		typeStr = "Light theme"
	}
	return fmt.Sprintf("%s • %s", typeStr, i.theme.Path)
}

// initialModel sets up the initial application state
func initialModel() Model {
	// Check platform support
	if err := validatePlatformSupport(); err != nil {
		log.Fatal(err)
	}
	
	// Discover VS Code themes
	themes, err := DiscoverVSCodeThemes()
	if err != nil {
		log.Fatal("Failed to discover VS Code themes:", err)
	}

	if len(themes) == 0 {
		log.Fatal("No VS Code themes found. Please ensure you have VS Code installed with some theme extensions.")
	}

	// Create list items
	items := make([]list.Item, len(themes))
	for i, theme := range themes {
		items[i] = item{theme: theme}
	}

	// Set up the list
	l := list.New(items, itemDelegate{}, 80, 20)
	l.Title = "VS Code Themes"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false) // Disable built-in filtering, we'll handle it ourselves
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	// Set up text input for filtering
	ti := textinput.New()
	ti.Placeholder = "Type to filter themes..."
	ti.CharLimit = 50
	ti.Width = 50

	return Model{
		list:           l,
		textInput:      ti,
		themes:         themes,
		filteredThemes: themes,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}


func (m Model) View() string {
	if m.quitting {
		return quitTextStyle.Render("Thanks for using VS Code to Warp theme converter!\n")
	}

	if m.errorMsg != "" {
		return fmt.Sprintf("\n  ❌ Error: %s\n\n  Press 'q' to quit.\n", m.errorMsg)
	}

	if m.converted {
		return fmt.Sprintf("\n  ✅ Successfully converted '%s' to Warp theme!\n  \n  The theme has been saved to ~/.warp/themes/\n  You can now select it in Warp's settings.\n\n  Press 'q' to quit.\n", m.choice)
	}

	if m.converting {
		return fmt.Sprintf("\n  🔄 Converting '%s' to Warp theme...\n", m.choice)
	}

	// Build the main view
	var content strings.Builder
	content.WriteString("\n")
	
	if m.filterMode {
		// Show filter input when in filter mode
		content.WriteString("🔍 Filter: ")
		content.WriteString(m.textInput.View())
		content.WriteString("\n")
		content.WriteString(fmt.Sprintf("📝 Found %d matching themes • Press Enter to navigate, Esc to cancel\n\n", len(m.filteredThemes)))
	} else {
		// Show filter hint when not in filter mode
		if m.filterText != "" {
			content.WriteString(fmt.Sprintf("🔍 Filtered by: \"%s\" (%d results) • Press / to change filter\n\n", m.filterText, len(m.filteredThemes)))
		} else {
			content.WriteString("💡 Press / to filter • j/k or ↑/↓ to navigate • Enter to convert\n\n")
		}
	}
	
	content.WriteString(m.list.View())
	return content.String()
}

// filterThemes filters the theme list based on the filter text
func (m *Model) filterThemes() {
	if m.filterText == "" {
		m.filteredThemes = m.themes
	} else {
		m.filteredThemes = make([]ThemeInfo, 0)
		filterLower := strings.ToLower(m.filterText)
		for _, theme := range m.themes {
			if strings.Contains(strings.ToLower(theme.DisplayName), filterLower) {
				m.filteredThemes = append(m.filteredThemes, theme)
			}
		}
	}
	
	// Update list items
	items := make([]list.Item, len(m.filteredThemes))
	for i, theme := range m.filteredThemes {
		items[i] = item{theme: theme}
	}
	m.list.SetItems(items)
}

// convertTheme handles the conversion process
func (m Model) convertTheme(themeInfo ThemeInfo) tea.Cmd {
	return func() tea.Msg {
		// Load the VS Code theme
		vscodeTheme, err := LoadVSCodeTheme(themeInfo.Path)
		if err != nil {
			return errorMsg{fmt.Sprintf("Failed to load theme: %v", err)}
		}

		// Try to load extension metadata for attribution
		var extensionMetadata *ExtensionMetadata
		if metadata, err := LoadExtensionMetadata(themeInfo.Path); err == nil {
			extensionMetadata = metadata
		}

		// Convert to Warp theme
		warpTheme, err := ConvertVSCodeToWarp(vscodeTheme, extensionMetadata)
		if err != nil {
			return errorMsg{fmt.Sprintf("Failed to convert theme: %v", err)}
		}

		// Save the Warp theme
		if err := SaveWarpTheme(warpTheme, vscodeTheme.Name); err != nil {
			return errorMsg{fmt.Sprintf("Failed to save theme: %v", err)}
		}

		return convertedMsg{}
	}
}

// Message types for async operations
type errorMsg struct {
	err string
}

type convertedMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errorMsg:
		m.converting = false
		m.errorMsg = msg.err
		return m, nil

	case convertedMsg:
		m.converting = false
		m.converted = true
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		// Handle different states
		if m.converting || m.converted || m.errorMsg != "" {
			// In conversion or end states, only handle q and ctrl+c
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			case "q":
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}

		if m.filterMode {
			// In filter mode - handle text input
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			case "esc":
				// Exit filter mode without applying
				m.filterMode = false
				m.filterText = ""
				m.textInput.SetValue("")
				m.filterThemes()
				return m, nil
			case "enter":
				// Exit filter mode and apply filter
				m.filterMode = false
				return m, nil
			case "backspace":
				// Handle backspace manually
				if len(m.filterText) > 0 {
					m.filterText = m.filterText[:len(m.filterText)-1]
					m.textInput.SetValue(m.filterText)
					m.filterThemes()
				}
				return m, nil
			default:
				// Add character to filter
				if len(msg.String()) == 1 {
					m.filterText += msg.String()
					m.textInput.SetValue(m.filterText)
					m.filterThemes()
				}
				return m, nil
			}
		} else {
			// Normal navigation mode
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			case "q":
				m.quitting = true
				return m, tea.Quit
			case "/":
				// Enter filter mode
				m.filterMode = true
				m.textInput.Focus()
				return m, nil
			case "enter":
				// Convert selected theme
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = i.theme.DisplayName
					m.converting = true
					return m, m.convertTheme(i.theme)
				}
				return m, nil
			case "up", "k":
				// Vim-style up navigation
				m.list.CursorUp()
				return m, nil
			case "down", "j":
				// Vim-style down navigation
				m.list.CursorDown()
				return m, nil
			case "g":
				// Go to top (vim-style)
				m.list.Select(0)
				return m, nil
			case "G":
				// Go to bottom (vim-style)
				m.list.Select(len(m.list.Items()) - 1)
				return m, nil
			}
		}
	}

	// Handle list updates for other keys (like page up/down)
	if !m.filterMode {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

// itemDelegate defines how items are rendered in the list
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Title())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		os, arch := getPlatformInfo()
		fmt.Println("VS Code to Warp Theme Converter")
		fmt.Println("===============================")
		fmt.Printf("Platform: %s/%s\n", os, arch)
		fmt.Println()
		fmt.Println("This CLI tool discovers VS Code themes installed on your system")
		fmt.Println("and converts them to Warp terminal themes.")
		fmt.Println()
		fmt.Println("✅ Warp is now supported on all platforms: Windows, macOS, and Linux!")
		fmt.Println()
		fmt.Println("Usage: vscode-to-warp")
		fmt.Println()
		fmt.Println("Instructions:")
		fmt.Println("  1. Run the command to see all available VS Code themes")
		fmt.Println("  2. Use arrow keys or vim bindings to navigate")
		fmt.Println("  3. Press / to filter, type to search, Enter to apply filter")
		fmt.Println("  4. Navigate filtered results with j/k or arrow keys")
		fmt.Println("  5. Press Enter to convert the selected theme")
		fmt.Println("  6. The converted theme will be saved to ~/.warp/themes/")
		fmt.Println()
		fmt.Println("Controls:")
		fmt.Println("  Navigation:")
		fmt.Println("    ↑/↓, j/k    Navigate themes")
		fmt.Println("    g           Go to top")
		fmt.Println("    G           Go to bottom")
		fmt.Println("  Filtering:")
		fmt.Println("    /           Enter filter mode")
		fmt.Println("    Enter       Apply filter and navigate results")
		fmt.Println("    Esc         Cancel filter")
		fmt.Println("  Actions:")
		fmt.Println("    Enter       Convert selected theme")
		fmt.Println("    q           Quit")
		fmt.Println("    Ctrl+C      Force quit")
		return
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
