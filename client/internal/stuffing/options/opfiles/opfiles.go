package opfiles

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	textGreenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#008000"))
	textRedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)

type Files struct {
	data         *model.Data
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
	logField     keeperlog.LogField
}

type Page interface {
	tea.Model
}

func NewOpFiles(model *model.Data) Page {
	lp := initialModel()
	lp.data = model

	// Логируем запросы на сервер в файл keeper.log
	logField := keeperlog.LogField{NextStepByName: lp.data.NextStep.NextStepByName}
	logField.Action = "NewOpFiles"
	logField.RequestByName = lp.data.NextStep.RequestByName

	// Установка UserID в логгер
	lp.data.Log.SetUserID(lp.data.User.IDUser)
	lp.logField = logField
	return lp
}

func initialModel() Files {

	fp := filepicker.New()
	fp.AllowedTypes = []string{".jpeg", ".jpg", ".zip", ".dmg", ".mp4", ".app"}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	m := Files{
		filepicker: fp,
	}
	return m
}

func (m Files) BackToMenuFiles() {
	m.data.NextStep.NextStepByName = "files"
}

func (m Files) NextStepToAfterSelected() {
	m.data.NextStep.NextStepByName = "opfilessenddialog"
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m Files) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m Files) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.BackToMenuFiles()
			m.quitting = true
			return m, tea.Quit
		case "ctrl+s", "s":
			m.NextStepToAfterSelected()
			m.quitting = true
			return m, tea.Quit
		}

	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m Files) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder

	s.WriteString("\n  ")
	send := textGreenStyle.Render("'Ctr+S - send'")
	s.WriteString("After selected file, click: ")
	s.WriteString(send)
	s.WriteString(" or ")
	exit := textRedStyle.Render("'Ctr+C - quit'")
	s.WriteString(exit)
	s.WriteString("\n  ")

	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	} else {
		m.data.OptionFiles.SelectedFile = m.selectedFile
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}
