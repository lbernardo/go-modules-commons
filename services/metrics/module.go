package metrics

import "go.uber.org/fx"

var Module = fx.Module("gmc.metrics",
	fx.Provide(NewRegister),
	fx.Invoke(NewServer))
