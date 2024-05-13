package test

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/taz03/monkeytui/config"
)

type Test struct {
    Words  []string
    Config config.Config

    typedWords []string
    pos        [2]int
}

func (this *Test) Init() tea.Cmd {
    this.typedWords = []string{""}
    return nil
}

func (this *Test) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case " ":
            this.typedWords = append(this.typedWords, "")
            this.pos[0]++
            this.pos[1] = 0

        case "backspace":
            if this.pos[1]--; this.pos[1] < 0 {
                if this.pos[0] > 0 {
                    this.pos[0]--
                    this.typedWords = this.typedWords[:len(this.typedWords) - 1]
                }
                this.pos[1] = len(this.typedWords[this.pos[0]])
            } else {
                this.typedWords[this.pos[0]] = this.typedWords[this.pos[0]][:this.pos[1]]
            }

        default:
            this.typedWords[len(this.typedWords) - 1] += msg.String()
            this.pos[1]++
        }
    }

    return this, nil
}

func (this *Test) View() string {
    styledWhitespace := lipgloss.NewStyle().Background(this.Config.BackgroundColor()).Render(" ")
    style := lipgloss.NewStyle().Background(this.Config.BackgroundColor()).Bold(true)

    test := strings.Builder{}
    for i := range (this.pos[0] + 1) {
        wordStyle := style.Copy()
        if i < this.pos[0] && this.Words[i] != this.typedWords[i] {
            wordStyle.Underline(true)
        }

        for j, char := range this.typedWords[i] {
            str := string(char)

            charStyle := wordStyle.Copy()
            switch {
            case j >= len(this.Words[i]):
                charStyle = this.Config.StyleErrorExtra(charStyle, str)
            case str != string(this.Words[i][j]):
                charStyle = this.Config.StyleError(charStyle, str, string(this.Words[i][j]))
            default:
                charStyle = this.Config.StyleCorrect(charStyle, str)
            }

            test.WriteString(charStyle.Render())
        }

        if i < this.pos[0] {
            if len(this.Words[i]) >= len(this.typedWords[i]) {
                remainingWord := this.Words[i][len(this.typedWords[i]):]
                test.WriteString(this.Config.StyleUntyped(wordStyle, remainingWord).Render())
            }

            test.WriteString(styledWhitespace)
        }
    }

    if len(this.Words[this.pos[0]]) > len(this.typedWords[this.pos[0]]) {
        test.WriteString(this.Config.StyleUntyped(style.Copy(), this.Words[this.pos[0]][this.pos[1]:]).Render())
    }
    test.WriteString(styledWhitespace)
    test.WriteString(this.Config.StyleUntyped(style.Copy(), strings.Join(this.Words[this.pos[0] + 1:], " ")).Render())

    return test.String()
}
