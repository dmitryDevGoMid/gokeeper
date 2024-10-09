package opfilessenddialog

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fogleman/ease"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

// General stuff for styling the view
var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)

	quitViewStyle = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))
	choiceStyle   = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("241"))

	// Gradient colors we'll use for the progress bar
	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)

type (
	tickMsg  struct{}
	frameMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

type FilesDialog struct {
	data     *model.Data
	Choice   int
	Chosen   bool
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
	logField keeperlog.LogField
}

type Page interface {
	tea.Model
}

func NewOpFilesDialog(model *model.Data) Page {
	fileDialog := initialModel()
	fileDialog.data = model

	// Логируем запросы на сервер в файл keeper.log
	logField := keeperlog.LogField{NextStepByName: fileDialog.data.NextStep.NextStepByName}
	logField.Action = "NewOpFilesDialog"
	logField.RequestByName = fileDialog.data.NextStep.RequestByName

	// Установка UserID в логгер
	fileDialog.data.Log.SetUserID(fileDialog.data.User.IDUser)
	fileDialog.logField = logField
	return fileDialog
}

func initialModel() FilesDialog {
	initialModel := FilesDialog{Choice: 0, Chosen: false, Frames: 0, Progress: 0, Loaded: false, Quitting: false}
	return initialModel
}

func (m FilesDialog) Init() tea.Cmd {
	return nil
}
func (m FilesDialog) NextStepToSend() {
	m.logField.Method = "NextStepToSend"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set Step: requesthttp and sendfiles")

	m.data.NextStep.NextStepByName = "requesthttp"
	m.data.NextStep.RequestByName = "sendfiles"
}

func (m FilesDialog) NextStepBackToSelectedFiles() {
	m.logField.Method = "NextStepBackToSelectedFiles"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Back Step: selectfiles")
	m.data.NextStep.NextStepByName = "selectfiles"
}

func (m FilesDialog) NextStepBackToFiles() {
	m.logField.Method = "NextStepBackToFiles"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Back Step: files")
	m.data.NextStep.NextStepByName = "files"
}

type funcUp func(msg tea.Msg, m FilesDialog) (tea.Model, tea.Cmd)

func Prompt(msg tea.Msg, m FilesDialog) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

// Main update function.
func (m FilesDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logField.Method = "Update"
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.NextStepBackToFiles()
			return m, tea.Quit
		}

		if k == "enter" {
			if m.Choice == 0 {
				m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Enter Set - Choice - 0")
				m.NextStepToSend()
				return m, tea.Quit
			}

			if m.Choice == 1 {
				m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Enter Set - Choice - 1")
				m.NextStepBackToSelectedFiles()
				return m, tea.Quit
			}

			if m.Choice == 2 {
				m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Enter Set - Choice - 2")
				m.NextStepBackToFiles()
				return m, tea.Quit
			}
		}
	}

	getUp := m.GetUpdate()

	return getUp(msg, m)
}

func (m FilesDialog) GetUpdate() funcUp {
	if m.Quitting {
		return Prompt
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		return updateChoices
	}
	return updateChosen
}

// The main view, which just calls the appropriate sub-view
func (m FilesDialog) View() string {
	var s string
	if m.Quitting {
		text := lipgloss.JoinHorizontal(lipgloss.Top, "You have unsaved changes. Quit without saving?", choiceStyle.Render("[yn]"))
		return quitViewStyle.Render(text)
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return mainStyle.Render("\n" + s + "\n\n")
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m FilesDialog) (tea.Model, tea.Cmd) {
	m.logField.Method = "updateChoices"

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set - down")
			m.Choice++
			if m.Choice > 3 {
				m.Choice = 3
			}
		case "k", "up":
			m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set - up")
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		}
	}

	return m, nil
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m FilesDialog) (tea.Model, tea.Cmd) {
	m.logField.Method = "updateChosen"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set - Chosen")

	switch msg.(type) {
	case frameMsg:
		if !m.Loaded {
			m.Frames++
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
			if m.Progress >= 1 {
				m.Progress = 1
				m.Loaded = true
				return m, tick()
			}
			return m, frame()
		}
	}

	return m, nil
}

// Sub-views

// The first view, where you're choosing a task
func choicesView(m FilesDialog) string {
	c := m.Choice

	tpl := "Selected file:"
	tpl += m.data.OptionFiles.SelectedFile
	tpl += "\n\n"
	tpl += "What to do?\n\n"
	tpl += "%s\n\n"
	tpl += "Program quits in %s seconds\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s",
		checkbox("Upload this file", c == 0),
		checkbox("Select another file", c == 1),
		checkbox("Back to menu", c == 2),
	)

	return fmt.Sprintf(tpl, choices) //ticksStyle.Render(strconv.Itoa(m.Ticks)))
}

// The second view, after a task has been chosen
func chosenView(m FilesDialog) string {
	var msg string

	switch m.Choice {
	case 0:
		msg = fmt.Sprintf("Carrot planting?\n\nCool, we'll need %s and %s...", keywordStyle.Render("libgarden"), keywordStyle.Render("vegeutils"))
	case 1:
		m.Quitting = true
	case 2:
		m.Quitting = true
	default:
		msg = ""
	}

	return msg
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

func progressbar(percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += ramp[i].Render(progressFullChar)
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// Generate a blend of colors.
func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
