# ToTUI

Minimal full-screen TUI for managing todo lists. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## Features

- Two-pane layout: lists on the left, todos on the right, active pane border highlighted.
- Keyboard-first: vim-style movement, `dd` to delete, `gg`/`G` to jump, `tab`/`shift+tab` to switch panes.
- Inline text input lives inside the active pane; todo input is prefixed with `[ ]` to match list items.
- Toggle todo completion with `enter` when in the todo pane.
- Resizes with your terminal to stay full-screen.
- Runs in the terminal alternate screen (like htop) so closing leaves no scrollback.

## Install & Run

```bash
go run main.go
```

## Keybindings

| Keys             | Action                                                     |
| ---------------- | ---------------------------------------------------------- |
| `q`, `ctrl+c`    | Quit                                                       |
| `tab`, `l`       | Switch to the other pane                                   |
| `shift+tab`, `h` | Switch to the other pane                                   |
| `j`, `down`      | Move selection down                                        |
| `k`, `up`        | Move selection up                                          |
| `enter`          | In lists pane: focus todos pane; in todos: toggle complete |
| `a`              | Add list (when in lists) or add todo (when in todos)       |
| `r`              | Rename selected list or todo                               |
| `dd`             | Delete selected list or todo                               |
| `gg`             | Jump to top                                                |
| `G`              | Jump to bottom                                             |
| `esc`            | Cancel current input                                       |

## Notes

- When adding or renaming, the input appears inside the active pane; type and press `enter` to confirm or `esc` to cancel.
- Deletion and jump use timed double-press detection (`d` then `d`, `g` then `g`).
- Layout adapts to terminal size; panels stretch to fill width and height.

## License

MIT
