// internal/cli/cli.go
package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Styles ---
var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	activeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// Screens/pages
type page int

const (
	pageMenu page = iota
	pageList
	pageSearch
	pageSubmit
)

// --- Model ---
type model struct {
	page      page
	cursor    int
	status    string
	options   []string
	client    *APIClient
	listItems []map[string]any
}

// --- Init ---
func initialModel(client *APIClient) model {
	return model{
		page:    pageMenu,
		cursor:  0,
		status:  "Select an option to continue...",
		options: []string{"List Items", "Submit Item", "Search Items", "Exit"},
		client:  client,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// --- Update ---
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		// Menu navigation
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}

		case "enter":
			switch m.page {
			case pageMenu:
				switch m.cursor {
				case 0: // List Items
					items, err := m.client.ListItems()
					if err != nil {
						m.status = "❌ Error fetching items: " + err.Error()
					} else {
						m.listItems = items
						m.page = pageList
					}
				case 1: // Submit Item
					m.page = pageSubmit
				case 2: // Search
					m.page = pageSearch
				case 3: // Exit
					return m, tea.Quit
				}
			case pageList, pageSubmit, pageSearch:
				// Return to menu
				m.page = pageMenu
			}
		}
	}
	return m, nil
}

// --- View ---
func (m model) View() string {
	switch m.page {

	case pageMenu:
		return m.viewMenu()

	case pageList:
		return m.viewList()

	case pageSubmit:
		return m.viewSubmit()

	case pageSearch:
		return m.viewSearch()

	default:
		return "Unknown page"
	}
}

// --- Menu view ---
func (m model) viewMenu() string {
	s := titleStyle.Render("📚 Repository CLI") + "\n\n"

	for i, opt := range m.options {
		cursor := "  "
		if m.cursor == i {
			cursor = cursorStyle.Render("➜")
			s += fmt.Sprintf("%s %s\n", cursor, activeStyle.Render(opt))
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, opt)
		}
	}

	s += "\n" + helpStyle.Render("↑/↓ to navigate • enter to select • q to quit") + "\n"
	s += "\nStatus: " + m.status + "\n"
	return s
}

// --- List Items view ---
func (m model) viewList() string {
	s := titleStyle.Render("📦 Items") + "\n\n"
	if len(m.listItems) == 0 {
		s += "No items found.\n"
	} else {
		for _, item := range m.listItems {
			s += fmt.Sprintf("• %s (ID: %v)\n", item["title"], item["id"])
		}
	}
	s += "\nPress Enter to return to Menu"
	return s
}

// --- Submit view ---
func (m model) viewSubmit() string {
	s := titleStyle.Render("📤 Submit Item") + "\n\n"
	s += "Feature not implemented yet (would upload file).\n"
	s += "\nPress Enter to return to Menu"
	return s
}

// --- Search view ---
func (m model) viewSearch() string {
	s := titleStyle.Render("🔍 Search Items") + "\n\n"
	s += "Feature not implemented yet (would prompt for query).\n"
	s += "\nPress Enter to return to Menu"
	return s
}

// --- Run ---
func Run(apiBase string) error {
	client := NewAPIClient(apiBase)
	p := tea.NewProgram(initialModel(client))
	return p.Start()
}
