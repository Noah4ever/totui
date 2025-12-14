package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type todoCursor struct {
	Item int
	Sub  int // -1 means top-level item
}

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
			m = m.submitInput()
			return m, m.saveCmd(), true
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
	case "o":
		return m.handleAddKey(true)
	case "O":
		return m.handleAddKey(false)
	case "a":
		return m.handleAddKey(true)
	case "s":
		return m.handleAddSubKey()
	case "g":
		return m.recordKey(key), nil
	case "G":
		return m.handleBottomKey()
	case "r":
		return m.handleRenameKey()
	case "J":
		return m.handleMoveDown()
	case "K":
		return m.handleMoveUp()
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
			idx := clampIndex(m.SelectedList, len(m.Lists))
			if m.LastKey == "addAfterList" {
				idx = clampIndex(m.SelectedList+1, len(m.Lists))
			} else if m.LastKey == "addBeforeList" {
				idx = clampIndex(m.SelectedList, len(m.Lists))
			} else if len(m.Lists) > 0 {
				idx = len(m.Lists)
			}
			m.Lists = insertList(m.Lists, idx, newList)
			m.SelectedList = idx
			m.SelectedItem = 0
			m.SelectedSub = -1
		case InputAddTodo:
			if len(m.Lists) > 0 {
				newItem := TodoItem{ID: newID("item"), Title: value, Completed: false}
				items := &m.Lists[m.SelectedList].Items
				if len(*items) == 0 {
					*items = append(*items, newItem)
					m.SelectedItem = 0
				} else {
					pos := clampIndex(m.SelectedItem, len(*items))
					if m.LastKey == "addAfter" {
						pos++
					} else if m.LastKey == "addBefore" {
						// pos stays
					} else {
						pos = len(*items)
					}
					*items = insertItem(*items, pos, newItem)
					m.SelectedItem = pos
				}
				m.SelectedSub = -1
			}
		case InputRenameList:
			if len(m.Lists) > 0 {
				m.Lists[m.SelectedList].Name = value
			}
		case InputRenameTodo:
			if len(m.Lists) > 0 && len(m.Lists[m.SelectedList].Items) > 0 {
				m.Lists[m.SelectedList].Items[m.SelectedItem].Title = value
			}
		case InputAddSubTodo:
			if len(m.Lists) > 0 && len(m.Lists[m.SelectedList].Items) > 0 {
				parent := &m.Lists[m.SelectedList].Items[m.SelectedItem]
				newItem := TodoItem{ID: newID("sub"), Title: value}
				parent.SubItems = append(parent.SubItems, newItem)
				m.SelectedSub = len(parent.SubItems) - 1
			}
		case InputRenameSubTodo:
			if len(m.Lists) > 0 && len(m.Lists[m.SelectedList].Items) > 0 && m.SelectedSub >= 0 {
				m.Lists[m.SelectedList].Items[m.SelectedItem].SubItems[m.SelectedSub].Title = value
			}
		}
	}

	m.InputMode = InputNone
	m.Input.Reset()
	m.Input.Blur()
	m.LastKey = ""
	return m
}

func (m Model) cancelInput() Model {
	m.InputMode = InputNone
	m.Input.Reset()
	m.Input.Blur()
	m.LastKey = ""
	return m
}

func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		m.ActiveWindow = Todos
	case Todos:
		if len(m.Lists) != 0 && len(m.Lists[m.SelectedList].Items) != 0 {
			if m.SelectedSub >= 0 {
				completed := &m.Lists[m.SelectedList].Items[m.SelectedItem].SubItems[m.SelectedSub].Completed
				*completed = !*completed
			} else {
				completed := &m.Lists[m.SelectedList].Items[m.SelectedItem].Completed
				*completed = !*completed
			}
			return m, m.saveCmd()
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
			m.SelectedSub = -1
		}
	case Todos:
		items := &m.Lists[m.SelectedList].Items
		if len(*items) > 0 {
			if m.SelectedSub >= 0 {
				subs := &(*items)[m.SelectedItem].SubItems
				if len(*subs) > 0 {
					*subs = append((*subs)[:m.SelectedSub], (*subs)[m.SelectedSub+1:]...)
					if m.SelectedSub >= len(*subs) {
						m.SelectedSub = len(*subs) - 1
					}
				}
			} else {
				*items = append((*items)[:m.SelectedItem], (*items)[m.SelectedItem+1:]...)
				if m.SelectedItem >= len(*items) && m.SelectedItem > 0 {
					m.SelectedItem--
				}
			}
			if m.SelectedItem >= len(*items) {
				m.SelectedItem = len(*items) - 1
			}
			m.SelectedSub = -1
			if len(*items) == 0 {
				m.SelectedItem = 0
				m.SelectedSub = -1
			}
		}
	}
	return m, m.saveCmd()
}

func (m Model) handleAddKey(after bool) (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		cmd := beginInput(&m, InputAddList, "List name", "")
		if after {
			m.LastKey = "addAfterList"
		} else {
			m.LastKey = "addBeforeList"
		}
		return m, cmd
	case Todos:
		if len(m.Lists) == 0 {
			return m, nil
		}
		cmd := beginInput(&m, InputAddTodo, "Todo title", "")
		// encode before/after intent in LastKey sentinel
		if after {
			m.LastKey = "addAfter"
		} else {
			m.LastKey = "addBefore"
		}
		return m, cmd
	default:
		return m, nil
	}
}

func (m Model) handleAddSubKey() (tea.Model, tea.Cmd) {
	if m.ActiveWindow != Todos {
		return m, nil
	}
	if len(m.Lists) == 0 || len(m.Lists[m.SelectedList].Items) == 0 {
		return m, nil
	}
	cmd := beginInput(&m, InputAddSubTodo, "Sub item", "")
	return m, cmd
}

func (m Model) handleTopKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		m.SelectedList = 0
		m.SelectedItem = 0
		m.SelectedSub = -1
	case Todos:
		order := m.todoOrder()
		if len(order) > 0 {
			m.SelectedItem = order[0].Item
			m.SelectedSub = order[0].Sub
		}
	}
	return m, nil
}

func (m Model) handleBottomKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		m.SelectedList = len(m.Lists) - 1
		m.SelectedItem = 0
		m.SelectedSub = -1
	case Todos:
		order := m.todoOrder()
		if len(order) > 0 {
			last := order[len(order)-1]
			m.SelectedItem = last.Item
			m.SelectedSub = last.Sub
		}
	}
	return m, nil
}

func (m Model) handleMoveDown() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		if m.SelectedList < len(m.Lists)-1 {
			m.Lists[m.SelectedList], m.Lists[m.SelectedList+1] = m.Lists[m.SelectedList+1], m.Lists[m.SelectedList]
			m.SelectedList++
			return m, m.saveCmd()
		}
	case Todos:
		if len(m.Lists) == 0 || len(m.Lists[m.SelectedList].Items) == 0 {
			return m, nil
		}
		if m.SelectedSub >= 0 {
			subs := m.Lists[m.SelectedList].Items[m.SelectedItem].SubItems
			if m.SelectedSub < len(subs)-1 {
				subs[m.SelectedSub], subs[m.SelectedSub+1] = subs[m.SelectedSub+1], subs[m.SelectedSub]
				m.Lists[m.SelectedList].Items[m.SelectedItem].SubItems = subs
				m.SelectedSub++
				return m, m.saveCmd()
			}
		} else {
			items := m.Lists[m.SelectedList].Items
			if m.SelectedItem < len(items)-1 {
				items[m.SelectedItem], items[m.SelectedItem+1] = items[m.SelectedItem+1], items[m.SelectedItem]
				m.Lists[m.SelectedList].Items = items
				m.SelectedItem++
				return m, m.saveCmd()
			}
		}
	}
	return m, nil
}

func (m Model) handleMoveUp() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		if m.SelectedList > 0 {
			m.Lists[m.SelectedList], m.Lists[m.SelectedList-1] = m.Lists[m.SelectedList-1], m.Lists[m.SelectedList]
			m.SelectedList--
			return m, m.saveCmd()
		}
	case Todos:
		if len(m.Lists) == 0 || len(m.Lists[m.SelectedList].Items) == 0 {
			return m, nil
		}
		if m.SelectedSub >= 0 {
			subs := m.Lists[m.SelectedList].Items[m.SelectedItem].SubItems
			if m.SelectedSub > 0 {
				subs[m.SelectedSub], subs[m.SelectedSub-1] = subs[m.SelectedSub-1], subs[m.SelectedSub]
				m.Lists[m.SelectedList].Items[m.SelectedItem].SubItems = subs
				m.SelectedSub--
				return m, m.saveCmd()
			}
		} else {
			items := m.Lists[m.SelectedList].Items
			if m.SelectedItem > 0 {
				items[m.SelectedItem], items[m.SelectedItem-1] = items[m.SelectedItem-1], items[m.SelectedItem]
				m.Lists[m.SelectedList].Items = items
				m.SelectedItem--
				return m, m.saveCmd()
			}
		}
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
		if m.SelectedSub >= 0 {
			current := m.Lists[m.SelectedList].Items[m.SelectedItem].SubItems[m.SelectedSub].Title
			cmd := beginInput(&m, InputRenameSubTodo, "Rename sub", current)
			return m, cmd
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

func (m Model) todoOrder() []todoCursor {
	order := []todoCursor{}
	if m.SelectedList < 0 || m.SelectedList >= len(m.Lists) {
		return order
	}
	for i, item := range m.Lists[m.SelectedList].Items {
		order = append(order, todoCursor{Item: i, Sub: -1})
		for s := range item.SubItems {
			order = append(order, todoCursor{Item: i, Sub: s})
		}
	}
	return order
}

func (m Model) todoCursorIndex(order []todoCursor) int {
	for i, c := range order {
		if c.Item == m.SelectedItem && c.Sub == m.SelectedSub {
			return i
		}
	}
	return -1
}

func (m Model) handleDownKey() (tea.Model, tea.Cmd) {
	switch m.ActiveWindow {
	case Lists:
		if m.SelectedList < len(m.Lists)-1 {
			m.SelectedList++
			m.SelectedItem = 0
			m.SelectedSub = -1
		}
	case Todos:
		order := m.todoOrder()
		idx := m.todoCursorIndex(order)
		if idx >= 0 && idx < len(order)-1 {
			next := order[idx+1]
			m.SelectedItem = next.Item
			m.SelectedSub = next.Sub
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
			m.SelectedSub = -1
		}
	case Todos:
		order := m.todoOrder()
		idx := m.todoCursorIndex(order)
		if idx > 0 {
			prev := order[idx-1]
			m.SelectedItem = prev.Item
			m.SelectedSub = prev.Sub
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

func insertList(lists []TodoList, idx int, v TodoList) []TodoList {
	if idx < 0 {
		idx = 0
	}
	if idx > len(lists) {
		idx = len(lists)
	}
	lists = append(lists, TodoList{})
	copy(lists[idx+1:], lists[idx:])
	lists[idx] = v
	return lists
}

func insertItem(items []TodoItem, idx int, v TodoItem) []TodoItem {
	if idx < 0 {
		idx = 0
	}
	if idx > len(items) {
		idx = len(items)
	}
	items = append(items, TodoItem{})
	copy(items[idx+1:], items[idx:])
	items[idx] = v
	return items
}

func clampIndex(idx, length int) int {
	if idx < 0 {
		return 0
	}
	if idx > length {
		return length
	}
	return idx
}

func newID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func (m Model) saveCmd() tea.Cmd {
	if m.SaveFunc == nil || m.StoragePath == "" {
		return nil
	}
	path := m.StoragePath
	lists := m.Lists
	return func() tea.Msg {
		_ = m.SaveFunc(path, lists)
		return nil
	}
}
