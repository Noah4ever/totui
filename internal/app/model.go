package app

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

type Mode int

const (
	Simple Mode = iota
	Advanced
)

type Window int

const (
	Lists Window = iota
	Todos
)

type InputMode int

const (
	InputNone InputMode = iota
	InputAddList
	InputAddTodo
	InputRenameList
	InputRenameTodo
	InputAddSubTodo
	InputRenameSubTodo
)

type TodoItem struct {
	ID        string
	Title     string
	Completed bool
	Due       string // ISO date string, empty if none
	SubItems  []TodoItem
}

type TodoList struct {
	ID    string
	Name  string
	Items []TodoItem
}

type Model struct {
	Lists        []TodoList
	SelectedList int
	SelectedItem int
	SelectedSub  int // -1 when no sub-item is selected
	Mode         Mode
	ActiveWindow Window
	Input        textinput.Model
	InputMode    InputMode
	LastKey      string
	LastKeyAt    time.Time
	Width        int
	Height       int
	StoragePath  string
	SaveFunc     func(string, []TodoList) error
	Running      bool
}

func NewModel() Model {
	return NewModelWithData(nil)
}

func NewModelWithData(lists []TodoList) Model {
	tInput := textinput.New()
	tInput.Prompt = " > "
	tInput.CharLimit = 140
	tInput.Placeholder = "Name"
	tInput.Blur()

	if len(lists) == 0 {
		lists = []TodoList{
			{
				ID:   "inbox",
				Name: "Inbox",
				Items: []TodoItem{
					{ID: "1", Title: "Welcome to toTUI", Completed: false},
				},
			},
		}
	}

	return Model{
		Lists:        lists,
		SelectedList: 0,
		SelectedItem: 0,
		SelectedSub:  -1,
		Mode:         Simple,
		ActiveWindow: Lists,
		Input:        tInput,
		InputMode:    InputNone,
		Width:        0,
		Height:       0,
		StoragePath:  "totui_data.json",
		Running:      true,
	}
}
