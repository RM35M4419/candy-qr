package main

import (
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func key(t tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: t}
}

func update(t *testing.T, m Model, msg tea.Msg) Model {
	t.Helper()
	nm, _ := m.Update(msg)
	mm, ok := nm.(Model)
	if !ok {
		t.Fatalf("Update returned %T, want Model", nm)
	}
	return mm
}

func TestUpdateTypeSelection(t *testing.T) {
	m := initialModel()

	if got := update(t, m, key(tea.KeyDown)).typeCursor; got != 1 {
		t.Errorf("down: cursor = %d, want 1", got)
	}
	if got := update(t, m, key(tea.KeyUp)).typeCursor; got != 2 {
		t.Errorf("up wraps: cursor = %d, want 2", got)
	}

	selected := update(t, m, key(tea.KeyEnter))
	if selected.screen != screenForm {
		t.Errorf("enter: screen = %v, want form", selected.screen)
	}
	if len(selected.fields) != len(vcardSpecs) {
		t.Errorf("enter: fields = %d, want %d", len(selected.fields), len(vcardSpecs))
	}
}

func TestUpdateFormNavigation(t *testing.T) {
	m := initialModel()
	m.screen = screenForm
	m.fields = newFields(vcardSpecs)
	m.fieldIndex = 0

	if got := update(t, m, key(tea.KeyTab)).fieldIndex; got != 1 {
		t.Errorf("tab: fieldIndex = %d, want 1", got)
	}
	if got := update(t, m, key(tea.KeyEnter)).fieldIndex; got != 1 {
		t.Errorf("enter: fieldIndex = %d, want 1", got)
	}
	if got := update(t, m, key(tea.KeyShiftTab)).fieldIndex; got != len(vcardSpecs)-1 {
		t.Errorf("shift+tab wraps: fieldIndex = %d, want %d", got, len(vcardSpecs)-1)
	}
}

func TestUpdateFormCommit(t *testing.T) {
	m := initialModel()
	m.screen = screenForm
	m.fields = newFields(vcardSpecs)

	committed := update(t, m, key(tea.KeyCtrlS))
	if committed.screen != screenPreview {
		t.Errorf("ctrl+s: screen = %v, want preview", committed.screen)
	}
}

func TestUpdateFormEsc(t *testing.T) {
	m := initialModel()
	m.screen = screenForm
	m.fields = newFields(vcardSpecs)

	if got := update(t, m, key(tea.KeyEsc)).screen; got != screenType {
		t.Errorf("esc: screen = %v, want type", got)
	}
}

func TestUpdatePreviewMenu(t *testing.T) {
	m := initialModel()
	m.screen = screenPreview
	m.fields = newFields(urlSpecs)
	m.fields[0].input.SetValue("https://example.com")

	m = update(t, m, key(tea.KeyDown)) // cursor -> 1, "Edit"
	if m.previewCursor != 1 {
		t.Errorf("down: previewCursor = %d, want 1", m.previewCursor)
	}

	// cursor on "Edit" returns to the form
	edited := update(t, m, key(tea.KeyEnter))
	if edited.screen != screenForm {
		t.Errorf("enter on Edit: screen = %v, want form", edited.screen)
	}
}

func TestUpdatePreviewStyle(t *testing.T) {
	m := initialModel()
	m.screen = screenPreview
	m.fields = newFields(urlSpecs)
	m.previewCursor = 2 // "Style"

	styled := update(t, m, key(tea.KeyEnter))
	if styled.screen != screenStyle {
		t.Errorf("style should switch to screenStyle, got %v", styled.screen)
	}
}

func TestUpdatePreviewEsc(t *testing.T) {
	m := initialModel()
	m.screen = screenPreview
	m.fields = newFields(urlSpecs)

	if got := update(t, m, key(tea.KeyEsc)).screen; got != screenForm {
		t.Errorf("esc: screen = %v, want form", got)
	}
}

func TestBuildContentDispatch(t *testing.T) {
	m := initialModel()

	m.typeCursor = 0
	m.fields = newFields(vcardSpecs)
	m.fields[0].input.SetValue("Jane")
	m.fields[1].input.SetValue("Doe")
	if got := m.buildContent(); !strings.Contains(got, "FN:Jane Doe") {
		t.Errorf("vcard content = %q", got)
	}

	m.typeCursor = 1
	m.fields = newFields(wifiSpecs)
	m.fields[0].input.SetValue("Net")
	if got := m.buildContent(); !strings.Contains(got, "S:Net") {
		t.Errorf("wifi content = %q", got)
	}

	m.typeCursor = 2
	m.fields = newFields(urlSpecs)
	m.fields[0].input.SetValue("https://x.com")
	if got := m.buildContent(); got != "https://x.com" {
		t.Errorf("url content = %q", got)
	}
}

func TestExportPNGEmpty(t *testing.T) {
	m := initialModel()
	m.screen = screenPreview
	m.fields = newFields(vcardSpecs)

	m.exportPNG()
	if m.message != "Nothing to export yet" {
		t.Errorf("message = %q", m.message)
	}
}

func TestExportPNGSuccess(t *testing.T) {
	m := initialModel()
	m.typeCursor = 2 // URL
	m.screen = screenPreview
	m.fields = newFields(urlSpecs)
	m.fields[0].input.SetValue("https://example.com")

	m.exportPNG()
	if !strings.HasPrefix(m.message, "Saved to ") || !strings.HasSuffix(m.message, "candy-qr.png") {
		t.Errorf("unexpected export message: %q", m.message)
	}
	if _, err := os.Stat("candy-qr.png"); err != nil {
		t.Errorf("expected candy-qr.png to exist: %v", err)
	}
	os.Remove("candy-qr.png")
}

func TestViews(t *testing.T) {
	m := initialModel()
	if v := m.View(); !strings.Contains(v, "What do you want to share?") {
		t.Errorf("typeView missing prompt: %s", v)
	}

	m.screen = screenForm
	m.fields = newFields(vcardSpecs)
	if v := m.View(); !strings.Contains(v, "First name") || !strings.Contains(v, "Preview") {
		t.Errorf("formView missing expected elements: %s", v)
	}

	m.screen = screenPreview
	m.fields[0].input.SetValue("Jane")
	m.fields[1].input.SetValue("Doe")
	if v := m.View(); !strings.Contains(v, "Export PNG") {
		t.Errorf("previewView missing menu item: %s", v)
	}
}

func TestFormViewScrolling(t *testing.T) {
	m := initialModel()
	m.screen = screenForm
	m.fields = newFields(vcardSpecs)
	m.height = 20
	m.width = 100

	// Cursor at index 0: should show "more below" but not "more above"
	m.fieldIndex = 0
	v := m.formView()
	if !strings.Contains(v, "more below") {
		t.Errorf("expected bottom scroll indicator when at start, got:\n%s", v)
	}
	if strings.Contains(v, "more above") {
		t.Errorf("unexpected top scroll indicator when at start, got:\n%s", v)
	}

	// Cursor at last index: should show "more above"
	m.fieldIndex = len(m.fields) - 1
	v = m.formView()
	if !strings.Contains(v, "more above") {
		t.Errorf("expected top scroll indicator when at end, got:\n%s", v)
	}
}

func TestFormViewNarrowLayout(t *testing.T) {
	m := initialModel()
	m.screen = screenForm
	m.fields = newFields(urlSpecs)
	m.width = 60
	m.height = 30

	v := m.formView()
	if !strings.Contains(v, "URL") || !strings.Contains(v, "Preview") {
		t.Errorf("narrow formView missing elements: %s", v)
	}
}

func TestWindowSizeMsg(t *testing.T) {
	m := initialModel()
	updated := update(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	if updated.width != 120 || updated.height != 40 {
		t.Errorf("expected 120x40, got %dx%d", updated.width, updated.height)
	}
}

func TestUpdateFormTextInput(t *testing.T) {
	m := initialModel()
	m.screen = screenForm
	m.fields = newFields(urlSpecs)

	// Send key 'a'
	updated := update(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if updated.fields[0].input.Value() != "a" {
		t.Errorf("expected input value 'a', got %q", updated.fields[0].input.Value())
	}
}

func TestUpdateStyleNavigation(t *testing.T) {
	m := initialModel()
	m.screen = screenStyle
	m.fields = newFields(urlSpecs)
	m.fields[0].input.SetValue("https://example.com")
	m.styleCursor = 0 // Preset

	// Right arrow on preset -> next preset
	m = update(t, m, key(tea.KeyRight))
	if m.style.PresetIndex != 1 {
		t.Errorf("preset index = %d, want 1", m.style.PresetIndex)
	}

	// Down arrow -> Shape
	m = update(t, m, key(tea.KeyDown))
	if m.styleCursor != 1 {
		t.Errorf("styleCursor = %d, want 1", m.styleCursor)
	}
	m = update(t, m, key(tea.KeyRight))
	if m.style.Shape != ShapeSquare {
		t.Errorf("shape = %v, want Square", m.style.Shape)
	}

	// Down arrow -> Gradient
	m = update(t, m, key(tea.KeyDown))
	if m.styleCursor != 2 {
		t.Errorf("styleCursor = %d, want 2", m.styleCursor)
	}
	m = update(t, m, key(tea.KeyRight))
	if m.style.Gradient != GradientVertical {
		t.Errorf("gradient = %v, want Vertical", m.style.Gradient)
	}

	// Down arrow -> Logo input
	m = update(t, m, key(tea.KeyDown))
	if m.styleCursor != 3 {
		t.Errorf("styleCursor = %d, want 3", m.styleCursor)
	}
	// Type 'l', 'o', 'g', 'o', '.', 'p', 'n', 'g'
	m = update(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if m.style.LogoPath != "l" {
		t.Errorf("logoPath = %q, want 'l'", m.style.LogoPath)
	}

	// Enter from logo input moves to Apply
	m = update(t, m, key(tea.KeyEnter))
	if m.styleCursor != 4 {
		t.Errorf("styleCursor = %d, want 4 (Apply)", m.styleCursor)
	}

	// Enter on Apply returns to Preview
	m = update(t, m, key(tea.KeyEnter))
	if m.screen != screenPreview {
		t.Errorf("enter on Apply should return to screenPreview, got %v", m.screen)
	}
}

func TestStyleView(t *testing.T) {
	m := initialModel()
	m.screen = screenStyle
	m.fields = newFields(urlSpecs)
	m.fields[0].input.SetValue("https://example.com")
	m.width = 100
	m.height = 30

	v := m.styleView()
	if !strings.Contains(v, "Theme Preset") || !strings.Contains(v, "Module Corners") || !strings.Contains(v, "Gradient Direction") {
		t.Errorf("styleView missing settings: %s", v)
	}
	if !strings.Contains(v, "Preview") {
		t.Errorf("styleView missing preview: %s", v)
	}

	// Test narrow width
	m.width = 60
	narrowView := m.styleView()
	if !strings.Contains(narrowView, "Theme Preset") {
		t.Errorf("narrow styleView missing content: %s", narrowView)
	}
}
