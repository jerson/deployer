package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jerson/deployer/pkg/deployer"
	"github.com/nsf/termbox-go"

	"github.com/VladimirMarkelov/clui"
)

// UI ...
type UI struct {
	deployer *deployer.Deployer
}

// NewUI ...
func NewUI(deployer *deployer.Deployer) *UI {
	return &UI{
		deployer: deployer,
	}
}

func updateProgress(value string, pb *clui.ProgressBar) {
	v, _ := strconv.Atoi(value)
	pb.SetValue(v)
}

func changeTheme(lb *clui.ListBox, btn *clui.Button, tp int) {
	items := clui.ThemeNames()
	dlgType := clui.SelectDialogRadio
	if tp == 1 {
		dlgType = clui.SelectDialogList
	}

	curr := -1
	for i, tName := range items {
		if tName == clui.CurrentTheme() {
			curr = i
			break
		}
	}

	selDlg := clui.CreateSelectDialog("Choose a theme", items, curr, dlgType)
	selDlg.OnClose(func() {
		switch selDlg.Result() {
		case clui.DialogButton1:
			idx := selDlg.Value()
			lb.AddItem(fmt.Sprintf("Selected item: %v", selDlg.Value()))
			lb.SelectItem(lb.ItemCount() - 1)
			if idx != -1 {
				clui.SetCurrentTheme(items[idx])
			}
		}

		btn.SetEnabled(true)
		// ask the composer to repaint all windows
		clui.PutEvent(clui.Event{Type: clui.EventRedraw})
	})
}

// Render ...
func (ui *UI) Render() (err error) {

	//go func() {
	//	go ui.runEvery(time.Second*10, ui.updateHealth)
	//	go ui.runEvery(time.Millisecond*100, ui.updateDeployments)
	//	go ui.runEvery(time.Millisecond*100, ui.updateMain)
	//	go ui.runEvery(time.Second*1, ui.updateProject)
	//}()

	clui.InitLibrary()
	defer clui.DeinitLibrary()

	ui.Layout()

	clui.MainLoop()
	return nil
}

// GetWidthPercent ...
func (ui *UI) GetWidthPercent(percentage float32) int {
	width, _ := clui.ScreenSize()
	return int(float32(width) * percentage)
}

// GetHeightPercent ...
func (ui *UI) GetHeightPercent(percentage float32) int {
	_, height := clui.ScreenSize()
	return int(float32(height) * percentage)
}

// Layout ...
func (ui *UI) Layout() {

	projectView := clui.AddWindow(0, 0, ui.GetWidthPercent(0.3), ui.GetHeightPercent(0.1), "Project")
	engineView := clui.AddWindow(0, ui.GetHeightPercent(0.1), ui.GetWidthPercent(0.3), ui.GetHeightPercent(0.3), "Engine")
	deploymentsView := clui.AddWindow(0, ui.GetHeightPercent(0.4), ui.GetWidthPercent(0.3), ui.GetHeightPercent(0.6), "Console")
	mainView := clui.AddWindow(ui.GetWidthPercent(0.3), 0, ui.GetWidthPercent(0.7), ui.GetHeightPercent(0.97), "Main")

	mainView.SetSizable(false)
	engineView.SetSizable(false)
	deploymentsView.SetSizable(false)
	projectView.SetSizable(false)

	engineView.SetGaps(0, 0)
	engineView.SetPaddings(0, 0)
	//view.SetMaximized(true)

	frmLeft := clui.CreateFrame(mainView, 8, 4, clui.BorderThin, 1)
	frmLeft.SetPack(clui.Vertical)
	frmLeft.SetGaps(clui.KeepValue, 1)
	frmLeft.SetPaddings(1, 1)
	frmLeft.SetBackColor(termbox.ColorGreen)

	reader := clui.CreateTextDisplay(mainView, 4, 4, 1)
	reader.SetBackColor(termbox.ColorBlue)
	reader.SetLineCount(50)
	reader.OnDrawLine(func(ind int) string {
		return fmt.Sprintf("%03d line line line", ind+1)
	})

	//frmTheme := clui.CreateFrame(frmLeft, 8, 1, clui.BorderNone, clui.Fixed)
	//frmTheme.SetGaps(1, clui.KeepValue)
	//checkBox := clui.CreateCheckBox(frmTheme, clui.AutoSize, "Use ListBox", clui.Fixed)
	btnTheme := clui.CreateButton(engineView, ui.GetWidthPercent(0.3), 1, "Exit", clui.Fixed)
	btnTheme.OnClick(func(event clui.Event) {
		go clui.Stop()
	})
	//clui.CreateFrame(frmLeft, 1, 1, clui.BorderNone, 1)
	//
	//frmPb := clui.CreateFrame(frmLeft, 8, 1, clui.BorderNone, clui.Fixed)
	//clui.CreateLabel(frmPb, 1, 1, "[", clui.Fixed)
	//pb := clui.CreateProgressBar(frmPb, 20, 1, 1)
	//pb.SetLimits(0, 10)
	//pb.SetTitle("{{value}} of {{max}}")
	//clui.CreateLabel(frmPb, 1, 1, "]", clui.Fixed)
	//
	//edit := clui.CreateEditField(frmLeft, 5, "0", clui.Fixed)
	//
	//frmEdit := clui.CreateFrame(frmLeft, 8, 1, clui.BorderNone, clui.Fixed)
	//frmEdit.SetPaddings(1, 1)
	//frmEdit.SetGaps(1, clui.KeepValue)
	//btnSet := clui.CreateButton(frmEdit, clui.AutoSize, 4, "Set", clui.Fixed)
	//btnStep := clui.CreateButton(frmEdit, clui.AutoSize, 4, "Step", clui.Fixed)
	//clui.CreateFrame(frmEdit, 1, 1, clui.BorderNone, 1)
	//btnQuit := clui.CreateButton(frmEdit, clui.AutoSize, 4, "Quit", clui.Fixed)
	//
	//logBox := clui.CreateListBox(view, 28, 5, clui.Fixed)
	//
	//clui.ActivateControl(view, edit)
	//
	//edit.OnKeyPress(func(key termbox.Key, ch rune) bool {
	//	if key == termbox.KeyCtrlM {
	//		v := edit.Title()
	//		logBox.AddItem(fmt.Sprintf("New PB value(KeyPress): %v", v))
	//		logBox.SelectItem(logBox.ItemCount() - 1)
	//		updateProgress(v, pb)
	//		return true
	//	}
	//	return false
	//})
	//btnTheme.OnClick(func(ev clui.Event) {
	//	btnTheme.SetEnabled(false)
	//	tp := checkBox.State()
	//	changeTheme(logBox, btnTheme, tp)
	//})
	//btnSet.OnClick(func(ev clui.Event) {
	//	v := edit.Title()
	//	logBox.AddItem(fmt.Sprintf("New ProgressBar value: %v", v))
	//	logBox.SelectItem(logBox.ItemCount() - 1)
	//	updateProgress(v, pb)
	//})
	//btnStep.OnClick(func(ev clui.Event) {
	//	go pb.Step()
	//	logBox.AddItem("ProgressBar step")
	//	logBox.SelectItem(logBox.ItemCount() - 1)
	//	clui.PutEvent(clui.Event{Type: clui.EventRedraw})
	//})
	//btnQuit.OnClick(func(ev clui.Event) {
	//	go clui.Stop()
	//})
}

func (ui *UI) runEvery(duration time.Duration, function func() error) {
	for {
		_ = function()
		time.Sleep(duration)
	}
}
