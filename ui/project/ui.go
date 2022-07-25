package ui

import (
	"fmt"
	"os"

	"github.com/76creates/stickers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/image"
)

var (
	styleBlank      = lipgloss.NewStyle()
	styleRow        = lipgloss.NewStyle().Align(lipgloss.Center).Foreground(lipgloss.Color("#000000")).Bold(true)
	styleBackground = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#ffffff"))

	logoRowIndex  = 0
	logoCellIndex = 0
)

func ProjectUI() {
	imageModel := image.New(true, true, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})

	initialModel := model{
		logo:    imageModel,
		flexBox: stickers.NewFlexBox(0, 0),
	}

	r1 := initialModel.flexBox.NewRow().AddCells(
		[]*stickers.FlexBoxCell{
			stickers.NewFlexBoxCell(1, 3).SetStyle(styleRow),
			stickers.NewFlexBoxCell(1, 9).SetStyle(styleBackground),
		},
	).SetStyle(styleRow)

	_rows := []*stickers.FlexBoxRow{r1}
	initialModel.flexBox.AddRows(_rows)

	err := tea.NewProgram(initialModel, tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{}
}

type model struct {
	logo    image.Bubble
	flexBox *stickers.FlexBox
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowHeight := msg.Height
		windowWidth := msg.Width
		m.flexBox.SetWidth(windowWidth)
		m.flexBox.SetHeight(windowHeight)
		m.logo.SetSize(windowWidth, windowHeight)
		cmds = append(cmds, m.logo.SetFileName("ui/images/EMERALD-CITY-3.png"))
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	m.logo, cmd = m.logo.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	m.flexBox.ForceRecalculate()
	_r := m.flexBox.Row(logoRowIndex)
	if _r == nil {
		panic("could not find the table row")
	}
	_c := _r.Cell(logoCellIndex)
	if _c == nil {
		panic("could not find the table cell")
	}
	m.logo.SetSize(_c.GetWidth(), _c.GetHeight())
	_c.SetContent(m.logo.View())

	return m.flexBox.Render()
}
