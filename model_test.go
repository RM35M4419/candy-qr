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
	if styled.message == "" {
		t.Error("style should set a message")
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
	if got := m.buildContent(); !strings.Contains(got, "FN:Jane") {
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
	if m.message != "Saved to candy-qr.png" {
		t.Errorf("message = %q", m.message)
	}
	if _, err := os.Stat("candy-qr.png"); err != nil {
		t.Errorf("expected candy-qr.png to exist: %v", err)
	}
	os.Remove("candy-qr.png")
}
