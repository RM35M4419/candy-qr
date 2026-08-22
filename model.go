package main

import (
	"fmt"
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
	screenStyle
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
	{"EMAIL2", "Email 2", "optional"},
	{"ORG", "Company", "Acme Inc."},
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

	style       QRStyle
	styleCursor int
	logoInput   textinput.Model

	width, height int
}

func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "logo.png (optional)"
	ti.Prompt = "  "
	ti.CharLimit = 250
	ti.Width = 34

	return Model{
		screen:     screenType,
		typeCursor: 0,
		style:      defaultStyle(),
		logoInput:  ti,
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
		case screenStyle:
			return m.updateStyle(msg)
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
		if len(m.fields) > 0 {
			return m, m.fields[0].input.Focus()
		}
	}
	return m, nil
}

func (m Model) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "enter":
		cmd := m.nextField()
		return m, cmd
	case "shift+tab":
		cmd := m.prevField()
		return m, cmd
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
			if len(m.fields) > 0 {
				return m, m.fields[m.fieldIndex].input.Focus()
			}
		case "Style":
			m.screen = screenStyle
			m.styleCursor = 0
			m.logoInput.SetValue(m.style.LogoPath)
			m.logoInput.Blur()
			m.message = ""
		case "Quit":
			return m, tea.Quit
		}
	case "esc":
		m.screen = screenForm
		if len(m.fields) > 0 {
			return m, m.fields[m.fieldIndex].input.Focus()
		}
	}
	return m, nil
}

func (m Model) updateStyle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	numItems := 5 // 0: Preset, 1: Shape, 2: Gradient, 3: Logo, 4: Done

	// When focused on Logo input
	if m.styleCursor == 3 && m.logoInput.Focused() {
		switch msg.String() {
		case "tab", "enter", "down":
			m.style.LogoPath = strings.TrimSpace(m.logoInput.Value())
			m.logoInput.Blur()
			m.styleCursor = 4
			return m, nil
		case "shift+tab", "up":
			m.style.LogoPath = strings.TrimSpace(m.logoInput.Value())
			m.logoInput.Blur()
			m.styleCursor = 2
			return m, nil
		case "esc":
			m.style.LogoPath = strings.TrimSpace(m.logoInput.Value())
			m.logoInput.Blur()
			m.screen = screenPreview
			return m, nil
		default:
			var cmd tea.Cmd
			m.logoInput, cmd = m.logoInput.Update(msg)
			m.style.LogoPath = strings.TrimSpace(m.logoInput.Value())
			return m, cmd
		}
	}

	switch msg.String() {
	case "up", "k":
		m.styleCursor = (m.styleCursor - 1 + numItems) % numItems
		if m.styleCursor == 3 {
			return m, m.logoInput.Focus()
		}
	case "down", "j", "tab":
		m.styleCursor = (m.styleCursor + 1) % numItems
		if m.styleCursor == 3 {
			return m, m.logoInput.Focus()
		}
	case "shift+tab":
		m.styleCursor = (m.styleCursor - 1 + numItems) % numItems
		if m.styleCursor == 3 {
			return m, m.logoInput.Focus()
		}
	case "left", "h":
		switch m.styleCursor {
		case 0:
			m.style.PresetIndex = (m.style.PresetIndex - 1 + len(presets)) % len(presets)
		case 1:
			m.style.Shape = ModuleShape((int(m.style.Shape) + 1) % 2)
		case 2:
			m.style.Gradient = GradientDirection((int(m.style.Gradient) - 1 + len(gradientNames)) % len(gradientNames))
		}
	case "right", "l":
		switch m.styleCursor {
		case 0:
			m.style.PresetIndex = (m.style.PresetIndex + 1) % len(presets)
		case 1:
			m.style.Shape = ModuleShape((int(m.style.Shape) + 1) % 2)
		case 2:
			m.style.Gradient = GradientDirection((int(m.style.Gradient) + 1) % len(gradientNames))
		}
	case "enter", " ":
		switch m.styleCursor {
		case 0:
			m.style.PresetIndex = (m.style.PresetIndex + 1) % len(presets)
		case 1:
			m.style.Shape = ModuleShape((int(m.style.Shape) + 1) % 2)
		case 2:
			m.style.Gradient = GradientDirection((int(m.style.Gradient) + 1) % len(gradientNames))
		case 3:
			return m, m.logoInput.Focus()
		case 4:
			m.style.LogoPath = strings.TrimSpace(m.logoInput.Value())
			m.screen = screenPreview
		}
	case "esc":
		m.style.LogoPath = strings.TrimSpace(m.logoInput.Value())
		m.screen = screenPreview
	}
	return m, nil
}

func (m *Model) nextField() tea.Cmd {
	if len(m.fields) == 0 {
		return nil
	}
	m.fields[m.fieldIndex].input.Blur()
	m.fieldIndex = (m.fieldIndex + 1) % len(m.fields)
	return m.fields[m.fieldIndex].input.Focus()
}

func (m *Model) prevField() tea.Cmd {
	if len(m.fields) == 0 {
		return nil
	}
	m.fields[m.fieldIndex].input.Blur()
	m.fieldIndex = (m.fieldIndex - 1 + len(m.fields)) % len(m.fields)
	return m.fields[m.fieldIndex].input.Focus()
}

func (m Model) View() string {
	switch m.screen {
	case screenType:
		return m.typeView()
	case screenForm:
		return m.formView()
	case screenPreview:
		return m.previewView()
	case screenStyle:
		return m.styleView()
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

	start, end := 0, len(m.fields)
	if m.height > 0 {
		avail := m.height - 6
		if avail < 6 {
			avail = 6
		}
		maxVisible := avail / 3
		if maxVisible < 2 {
			maxVisible = 2
		}
		if len(m.fields) > maxVisible {
			start = m.fieldIndex - maxVisible/2
			if start < 0 {
				start = 0
			}
			end = start + maxVisible
			if end > len(m.fields) {
				end = len(m.fields)
				start = end - maxVisible
				if start < 0 {
					start = 0
				}
			}
		}
	}

	var form strings.Builder
	if start > 0 {
		form.WriteString(fmt.Sprintf("  ▲ %d more above\n\n", start))
	}
	for i := start; i < end; i++ {
		f := m.fields[i]
		cursor := "  "
		if i == m.fieldIndex {
			cursor = "> "
		}
		form.WriteString(cursor + f.label + "\n")
		form.WriteString("  " + f.input.View() + "\n\n")
	}
	if end < len(m.fields) {
		form.WriteString(fmt.Sprintf("  ▼ %d more below\n", len(m.fields)-end))
	}

	preview := m.renderPreview()

	left := lipgloss.NewStyle().Width(40).Render(form.String())
	right := lipgloss.NewStyle().PaddingLeft(2).Render(preview)

	var body string
	if m.width > 0 && m.width < 75 {
		body = lipgloss.JoinVertical(lipgloss.Left, left, right)
	} else {
		body = lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	}

	return title + "\n\n" + body + "\n\n" + m.footer()
}

func (m Model) previewView() string {
	var b strings.Builder
	b.WriteString("🍬 Candy QR — " + qrTypes[m.typeCursor] + "\n\n")
	b.WriteString(renderStyledQRString(m.buildContent(), m.style) + "\n\n")

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

func (m Model) styleView() string {
	title := "🍬 Candy QR — QR Style & Theme"

	var left strings.Builder
	p := m.style.CurrentPreset()

	// 0: Preset
	c0 := "  "
	if m.styleCursor == 0 {
		c0 = "> "
	}
	left.WriteString(c0 + "Theme Preset:\n")
	left.WriteString("    ◀ " + p.Name + " ▶\n")
	left.WriteString("    " + lipgloss.NewStyle().Faint(true).Render(p.Description) + "\n")
	swatch := lipgloss.NewStyle().Foreground(lipgloss.Color(p.FgStart)).Render("■ "+p.FgStart) +
		" ➔ " + lipgloss.NewStyle().Foreground(lipgloss.Color(p.FgEnd)).Render("■ "+p.FgEnd)
	left.WriteString("    " + swatch + "\n\n")

	// 1: Shape / Corners
	c1 := "  "
	if m.styleCursor == 1 {
		c1 = "> "
	}
	left.WriteString(c1 + "Module Corners:\n")
	left.WriteString("    ◀ " + shapeNames[m.style.Shape] + " ▶\n\n")

	// 2: Gradient
	c2 := "  "
	if m.styleCursor == 2 {
		c2 = "> "
	}
	left.WriteString(c2 + "Gradient Direction:\n")
	left.WriteString("    ◀ " + gradientNames[m.style.Gradient] + " ▶\n\n")

	// 3: Logo Path
	c3 := "  "
	if m.styleCursor == 3 {
		c3 = "> "
	}
	left.WriteString(c3 + "Center Logo (PNG):\n")
	left.WriteString("  " + m.logoInput.View() + "\n\n")

	// 4: Apply button
	c4 := "  "
	if m.styleCursor == 4 {
		c4 = "> "
	}
	left.WriteString(c4 + "[ Apply & Return to Preview ]\n")

	preview := m.renderPreview()

	leftPanel := lipgloss.NewStyle().Width(42).Render(left.String())
	rightPanel := lipgloss.NewStyle().PaddingLeft(2).Render(preview)

	var body string
	if m.width > 0 && m.width < 78 {
		body = lipgloss.JoinVertical(lipgloss.Left, leftPanel, rightPanel)
	} else {
		body = lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	}

	return title + "\n\n" + body + "\n\n" + m.footer()
}

func (m Model) footer() string {
	switch m.screen {
	case screenType:
		return "↑/↓ select · enter choose · ctrl+c quit"
	case screenForm:
		return "tab/shift+tab move · enter next · ctrl+s generate · esc back · ctrl+c quit"
	case screenPreview:
		return "↑/↓ select · enter confirm · esc back · ctrl+c quit"
	case screenStyle:
		return "↑/↓ move · ←/→ cycle · enter select · esc done · ctrl+c quit"
	}
	return ""
}
