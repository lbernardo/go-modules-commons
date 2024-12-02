package apiserver

import "go.uber.org/fx"

var Module = fx.Module("gmc.api_server",
	fx.Provide(NewServer),
	fx.Invoke(Health),
	fx.Invoke(Start))
