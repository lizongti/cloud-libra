package component

type Component interface {
	OnInit()
	OnStart()
	OnStop()
	OnClose()
}

type ComponentBase struct{}

func (*ComponentBase) OnInit() {}

func (*ComponentBase) OnStart() {}

func (*ComponentBase) OnStop() {}

func (*ComponentBase) OnClose() {}
