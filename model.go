package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	screenType screen = iota
	screenForm
	screenPreview
)

var qrTypes = []string{"vCard", "Wi-Fi", "URL"}

var previewMenu = []string{"Export PNG", "Edit", "Style", "Quit"}

type fieldSpec struct {
	key         string
	label       string
	placeholder string
}

var vcardSpecs = []fieldSpec{
	{"GIVEN", "First name", "Jane"},
	{"FAMILY", "Last name", "Doe"},
	{"TEL", "Phone", "+1 555 123 4567"},
	{"TEL2", "Phone 2", "optional"},
	{"EMAIL", "Email", "jane@example.com"},
	{"EMAIL2", "Email 2", "optional"}, {"ORG", "Company", "Acme Inc."},
	{"TITLE", "Title", "Engineer"},
	{"ADR", "Address", "123 Main St, City"},
	{"URL", "Website", "https://example.com"},
	{"NOTE", "Notes", "Say hi!"},
}

var wifiSpecs = []fieldSpec{
	{"S", "Network name (SSID)", "MyWiFi"},
	{"P", "Password", "hunter2"},
	{"T", "Encryption", "WPA"},
}

var urlSpecs = []fieldSpec{
	{"URL", "URL", "https://example.com"},
}

func fieldSpecsFor(typ string) []fieldSpec {
	switch typ {
	case "vCard":
		return vcardSpecs
	case "Wi-Fi":
		return wifiSpecs
	case "URL":
		return urlSpecs
	}
	return nil
}

type field struct {
	key   string
	label string
	input textinput.Model
}

type Model struct {
	screen screen

	typeCursor int

	fields     []field
	fieldIndex int

	previewCursor int
	message       string

	width, height int
}

func initialModel() Model {
	return Model{
		screen:     screenType,
		typeCursor: 0,
	}
}

func newFields(specs []fieldSpec) []field {
	fields := make([]field, len(specs))
	for i, s := range specs {
		ti := textinput.New()
		ti.Placeholder = s.placeholder
		ti.Prompt = ""
		ti.CharLimit = 200
		ti.Width = 36
		fields[i] = field{key: s.key, label: s.label, input: ti}
	}
	if len(fields) > 0 {
		fields[0].input.Focus()
	}
	return fields
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch m.screen {
		case screenType:
			return m.updateType(msg)
		case screenForm:
			return m.updateForm(msg)
		case screenPreview:
			return m.updatePreview(msg)
		}
	}
	return m, nil
}

func (m Model) updateType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.typeCursor = (m.typeCursor - 1 + len(qrTypes)) % len(qrTypes)
	case "down", "j":
		m.typeCursor = (m.typeCursor + 1) % len(qrTypes)
	case "enter":
		m.screen = screenForm
		m.fields = newFields(fieldSpecsFor(qrTypes[m.typeCursor]))
		m.fieldIndex = 0
	}
	return m, nil
}

func (m Model) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "enter":
		m.nextField()
	case "shift+tab":
		m.prevField()
	case "ctrl+s":
		m.screen = screenPreview
		m.previewCursor = 0
		m.message = ""
	case "esc":
		m.screen = screenType
	default:
		var cmd tea.Cmd
		f := &m.fields[m.fieldIndex]
		f.input, cmd = f.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) updatePreview(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.previewCursor = (m.previewCursor - 1 + len(previewMenu)) % len(previewMenu)
	case "down", "j":
		m.previewCursor = (m.previewCursor + 1) % len(previewMenu)
	case "enter":
		switch previewMenu[m.previewCursor] {
		case "Export PNG":
			m.exportPNG()
		case "Edit":
			m.screen = screenForm
		case "Style":
			m.message = "Style editing is coming soon ✨"
		case "Quit":
			return m, tea.Quit
		}
	case "esc":
		m.screen = screenForm
	}
	return m, nil
}

func (m *Model) nextField() {
	if len(m.fields) == 0 {
		return
	}
	m.fields[m.fieldIndex].input.Blur()
	m.fieldIndex = (m.fieldIndex + 1) % len(m.fields)
	m.fields[m.fieldIndex].input.Focus()
}

func (m *Model) prevField() {
	if len(m.fields) == 0 {
		return
	}
	m.fields[m.fieldIndex].input.Blur()
	m.fieldIndex = (m.fieldIndex - 1 + len(m.fields)) % len(m.fields)
	m.fields[m.fieldIndex].input.Focus()
}

func (m Model) View() string {
	switch m.screen {
	case screenType:
		return m.typeView()
	case screenForm:
		return m.formView()
	case screenPreview:
		return m.previewView()
	}
	return ""
}

func (m Model) typeView() string {
	var b strings.Builder
	b.WriteString("🍬 Candy QR\n\n")
	b.WriteString("What do you want to share?\n\n")
	for i, t := range qrTypes {
		cursor := "  "
		if i == m.typeCursor {
			cursor = "> "
		}
		b.WriteString(cursor + t + "\n")
	}
	b.WriteString("\n" + m.footer())
	return b.String()
}

func (m Model) formView() string {
	title := "🍬 Candy QR — " + qrTypes[m.typeCursor]

	var form strings.Builder
	for i, f := range m.fields {
		cursor := "  "
		if i == m.fieldIndex {
			cursor = "> "
		}
		form.WriteString(cursor + f.label + "\n")
		form.WriteString("  " + f.input.View() + "\n\n")
	}

	preview := m.renderPreview()

	left := lipgloss.NewStyle().Width(40).Render(form.String())
	right := lipgloss.NewStyle().PaddingLeft(2).Render(preview)

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	return title + "\n\n" + body + "\n\n" + m.footer()
}

func (m Model) previewView() string {
	var b strings.Builder
	b.WriteString("🍬 Candy QR — " + qrTypes[m.typeCursor] + "\n\n")
	b.WriteString(renderQRString(m.buildContent()) + "\n\n")

	for i, item := range previewMenu {
		cursor := "  "
		if i == m.previewCursor {
			cursor = "> "
		}
		b.WriteString(cursor + item + "\n")
	}

	if m.message != "" {
		b.WriteString("\n" + m.message + "\n")
	}

	b.WriteString("\n" + m.footer())
	return b.String()
}

func (m Model) footer() string {
	switch m.screen {
	case screenType:
		return "↑/↓ select · enter choose · ctrl+c quit"
	case screenForm:
		return "tab/shift+tab move · enter next · ctrl+s generate · esc back · ctrl+c quit"
	case screenPreview:
		return "↑/↓ select · enter confirm · esc back · ctrl+c quit"
	}
	return ""
}
