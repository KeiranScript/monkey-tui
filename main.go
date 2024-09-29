package main

import (
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/taz03/monkeytui/config"
	"github.com/taz03/monkeytui/test"
)

type model struct {
	Test       *test.Model
	Config     *config.Model
	TotalChars int
	TotalWords int
	WPM        string
	StartTime  time.Time
	ShowStats  bool
}

var width, height int

func main() {
	userConfigPath := "config.json"
	if _, err := os.Stat(userConfigPath); os.IsNotExist(err) {
		userConfigPath = "config/default.json"
	}
	userConfig := config.New(userConfigPath)

	app := tea.NewProgram(model{
		Test:   test.New(userConfig),
		Config: userConfig,
	}, tea.WithAltScreen())
	go userConfig.MonkeyTheme.Update(app)

	app.Run()
}

func (m model) calculateTestWidth() int {
	if m.Config.MaxLineWidth == 0 {
		return width - 10
	}

	if m.Config.MaxLineWidth > width {
		return width
	}

	return m.Config.MaxLineWidth
}

func (m model) Init() tea.Cmd {
	m.StartTime = time.Now()
	_ = m.StartTime
	return m.Test.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width, height = msg.Width, msg.Height

	case tea.KeyMsg:
		if msg.String() == m.Config.RestartKey() {
			m.Test = test.New(m.Config)
			m.Test.Width = m.calculateTestWidth()
			m.TotalChars = 0
			m.TotalWords = 0
			m.StartTime = time.Now()
			return m, m.Test.Init()
		}

		switch msg.String() {
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		case tea.KeyTab.String():
			m.ShowStats = !m.ShowStats
			return m, nil
		}
	}

	_, cmd := m.Test.Update(msg)
	return m, cmd
}

func (m model) View() string {
	m.Test.Width = m.calculateTestWidth()
	m.Test.ProgressBar.Width = width

	// Calculate total characters and words typed
	for _, word := range m.Test.TypedWords {
		m.TotalChars += len(word)
		m.TotalWords += 1
	}

	// Calculate WPM
	elapsed := time.Since(m.StartTime).Seconds()
	if elapsed > 0 {
		wpm := float64(m.TotalWords) / (elapsed / 60)
		m.WPM = strconv.FormatFloat(wpm, 'f', 2, 64) + " WPM"
	} else {
		m.WPM = "0.00 WPM"
	}

	// Display total characters, words, and WPM directly above the top row of words
	statsView := "Total Characters: " + strconv.Itoa(m.TotalChars) + "\n" +
		"Total Words: " + strconv.Itoa(m.TotalWords) + "\n" +
		"WPM: " + m.WPM + "\n"

	view := m.Test.View()
	if m.ShowStats {
		view = statsView + view
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.Place(
			width,
			height-3,
			lipgloss.Center,
			lipgloss.Center,
			view,
			lipgloss.WithWhitespaceBackground(m.Config.BackgroundColor()),
		),
		m.Test.ProgressBar.View(),
	)
}
