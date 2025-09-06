package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		choices:  []string{"List Items", "Submit Item", "Search Items", "Exit"},
		selected: make(map[int]struct{}),
	}
}

// --- Bubbletea interface methods ---

// Init runs when the program starts
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles keypresses and messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Quit on ctrl+c or q
		case "ctrl+c", "q":
			return m, tea.Quit

		// Navigate up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// Navigate down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// Select with Enter
		case "enter":
			choice := m.choices[m.cursor]
			if choice == "Exit" {
				return m, tea.Quit
			}
			// For now just print action
			fmt.Println("You chose:", choice)
		}
	}
	return m, nil
}

// View renders the UI
func (m model) View() string {
	s := "Main Menu:\n\n"

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\n(↑/↓ to navigate, enter to select, q to quit)\n"
	return s
}

// Run starts the Bubbletea program
func Run() error {
	p := tea.NewProgram(initialModel())
	return p.Start()
}
