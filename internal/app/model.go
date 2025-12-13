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
)

type TodoItem struct {
	ID        string
	Title     string
	Completed bool
	Due       string // ISO date string, empty if none
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
	Mode         Mode
	ActiveWindow Window
	Input        textinput.Model
	InputMode    InputMode
	LastKey      string
	LastKeyAt    time.Time
	Width        int
	Height       int
	Running      bool
}

func NewModel() Model {
	tInput := textinput.New()
	tInput.Prompt = " > "
	tInput.CharLimit = 140
	tInput.Placeholder = "Name"
	tInput.Blur()

	return Model{
		Lists: []TodoList{
			{
				ID:   "inbox",
				Name: "Inbox",
				Items: []TodoItem{
					{ID: "1", Title: "Welcome to toTUI", Completed: false},
				},
			},
		},
		SelectedList: 0,
		SelectedItem: 0,
		Mode:         Simple,
		ActiveWindow: Lists,
		Input:        tInput,
		InputMode:    InputNone,
		Width:        0,
		Height:       0,
		Running:      true,
	}
}
