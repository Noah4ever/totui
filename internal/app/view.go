package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/noah4ever/totui/internal/ui"
)

func (m Model) View() string {
	if len(m.Lists) == 0 {
		return "No todo lists available."
	}

	listsPanel := m.renderLists()
	todosPanel := m.renderTodos()

	main := lipgloss.JoinHorizontal(lipgloss.Top, listsPanel, todosPanel)

	if m.InputMode == InputNone {
		return ui.Header.Render("ToTUI") + "\n\n" + main
	}

	return ui.Header.Render("ToTUI") + "\n\n" + main + "\n\n" + m.renderInput()
}

func (m Model) renderLists() string {
	rows := make([]string, 0, len(m.Lists))
	for i, list := range m.Lists {
		style := ui.Item
		if i == m.SelectedList {
			style = ui.ItemSelected
		}
		rows = append(rows, style.Render(list.Name))
	}

	body := strings.Join(rows, "\n")
	w := m.sectionWidthLists()
	panel := ui.Panel(w, m.ActiveWindow == Lists)
	content := ui.Title.Render("Lists") + "\n" + body
	return panel.Render(content)
}

func (m Model) renderTodos() string {
	items := []string{}
	if len(m.Lists) > 0 {
		for i, item := range m.Lists[m.SelectedList].Items {
			status := "[ ]"
			style := ui.Item
			if item.Completed {
				status = "[x]"
				style = ui.Completed
			}
			if i == m.SelectedItem {
				style = ui.ItemSelected
			}
			items = append(items, fmt.Sprintf("%s %s", status, style.Render(item.Title)))
		}
	}

	body := strings.Join(items, "\n")
	w := m.sectionWidthTodos()
	panel := ui.Panel(w, m.ActiveWindow == Todos)
	content := ui.Title.Render("ToDo's") + "\n" + body
	return panel.Render(content)
}

func (m Model) renderInput() string {
	label := ""
	switch m.InputMode {
	case InputAddList:
		label = "Add list"
	case InputAddTodo:
		label = "Add todo"
	case InputRenameList:
		label = "Rename list"
	case InputRenameTodo:
		label = "Rename todo"
	}

	w := m.sectionWidthInput()
	return ui.Panel(w, true).Render(ui.Title.Render(label) + "\n" + m.Input.View())
}

func (m Model) sectionWidthLists() int {
	if m.Width <= 0 {
		return 30
	}
	padding := 4 // panel padding left+right from styles (2 each)
	available := m.Width - padding
	min := 24
	alloc := available / 3
	if alloc < min {
		alloc = min
	}
	return alloc
}

func (m Model) sectionWidthTodos() int {
	if m.Width <= 0 {
		return 60
	}
	padding := 4
	listsWidth := m.sectionWidthLists()
	remaining := m.Width - listsWidth - padding
	if remaining < 30 {
		remaining = 30
	}
	return remaining
}

func (m Model) sectionWidthInput() int {
	if m.Width <= 0 {
		return 60
	}
	return m.Width - 4
}
