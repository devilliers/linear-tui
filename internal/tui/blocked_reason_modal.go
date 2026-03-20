package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// BlockedReasonModal manages a text input prompt for entering a blocked reason.
type BlockedReasonModal struct {
	app       *App
	modal     *tview.Flex
	form      *tview.Form
	bodyField *tview.TextArea
	onSubmit  func(reason string)
}

// NewBlockedReasonModal creates a new blocked reason modal.
func NewBlockedReasonModal(app *App) *BlockedReasonModal {
	brm := &BlockedReasonModal{
		app: app,
	}

	brm.form = tview.NewForm()
	brm.form.SetBackgroundColor(app.theme.HeaderBg)
	brm.form.SetFieldBackgroundColor(app.theme.InputBg)
	brm.form.SetFieldTextColor(app.theme.Foreground)
	brm.form.SetButtonBackgroundColor(app.theme.Accent)
	brm.form.SetButtonTextColor(app.theme.SelectionText)
	brm.form.SetLabelColor(app.theme.Foreground)

	brm.form.AddTextArea("Reason", "", 60, 6, 0, nil)
	if item := brm.form.GetFormItemByLabel("Reason"); item != nil {
		if textArea, ok := item.(*tview.TextArea); ok {
			brm.bodyField = textArea
		}
	}

	brm.form.AddButton("Submit", func() {
		reason := brm.bodyField.GetText()
		brm.Hide()
		if brm.onSubmit != nil {
			brm.onSubmit(reason)
		}
	})
	brm.form.AddButton("Skip", func() {
		brm.Hide()
		if brm.onSubmit != nil {
			brm.onSubmit("")
		}
	})
	brm.form.AddButton("Cancel", func() {
		brm.Hide()
	})

	headerView := tview.NewTextView()
	headerView.SetText("Why is this blocked?")
	headerView.SetTextColor(app.theme.Accent)
	headerView.SetBackgroundColor(app.theme.HeaderBg)

	helpView := tview.NewTextView()
	helpView.SetText("Esc: cancel • Ctrl+Enter: submit • Tab: navigate")
	helpView.SetTextColor(app.theme.SecondaryText)
	helpView.SetBackgroundColor(app.theme.HeaderBg)
	helpView.SetTextAlign(tview.AlignCenter)

	modalContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(headerView, 1, 0, false).
		AddItem(brm.form, 0, 1, true).
		AddItem(helpView, 1, 0, false)
	modalContent.Box = tview.NewBox().SetBackgroundColor(app.theme.HeaderBg)
	modalContent.SetBackgroundColor(app.theme.HeaderBg).
		SetBorder(true).
		SetBorderColor(app.theme.Accent).
		SetTitle(" Blocked Reason ").
		SetTitleColor(app.theme.Foreground)
	padding := app.density.ModalPadding
	modalContent.SetBorderPadding(padding.Top, padding.Bottom, padding.Left, padding.Right)

	brm.modal = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(modalContent, 16, 0, true).
			AddItem(nil, 0, 1, false), 75, 0, true).
		AddItem(nil, 0, 1, false)
	brm.modal.SetBackgroundColor(app.theme.Background)

	return brm
}

// Show displays the blocked reason modal.
func (brm *BlockedReasonModal) Show(onSubmit func(reason string)) {
	brm.onSubmit = onSubmit
	brm.bodyField.SetText("", true)

	brm.app.pages.AddPage("blocked_reason", brm.modal, true, true)
	brm.app.pages.SendToFront("blocked_reason")
	brm.app.app.SetFocus(brm.form)
}

// Hide hides the blocked reason modal.
func (brm *BlockedReasonModal) Hide() {
	brm.app.pages.RemovePage("blocked_reason")
	brm.app.updateFocus()
}

// HandleKey handles keyboard input for the blocked reason modal.
func (brm *BlockedReasonModal) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		brm.Hide()
		return nil
	case tcell.KeyEnter:
		mod := event.Modifiers()
		if mod&tcell.ModCtrl != 0 || mod&tcell.ModMeta != 0 {
			reason := brm.bodyField.GetText()
			brm.Hide()
			if brm.onSubmit != nil {
				brm.onSubmit(reason)
			}
			return nil
		}
	}
	return event
}
