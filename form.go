package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// フォームの入力フィールド
const (
	nameInput = iota
	emailInput
	submitButton
)

// フォームの状態
type formState int

const (
	formInput formState = iota
	formSubmitted
)

// フォームモデル
type formModel struct {
	inputs       []textinput.Model
	focusIndex   int
	state        formState
	errorMessage string
	submitted    bool
}

// 入力データの構造体
type formData struct {
	name  string
	email string
}

// コンストラクタ
func NewFormModel() formModel {
	// 名前入力フィールド
	nameField := textinput.New()
	nameField.Placeholder = "例: 山田太郎"
	nameField.Focus()
	nameField.CharLimit = 50
	nameField.Width = 30

	// メール入力フィールド
	emailField := textinput.New()
	emailField.Placeholder = "例: taro@example.com"
	emailField.CharLimit = 100
	emailField.Width = 30

	inputs := []textinput.Model{nameField, emailField}

	return formModel{
		inputs:     inputs,
		focusIndex: nameInput,
		state:      formInput,
	}
}

// Init - 初期化
func (m formModel) Init() tea.Cmd {
	// 最初のフィールドにフォーカスを当てる
	return textinput.Blink
}

// Update - メッセージ処理
func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focusIndex == submitButton {
				// 送信処理
				if err := m.validate(); err != nil {
					m.errorMessage = err.Error()
					return m, nil
				}
				m.state = formSubmitted
				m.submitted = true
				return m, nil
			}
			// Enterキーで次のフィールドへ
			m = m.nextField()
			return m, nil

		case tea.KeyTab, tea.KeyShiftTab:
			// Tab/Shift+Tabでフィールド移動
			if msg.Type == tea.KeyTab {
				m = m.nextField()
			} else {
				m = m.prevField()
			}
			return m, nil

		case tea.KeyUp:
			m = m.prevField()
			return m, nil

		case tea.KeyDown:
			m = m.nextField()
			return m, nil

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				if m.state == formSubmitted {
					return m, tea.Quit
				}
			}

		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}

		// 入力中の状態でエラーメッセージをクリア
		if m.focusIndex < len(m.inputs) && m.errorMessage != "" {
			m.errorMessage = ""
		}
	}

	// textinputの更新
	cmd := m.updateInputs(msg)
	return m, cmd
}

// 入力フィールドの更新
func (m *formModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// 現在フォーカスがあるフィールドのみ更新
	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
	}

	return tea.Batch(cmds...)
}

// 次のフィールドへ移動
func (m formModel) nextField() formModel {
	m.focusIndex++
	if m.focusIndex > submitButton {
		m.focusIndex = nameInput
	}

	// フォーカスを更新
	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}

	return m
}

// 前のフィールドへ移動
func (m formModel) prevField() formModel {
	m.focusIndex--
	if m.focusIndex < nameInput {
		m.focusIndex = submitButton
	}

	// フォーカスを更新
	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}

	return m
}

// バリデーション
func (m formModel) validate() error {
	// 名前のチェック
	name := strings.TrimSpace(m.inputs[nameInput].Value())
	if name == "" {
		return fmt.Errorf("名前を入力してください")
	}

	// メールアドレスのチェック
	email := strings.TrimSpace(m.inputs[emailInput].Value())
	if email == "" {
		return fmt.Errorf("メールアドレスを入力してください")
	}

	// メール形式のチェック
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("正しいメールアドレスの形式で入力してください")
	}

	return nil
}

// フォームデータの取得
func (m formModel) getFormData() formData {
	return formData{
		name:  strings.TrimSpace(m.inputs[nameInput].Value()),
		email: strings.TrimSpace(m.inputs[emailInput].Value()),
	}
}

// View - UIの描画
func (m formModel) View() string {
	// スタイル定義
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Width(15)

	focusedLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true).
		Width(15)

	buttonStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("12")).
		Foreground(lipgloss.Color("15")).
		Padding(0, 3).
		MarginTop(1)

	focusedButtonStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("14")).
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Padding(0, 3).
		MarginTop(1)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		MarginTop(1)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true).
		MarginTop(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true).
		MarginTop(1)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2)

	// 送信済み画面
	if m.state == formSubmitted {
		data := m.getFormData()
		content := titleStyle.Render("📨 フォーム送信完了") + "\n\n"
		content += successStyle.Render("✅ 正常に送信されました！") + "\n\n"
		content += labelStyle.Render("名前:") + " " + data.name + "\n"
		content += labelStyle.Render("メール:") + " " + data.email + "\n\n"
		content += helpStyle.Render("q: 終了")
		return borderStyle.Render(content)
	}

	// フォーム入力画面
	var content strings.Builder
	content.WriteString(titleStyle.Render("📝 ユーザー登録フォーム"))
	content.WriteString("\n\n")

	// 名前フィールド
	label := "名前:"
	if m.focusIndex == nameInput {
		content.WriteString(focusedLabelStyle.Render(label))
	} else {
		content.WriteString(labelStyle.Render(label))
	}
	content.WriteString("\n")
	content.WriteString(m.inputs[nameInput].View())
	content.WriteString("\n\n")

	// メールフィールド
	label = "メールアドレス:"
	if m.focusIndex == emailInput {
		content.WriteString(focusedLabelStyle.Render(label))
	} else {
		content.WriteString(labelStyle.Render(label))
	}
	content.WriteString("\n")
	content.WriteString(m.inputs[emailInput].View())
	content.WriteString("\n\n")

	// 送信ボタン
	button := "[ 送信 ]"
	if m.focusIndex == submitButton {
		content.WriteString(focusedButtonStyle.Render(button))
	} else {
		content.WriteString(buttonStyle.Render(button))
	}

	// エラーメッセージ
	if m.errorMessage != "" {
		content.WriteString("\n")
		content.WriteString(errorStyle.Render("❌ " + m.errorMessage))
	}

	// ヘルプテキスト
	help := "\nTab: 次へ  Shift+Tab: 前へ  Enter: 決定  Esc: 終了"
	content.WriteString(helpStyle.Render(help))

	return borderStyle.Render(content.String())
}
