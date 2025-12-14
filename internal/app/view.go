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

	return ui.Header.Render("ToTUI") + "\n\n" + main
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

	// Inline input when adding or renaming lists.
	if m.InputMode == InputAddList || m.InputMode == InputRenameList {
		rows = append(rows, ui.ItemSelected.Render(m.Input.View()))
	}

	body := ui.Title.Render("Lists") + "\n" + strings.Join(rows, "\n")
	w := m.sectionWidthLists()
	h := m.sectionHeight()
	panel := ui.Panel(w, h, m.ActiveWindow == Lists)
	return panel.Render(body)
}

func (m Model) renderTodos() string {
	items := []string{}
	if len(m.Lists) > 0 {
		listIdx := m.SelectedList
		if listIdx < 0 {
			listIdx = 0
		}
		if listIdx >= len(m.Lists) {
			listIdx = len(m.Lists) - 1
		}
		for i, item := range m.Lists[listIdx].Items {
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

	// Inline input when adding or renaming todos.
	if m.InputMode == InputAddTodo || m.InputMode == InputRenameTodo {
		items = append(items, fmt.Sprintf("[ ] %s", m.Input.View()))
	}

	body := ui.Title.Render("ToDo's") + "\n" + strings.Join(items, "\n")
	w := m.sectionWidthTodos()
	h := m.sectionHeight()
	panel := ui.Panel(w, h, m.ActiveWindow == Todos)
	return panel.Render(body)
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

func (m Model) sectionHeight() int {
	if m.Height <= 0 {
		return 0
	}
	// Subtract header (1) + blank line (1) for spacing.
	height := m.Height - 2
	if height < 6 {
		height = 6
	}
	return height
}
