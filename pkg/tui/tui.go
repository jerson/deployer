package tui

import (
	"bytes"
	"io"
	"time"

	"github.com/jerson/deployer/pkg/deployer"
	"github.com/jerson/deployer/pkg/formatter"
	"github.com/jerson/deployer/pkg/tui/binding"
	"github.com/jerson/deployer/pkg/tui/common"
	"github.com/jerson/deployer/pkg/tui/focus"
	"github.com/jerson/deployer/pkg/tui/state"
	"github.com/jerson/deployer/pkg/tui/utils"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/sirupsen/logrus"
)

// Config ...
type Config struct {
	Name string
}

// GUI ...
type GUI struct {
	config  Config
	g       *gocui.Gui
	cli     *deployer.Deployer
	state   *state.State
	output  *bytes.Buffer
	log     *logrus.Logger
	focus   *focus.Manager
	binding *binding.Manager
	debug   bool
}

// NewGUI ...
func NewGUI(cli *deployer.Deployer, config Config, opts ...Option) *GUI {

	instance := &GUI{
		g:       nil,
		cli:     cli,
		config:  config,
		state:   state.NewState(),
		output:  nil,
		log:     nil,
		focus:   nil,
		binding: nil,
		debug:   false,
	}

	for _, opt := range opts {
		opt(instance)
	}

	if instance.output == nil {
		instance.output = &bytes.Buffer{}
	}

	if instance.log == nil {
		instance.log = common.CustomLogger(instance.output, &formatter.Console{
			FieldsOrder: []string{
				common.LogFieldModule,
				common.LogFieldWidget,
			},
			IgnoreNotFoundFields: false,
		})
	}

	return instance
}

// Render ...
func (gui *GUI) Render() (err error) {

	g, err := gocui.NewGui(gocui.OutputNormal, false, gui.log.WithField(common.LogFieldModule, "gocui"))
	if err != nil {
		return err
	}
	defer g.Close()

	mainViewName := "main"
	gui.g = g
	gui.focus = focus.NewManager(
		gui.g,
		gui.log.WithField(common.LogFieldModule, "focus"),
		mainViewName,
		gui.isModal,
	)
	gui.binding = binding.NewManager(
		gui.g,
		gui.log.WithField(common.LogFieldModule, "binding"),
		gui.isModal,
	)

	g.Mouse = true

	err = gui.setColorScheme()
	if err != nil {
		return err
	}

	g.SetManager(gocui.ManagerFunc(gui.layout), gocui.ManagerFunc(gui.focus.Layout()))

	err = gui.registerBindings()
	if err != nil {
		return err
	}

	go func() {
		// this delay avoid invalid view errors
		time.Sleep(time.Millisecond * 10)
		err := gui.init()
		if err != nil {
			panic(err)
		}
	}()
	err = g.MainLoop()
	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}
	return nil
}

func (gui *GUI) init() error {
	err := gui.initInternal()
	if err != nil {
		return err
	}
	err = gui.initProject()
	if err != nil {
		return err
	}
	go utils.RunEvery(time.Second*1, gui.updateEngine)
	go utils.RunEvery(time.Second*3, gui.updateDeployments)
	go utils.RunEvery(time.Second*1, gui.updateHealth)
	go utils.RunEvery(time.Second*1, gui.updateOptions)
	go utils.RunEvery(time.Millisecond*500, gui.updateMain)

	return nil
}

func (gui *GUI) initInternal() error {

	view, err := gui.getDebugView()
	if err != nil {
		return err
	}
	view.Autoscroll = true

	_, err = io.Copy(view, gui.output)
	if err != nil {
		return err
	}

	gui.log.SetOutput(io.MultiWriter(view, gui.output))

	//mainView, err := gui.getMainView()
	//if err != nil {
	//	return err
	//}
	//mainView.Wrap = true

	return nil
}

func (gui *GUI) setColorScheme() error {
	gui.g.FgColor = gocui.ColorDefault
	gui.g.SelFgColor = gocui.ColorGreen
	return nil
}

func (gui *GUI) isModal(name string) bool {

	modalViewNames := []string{
		"confirmation",
		"engine_modal",
		"engine_modal_start",
		"engine_modal_stop",
		"error_modal",
		"options_modal",
	}
	for _, viewName := range modalViewNames {
		if viewName == name {
			return true
		}
	}
	return false
}
