package cli

import (
	"github.com/elves/elvish/cli/el"
	"github.com/elves/elvish/cli/el/codearea"
	"github.com/elves/elvish/styled"
)

// AppSpec specifies the configuration and initial state for an App.
type AppSpec struct {
	TTY               TTY
	MaxHeight         func() int
	RPromptPersistent func() bool
	BeforeReadline    func()
	AfterReadline     func(string)

	Highlighter Highlighter
	Prompt      Prompt
	RPrompt     Prompt

	OverlayHandler el.Handler
	Abbreviations  func(f func(abbr, full string))
	QuotePaste     func() bool

	CodeAreaState codearea.State
	State         State
}

// Highlighter represents a code highlighter whose result can be delivered
// asynchronously.
type Highlighter interface {
	// Get returns the highlighted code and any static errors.
	Get(code string) (styled.Text, []error)
	// LateUpdates returns a channel for delivering late updates.
	LateUpdates() <-chan styled.Text
}

// A Highlighter implementation that always returns plain text.
type dummyHighlighter struct{}

func (dummyHighlighter) Get(code string) (styled.Text, []error) {
	return styled.Plain(code), nil
}

func (dummyHighlighter) LateUpdates() <-chan styled.Text { return nil }

// A Highlighter implementation useful for testing.
type testHighlighter struct {
	get         func(code string) (styled.Text, []error)
	lateUpdates chan styled.Text
}

func (hl testHighlighter) Get(code string) (styled.Text, []error) {
	return hl.get(code)
}

func (hl testHighlighter) LateUpdates() <-chan styled.Text {
	return hl.lateUpdates
}

// Prompt represents a prompt whose result can be delivered asynchronously.
type Prompt interface {
	// Trigger requests a re-computation of the prompt. The force flag is set
	// when triggered for the first time during a ReadCode session or after a
	// SIGINT that resets the editor.
	Trigger(force bool)
	// Get returns the current prompt.
	Get() styled.Text
	// LastUpdates returns a channel for delivering late updates.
	LateUpdates() <-chan styled.Text
}

// A Prompt implementation that always return the same styled.Text.
type constPrompt struct{ t styled.Text }

func (constPrompt) Trigger(force bool)              {}
func (p constPrompt) Get() styled.Text              { return p.t }
func (constPrompt) LateUpdates() <-chan styled.Text { return nil }

// A Prompt implementation useful for testing.
type testPrompt struct {
	trigger     func(force bool)
	get         func() styled.Text
	lateUpdates chan styled.Text
}

func (p testPrompt) Trigger(force bool) {
	if p.trigger != nil {
		p.trigger(force)
	}
}

func (p testPrompt) Get() styled.Text {
	if p.get != nil {
		return p.get()
	}
	return nil
}

func (p testPrompt) LateUpdates() <-chan styled.Text {
	return p.lateUpdates
}