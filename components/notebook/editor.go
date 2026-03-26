package notebook

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type EditorFinishedMsg struct {
	Content string
	Err     error
	Idx     int
	Note    *Item
}

func openMarkdownInEditor(md string, idx int, note Item) tea.Cmd {
	file, _ := os.CreateTemp("", "*.md") // create temp md file
	file.WriteString(md)
	file.Close()

	editor := os.Getenv("EDITOR") // expect user to define
	if editor == "" {
		editor = "nvim"
	}

	cmd := exec.Command(editor, file.Name())
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		edited, _ := os.ReadFile(file.Name()) // read file contents on close
		return EditorFinishedMsg{Content: string(edited), Err: err, Idx: idx, Note: &note}
	})
}
