package tui

func (gui *GUI) updateHealth() error {

	state := gui.state.Background
	//log := state.Logger.Log
	//
	//view, err := gui.getDebugView()
	//if err != nil {
	//	return err
	//}
	//
	//log.SetOutput(io.MultiWriter(view, state.Logger.Output))
	//state.Ctx = context.WithValue(state.Ctx, deployer.ContextKeyLog, log.WithFields(nil))

	return gui.cli.HealthAll(state.Ctx)
}
