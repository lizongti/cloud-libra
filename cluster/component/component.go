package component

type Component interface {
	OnInit()
	OnStart()
	OnStop()
	OnClose()
	String() string
}

type ComponentBase struct{}

func (*ComponentBase) OnInit() {}

func (*ComponentBase) OnStart() {}

func (*ComponentBase) OnStop() {}

func (*ComponentBase) OnClose() {}

func (*ComponentBase) String() string {
	return ""
}
