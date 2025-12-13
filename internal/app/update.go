package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.Width = ws.Width
		m.Height = ws.Height
		return m, nil
	}

	if updated, cmd, handled := m.updateInputMode(msg); handled {
		return updated, cmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	return m.handleKeyMsg(keyMsg)
}

func (m Model) updateInputMode(msg tea.Msg) (Model, tea.Cmd, bool) {
	if m.InputMode == InputNone {
		return m, nil, false
	}

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			return m.submitInput(), nil, true
		case "esc":
			return m.cancelInput(), nil, true
		}
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd, true
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if double, handled := m.checkDoubleKey(key); handled {
		return double, nil
	}

	if key != "d" && key != "g" {
		m.clearRecordedKey()
	}

	switch key {
	case "q", "ctrl+c":
		m.Running = false
		return m, tea.Quit
	case "enter":
		return m.handleEnterKey()
	case "d":
		return m.recordKey(key), nil
	case "a":
		return m.handleAddKey()
	case "g":
		return m.recordKey(key), nil
	case "G":
		return m.handleBottomKey()
	case "r":
		return m.handleRenameKey()
	case "j", "down":
		return m.handleDownKey()
	case "k", "up":
		return m.handleUpKey()
	case "tab", "l":
		return m.toggleWindow(), nil
	case "shift+tab", "h":
		return m.toggleWindow(), nil
	default:
		return m, nil
	}
}

func (m Model) submitInput() Model {
	value := strings.TrimSpace(m.Input.Value())
	if value != "" {
		switch m.InputMode {
		case InputAddList:
			newList := TodoList{ID: newID("list"), Name: value, Items: []TodoItem{}}
			m.Lists = append(m.Lists, newList)
			m.SelectedList = len(m.Lists) - 1
			m.SelectedItem = 0
		case InputAddTodo:
			if len(m.Lists) > 0 {
				newItem := TodoItem{ID: newID("item"), Title: value, Completed: false}
				m.Lists[m.SelectedList].Items = append(m.Lists[m.SelectedList].Items, newItem)
				m.SelectedItem = len(m.Lists[m.SelectedList].Items) - 1
			}
		case InputRenameList:
			if len(m.Lists) > 0 {
				m.Lists[m.SelectedList].Name = value
			}
		case InputRenameTodo:
			if len(m.Lists) > 0 && len(m.Lists[m.SelectedList].Items) > 0 {
				m.Lists[m.SelectedList].Items[m.SelectedItem].Title = value
			}
		}
	}

	m.InputMode = InputNone
	m.Input.Reset()
	m.Input.Blur()
	return m
}

func (m Model) cancelInput() Model {
	m.InputMode = InputNone
	m.Input.Reset()
	m.Input.Blur()
	return m
}

func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		m.ActiveWindow = Todos
	case Todos:
		if len(m.Lists) != 0 && len(m.Lists[m.SelectedList].Items) != 0 {
			completed := &m.Lists[m.SelectedList].Items[m.SelectedItem].Completed
			*completed = !*completed
		}
	}
	return m, nil
}

func (m Model) handleDeleteKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		if len(m.Lists) > 0 {
			m.Lists = append(m.Lists[:m.SelectedList], m.Lists[m.SelectedList+1:]...)
			if m.SelectedList >= len(m.Lists) && m.SelectedList > 0 {
				m.SelectedList--
			}
			m.SelectedItem = 0
		}
	case Todos:
		items := &m.Lists[m.SelectedList].Items
		if len(*items) > 0 {
			*items = append((*items)[:m.SelectedItem], (*items)[m.SelectedItem+1:]...)
			if m.SelectedItem >= len(*items) && m.SelectedItem > 0 {
				m.SelectedItem--
			}
		}
	}
	return m, nil
}

func (m Model) handleAddKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		cmd := beginInput(&m, InputAddList, "List name", "")
		return m, cmd
	case Todos:
		if len(m.Lists) == 0 {
			return m, nil
		}
		cmd := beginInput(&m, InputAddTodo, "Todo title", "")
		return m, cmd
	default:
		return m, nil
	}
}

func (m Model) handleTopKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		m.SelectedList = 0
		m.SelectedItem = 0
	case Todos:
		m.SelectedItem = 0
	}
	return m, nil
}

func (m Model) handleBottomKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		m.SelectedList = len(m.Lists) - 1
		m.SelectedItem = 0
	case Todos:
		m.SelectedItem = len(m.Lists[m.SelectedList].Items) - 1
	}
	return m, nil
}

func (m Model) handleRenameKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		if len(m.Lists) == 0 {
			return m, nil
		}
		current := m.Lists[m.SelectedList].Name
		cmd := beginInput(&m, InputRenameList, "Rename list", current)
		return m, cmd
	case Todos:
		if len(m.Lists) == 0 || len(m.Lists[m.SelectedList].Items) == 0 {
			return m, nil
		}
		current := m.Lists[m.SelectedList].Items[m.SelectedItem].Title
		cmd := beginInput(&m, InputRenameTodo, "Rename todo", current)
		return m, cmd
	default:
		return m, nil
	}
}

func (m Model) checkDoubleKey(key string) (Model, bool) {
	const threshold = 450 * time.Millisecond

	if m.LastKey == key && !m.LastKeyAt.IsZero() && time.Since(m.LastKeyAt) <= threshold {
		switch key {
		case "d":
			updated, _ := m.handleDeleteKey()
			if model, ok := updated.(Model); ok {
				m = model
			}
		case "g":
			updated, _ := m.handleTopKey()
			if model, ok := updated.(Model); ok {
				m = model
			}
		}
		m.LastKey = ""
		m.LastKeyAt = time.Time{}
		return m, true
	}

	// No double detected; leave handling to caller.
	return m, false
}

func (m Model) recordKey(key string) Model {
	m.LastKey = key
	m.LastKeyAt = time.Now()
	return m
}

func (m Model) clearRecordedKey() {
	m.LastKey = ""
	m.LastKeyAt = time.Time{}
}

func (m Model) handleDownKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		if m.SelectedList < len(m.Lists)-1 {
			m.SelectedList++
			m.SelectedItem = 0
		}
	case Todos:
		if m.SelectedItem < len(m.Lists[m.SelectedList].Items)-1 {
			m.SelectedItem++
		}
	}
	return m, nil
}

func (m Model) handleUpKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		if m.SelectedList > 0 {
			m.SelectedList--
			m.SelectedItem = 0
		}
	case Todos:
		if m.SelectedItem > 0 {
			m.SelectedItem--
		}
	}
	return m, nil
}

func (m Model) switchWindow(target Window) Model {
	m.ActiveWindow = target
	return m
}

func (m Model) toggleWindow() Model {
	if m.ActiveWindow == Lists {
		m.ActiveWindow = Todos
	} else {
		m.ActiveWindow = Lists
	}
	return m
}

func beginInput(m *Model, mode InputMode, placeholder, value string) tea.Cmd {
	m.InputMode = mode
	m.Input.Placeholder = placeholder
	m.Input.SetValue(value)
	m.Input.CursorEnd()
	return m.Input.Focus()
}

func newID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
