package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	selectTip        string
	choices          []string         // items on the to-do list
	cursor           int              // which to-do list item our cursor is pointing at
	selected         map[int]struct{} // which to-do items are selected
	metaData         map[int]interface{}
	selectedMetaData interface{}
}

func NewListModel(title string, metaData map[string]interface{}) *model {
	m := &model{
		selectTip: title,
		selected:  make(map[int]struct{}),
		metaData:  make(map[int]interface{}),
	}
	var i int
	for choice, data := range metaData {
		m.choices = append(m.choices, choice)
		m.metaData[i] = data
		i++
	}
	return m
}

func (m *model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "Q":
			m.selectedMetaData = nil
			return m, tea.Quit
		case "enter":
			if m.selectedMetaData != nil {
				return m, tea.Quit
			}

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case " ":
			_, ok := m.selected[m.cursor]
			if !ok {
				m.selected = make(map[int]struct{})
				m.selected[m.cursor] = struct{}{}
				m.selectedMetaData = m.metaData[m.cursor]
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := m.selectTip + "\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\n↑↓ 移动光标, 空格切换选择, 回车保存, Q 退出.\n"

	// Send the UI for rendering
	return s
}

func (m *model) Select() interface{} {
	defer tea.Quit()
	p := tea.NewProgram(m)
	if teaModel, err := p.StartReturningModel(); err != nil {
		return nil
	} else {
		md := teaModel.(*model)
		return md.selectedMetaData
	}
}
